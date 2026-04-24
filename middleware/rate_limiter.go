package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	r        rate.Limit
	burst    int
}

func newRateLimiterStore(r rate.Limit, burst int) *RateLimiterStore {
	store := &RateLimiterStore{
		limiters: make(map[string]*ipLimiter),
		r:        r,
		burst:    burst,
	}
	// Goroutine de limpeza: remove IPs inativos a cada 5 minutos
	go store.cleanup(5 * time.Minute)
	return store
}

func (s *RateLimiterStore) get(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, exists := s.limiters[ip]; exists {
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	l := rate.NewLimiter(s.r, s.burst)
	s.limiters[ip] = &ipLimiter{limiter: l, lastSeen: time.Now()}
	return l
}

func (s *RateLimiterStore) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		for ip, entry := range s.limiters {
			if time.Since(entry.lastSeen) > interval {
				delete(s.limiters, ip)
			}
		}
		s.mu.Unlock()
	}
}

// globalStore: 60 requisições por minuto por IP (1 req/s com burst de 10)
var globalStore = newRateLimiterStore(rate.Every(time.Second), 10)

// authStore: 10 requisições por minuto por IP em rotas de autenticação
var authStore = newRateLimiterStore(rate.Every(6*time.Second), 5)

// RateLimiter aplica limite genérico de 60 req/min por IP (burst 10)
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !globalStore.get(ip).Allow() {
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "muitas requisições, tente novamente em instantes",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// StrictRateLimiter aplica limite restrito de 10 req/min por IP (burst 5)
// Usado nas rotas de autenticação para prevenir brute force
func StrictRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !authStore.get(ip).Allow() {
			c.Header("Retry-After", "6")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "muitas tentativas de autenticação, aguarde alguns segundos",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
