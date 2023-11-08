package handlers

import (
	"context"
	"net/http"

	"github.com/Jason-CKY/htmx-todo-app/pkg/components"
	"github.com/Jason-CKY/htmx-todo-app/pkg/core"
	"github.com/Jason-CKY/htmx-todo-app/pkg/schemas"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func HomePage(c echo.Context) error {
	component := components.HomePage(4)
	return component.Render(context.Background(), c.Response().Writer)
}

func TasksView(c echo.Context) error {
	backlogTaskList, progressTaskList, doneTaskList, err := core.GetTasksInOrder()
	if err != nil {
		return err
	}

	component := components.TaskView(backlogTaskList, progressTaskList, doneTaskList)
	return component.Render(context.Background(), c.Response().Writer)
}

func DeleteTaskView(c echo.Context) error {
	task_id := c.Param("id")
	err := core.DeleteTaskById(task_id)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "")
}

func EditTaskView(c echo.Context) error {
	task_id := c.Param("id")
	task, err := core.GetTaskById(task_id)
	if err != nil {
		return err
	}
	component := components.EditTask(task)
	return component.Render(context.Background(), c.Response().Writer)
}

func EmptyEditTaskView(c echo.Context) error {
	task_status := c.Param("status")
	task := schemas.Task{
		Id:          uuid.New().String(),
		Status:      task_status,
		Title:       "",
		Description: "",
	}
	component := components.EditTask(task)
	return component.Render(context.Background(), c.Response().Writer)
}

func CancelEditTaskView(c echo.Context) error {
	task_id := c.Param("id")
	task, err := core.GetTaskById(task_id)
	if err != nil {
		if err.Code == http.StatusNotFound {
			return c.String(http.StatusOK, "")
		}
		return err
	}
	component := components.TaskSingleton(task)
	return component.Render(context.Background(), c.Response().Writer)
}

func UpdateTaskView(c echo.Context) error {
	task_id := c.Param("id")
	status := c.FormValue("status")
	title := c.FormValue("title")
	description := c.FormValue("description")
	new_task := schemas.Task{
		Id:          task_id,
		Title:       title,
		Description: description,
		Status:      status,
	}
	task, err := core.UpsertTask(new_task)
	if err != nil {
		return err
	}
	component := components.TaskSingleton(task)
	return component.Render(context.Background(), c.Response().Writer)
}
