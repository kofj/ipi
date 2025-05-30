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

	// Initialize the router
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
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
