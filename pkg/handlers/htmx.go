package handlers

import (
	"context"

	"github.com/Jason-CKY/htmx-todo-app/pkg/components"
	"github.com/Jason-CKY/htmx-todo-app/pkg/core"
	"github.com/labstack/echo/v4"
)

func HomePage(c echo.Context) error {
	component := components.HomePage(4)
	return component.Render(context.Background(), c.Response().Writer)
}

func TasksView(c echo.Context) error {
	backlogTaskList, progressTaskList, doneTaskList := core.GetTasks()

	component := components.TaskView(backlogTaskList, progressTaskList, doneTaskList)
	return component.Render(context.Background(), c.Response().Writer)
}
