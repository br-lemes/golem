package api

import "testing"

func TestTasksListPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/tasks/list", TasksListSize, func(client *Client) (int, error) {
		items, err := client.TasksList(TasksListOptions{})
		return len(items), err
	})
}
