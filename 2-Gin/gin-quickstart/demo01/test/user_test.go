package test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkLogin(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default() // 初始化路由和依赖
	r.GET("/users/:id", getUserById)

	// body := `{"email":"test@example.com","password":"123456"}`

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}
	})
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var users map[int]*User = make(map[int]*User)

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
