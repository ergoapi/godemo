package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ysicing/ext/e"
	"github.com/ysicing/ext/ginmid"
	"github.com/ysicing/ext/httputil"
	"github.com/ysicing/ext/logger"
)

var (
	reqCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "default_req_counter",
		Help: "访问统计",
	})
)

func init() {
	logcfg := logger.Config{Simple: true, ConsoleOnly: true}
	logger.InitLogger(&logcfg)
	prometheus.MustRegister(reqCounter)
}

func prom() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCounter.Inc()
		c.Next()
	}
}

func main() {
	gin.SetMode(gin.DebugMode)
	gin.DisableConsoleColor()
	r := gin.New()
	r.Use(ginmid.RequestID(), ginmid.Log(), prom())

	// Example ping request.
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Example / request.
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, e.Done(map[string]interface{}{
			"message": "power by godemo",
			"id":      ginmid.GetRequestID(c),
			"method":  c.Request.Method,
			"url":     c.Request.Host,
			"client":  c.ClientIP(),
			"ua":      c.Request.UserAgent(),
		}))
	})

	// Example /metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	addr := ":8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go start2svc(r)
	go func() {
		logger.Slog.Infof("http listen to %v, pid is %v", addr, os.Getpid())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Slog.Fatal(err)
		}
	}()
	httputil.SetupGracefulStop(srv)
}

func start2svc(e *gin.Engine) {
	addr := ":8081"
	srv := &http.Server{
		Addr:    addr,
		Handler: e,
	}
	go func() {
		logger.Slog.Infof("http listen to %v, pid is %v", addr, os.Getpid())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Slog.Fatal(err)
		}
	}()
	httputil.SetupGracefulStop(srv)
}
