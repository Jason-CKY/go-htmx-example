package handlers

import (
	"context"

	"github.com/Jason-CKY/htmx-todo-app/pkg/components"
	"github.com/Jason-CKY/htmx-todo-app/pkg/schemas"
	"github.com/labstack/echo/v4"
)

func HomePage(c echo.Context) error {
	component := components.HomePage(4)
	return component.Render(context.Background(), c.Response().Writer)
}

func TasksView(c echo.Context) error {
	backlogTaskList := []schemas.Task{
		{
			Id:          "1234",
			Title:       "This is a test",
			Description: "This is a test description",
			Status:      "backlog",
		},
		{
			Id:          "123674",
			Title:       "This is a test",
			Description: "This is a test description",
			Status:      "backlog",
		},
	}
	progressTaskList := []schemas.Task{
		{
			Id:          "5436",
			Title:       "This is a test",
			Description: "This is a test description",
			Status:      "progress",
		},
	}
	doneTaskList := []schemas.Task{
		{
			Id:          "7888",
			Title:       "This is a test",
			Description: "This is a test description",
			Status:      "done",
		},
	}
	component := components.TaskView(backlogTaskList, progressTaskList, doneTaskList)
	return component.Render(context.Background(), c.Response().Writer)
}
