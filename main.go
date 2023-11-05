package main

import (
	"flag"
	"fmt"

	"github.com/Jason-CKY/htmx-todo-app/pkg/core"
	"github.com/Jason-CKY/htmx-todo-app/pkg/handlers"
	"github.com/Jason-CKY/htmx-todo-app/pkg/utils"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func main() {
	flag.StringVar(&core.DirectusHost, "fpath", utils.LookupEnvOrString("DIRECTUS_HOST", core.DirectusHost), "Path to routing json file")
	flag.IntVar(&core.WebPort, "port", utils.LookupEnvOrInt("PORT", core.WebPort), "Port for echo web server")

	flag.Parse()

	// setup logrus
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	log.Infof("connecting to directus at: %v", core.DirectusHost)

	e := echo.New()
	e.Static("/static", "static")
	e.GET("/", handlers.HomePage)
	e.GET("/htmx", handlers.TasksView)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", core.WebPort)))
}
