package core

import (
	"github.com/Jason-CKY/htmx-todo-app/pkg/schemas"
	log "github.com/sirupsen/logrus"
)

func GetTasks() ([]schemas.Task, []schemas.Task, []schemas.Task) {
	// res, err := http.Get()
	log.Info(DirectusHost)
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
	return backlogTaskList, progressTaskList, doneTaskList
}
