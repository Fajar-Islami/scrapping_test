package main

import (
	"github.com/Fajar-Islami/scrapping_test/app"
	"github.com/Fajar-Islami/scrapping_test/infrastructure/container"
	"github.com/gin-gonic/gin"
)

func main() {
	cont := container.New(".env")
	serve := gin.Default()
	app.Start(cont, serve)
}
