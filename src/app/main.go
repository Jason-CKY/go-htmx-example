package main

import (
	"flag"
	"os"

	"github.com/Jason-CKY/htmx-todo-app/pkg/handlers"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var directusHost = ""

func LookupEnvOrString(key string, defaultValue string) string {
	envVariable, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return envVariable
}

func main() {
	flag.StringVar(&directusHost, "fpath", LookupEnvOrString("directus_host", directusHost), "Path to routing json file")

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

	e.Logger.Fatal(e.Start(":8080"))
	// router := gin.Default()
	// router.LoadHTMLGlob("templates/*")

	// router.GET("/", handlers.HomePage)
	// router.Static("/static", "./static")
	// s := &http.Server{
	// 	Addr:           ":8080",
	// 	Handler:        router,
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// s.ListenAndServe()
}
