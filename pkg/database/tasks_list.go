package database

import (
	_ "embed"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

//go:embed tasks_list.json
var tasksList []byte

var Tasks = newStore(jsonLoader[schemas.TaskFullSchema](tasksList), func(task *schemas.TaskFullSchema) string {
	return task.Code
})

var taskSkills = sync.OnceValue(func() []string {
	seen := make(map[string]struct{})
	var skills []string

	for _, task := range Tasks.All() {
		if task.Skill == nil {
			continue
		}
		_, exists := seen[*task.Skill]
		if exists {
			continue
		}
		seen[*task.Skill] = struct{}{}
		skills = append(skills, *task.Skill)
	}
	return skills
})

func TaskSkills() []string {
	return taskSkills()
}
