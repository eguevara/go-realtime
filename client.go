package realtime

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	maxErrMsgLen   = 256
	defaultTimeout = 2 * time.Second
)

var (
	defaultBaseURL = func() *url.URL {
		u, err := url.Parse("https://www.googleapis.com/analytics/v3/data/")
		if err != nil {
			panic(err)
		}
		return u
	}()
)

// Response stores the google analytics response.
type Response struct {
	TotalsForAllResults *ResponseTotals `json:"totalsForAllResults,omitempty"`
}

// ResponseTotals stores the total active users response.
type ResponseTotals struct {
	RtActiveUsers *string `json:"rt:activeUsers,omitempty"`
}

// Options specifies the parameters required to make a valid request.
type Options struct {
	IDs     string `url:"ids"`
	Metrics string `url:"metrics"`
}

// ClientOption modifies the default behavior of a realtime client.
type ClientOption func(*Client)

// WithHTTPClient makes the realtime client use the given HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) { c.client = client }
}

// Client to interact with the realtime API
type Client struct {
	client  *http.Client
	baseURL *url.URL
}

// NewClient creates a client for the realtime API.
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		client:  &http.Client{Timeout: defaultTimeout},
		baseURL: defaultBaseURL,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// GetRealTime returns the realtime analytics api.
// https://developers.google.com/apis-explorer/#p/analytics/v3/analytics.data.realtime.get
func (c *Client) GetRealTime(opt *Options) (*Response, error) {
	url, errOptions := addOptions("realtime", opt)
	if errOptions != nil {
		return nil, errOptions
	}

	response := new(Response)
	err := c.doGet(url, func(r io.Reader) error {
		return json.NewDecoder(r).Decode(response)
	})
	return response, err
}

func (c *Client) doGet(resource string, decoder func(r io.Reader) error) error {
	return c.doGetURL(c.resolve(resource), decoder)
}

func (c *Client) doGetURL(url string, decoder func(r io.Reader) error) error {
	resp, err := c.client.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c.makeError(resp)
	}

	return decoder(resp.Body)
}

func (c *Client) makeError(resp *http.Response) error {
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, maxErrMsgLen))
	if err != nil {
		return err
	}
	if len(body) >= maxErrMsgLen {
		body = append(body[:maxErrMsgLen], []byte("... (elided)")...)
	} else if len(body) == 0 {
		body = []byte(resp.Status)
	}
	return fmt.Errorf("unexpected response from realtime API, status %d: %s",
		resp.StatusCode, string(body))
}

func (c *Client) resolve(urlStr string) string {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	u := *c.baseURL.ResolveReference(rel)

	return u.String()
}

// addOptions adds the parameters in opt as URL query parameters to s.
// opt must be a struct whose fields may contain "url" tags.
// all parameters passed as s will be replaced by options in opt.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
