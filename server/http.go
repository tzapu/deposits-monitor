package server

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/tzapu/deposits-monitor/helper"

	"github.com/gin-gonic/gin"
)

var log = logrus.WithField("module", "server")

func Serve() {
	hub := newHub()
	go hub.run()

	watcher := NewWatcher(hub)
	go func() {
		err := watcher.Watch("./web/assets", "./web/templates")
		helper.FatalIfError(err, "watcher")
	}()

	r := gin.Default()

	r.LoadHTMLGlob("./web/templates/*")

	r.Static("/assets", "./web/assets")

	r.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/reload", func(c *gin.Context) {
		hub.Broadcast([]byte(`{"type":"build_complete"}`))
		c.JSON(200, gin.H{
			"message": "broadcasted reload",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	err := r.Run()
	helper.FatalIfError(err, "gin run")
}
