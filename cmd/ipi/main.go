package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kofj/ipi/pkg/handlers"
	"github.com/kofj/ipi/pkg/ipdb"
	"github.com/kofj/ipi/pkg/otel"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var router = gin.New()

var formatter = func(p gin.LogFormatterParams) string {
	// nginx style log format
	return fmt.Sprintf(
		"%s - [%s] \"%s %s %s\" %d %s \"%s\" \"%d\" \"%s\" \"%s\" \n",
		p.ClientIP,
		p.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
		p.Method,
		p.Path,
		p.Request.Proto,
		p.StatusCode,
		p.Latency,
		p.Request.UserAgent(),
		p.BodySize,
		// otel trace ID
		p.Request.Header.Get("traceparent"),
		p.Request.Header.Get("tracestate"),
	)
}

func init() {
	router.LoadHTMLGlob("templates/*")

	var OTEL_OTLP_HTTP_ENDPOINT = os.Getenv("OTEL_OTLP_HTTP_ENDPOINT")
	var OTEL_OTLP_HTTP_HEADERS = os.Getenv("OTEL_OTLP_HTTP_AUTH_HEADER")
	var OTEL_OTLP_HTTP_STREAM_NAME = os.Getenv("OTEL_OTLP_HTTP_STREAM_NAME")
	var tp = otel.InitTracerHTTP(OTEL_OTLP_HTTP_ENDPOINT, map[string]string{
		"Authorization": OTEL_OTLP_HTTP_HEADERS,
		"stream-name":   OTEL_OTLP_HTTP_STREAM_NAME,
	})
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Println("Error shutting down tracer provider: ", err)
		}
	}()

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	var conf = gin.LoggerConfig{
		Formatter: formatter,
		Output:    os.Stdout,              // Log to standard output
		SkipPaths: []string{"/public/.*"}, // Skip logging for static files
	}

	// Initialize the router
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(conf))
	router.Use(otelgin.Middleware("ipi-server"))

	// Load the routes
	loadRoutes(router)
}

func loadRoutes(r *gin.Engine) {
	// Define your routes here
	r.GET("/", handlers.IpiPage)
	r.GET("/:ip", handlers.IpiPage)

	// assets static files
	r.Static("/public", "./public")

	// Add more routes as needed
}

func main() {
	var err = ipdb.Singleton()
	if err != nil {
		logrus.WithError(err).Error("ipdb init failed")
		return
	}
	// Start the server
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
