package handlers

import (
	"context"

	"github.com/Jason-CKY/htmx-todo-app/pkg/views"
	"github.com/labstack/echo/v4"
)

func HomePage(c echo.Context) error {
	component := views.HomePage(5)
	return component.Render(context.Background(), c.Response().Writer)
}
