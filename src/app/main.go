package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/Jason-CKY/htmx-todo-app/handlers"
	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handlers.HomePage)
	router.Static("/static", "./static")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
