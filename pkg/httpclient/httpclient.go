package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/projectdiscovery/ratelimit"
)

type HttpClient struct {
	client  *http.Client
	headers map[string][]string
	limiter *ratelimit.Limiter
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client:  &http.Client{},
		headers: map[string][]string{},
		limiter: nil,
	}
}

func (c *HttpClient) AddLimiter(ctx context.Context, reqPerSecond uint) {
	c.limiter = ratelimit.New(ctx, reqPerSecond, time.Second)
}

func (c *HttpClient) AddHeaders(headers map[string][]string) {
	for k, v := range headers {
		c.headers[k] = v
	}
}

func (c *HttpClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	c.beforeRequest(req)

	return c.do(req)
}

func (c *HttpClient) do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *HttpClient) beforeRequest(req *http.Request) {
	if c.limiter != nil {
		c.limiter.Take()
	}
	c.setHeaders(req)
}

func (c *HttpClient) setHeaders(req *http.Request) {
	for k, v := range c.headers {
		req.Header[k] = v
	}
}
