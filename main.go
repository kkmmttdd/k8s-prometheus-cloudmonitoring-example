package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var SummaryOpts = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "go_app_response_latency_seconds",
		Help: "Response latency in seconds.",
	}, []string{"path"})


func observeMiddleware(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	timePassed := time.Since(startTime)
	SummaryOpts.WithLabelValues(c.Request.URL.Path).Observe(timePassed.Seconds())
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}


func main() {
	run()
}

func run() {
	r := gin.Default()
	r.Use(observeMiddleware)
	r.GET("/immediate", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "immediate",
		})
	})
	r.GET("/delay", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		c.JSON(200, gin.H{
			"message": "delay",
		})
	})
	r.GET("/metrics", prometheusHandler())
	_ = r.Run()
}
