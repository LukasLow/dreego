package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimitBlocktNachBurst(t *testing.T) {
	limiter := New(Config{Burst: 3, Every: time.Minute})

	durch := 0
	blockiert := 0
	for index := 0; index < 5; index++ {
		schreiber := httptest.NewRecorder()
		anfrage := httptest.NewRequest("GET", "/", nil)
		anfrage.RemoteAddr = "1.2.3.4:1111"
		limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			durch++
		})).ServeHTTP(schreiber, anfrage)
		if schreiber.Code == http.StatusTooManyRequests {
			blockiert++
		}
	}

	if durch != 3 {
		t.Fatalf("durch = %d, erwartet 3", durch)
	}
	if blockiert != 2 {
		t.Fatalf("blockiert = %d, erwartet 2", blockiert)
	}
}

func TestVerschiedeneIPsGetrennt(t *testing.T) {
	limiter := New(Config{Burst: 1, Every: time.Minute})
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	erst := httptest.NewRequest("GET", "/", nil)
	erst.RemoteAddr = "1.1.1.1:1"
	zweit := httptest.NewRequest("GET", "/", nil)
	zweit.RemoteAddr = "2.2.2.2:2"

	s1 := httptest.NewRecorder()
	handler.ServeHTTP(s1, erst)
	s2 := httptest.NewRecorder()
	handler.ServeHTTP(s2, zweit)

	if s1.Code != 200 || s2.Code != 200 {
		t.Fatalf("verschiedene IPs sollen getrennt zählen: %d %d", s1.Code, s2.Code)
	}
}

func TestXForwardedForNurVonTrusted(t *testing.T) {
	limiter := New(Config{
		Burst:          1,
		Every:          time.Minute,
		TrustedProxies: []string{"10.0.0.1"},
	})
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// Vertrauenswürdiger Proxy: X-Forwarded-For zählt.
	a := httptest.NewRequest("GET", "/", nil)
	a.RemoteAddr = "10.0.0.1:1"
	a.Header.Set("X-Forwarded-For", "9.9.9.9")
	sa := httptest.NewRecorder()
	handler.ServeHTTP(sa, a)

	// Gleicher echter Client nochmal → blockiert.
	b := httptest.NewRequest("GET", "/", nil)
	b.RemoteAddr = "10.0.0.1:2"
	b.Header.Set("X-Forwarded-For", "9.9.9.9")
	sb := httptest.NewRecorder()
	handler.ServeHTTP(sb, b)

	if sb.Code != http.StatusTooManyRequests {
		t.Fatalf("zweiter Request desselben Clients soll 429 sein, war %d", sb.Code)
	}

	// Nicht vertrauenswürdiger Client: Header wird ignoriert, eigene IP zählt.
	c := httptest.NewRequest("GET", "/", nil)
	c.RemoteAddr = "5.5.5.5:3"
	c.Header.Set("X-Forwarded-For", "9.9.9.9")
	sc := httptest.NewRecorder()
	handler.ServeHTTP(sc, c)

	if sc.Code != 200 {
		t.Fatalf("gespoofte XFF darf nicht zählen: %d", sc.Code)
	}
}

func TestDeaktiviert(t *testing.T) {
	limiter := New(Config{})
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	schreiber := httptest.NewRecorder()
	handler.ServeHTTP(schreiber, httptest.NewRequest("GET", "/", nil))
	if schreiber.Code != 200 {
		t.Fatalf("deaktivierter Limiter soll durchlassen: %d", schreiber.Code)
	}
}
