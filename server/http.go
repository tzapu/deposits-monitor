package server

import (
	"html/template"
	"net/http"

	"github.com/tzapu/deposits-monitor/importer"

	"github.com/sirupsen/logrus"

	"github.com/tzapu/deposits-monitor/helper"

	"github.com/gin-gonic/gin"
)

var log = logrus.WithField("module", "server")

func Serve(imp *importer.Importer) {
	hub := newHub()
	go hub.run()

	watcher := NewWatcher(hub)
	go func() {
		err := watcher.Watch("./web/assets", "./web/templates")
		helper.FatalIfError(err, "watcher")
	}()

	r := gin.New()
	r.Use(gin.Logger()) // 25%-50%  extra performace if we disable this
	r.Use(gin.Recovery())

	r.SetFuncMap(template.FuncMap{
		"formatDate":   importer.FormatDate,
		"formatStart":  importer.FormatStart,
		"formatEnd":    importer.FormatEnd,
		"formatMiddle": importer.FormatMiddle,
	})

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
		transfers := imp.TransfersList()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":     "Deposits Monitor",
			"Transfers": transfers,
		})
	})

	err := r.Run()
	helper.FatalIfError(err, "gin run")
}
