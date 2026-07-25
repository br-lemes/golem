package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/tidwall/gjson"
)

var (
	//go:embed tasks_list.json
	tasksList []byte

	tasksListList  []schemas.TaskFullSchema
	tasksListCache map[string]schemas.TaskFullSchema
	tasksCodesList []string
	tasksSkills    []string
)

func GetTasksList() []schemas.TaskFullSchema {
	initTasksListCache()
	return tasksListList
}

func GetTaskList(code string) (schemas.TaskFullSchema, bool) {
	initTasksListCache()
	task, exists := tasksListCache[code]
	return task, exists
}

func GetTasksCodes() []string {
	initTasksCodes()
	return tasksCodesList
}

func GetTasksSkills() []string {
	initTasksSkills()
	return tasksSkills
}

var initTasksListCache = sync.OnceFunc(func() {
	tasksListCache = make(map[string]schemas.TaskFullSchema)
	err := json.Unmarshal(tasksList, &tasksListList)
	if err != nil {
		panic("failed to unmarshal tasks list: " + err.Error())
	}
	for _, task := range tasksListList {
		tasksListCache[task.Code] = task
	}
})

var initTasksCodes = sync.OnceFunc(func() {
	gjson.GetBytes(tasksList, "#.code").ForEach(func(_, value gjson.Result) bool {
		tasksCodesList = append(tasksCodesList, value.String())
		return true
	})
})

var initTasksSkills = sync.OnceFunc(func() {
	gjson.GetBytes(tasksList, "#.skill").ForEach(func(_, value gjson.Result) bool {
		tasksSkills = append(tasksSkills, value.String())
		return true
	})
})
