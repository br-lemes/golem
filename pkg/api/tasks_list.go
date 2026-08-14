package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

const TasksListSize = 10000

func TasksList() ([]schemas.TaskFullSchema, error) {
	result := []schemas.TaskFullSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/tasks/list?page=%d&size=%d", page, TasksListSize), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageTaskFullSchema
		err = json.Unmarshal(resp, &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data.Data...)
		if page >= data.Pages {
			break
		}
		page++
	}
	return result, nil
}
