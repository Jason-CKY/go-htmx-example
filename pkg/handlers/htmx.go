package handlers

import (
	"context"
	"net/http"

	"github.com/Jason-CKY/htmx-todo-app/pkg/components"
	"github.com/Jason-CKY/htmx-todo-app/pkg/core"
	"github.com/labstack/echo/v4"
)

func HomePage(c echo.Context) error {
	component := components.HomePage(4)
	return component.Render(context.Background(), c.Response().Writer)
}

func TasksView(c echo.Context) error {
	backlogTaskList, progressTaskList, doneTaskList, err := core.GetTasks()
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
