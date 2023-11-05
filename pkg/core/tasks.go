package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jason-CKY/htmx-todo-app/pkg/schemas"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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

func DeleteTaskById(task_id string) error {
	log.Debugf("Deleting task id: %v...", task_id)
	endpoint := fmt.Sprintf("%v/items/task/%v", DirectusHost, task_id)

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	client := &http.Client{}
	res, err := client.Do(req)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if res.StatusCode != 204 {
		return echo.NewHTTPError(res.StatusCode, string(body))
	}

	return nil
}
