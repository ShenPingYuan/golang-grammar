package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello,gin!")
	router := gin.Default()
	router.GET("/ping", get)
	router.POST("/ping", post)
	router.PUT("/ping", update)
	router.DELETE("/ping", delete)
	router.PATCH("/ping", patch)
	router.OPTIONS("/ping", options)
	router.HEAD("/ping", head)
	router.Run() // listen and serve on 0.0.0.0:8080 by default
}

func get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get!"})
}

func post(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "post!"})
}

func update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update"})
}

func delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "delete!"})
}

func patch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "patch!"})
}

func options(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "options"})
}

func head(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "head!"})
}
