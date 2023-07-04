package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/exgin"
	"github.com/ergoapi/util/exhttp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	g := exgin.Init(&exgin.Config{
		Debug:   true,
		Cors:    true,
		Metrics: true,
	})
	g.Use(exgin.ExLog(), exgin.ExRecovery(), exgin.ExTraceID())

	// Example ping request.
	g.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Example / request.
	g.GET("/", func(c *gin.Context) {
		exgin.GinsData(c, map[string]interface{}{
			"message": "power by godemo",
			"method":  c.Request.Method,
			"url":     c.Request.Host,
			"client":  c.ClientIP(),
			"ua":      c.Request.UserAgent(),
		}, nil)
	})

	g.GET("/:id", func(c *gin.Context) {
		id := exgin.GinsParamStr(c, "id")
		exgin.GinsData(c, map[string]interface{}{
			"message": fmt.Sprintf("power by %s", id),
			"method":  c.Request.Method,
			"url":     c.Request.Host,
			"ns":      environ.GetEnv("NAMESPACE", "default"),
			"client":  c.ClientIP(),
			"ua":      c.Request.UserAgent(),
		}, nil)
	})

	addr := ":8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: g,
	}

	go func() {
		logrus.Infof("http listen to %v, pid is %v", addr, os.Getpid())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatal(err)
		}
	}()
	exhttp.SetupGracefulStop(srv)
}
