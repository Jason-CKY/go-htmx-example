package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jason-CKY/htmx-todo-app/pkg/schemas"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func GetTasks() ([]schemas.Task, []schemas.Task, []schemas.Task, *echo.HTTPError) {
	backlogTaskList, progressTaskList, doneTaskList := []schemas.Task{}, []schemas.Task{}, []schemas.Task{}
	endpoint := fmt.Sprintf("%v/items/task", DirectusHost)
	res, err := http.Get(endpoint)
	// error handling for http request
	if err != nil {
		return backlogTaskList, progressTaskList, doneTaskList, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	// error handling for anything above 2xx response
	if res.StatusCode > 299 {
		return backlogTaskList, progressTaskList, doneTaskList, echo.NewHTTPError(res.StatusCode, string(body))
	}
	var tasksResponse map[string][]schemas.Task
	err = json.Unmarshal(body, &tasksResponse)
	// error handling for json unmarshaling
	if err != nil {
		return backlogTaskList, progressTaskList, doneTaskList, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

func GetTaskById(task_id string) (schemas.Task, *echo.HTTPError) {
	endpoint := fmt.Sprintf("%v/items/task?filter[id][_eq]=%v", DirectusHost, task_id)
	res, err := http.Get(endpoint)
	// error handling for http request
	if err != nil {
		return schemas.Task{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return schemas.Task{}, echo.NewHTTPError(res.StatusCode, string(body))
	}

	var taskResponse map[string][]schemas.Task
	err = json.Unmarshal(body, &taskResponse)
	// error handling for json unmarshaling
	if err != nil {
		return schemas.Task{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if len(taskResponse["data"]) == 0 {
		return schemas.Task{}, echo.NewHTTPError(http.StatusNotFound, "task not found")
	}

	return taskResponse["data"][0], nil
}

func DeleteTaskById(task_id string) *echo.HTTPError {
	log.Debugf("Deleting task id: %v...", task_id)
	endpoint := fmt.Sprintf("%v/items/task/%v", DirectusHost, task_id)

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 204 {
		return echo.NewHTTPError(res.StatusCode, string(body))
	}

	return nil
}

func UpdateTask(task schemas.Task) (schemas.Task, *echo.HTTPError) {
	endpoint := fmt.Sprintf("%v/items/task/%v", DirectusHost, task.Id)
	reqBody, _ := json.Marshal(task)
	req, err := http.NewRequest(http.MethodPatch, endpoint, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return task, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return task, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return task, echo.NewHTTPError(res.StatusCode, string(body))
	}
	var taskResponse map[string]schemas.Task
	err = json.Unmarshal(body, &taskResponse)
	// error handling for json unmarshaling
	if err != nil {
		return task, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return taskResponse["data"], nil
}

func CreateTask(task schemas.Task) (schemas.Task, *echo.HTTPError) {
	endpoint := fmt.Sprintf("%v/items/task", DirectusHost)
	reqBody, _ := json.Marshal(task)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return task, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return task, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return task, echo.NewHTTPError(res.StatusCode, string(body))
	}
	var taskResponse map[string]schemas.Task
	err = json.Unmarshal(body, &taskResponse)
	// error handling for json unmarshaling
	if err != nil {
		return task, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return taskResponse["data"], nil
}

func UpsertTask(task schemas.Task) (schemas.Task, *echo.HTTPError) {
	_, err := GetTaskById(task.Id)
	if err != nil {
		if err.Code == http.StatusNotFound {
			newTask, err := CreateTask(task)
			return newTask, err
		}
		return task, err
	}
	newTask, err := UpdateTask(task)
	return newTask, err
}
