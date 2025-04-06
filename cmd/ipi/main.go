package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kofj/ipi/pkg/handlers"
	"github.com/kofj/ipi/pkg/ipdb"
	"github.com/sirupsen/logrus"
)

var router = gin.New()

func init() {
	router.LoadHTMLGlob("templates/*")

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Initialize the router
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

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
