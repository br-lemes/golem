package database

import (
	_ "embed"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

//go:embed tasks_list.json
var tasksList []byte

type taskCatalog struct {
	*store[schemas.TaskFullSchema, string]
}

func (c *taskCatalog) Items() []string {
	return taskItemCodes()
}

func (c *taskCatalog) Monsters() []string {
	return taskMonsterCodes()
}

func Tasks() *taskCatalog {
	return tasksCatalog
}

func (c *taskCatalog) Skills() []string {
	return taskSkills()
}

var tasksCatalog = func() *taskCatalog {
	return &taskCatalog{
		store: newStore(jsonLoader[schemas.TaskFullSchema](tasksList), func(task *schemas.TaskFullSchema) string {
			return task.Code
		}),
	}
}()

var taskItemCodes = sync.OnceValue(func() []string {
	return taskCodes(schemas.Items)
})

var taskMonsterCodes = sync.OnceValue(func() []string {
	return taskCodes(schemas.Monsters)
})

var taskSkills = sync.OnceValue(func() []string {
	seen := make(map[string]struct{})
	var skills []string

	for _, task := range Tasks().All() {
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

func taskCodes(taskType schemas.TaskType) []string {
	var codes []string
	for _, task := range Tasks().All() {
		if task.Type == taskType {
			codes = append(codes, task.Code)
		}
	}
	return codes
}
