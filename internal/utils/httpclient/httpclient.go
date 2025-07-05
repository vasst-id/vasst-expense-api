package httpclient

import (
	"crypto/tls"
	"net/http"
	"time"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

type (
	client struct {
		Doer *http.Client
	}
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Config struct {
	Timeout       int    // in millisecond
	ServiceName   string // trace serviceName
	ResourceNamer func(req *http.Request) string
	Transport     struct {
		DisableKeepAlives   bool
		MaxIdleConns        int
		MaxConnsPerHost     int
		MaxIdleConnsPerHost int
		IdleConnTimeout     time.Duration
	}
	DataDogTracer bool //trace to be on or off
}

// DefaultResourceNamer provides default resource namer function for tracing resource tag
func DefaultResourceNamer() func(req *http.Request) string {
	return func(req *http.Request) string {
		return req.Method + " " + req.URL.Path
	}
}

// New provides a http client with preconfigured configs to be used on making http calls on go app
func New(cfg *Config) Client {
	c := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Millisecond,
		Transport: &http.Transport{
			DisableKeepAlives:   cfg.Transport.DisableKeepAlives,
			MaxIdleConns:        cfg.Transport.MaxIdleConns,
			MaxConnsPerHost:     cfg.Transport.MaxConnsPerHost,
			MaxIdleConnsPerHost: cfg.Transport.MaxIdleConnsPerHost,
			IdleConnTimeout:     cfg.Transport.IdleConnTimeout,
		},
	}

	if cfg.ResourceNamer == nil {
		cfg.ResourceNamer = DefaultResourceNamer()
	}

	if cfg.DataDogTracer {
		httptrace.WrapClient(c,
			httptrace.RTWithServiceName(cfg.ServiceName),
			httptrace.RTWithResourceNamer(cfg.ResourceNamer),
		)
	}

	return &client{c}
}

// NewWithInsecureTLS provides a http client with preconfigured configs and TLS ClientConfig to Insecure Skip Verify, to be used on making http calls on go app
func NewWithInsecureTLS(cfg *Config) Client {
	c := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Millisecond,
		Transport: &http.Transport{
			DisableKeepAlives:   cfg.Transport.DisableKeepAlives,
			MaxIdleConns:        cfg.Transport.MaxIdleConns,
			MaxConnsPerHost:     cfg.Transport.MaxConnsPerHost,
			MaxIdleConnsPerHost: cfg.Transport.MaxIdleConnsPerHost,
			IdleConnTimeout:     cfg.Transport.IdleConnTimeout,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
	}

	if cfg.ResourceNamer == nil {
		cfg.ResourceNamer = DefaultResourceNamer()
	}

	if cfg.DataDogTracer {
		httptrace.WrapClient(c,
			httptrace.RTWithServiceName(cfg.ServiceName),
			httptrace.RTWithResourceNamer(cfg.ResourceNamer),
		)
	}

	return &client{c}
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	return c.Doer.Do(req)
}
