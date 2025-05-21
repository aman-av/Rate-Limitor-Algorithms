package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type RateLimiter struct {
	requests    map[string]int
	window      time.Duration
	maxRequests int
	mutex       sync.Mutex
}

func NewRateLimiter(window time.Duration, maxRequests int) *RateLimiter {
	rl := &RateLimiter{
		requests:    make(map[string]int),
		window:      window,
		maxRequests: maxRequests,
	}
	go rl.resetRequests()
	return rl
}

func (rl *RateLimiter) resetRequests() {
	for {
		time.Sleep(rl.window)
		rl.mutex.Lock()
		rl.requests = make(map[string]int)
		rl.mutex.Unlock()
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.requests[ip]++
	return rl.requests[ip] <= rl.maxRequests
}

// secureRandomString generates a random string of length n using crypto/rand
func secureRandomString(n int) string {
	letters := "123456789"
	result := make([]byte, n)

	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[num.Int64()]
	}

	return string(result)
}

func main() {
	limiter := NewRateLimiter(1*time.Minute, 1)
	for {
		time.Sleep(1 * time.Second) // Simulate a request every second
		ip := secureRandomString(4) // Simulate a random IP address
		if limiter.Allow(ip) {
			// Allow request
			fmt.Printf("Request for IP %s allowed\n", ip)
		} else {
			// Deny request
			fmt.Printf("Request for IP %s denied\n", ip)
		}
	}
}
