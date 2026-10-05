package dreego

import (
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
)

// staticMount is one embedded file tree served under a URL prefix.
type staticMount struct {
	fs     fs.FS
	prefix string
}

// Static serves files from an embedded filesystem (go:embed). The prefix is
// BOTH the URL prefix and the folder inside the filesystem:
//
//	//go:embed public/*
//	var publicFS embed.FS
//
//	app.Static(publicFS, "public")   // public/design.css -> /public/design.css
//	app.Static(publicFS, "")         // favicon.svg      -> /favicon.svg
//
// Static is consulted only when no page matches, so a page always wins.
// Unknown paths fall through to 404.
func (app *App) Static(dateiSystem fs.FS, prefix string) {
	mount := staticMount{fs: dateiSystem, prefix: normalizePrefix(prefix)}
	app.statics = append(app.statics, mount)

	// A prefixed mount registers its own subtree. The empty prefix uses the
	// app's "/" catch-all, which already calls serveStatic.
	if mount.prefix == "" {
		return
	}
	app.mux.HandleFunc(mount.prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		// Static files are read-only: anything but GET/HEAD is a method error.
		if r.Method != "GET" && r.Method != "HEAD" {
			w.Header().Set("Allow", "GET, HEAD")
			app.serveError(w, r, http.StatusMethodNotAllowed)
			return
		}
		if app.serveStatic(w, r) {
			return
		}
		app.serveError(w, r, http.StatusNotFound)
	})
}

// normalizePrefix turns "public" or "/public" into "/public", and "" into "".
func normalizePrefix(prefix string) string {
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return ""
	}
	return "/" + prefix
}

// serveStatic tries to serve an embedded file for the request path. It reports
// whether it handled the request.
func (app *App) serveStatic(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "GET" && r.Method != "HEAD" {
		return false
	}

	for _, mount := range app.statics {
		if !passesPrefix(r.URL.Path, mount.prefix) {
			continue
		}

		// The path inside the filesystem is the URL path without the leading
		// slash: the prefix is the folder inside the filesystem.
		relativ := strings.TrimPrefix(r.URL.Path, "/")

		inhalt, err_lesen := fs.ReadFile(mount.fs, relativ)
		if err_lesen != nil {
			continue
		}

		w.Header().Set("Content-Type", contentTypeFor(relativ, inhalt))
		w.Header().Set("Cache-Control", "public, max-age=300")
		w.WriteHeader(http.StatusOK)
		if r.Method == "GET" {
			_, _ = w.Write(inhalt)
		}
		return true
	}

	return false
}

// passesPrefix reports whether the request path lies under the mount prefix.
// An empty prefix matches everything.
func passesPrefix(pfad string, prefix string) bool {
	if prefix == "" {
		return true
	}
	return pfad == prefix || strings.HasPrefix(pfad, prefix+"/")
}

// webTypes covers the common web assets explicitly. The system MIME database is
// not always available (e.g. in a minimal container), so .css would otherwise
// be served as text/plain.
var webTypes = map[string]string{
	".css":         "text/css; charset=utf-8",
	".js":          "text/javascript; charset=utf-8",
	".mjs":         "text/javascript; charset=utf-8",
	".json":        "application/json; charset=utf-8",
	".svg":         "image/svg+xml",
	".html":        "text/html; charset=utf-8",
	".txt":         "text/plain; charset=utf-8",
	".xml":         "application/xml; charset=utf-8",
	".woff":        "font/woff",
	".woff2":       "font/woff2",
	".png":         "image/png",
	".jpg":         "image/jpeg",
	".jpeg":        "image/jpeg",
	".webp":        "image/webp",
	".ico":         "image/x-icon",
	".webmanifest": "application/manifest+json",
}

// contentTypeFor resolves the Content-Type for a file: explicit web types
// first, then the system MIME database, then content sniffing.
func contentTypeFor(name string, inhalt []byte) string {
	erweiterung := strings.ToLower(path.Ext(name))
	if typ, bekannt := webTypes[erweiterung]; bekannt {
		return typ
	}
	if typ := mime.TypeByExtension(erweiterung); typ != "" {
		return typ
	}
	return http.DetectContentType(inhalt)
}
