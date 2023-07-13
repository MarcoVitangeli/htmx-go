package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())

	r.LoadHTMLGlob("./html/*.html")
	r.StaticFile("css/pico.min.css", "./pico-1.5.9/css/pico.min.css")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"timestamp": time.Now().UTC().String(),
		})
	})

	r.Run(":8080")
}
