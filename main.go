package main

import (
	"example/web-service-gin/api"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// Simple Gin Server
	fmt.Println("Starting gin server..")
	router := gin.Default()
	// router.GET("/order", api.OrderSimple)
	router.GET("/catalog", api.ListCatalog)
	router.POST("/catalog/update", api.UpdateBatchCatalogObject)

	router.Run("localhost:8080")

}
