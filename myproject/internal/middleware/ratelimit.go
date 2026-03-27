package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimit 基于 IP 的简易限流中间件
func RateLimit(rps int) func(http.Handler) http.Handler {
	var mu sync.Mutex
	limiters := make(map[string]*rate.Limiter)

	getLimiter := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()
		if lim, ok := limiters[ip]; ok {
			return lim
		}
		lim := rate.NewLimiter(rate.Limit(rps), rps*2)
		limiters[ip] = lim
		return lim
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := getLimiter(r.RemoteAddr)
			if !limiter.Allow() {
				http.Error(w, `{"code":429,"message":"too many requests"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}