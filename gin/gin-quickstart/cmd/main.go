package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var users map[int]*User = make(map[int]*User)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	fmt.Println("hello,gin!")
	users[1] = &User{
		Id:   1,
		Name: "admin",
	}
	users[2] = &User{
		Id:   2,
		Name: "spy",
	}
	router := gin.Default()
	router.LoadHTMLGlob("../templates/*")
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUserById)
	router.GET("/users/:id/*action", doSomethingForUser)
	router.GET("/ping", get)
	router.GET("/index", getIndex)
	router.POST("/ping", post)
	router.PUT("/ping", update)
	router.DELETE("/ping", delete)
	router.PATCH("/ping", patch)
	router.OPTIONS("/ping", options)
	router.HEAD("/ping", head)
	router.GET("/", get)
	router.Run() // listen and serve on 0.0.0.0:8080 by default
}

func getUsers(c *gin.Context) {
	userIdStr := c.DefaultQuery("id", "0")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求url错误"})
		return
	}
	result := []User{}
	for _, user := range users {
		result = append(result, *user)
	}
	if userId != 0 {
		newResult := []User{}
		for _, user := range result {
			if user.Id == userId {
				newResult = append(newResult, user)
			}
		}
		c.JSON(http.StatusOK, newResult)
		return
	}
	c.JSON(http.StatusOK, result)
}

func getUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求url错误"})
		return
	}
	if user, ok := users[id]; ok {
		c.JSON(http.StatusOK, user)
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "未找到"})
		return
	}
}

func doSomethingForUser(c *gin.Context) {
	idStr := c.Param("id")
	action := c.Param("action")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求url错误",
			"action":  action,
		})
		return
	}
	if user, ok := users[id]; ok {
		c.JSON(http.StatusOK, user)
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "未找到", "action": action})
		return
	}
}

func get(c *gin.Context) {
	p := Person{
		Name: "hkpf1729",
		Age:  20,
	}
	c.JSON(http.StatusOK, p)
}

func getIndex(c *gin.Context) {
	p := Person{
		Name: "hkpf1729",
		Age:  20,
	}
	c.HTML(http.StatusOK, "index.html", p)
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
