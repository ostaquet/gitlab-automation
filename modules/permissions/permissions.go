package permissions

import (
	"flag"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"strconv"
	"time"
)

type arguments struct {
	token string
	gid   int
}

// CommandPermissions start the subcommand "permissions" based on arguments received by the program.
// Stop the program in case of error.
func CommandPermissions() {
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
	myGID := flagSet.String("gid", "", "Gitlab group id")

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

	// Create new gitlab REST client
	git, err := gitlab.NewClient(args.token)
	if err != nil {
		return err
	}

	fmt.Print("Permissions audit in progress...")

	// Process the group recursively
	err = processGroup(args.gid, git, users)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Printf("Permissions audit on %v\n", time.Now().Format("2006-01-02"))

	// Show the results
	for _, itMember := range users {
		fmt.Printf("%v (ID: %v) has access to:\n", itMember.username, itMember.id)
		for itGroup := range itMember.groupsID {
			fmt.Printf(" - group %v (ID: %v)\n", groupsDict[itGroup].name, itGroup)
		}
		for itProject := range itMember.projectsID {
			fmt.Printf(" - project %v (ID: %v)\n", projectsDict[itProject].name, itProject)
		}
	}

	fmt.Println("END")

	return nil
}

func processGroup(gid int, git *gitlab.Client, users map[int]user) error {
	fmt.Print(".")

	// Get details about the current group
	group, _, err := git.Groups.GetGroup(gid, &gitlab.GetGroupOptions{})
	if err != nil {
		return err
	}

	// Add details in the dictionary to avoid multiple calls to the API
	addGroupDictEntry(group)

	// Store group's members
	err = considerGroupMembers(gid, git, users)
	if err != nil {
		return err
	}

	for _, project := range group.Projects {
		// Add details in the dictionary
		addProjectDictEntry(project)
		// Store project's members
		err = considerProjectMembers(project.ID, git, users)
		if err != nil {
			return err
		}
	}

	// List subgroups
	subgroups, _, err := git.Groups.ListSubgroups(gid, &gitlab.ListSubgroupsOptions{})

	// Process subgroups recursively
	for _, subgroup := range subgroups {
		err = processGroup(subgroup.ID, git, users)
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
		// Check the existence of the member in the list, add it if not found
		if _, found := users[projectMember.ID]; !found {
			users[projectMember.ID] = newUserFromProject(*projectMember)
		}

		// Store the info of the current group
		users[projectMember.ID].projectsID[pid] = true
	}

	return nil
}

func considerGroupMembers(gid int, git *gitlab.Client, users map[int]user) error {
	groupMembers, _, err := git.Groups.ListGroupMembers(gid, &gitlab.ListGroupMembersOptions{})
	if err != nil {
		return err
	}

	for _, groupMember := range groupMembers {
		// Check the existence of the member in the list, add it if not found
		if _, found := users[groupMember.ID]; !found {
			users[groupMember.ID] = newUserFromGroup(*groupMember)
		}

		// Store the info of the current group
		users[groupMember.ID].groupsID[gid] = true
	}

	return nil
}
