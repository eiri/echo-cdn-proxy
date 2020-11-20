package cdnproxy_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eiri/echo-cdn-proxy"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestProxy verifies that proxy works with different router configs
func TestProxy(t *testing.T) {
	cdn := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", echo.MIMEApplicationJavaScriptCharsetUTF8)
		fmt.Fprint(w, `/* fakelib v0.20.0 */`)
	}))
	defer cdn.Close()

	var req *http.Request
	var rec *httptest.ResponseRecorder
	echoes := map[string]*echo.Echo{
		"echo router simple":            echoSimple(cdn),
		"echo router with static":       echoStatic(cdn),
		"echo router with prefix clash": echoClash(cdn),
	}

	for tName, e := range echoes {
		t.Run(tName, func(t *testing.T) {
			// does proxy to our fake CDN
			req = httptest.NewRequest(http.MethodGet, "/npm/fakelib.min.js", nil)
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, echo.MIMEApplicationJavaScriptCharsetUTF8, rec.Header().Get("Content-Type"))
			assert.Equal(t, `/* fakelib v0.20.0 */`, rec.Body.String())

			// still can access other routes
			req = httptest.NewRequest(http.MethodGet, "/api", nil)
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get("Content-Type"))
			assert.Equal(t, `{"ok":true}`, strings.TrimRight(rec.Body.String(), "\n"))
		})
	}
}

// TestProxyNotFound verifies that error codes proxied through
func TestProxyNotFound(t *testing.T) {
	cdn := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", echo.MIMETextPlain)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `not found`)
	}))
	defer cdn.Close()

	e := echoSimple(cdn)

	req := httptest.NewRequest(http.MethodGet, "/npm/fakelib.min.js", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, echo.MIMETextPlain, rec.Header().Get("Content-Type"))
	assert.Equal(t, `not found`, strings.TrimRight(rec.Body.String(), "\n"))

}

// helper functions

func echoSimple(cdn *httptest.Server) *echo.Echo {
	e := echo.New()
	cfg := cdnproxy.NewConfig(cdn.URL, "/npm")
	e.Use(cfg.Proxy)
	e.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"ok": true})
	})
	return e
}

func echoStatic(cdn *httptest.Server) *echo.Echo {
	e := echoSimple(cdn)
	e.Static("/", ".")
	return e
}

func echoClash(cdn *httptest.Server) *echo.Echo {
	e := echoSimple(cdn)
	e.GET("/npm/*", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"ok": true})
	})
	return e
}
