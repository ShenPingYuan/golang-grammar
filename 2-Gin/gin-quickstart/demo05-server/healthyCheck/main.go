package main

import (
	"context"
	"database/sql"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var isReady atomic.Bool

func main() {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost/dbname?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()

	// Liveness: is the process alive?
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	// Readiness: can we serve traffic?
	r.GET("/readyz", func(c *gin.Context) {
		if !isReady.Load() {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}

		// Check database connectivity
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"reason": "database unreachable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// Mark as ready after initialization is complete
	isReady.Store(true)

	r.Run(":8080")
}

type HealthChecker struct {
	DB    *sql.DB
	Redis *redis.Client
}

func (h *HealthChecker) CheckHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	checks := gin.H{}
	healthy := true

	// Check database
	if err := h.DB.PingContext(ctx); err != nil {
		checks["database"] = gin.H{"status": "unhealthy", "error": err.Error()}
		healthy = false
	} else {
		checks["database"] = gin.H{"status": "healthy"}
	}

	// Check Redis
	if err := h.Redis.Ping(ctx).Err(); err != nil {
		checks["redis"] = gin.H{"status": "unhealthy", "error": err.Error()}
		healthy = false
	} else {
		checks["redis"] = gin.H{"status": "healthy"}
	}

	status := http.StatusOK
	if !healthy {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"status": map[bool]string{true: "healthy", false: "unhealthy"}[healthy],
		"checks": checks,
	})
}
