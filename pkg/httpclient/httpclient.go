package httpclient

import (
	"context"
	"crypto/tls"
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

func (c *HttpClient) FollowRedirects(value bool) {
	c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if !value {
			return http.ErrUseLastResponse
		}
		return nil
	}
}

// AllowInsecure allows insecure connections to be made by the client.
// It skips TLS verification.
func (c *HttpClient) AllowInsecure() error {
	currentTransport := c.client.Transport
	if currentTransport == nil {
		c.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return nil
	} else {
		transport, ok := currentTransport.(*http.Transport)
		if !ok {
			return fmt.Errorf("transport is not an *http.Transport")
		}
		currentTLSClientConfig := transport.TLSClientConfig
		if currentTLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			currentTLSClientConfig.InsecureSkipVerify = true
		}
		return nil
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
