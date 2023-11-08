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

// returns backlogSortOrder, progressSortOrder, doneSortOrder, error
func GetTaskSort() ([]string, []string, []string, *echo.HTTPError) {
	endpoint := fmt.Sprintf("%v/items/task_sorting", DirectusHost)
	res, err := http.Get(endpoint)
	// error handling for http request
	if err != nil {
		return []string{}, []string{}, []string{}, echo.NewHTTPError(res.StatusCode, err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	// error handling for anything above 2xx response
	if res.StatusCode > 299 {
		return []string{}, []string{}, []string{}, echo.NewHTTPError(res.StatusCode, err.Error())
	}
	var httpResponse map[string][]schemas.TaskSort
	err = json.Unmarshal(body, &httpResponse)
	// error handling for json unmarshaling
	if err != nil {
		return []string{}, []string{}, []string{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	backlogSortOrder, progressSortOrder, doneSortOrder := []string{}, []string{}, []string{}
	for _, taskSort := range httpResponse["data"] {
		if taskSort.Status == "backlog" {
			backlogSortOrder = taskSort.Sorting_order
		} else if taskSort.Status == "progress" {
			progressSortOrder = taskSort.Sorting_order
		} else {
			doneSortOrder = taskSort.Sorting_order
		}
	}
	return backlogSortOrder, progressSortOrder, doneSortOrder, nil
}

func GetTaskSortByStatus(status string) (schemas.TaskSort, *echo.HTTPError) {
	endpoint := fmt.Sprintf("%v/items/task_sorting?filter[status][_eq]=%v", DirectusHost, status)
	res, err := http.Get(endpoint)
	// error handling for http request
	if err != nil {
		return schemas.TaskSort{}, echo.NewHTTPError(res.StatusCode, err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	// error handling for anything above 2xx response
	if res.StatusCode > 299 {
		return schemas.TaskSort{}, echo.NewHTTPError(res.StatusCode, err.Error())
	}
	var httpResponse map[string][]schemas.TaskSort
	err = json.Unmarshal(body, &httpResponse)
	// error handling for json unmarshaling
	if err != nil {
		return schemas.TaskSort{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if len(httpResponse["data"]) == 0 {
		return schemas.TaskSort{}, echo.NewHTTPError(http.StatusNotFound, "Item not found")
	}

	return httpResponse["data"][0], nil
}

func UpdateTaskSort(taskSort schemas.TaskSort) (schemas.TaskSort, *echo.HTTPError) {
	endpoint := fmt.Sprintf("%v/items/task_sorting/%v", DirectusHost, taskSort.Id)
	reqBody, _ := json.Marshal(taskSort)
	req, err := http.NewRequest(http.MethodPatch, endpoint, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return taskSort, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return taskSort, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return taskSort, echo.NewHTTPError(res.StatusCode, string(body))
	}
	var taskResponse map[string]schemas.TaskSort
	err = json.Unmarshal(body, &taskResponse)
	// error handling for json unmarshaling
	if err != nil {
		return taskSort, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return taskResponse["data"], nil
}

func UpdateTaskSortByTasks(status string, tasks []schemas.Task) (schemas.TaskSort, *echo.HTTPError) {
	taskSort, httpErr := GetTaskSortByStatus(status)
	if httpErr != nil {
		return taskSort, httpErr
	}
	sortOrder := []string{}
	for _, task := range tasks {
		sortOrder = append(sortOrder, task.Id)
	}
	taskSort.Sorting_order = sortOrder
	updatedTaskSort, httpErr := UpdateTaskSort(taskSort)
	return updatedTaskSort, httpErr
}

func GetTasks() ([]schemas.Task, *echo.HTTPError) {
	endpoint := fmt.Sprintf("%v/items/task", DirectusHost)
	res, err := http.Get(endpoint)
	// error handling for http request
	if err != nil {
		return []schemas.Task{}, echo.NewHTTPError(res.StatusCode, err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	// error handling for anything above 2xx response
	if res.StatusCode > 299 {
		return []schemas.Task{}, echo.NewHTTPError(res.StatusCode, string(body))
	}
	var tasksResponse map[string][]schemas.Task
	err = json.Unmarshal(body, &tasksResponse)
	// error handling for json unmarshaling
	if err != nil {
		return []schemas.Task{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return tasksResponse["data"], nil
}

func FilterTaskById(task_id string, tasks []schemas.Task) schemas.Task {
	var filteredTask schemas.Task

	for _, task := range tasks {
		if task.Id == task_id {
			return task
		}
	}

	return filteredTask
}

func GetTasksInOrder() ([]schemas.Task, []schemas.Task, []schemas.Task, *echo.HTTPError) {
	tasks, err := GetTasks()
	if err != nil {
		return []schemas.Task{}, []schemas.Task{}, []schemas.Task{}, err
	}
	backlogTaskSort, progressTaskSort, doneTaskSort, err := GetTaskSort()
	if err != nil {
		return []schemas.Task{}, []schemas.Task{}, []schemas.Task{}, err
	}

	backlogTasks, progressTasks, doneTasks := []schemas.Task{}, []schemas.Task{}, []schemas.Task{}

	if len(backlogTaskSort)+len(progressTaskSort)+len(doneTaskSort) != len(tasks) {
		log.Error("Task Sorting not the same length as the total number of tasks!")
		for _, task := range tasks {
			if task.Status == "backlog" {
				backlogTasks = append(backlogTasks, task)
			} else if task.Status == "progress" {
				progressTasks = append(progressTasks, task)
			} else {
				doneTasks = append(doneTasks, task)
			}
		}
		UpdateTaskSortByTasks("backlog", backlogTasks)
		UpdateTaskSortByTasks("progress", progressTasks)
		UpdateTaskSortByTasks("done", doneTasks)
		return backlogTasks, progressTasks, doneTasks, nil

	}

	for _, taskId := range backlogTaskSort {
		backlogTasks = append(backlogTasks, FilterTaskById(taskId, tasks))
	}
	for _, taskId := range progressTaskSort {
		progressTasks = append(progressTasks, FilterTaskById(taskId, tasks))
	}
	for _, taskId := range doneTaskSort {
		doneTasks = append(doneTasks, FilterTaskById(taskId, tasks))
	}
	return backlogTasks, progressTasks, doneTasks, nil
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
	task, httpErr := GetTaskById(task_id)
	if httpErr != nil {
		return httpErr
	}

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

	taskSort, httpErr := GetTaskSortByStatus(task.Status)
	sortOrder := []string{}
	if httpErr != nil {
		return httpErr
	}
	for _, taskId := range taskSort.Sorting_order {
		if taskId != task_id {
			sortOrder = append(sortOrder, taskId)
		}
	}
	taskSort.Sorting_order = sortOrder
	_, httpErr = UpdateTaskSort(taskSort)
	if httpErr != nil {
		return httpErr
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

func UpdateTasksStatusById(task_ids []string, status string) ([]schemas.Task, *echo.HTTPError) {
	endpoint := fmt.Sprintf("%v/items/task", DirectusHost)
	data := map[string]interface{}{
		"keys": task_ids,
		"data": map[string]interface{}{
			"status": status,
		},
	}
	reqBody, err := json.Marshal(data)
	if err != nil {
		return []schemas.Task{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	req, err := http.NewRequest(http.MethodPatch, endpoint, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return []schemas.Task{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []schemas.Task{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return []schemas.Task{}, echo.NewHTTPError(res.StatusCode, string(body))
	}
	var taskResponse map[string][]schemas.Task
	err = json.Unmarshal(body, &taskResponse)
	// error handling for json unmarshaling
	if err != nil {
		return []schemas.Task{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	taskMapping := map[string]schemas.Task{}
	for _, task := range taskResponse["data"] {
		taskMapping[task.Id] = task
	}
	updatedTasks := []schemas.Task{}
	for _, taskId := range task_ids {
		updatedTasks = append(updatedTasks, taskMapping[taskId])
	}
	return updatedTasks, nil

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

	taskSort, httpErr := GetTaskSortByStatus(task.Status)
	sortOrder := []string{task.Id}
	if httpErr != nil {
		return task, httpErr
	}
	taskSort.Sorting_order = append(sortOrder, taskSort.Sorting_order[:]...)
	_, httpErr = UpdateTaskSort(taskSort)

	if httpErr != nil {
		return task, httpErr
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
