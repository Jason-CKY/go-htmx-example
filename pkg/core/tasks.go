package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jason-CKY/htmx-todo-app/pkg/schemas"
	"github.com/labstack/echo/v4"
)

func GetTasks() ([]schemas.Task, []schemas.Task, []schemas.Task, error) {
	backlogTaskList, progressTaskList, doneTaskList := []schemas.Task{}, []schemas.Task{}, []schemas.Task{}
	endpoint := fmt.Sprintf("%v/items/task", DirectusHost)
	res, err := http.Get(endpoint)
	// error handling for http request
	if err != nil {
		return backlogTaskList, progressTaskList, doneTaskList, echo.NewHTTPError(500, err.Error())
	}
	body, _ := io.ReadAll(res.Body)
	// error handling for anything above 2xx response
	if res.StatusCode > 299 {
		return backlogTaskList, progressTaskList, doneTaskList, echo.NewHTTPError(res.StatusCode, string(body))
	}
	var tasksResponse map[string][]schemas.Task
	defer res.Body.Close()
	err = json.Unmarshal(body, &tasksResponse)
	// error handling for json unmarshaling
	if err != nil {
		return backlogTaskList, progressTaskList, doneTaskList, echo.NewHTTPError(500, err.Error())
	}
	for _, task := range tasksResponse["data"] {
		if task.Status == "backlog" {
			backlogTaskList = append(backlogTaskList, task)
		} else if task.Status == "progress" {
			progressTaskList = append(progressTaskList, task)
		} else {
			doneTaskList = append(doneTaskList, task)
		}
	}
	return backlogTaskList, progressTaskList, doneTaskList, nil
}
