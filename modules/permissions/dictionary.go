package permissions

import "github.com/xanzy/go-gitlab"

type dictEntry struct {
	id   int
	name string
}

var projectsDict = make(map[int]dictEntry)
var groupsDict = make(map[int]dictEntry)

func addGroupDictEntry(group *gitlab.Group) {
	groupsDict[group.ID] = dictEntry{
		id:   group.ID,
		name: group.Name,
	}
}

func addProjectDictEntry(project *gitlab.Project) {
	projectsDict[project.ID] = dictEntry{
		id:   project.ID,
		name: project.Name,
	}
}
