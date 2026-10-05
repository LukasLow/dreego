package dreego

import (
	"errors"
	"net"
	"net/http"
	"sync"
)

// Server is one running HTTP server started by Start.
//
// Start returns immediately; the server runs in its own goroutine. Wait blocks
// until the server stops and returns the first non-graceful error. Close asks
// the server to shut down gracefully.
//
// This is the small piece that lets one process run several apps on several
// ports (e.g. www, app, link) without any ad-hoc net/http glue in main:
//
//	www := dreego.Start(":3000", wwwApp.Handler())
//	app := dreego.Start(":3001", appApp.Handler())
//	err := errors.Join(www.Wait(), app.Wait())
type Server struct {
	httpServer *http.Server
	listener   net.Listener

	// done is closed once Serve has returned (or Start failed to bind).
	done chan struct{}

	// errOnce guards the recorded error: the serve loop and a possible Start
	// failure both funnel into the same value.
	errOnce sync.Once
	err     error
}

// Start binds addr and serves handler in a background goroutine. It returns a
// *Server at once; use Wait to block. A bind failure is not fatal here — it is
// recorded and returned by Wait, so a caller can start several servers and
// collect all errors in one place.
func Start(addr string, handler http.Handler) *Server {
	server := &Server{
		httpServer: &http.Server{Addr: addr, Handler: handler},
		done:       make(chan struct{}),
	}

	listener, err_listen := net.Listen("tcp", addr)
	if err_listen != nil {
		server.recordError(err_listen)
		close(server.done)
		return server
	}
	server.listener = listener

	go func() {
		err_serve := server.httpServer.Serve(listener)
		// http.ErrServerClosed is the normal result of Close, not a failure.
		if errors.Is(err_serve, http.ErrServerClosed) {
			err_serve = nil
		}
		server.recordError(err_serve)
		close(server.done)
	}()

	return server
}

// Wait blocks until the server stops (clean shutdown, Close, or a serve error)
// and returns the recorded error, or nil on a clean shutdown.
func (server *Server) Wait() error {
	<-server.done
	return server.currentError()
}

// Addr returns the actual listen address ("host:port"). It is useful when addr
// used port 0. An empty string means the server never bound.
func (server *Server) Addr() string {
	if server.listener == nil {
		return ""
	}
	return server.listener.Addr().String()
}

// Close shuts the server down gracefully. It is safe to call more than once and
// makes a pending Wait return.
func (server *Server) Close() error {
	if server.httpServer == nil {
		return nil
	}
	return server.httpServer.Close()
}

// recordError stores the first non-nil error.
func (server *Server) recordError(err error) {
	if err == nil {
		return
	}
	server.errOnce.Do(func() { server.err = err })
}

// currentError reads the recorded error.
func (server *Server) currentError() error {
	return server.err
}
