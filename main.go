package main

import (
	"flag"
	"fmt"

	"github.com/Jason-CKY/htmx-todo-app/pkg/handlers"
	"github.com/Jason-CKY/htmx-todo-app/pkg/utils"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var (
	directusHost = "http://localhost:8055"
	webPort      = 8080
)

func main() {
	flag.StringVar(&directusHost, "fpath", utils.LookupEnvOrString("DIRECTUS_HOST", directusHost), "Path to routing json file")
	flag.IntVar(&webPort, "port", utils.LookupEnvOrInt("PORT", webPort), "Port for echo web server")

	flag.Parse()

	// setup logrus
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	log.Infof("connecting to directus at: %v", directusHost)

	e := echo.New()
	e.Static("/static", "static")
	e.GET("/", handlers.HomePage)
	e.GET("/htmx", handlers.TasksView)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", webPort)))
}
