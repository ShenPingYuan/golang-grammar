package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        int    `json:"id" xml:"id"`
	Name      string `json:"name" xml:"name"`
	Passworld string `json:"password" xml:"password"`
	Age       int    `json:"age" xml:"age"`
}

type Response struct {
	Code    int `json:"code" xml:"code" form:"code"`
	Message string
	Data    any
}

type UserRequest struct {
	PageIndex int `json:"pageIndex" form:"pageIndex" xml:"pageIndex" binding:"required"`
	PageSize  int `json:"pageSize" form:"pageSize" xml:"pageSize" binding:"-"`
}

type PagedResult struct {
	List  []any `json:"list" xml:"list"`
	Total int   `json:"total"`
}

var users []User = []User{
	{
		ID:        1,
		Name:      "Admin",
		Passworld: "admin123",
		Age:       30,
	},
}

func main() {
	router := gin.Default()

	router.GET("/users", func(ctx *gin.Context) {
		var request UserRequest
		err := ctx.ShouldBindQuery(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(200, Response{
			Code:    http.StatusOK,
			Message: "",
			Data:    request,
		})
	})

	router.Run()
}
