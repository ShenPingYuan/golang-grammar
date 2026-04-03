package main

import (
	"github.com/danielkov/gin-helmet/ginhelmet"
	"github.com/gin-gonic/gin"
)

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(ginhelmet.Default())
	r.Use(gin.Recovery())
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, "pang!")
	})
	r.GET("profile", gin.BasicAuth(gin.Accounts{
		"admin": "123",
	}), func(ctx *gin.Context) {
		user := ctx.MustGet(gin.AuthUserKey)
		ctx.JSON(200, user)
		ctx.Next()
	})
	r.Run()
}
