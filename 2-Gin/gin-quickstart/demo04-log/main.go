package main

import (
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func SlogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		requestID, _ := c.Get("request_id")
		logger.Info("request",
			slog.String("request_id", requestID.(string)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
		)
	}
}

func main() {
	// 写入文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Disable log's color
	gin.DisableConsoleColor()

	// Force log's color
	gin.ForceConsoleColor()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	// SkipQueryString indicates that the logger should not log the query string.
	// For example, /path?q=1 will be logged as /path
	// loggerConfig := gin.LoggerConfig{SkipQueryString: true, Formatter: func(param gin.LogFormatterParams) string {
	// 	// your custom format
	// 	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
	// 		param.ClientIP,
	// 		param.TimeStamp.Format(time.RFC1123),
	// 		param.Method,
	// 		param.Path,
	// 		param.Request.Proto,
	// 		param.StatusCode,
	// 		param.Latency,
	// 		param.Request.UserAgent(),
	// 		param.ErrorMessage,
	// 	)
	// }}

	r := gin.New()

	r.Use(RequestIDMiddleware())
	r.Use(SlogMiddleware(logger))
	// router.Use(gin.LoggerWithConfig(loggerConfig))
	r.GET("/ping", func(c *gin.Context) {
		// c.SetCookie()
		c.SetCookieData(&http.Cookie{
			
		})
		c.String(200, "pong")
	})

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// log.Fatal(autotls.Run(r, "example1.com", "example2.com"))
	s.ListenAndServe()
	// r.Run(":8080")
}
