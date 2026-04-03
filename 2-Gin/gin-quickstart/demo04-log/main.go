package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 写入文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// Disable log's color
	gin.DisableConsoleColor()

	// Force log's color
	gin.ForceConsoleColor()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	// SkipQueryString indicates that the logger should not log the query string.
	// For example, /path?q=1 will be logged as /path
	loggerConfig := gin.LoggerConfig{SkipQueryString: true, Formatter: func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}}

	router := gin.Default()

	router.Use(gin.LoggerWithConfig(loggerConfig))
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(":8080")
}
