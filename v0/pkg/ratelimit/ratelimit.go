// Package ratelimit is a small token-bucket rate limiter for dreego, as
// per-route or app-wide middleware.
//
//	limiter := ratelimit.New(ratelimit.Config{Burst: 20, Every: time.Minute})
//	app.Use(limiter.Middleware)          // whole app
//	// or per route:
//	Page{ …, Middleware: []func(http.Handler) http.Handler{limiter.Middleware} }
//
// State lives in memory, keyed by client IP (and an optional discriminator).
// That is enough for a single process; a multi-instance deployment needs a
// shared store (e.g. Redis) and is out of scope here.
package ratelimit

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// Config describes the limit: Burst requests per Every, refilled continuously.
type Config struct {
	// Burst is the bucket size — how many requests may pass at once.
	Burst int
	// Every is the window the burst refills over.
	Every time.Duration
	// Key returns the bucket key for a request. Default: client IP.
	Key func(r *http.Request) string
	// TrustedProxies are addresses allowed to set the real client IP via
	// X-Forwarded-For / X-Real-IP. Without them the socket IP is used, so an
	// attacker cannot spoof the key.
	TrustedProxies []string
}

// bucket is one client's token bucket.
type bucket struct {
	tokens   float64
	lastFill time.Time
}

// Limiter holds the buckets. Safe for concurrent use.
type Limiter struct {
	cfg        Config
	ratePerSec float64

	mu      sync.Mutex
	buckets map[string]*bucket

	trusted map[string]bool
}

// New builds a limiter. Burst <= 0 or Every <= 0 disables limiting (the
// middleware then just passes everything through).
func New(cfg Config) *Limiter {
	limiter := &Limiter{
		cfg:     cfg,
		buckets: map[string]*bucket{},
		trusted: map[string]bool{},
	}
	if cfg.Every > 0 {
		limiter.ratePerSec = float64(cfg.Burst) / cfg.Every.Seconds()
	}
	for _, adresse := range cfg.TrustedProxies {
		limiter.trusted[adresse] = true
	}
	return limiter
}

// Middleware enforces the limit. Over-limit requests get 429 with Retry-After.
func (limiter *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.cfg.Burst <= 0 || limiter.cfg.Every <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		key := limiter.key(r)
		if !limiter.allow(key) {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "429 too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// allow consumes one token for key, refilling by elapsed time. It reports
// whether the request may pass.
func (limiter *Limiter) allow(key string) bool {
	jetzt := time.Now()

	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	topf, vorhanden := limiter.buckets[key]
	if !vorhanden {
		topf = &bucket{tokens: float64(limiter.cfg.Burst), lastFill: jetzt}
		limiter.buckets[key] = topf
	}

	nachfuell := jetzt.Sub(topf.lastFill).Seconds() * limiter.ratePerSec
	topf.tokens = min(topf.tokens+nachfuell, float64(limiter.cfg.Burst))
	topf.lastFill = jetzt

	if topf.tokens < 1 {
		return false
	}

	topf.tokens--
	return true
}

// key returns the bucket key for a request.
func (limiter *Limiter) key(r *http.Request) string {
	if limiter.cfg.Key != nil {
		return limiter.cfg.Key(r)
	}
	return clientIP(r, limiter.trusted)
}

// clientIP returns the client address: the socket address, or — only from a
// trusted proxy — the first X-Forwarded-For entry.
func clientIP(r *http.Request, trusted map[string]bool) string {
	host, _, err_split := net.SplitHostPort(r.RemoteAddr)
	if err_split != nil {
		host = r.RemoteAddr
	}

	if trusted[host] {
		if weiter := r.Header.Get("X-Forwarded-For"); weiter != "" {
			return firstCSV(weiter)
		}
		if echt := r.Header.Get("X-Real-IP"); echt != "" {
			return echt
		}
	}

	return host
}

// firstCSV returns the first comma-separated token, trimmed.
func firstCSV(wert string) string {
	for index := 0; index < len(wert); index++ {
		if wert[index] == ',' {
			return trimSpace(wert[:index])
		}
	}
	return trimSpace(wert)
}

func trimSpace(wert string) string {
	start := 0
	ende := len(wert)
	for start < ende && (wert[start] == ' ' || wert[start] == '\t') {
		start++
	}
	for ende > start && (wert[ende-1] == ' ' || wert[ende-1] == '\t') {
		ende--
	}
	return wert[start:ende]
}
