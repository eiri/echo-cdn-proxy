package cdnproxy

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

// Entry is a Cache entry for CDN file
type Entry struct {
	CType string
	Data  []byte
}

// Config keeps proxy configuration
type Config struct {
	CDN    *url.URL
	Prefix string
	Client *http.Client
	Cache  *cache.Cache
}

// NewConfig returns new proxy configured with given CDN and prefix
func NewConfig(baseURL, prefix string) Config {
	cdn, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	// add a tracer to the round-tripper later
	client := &http.Client{Transport: http.DefaultTransport}
	cache := cache.New(cache.NoExpiration, 10*time.Minute)
	return Config{
		CDN:    cdn,
		Prefix: prefix,
		Client: client,
		Cache:  cache,
	}
}

// Proxy is the middleware echo handler
func (cfg Config) Proxy(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		p := c.Path()

		// we have the static route
		if p == "/*" {
			pp, err := url.PathUnescape(c.Param("*"))
			if err != nil {
				return err
			}
			p = path.Clean("/" + pp)
		}

		if !strings.HasPrefix(p, cfg.Prefix) {
			return next(c)
		}

		cfg.CDN.Path = p

		if e, ok := cfg.Cache.Get(p); ok {
			c.Logger().Debugf("cache hit for %s", p)
			entry := e.(Entry)
			return c.Blob(http.StatusOK, entry.CType, entry.Data)
		}

		c.Logger().Debugf("request %s", cfg.CDN.String())

		resp, err := cfg.Client.Get(cfg.CDN.String())
		if err != nil {
			return err
		}
		c.Logger().Debugf("response %s", resp.Status)
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}

		entry := Entry{Data: data}
		if contentType, ok := resp.Header["Content-Type"]; ok {
			entry.CType = contentType[0]
		} else {
			entry.CType = echo.MIMETextPlain
		}
		cfg.Cache.Set(p, entry, cache.NoExpiration)

		return c.Blob(resp.StatusCode, entry.CType, entry.Data)
	}
}
