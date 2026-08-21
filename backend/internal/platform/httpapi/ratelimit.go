package httpapi

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mutex  sync.Mutex
	limit  int
	window time.Duration
	now    func() time.Time
	visits map[string]visit
}

type visit struct {
	started time.Time
	count   int
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{limit: limit, window: window, now: time.Now, visits: make(map[string]visit)}
}

func (limiter *rateLimiter) allow(request *http.Request) bool {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		host = request.RemoteAddr
	}
	now := limiter.now()
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()
	current := limiter.visits[host]
	if current.started.IsZero() || now.Sub(current.started) >= limiter.window {
		limiter.visits[host] = visit{started: now, count: 1}
		return true
	}
	if current.count >= limiter.limit {
		return false
	}
	current.count++
	limiter.visits[host] = current
	return true
}
