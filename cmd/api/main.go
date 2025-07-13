package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visiab/appointment-calculator/internal/infrastructure/web"
)

func main() {
	router := gin.Default()
	
	// Setup routes
	web.SetupRoutes(router)
	
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
