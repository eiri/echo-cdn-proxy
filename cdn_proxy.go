package cdnproxy

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
)

// Config keeps proxy configuration
type Config struct {
	CDN    *url.URL
	Prefix string
	Client *http.Client
}

// NewConfig returns new proxy configured with given CDN and prefix
func NewConfig(baseURL, prefix string) Config {
	cdn, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	// add a tracer to the round-tripper later
	client := &http.Client{Transport: http.DefaultTransport}
	return Config{
		CDN:    cdn,
		Prefix: prefix,
		Client: client,
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

		contentType, ok := resp.Header["Content-Type"]
		if !ok {
			contentType = []string{echo.MIMETextPlain}
		}

		return c.Blob(resp.StatusCode, contentType[0], data)
	}
}
