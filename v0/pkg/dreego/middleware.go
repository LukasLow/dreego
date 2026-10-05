package dreego

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
	"sync"
)

// SecurityHeaders sets the response security headers that are always on. The
// Content-Security-Policy is set per request in ServeHTTP, because it carries a
// fresh nonce.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		kopf := w.Header()
		kopf.Set("X-Content-Type-Options", "nosniff")
		kopf.Set("X-Frame-Options", "DENY")
		kopf.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		kopf.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		next.ServeHTTP(w, r)
	})
}

// Compress enables gzip when the client accepts it. Text responses (HTML, CSS,
// JS, SVG, JSON) shrink to roughly a third, which matters most on slow links.
//
// Already-compressed content (images, woff/woff2) is never gzipped, and HEAD,
// 204 and 304 responses are left untouched. The output always carries
// Vary: Accept-Encoding so caches keep compressed and plain copies apart.
func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !acceptsGzip(r.Header.Get("Accept-Encoding")) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Vary", "Accept-Encoding")

		gzWrite := &gzipResponseWriter{ResponseWriter: w}
		defer gzWrite.Close()

		next.ServeHTTP(gzWrite, r)
	})
}

// acceptsGzip reports whether the Accept-Encoding header allows gzip. A
// "gzip;q=0" explicitly forbids it; "*" allows it.
func acceptsGzip(header string) bool {
	for _, teil := range strings.Split(header, ",") {
		teil = strings.TrimSpace(teil)
		name, q, hatQ := strings.Cut(teil, ";")
		name = strings.ToLower(strings.TrimSpace(name))

		if name != "gzip" && name != "*" {
			continue
		}

		if hatQ {
			wert := strings.TrimSpace(strings.TrimPrefix(q, "q="))
			if wert == "0" || wert == "0.0" || wert == "0.00" || wert == "0.000" {
				continue
			}
		}
		return true
	}
	return false
}

// gzipResponseWriter compresses the body. The decision is made at WriteHeader
// time — that is the only correct moment, because after WriteHeader the header
// is on the wire. Handlers that call WriteHeader before writing (c.JSON,
// http.Error) would otherwise get a gzipped body without a Content-Encoding
// header, which browsers show as binary garbage.
type gzipResponseWriter struct {
	http.ResponseWriter
	gz          *gzip.Writer
	geschrieben bool
}

var gzipPool = sync.Pool{New: func() any { return gzip.NewWriter(nil) }}

// WriteHeader decides whether to compress, based on the Content-Type that is
// already set, then writes the header.
func (gw *gzipResponseWriter) WriteHeader(status int) {
	if gw.geschrieben {
		return
	}
	gw.geschrieben = true

	// No body: pass through untouched.
	if status == http.StatusNoContent || status == http.StatusNotModified || status < 200 {
		gw.ResponseWriter.WriteHeader(status)
		return
	}

	if compressibleContentType(gw.Header().Get("Content-Type")) {
		gw.Header().Set("Content-Encoding", "gzip")
		gw.Header().Del("Content-Length")
		writer := gzipPool.Get().(*gzip.Writer)
		writer.Reset(gw.ResponseWriter)
		gw.gz = writer
	}
	gw.ResponseWriter.WriteHeader(status)
}

// Write compresses through gzip when it was enabled at header time.
func (gw *gzipResponseWriter) Write(daten []byte) (int, error) {
	if !gw.geschrieben {
		// No explicit WriteHeader: implicit 200. Decide now.
		gw.WriteHeader(http.StatusOK)
	}
	if gw.gz != nil {
		return gw.gz.Write(daten)
	}
	return gw.ResponseWriter.Write(daten)
}

// Flush makes the buffered gzip data visible (needed for streaming responses).
func (gw *gzipResponseWriter) Flush() {
	if gw.gz != nil {
		_ = gw.gz.Flush()
	}
	if flusher, ok := gw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Close finishes the gzip member and returns the writer to the pool.
func (gw *gzipResponseWriter) Close() {
	if gw.gz == nil {
		return
	}
	_ = gw.gz.Close()
	gzipPool.Put(gw.gz)
	gw.gz = nil
}

// Unwrap exposes the underlying writer (for http.ResponseController).
func (gw *gzipResponseWriter) Unwrap() http.ResponseWriter { return gw.ResponseWriter }

// compressibleContentType reports whether a Content-Type benefits from gzip.
// Empty means "unknown": we then compress, which is the safe default for HTML.
func compressibleContentType(typ string) bool {
	if typ == "" {
		return true
	}
	typ = strings.ToLower(typ)
	if strings.HasPrefix(typ, "text/") {
		return true
	}
	for _, teil := range []string{
		"json", "javascript", "xml", "svg", "wasm",
	} {
		if strings.Contains(typ, teil) {
			return true
		}
	}
	// Images, audio/video, fonts, zip/gzip are already compressed.
	return false
}

// Recover turns a panic in any handler into a generic 500 response. The internal
// cause is logged server-side and never sent to the client.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			panik := recover()
			if panik == nil {
				return
			}
			log.Printf("dreego: panic bei %s %s: %v", r.Method, r.URL.Path, panik)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}()
		next.ServeHTTP(w, r)
	})
}

// recoverMiddleware is Recover bound to the app, so a panic renders the app's
// own 500 page instead of bare text.
func (app *App) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			panik := recover()
			if panik == nil {
				return
			}
			log.Printf("dreego: panic bei %s %s: %v", r.Method, r.URL.Path, panik)
			app.serveError(w, r, http.StatusInternalServerError)
		}()
		next.ServeHTTP(w, r)
	})
}
