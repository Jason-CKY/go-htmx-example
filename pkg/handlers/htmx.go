package handlers

import (
	"context"

	"github.com/Jason-CKY/htmx-todo-app/pkg/components"
	"github.com/labstack/echo/v4"
)

func HomePage(c echo.Context) error {
	component := components.HomePage(4)
	return component.Render(context.Background(), c.Response().Writer)
}
