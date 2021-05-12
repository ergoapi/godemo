package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reqCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "default_req_counter",
		Help: "访问统计",
	})
)

func init() {
	prometheus.MustRegister(reqCounter)
}

func prom() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCounter.Inc()
		c.Next()
	}
}

func main() {
	r := gin.New()
	r.Use(prom())
	// Example ping request.
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Example / request.
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "power by godemo",
			"method":  c.Request.Method,
			"url":     c.Request.Host,
			"client":  c.ClientIP(),
			"ua":      c.Request.UserAgent(),
		})
	})

	// Example /metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.Run(":9090")
}
