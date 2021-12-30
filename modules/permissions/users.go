package permissions

import "github.com/xanzy/go-gitlab"

type user struct {
	id         int
	username   string
	name       string
	groupsID   map[int]bool
	projectsID map[int]bool
}

func newUser(userId int, username string, name string) user {
	return user{
		id:         userId,
		username:   username,
		name:       name,
		groupsID:   make(map[int]bool),
		projectsID: make(map[int]bool),
	}
}

func newUserFromGroup(member gitlab.GroupMember) user {
	return newUser(member.ID, member.Username, member.Name)
}

func newUserFromProject(member gitlab.ProjectMember) user {
	return newUser(member.ID, member.Username, member.Name)
}
