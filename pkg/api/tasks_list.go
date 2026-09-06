package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const TasksListSize = 10000

type TasksListOptions struct {
	MaxLevel int    `url:"max_level,omitempty"`
	MinLevel int    `url:"min_level,omitempty"`
	Skill    string `url:"skill,omitempty"`
	Type     string `url:"type,omitempty"`
}

func TasksList(options TasksListOptions) ([]schemas.TaskFullSchema, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	result := []schemas.TaskFullSchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(TasksListSize))
		path := fmt.Sprintf("/tasks/list?%s", params.Encode())
		resp, err := Get(path, nil)
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
