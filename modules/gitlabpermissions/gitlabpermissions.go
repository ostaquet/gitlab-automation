package gitlabpermissions

import (
	"flag"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"strconv"
)

type arguments struct {
	token string
	gid   int
}

type user struct {
	id       int
	username string
	name     string
	groups   map[int]int
	projects map[int]int
}

func CommandGitlabPermissions() {
	args, err := extractArgs()

	if err != nil {
		log.Fatal(err)
	}

	err = process(args)
	if err != nil {
		log.Fatal(err)
	}
}

func extractArgs() (args arguments, err error) {
	flagSet := flag.NewFlagSet("permissions", flag.ContinueOnError)

	myToken := flagSet.String("token", "", "Gitlab token")
	myGID := flagSet.String("gid", "", "Gitlab group id to check")

	err = flagSet.Parse(os.Args[2:])
	if err != nil {
		return
	}

	if *myToken == "" {
		flagSet.Usage()
		err = fmt.Errorf("mandatory Gitlab token required")
		return
	}
	args.token = *myToken

	if *myGID == "" {
		flagSet.Usage()
		err = fmt.Errorf("mandatory Gitlab group identifier required")
		return
	}

	args.gid, err = strconv.Atoi(*myGID)

	return
}

func process(args arguments) (err error) {
	users := make(map[int]user)
	projectsDict := make(map[int]string)
	groupsDict := make(map[int]string)

	git, err := gitlab.NewClient(args.token)
	if err != nil {
		return err
	}

	err = processGroup(args.gid, git, users, projectsDict, groupsDict)
	if err != nil {
		return err
	}

	for _, iUser := range users {
		fmt.Printf("%v (ID: %v) has access to:\n", iUser.username, iUser.id)
		for _, iGroup := range iUser.groups {
			fmt.Printf(" - group %v (ID: %v)\n", groupsDict[iGroup], iGroup)
		}
		for _, iProject := range iUser.projects {
			fmt.Printf(" - project %v (ID: %v)\n", projectsDict[iProject], iProject)
		}
	}

	return nil
}

func processGroup(gid int, git *gitlab.Client, users map[int]user, projectsDict map[int]string, groupsDict map[int]string) error {
	group, _, err := git.Groups.GetGroup(gid, &gitlab.GetGroupOptions{})
	if err != nil {
		return err
	}

	groupsDict[group.ID] = group.Name

	err = considerGroupMembers(gid, git, users)
	if err != nil {
		return err
	}

	for _, project := range group.Projects {
		projectsDict[project.ID] = project.Name

		err = considerProjectMembers(project.ID, git, users)
		if err != nil {
			return err
		}
	}

	subgroups, _, err := git.Groups.ListSubgroups(gid, &gitlab.ListSubgroupsOptions{})

	for _, subgroup := range subgroups {
		err = processGroup(subgroup.ID, git, users, projectsDict, groupsDict)
		if err != nil {
			return err
		}
	}

	return nil
}

func considerProjectMembers(pid int, git *gitlab.Client, users map[int]user) error {
	projectMembers, _, err := git.ProjectMembers.ListProjectMembers(pid, &gitlab.ListProjectMembersOptions{})
	if err != nil {
		return err
	}

	for _, projectMember := range projectMembers {
		// Check the existence of the member in the list
		_, ok := users[projectMember.ID]
		if !ok {
			users[projectMember.ID] = user{
				id:       projectMember.ID,
				username: projectMember.Username,
				name:     projectMember.Name,
				groups:   make(map[int]int),
				projects: make(map[int]int),
			}
		}

		// Store the info of the current group
		users[projectMember.ID].projects[pid] = pid
	}

	return nil
}

func considerGroupMembers(gid int, git *gitlab.Client, users map[int]user) error {
	groupMembers, _, err := git.Groups.ListGroupMembers(gid, &gitlab.ListGroupMembersOptions{})
	if err != nil {
		return err
	}

	for _, groupMember := range groupMembers {
		// Check the existence of the member in the list
		_, ok := users[groupMember.ID]
		if !ok {
			users[groupMember.ID] = user{
				id:       groupMember.ID,
				username: groupMember.Username,
				name:     groupMember.Name,
				groups:   make(map[int]int),
				projects: make(map[int]int),
			}
		}

		// Store the info of the current group
		users[groupMember.ID].groups[gid] = gid
	}

	return nil
}
