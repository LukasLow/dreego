// Package dreegotest offers helpers for testing a dreego application without a
// network: a Client that drives an *dreego.App through requests, keeps cookies
// and follows redirects, and a Response with convenient assertions.
//
//	app := dreego.NewApp()
//	app.Page(Home)
//	client := dreegotest.New(app)
//	resp := client.Get("/")
//	resp.AssertStatus(t, 200)
//	resp.AssertContains(t, "Hallo")
package dreegotest

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
)

// Client drives an application over its http.Handler, keeping cookies between
// requests (so sessions and CSRF work as in a browser).
type Client struct {
	handler http.Handler
	jar     http.CookieJar
	base    string
}

// New builds a client for the given handler (usually app.Handler()).
func New(handler http.Handler) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{handler: handler, jar: jar, base: "http://test.local"}
}

// Get performs a GET and follows redirects.
func (c *Client) Get(pfad string) *Response {
	return c.do(http.MethodGet, pfad, "", "")
}

// PostForm performs a POST with urlencoded form values (no redirect following;
// call Follow to follow manually, or use PostFormFollow).
func (c *Client) PostForm(pfad string, werte map[string]string) *Response {
	form := url.Values{}
	for name, wert := range werte {
		form.Set(name, wert)
	}
	return c.do(http.MethodPost, pfad, form.Encode(), "application/x-www-form-urlencoded")
}

// PostJSON performs a POST with a JSON body.
func (c *Client) PostJSON(pfad string, jsonBody string) *Response {
	return c.do(http.MethodPost, pfad, jsonBody, "application/json")
}

// Response is one response plus the request that produced it.
type Response struct {
	Status int
	Header http.Header
	Body   string
}

// do performs one request (no redirect following) and records the result.
func (c *Client) do(methode string, pfad string, body string, contentType string) *Response {
	var leser *strings.Reader
	if body != "" {
		leser = strings.NewReader(body)
	} else {
		leser = strings.NewReader("")
	}

	anfrage := httptest.NewRequest(methode, c.base+pfad, leser)
	if contentType != "" {
		anfrage.Header.Set("Content-Type", contentType)
	}
	for _, cookie := range c.jar.Cookies(mussURL(c.base + pfad)) {
		anfrage.AddCookie(cookie)
	}

	schreiber := httptest.NewRecorder()
	c.handler.ServeHTTP(schreiber, anfrage)

	// Keep any cookies the response set.
	antwort := schreiber.Result()
	c.jar.SetCookies(mussURL(c.base+pfad), antwort.Cookies())

	return &Response{
		Status: schreiber.Code,
		Header: schreiber.Header(),
		Body:   schreiber.Body.String(),
	}
}

// mussURL parst eine URL; bei Fehler ein leerer Wert (kann hier nicht passieren).
func mussURL(roh string) *url.URL {
	geparst, err := url.Parse(roh)
	if err != nil {
		return &url.URL{}
	}
	return geparst
}
