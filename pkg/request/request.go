package request

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

const (
	QueryParam  = "queries"
	HeaderParam = "headers"
)

// Client provided for mock
type Client interface {
	Get(targetURL string, opts SendOptions) (*Response, error)
	Put(targetURL string, opts SendOptions, body []byte) (*Response, error)
	Post(targetURL string, opts SendOptions, body []byte) (*Response, error)
	Patch(targetURL string, opts SendOptions, body []byte) (*Response, error)
	Delete(targetURL string, opts SendOptions, body []byte) (*Response, error)
	Send(method, path string, opts SendOptions, body []byte) (*Response, error)
	GetStandardClient() *http.Client
}

// SendOptions for attached data through a request
// Example should be add query params (http://abcd.com?user=ec&limit=5) or header
//
//	SendOptions {
//	 "queries": map[string]interface{}{
//	    "user": "ec",
//	    "limit": "5",
//	 },
//		"headers": map[string]interface{}{
//	    "xxx":         	"abc",
//	    "Content-Type":  "application/json",
//		},
//	}
type SendOptions map[string]map[string]interface{}

// WithTimeout sets timeout of the client options
func WithTimeout(timeout time.Duration) OptionClient {
	return func(r *retryablehttp.Client) {
		r.HTTPClient.Timeout = timeout
	}
}

// WithRetryMax sets max retry of the client options
func WithRetryMax(retryMax int) OptionClient {
	return func(r *retryablehttp.Client) {
		r.RetryMax = retryMax
	}
}

// OptionClient represents an option for the http client
type OptionClient func(*retryablehttp.Client)

type client struct {
	debugEnable bool
	logger      *logrus.Logger
	HTTPClient  *http.Client
}

// Response struct
type Response struct {
	Body       []byte
	Header     http.Header
	StatusCode int
}

// NewClient init http client
func NewClient(optsClient ...OptionClient) Client {
	return NewClientWithDebug(false, optsClient...)
}

// NewClientWithDebug init http client with debug config
func NewClientWithDebug(debugEnable bool, optsClient ...OptionClient) Client {
	clientlogger := logrus.New()
	httpClient := retryablehttp.NewClient()
	for _, optClient := range optsClient {
		optClient(httpClient)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient.HTTPClient.Transport = tr
	httpClient.Logger = log.New(io.Discard, "", log.LstdFlags)
	if debugEnable {
		clientlogger.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: false,
			DisableColors:             false,
			TimestampFormat:           "2006-01-02 15:04:05.000",
			FullTimestamp:             true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyLevel: "log_level",
			},
		})
		clientlogger.SetLevel(logrus.DebugLevel)
		httpClient.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
			clientlogger.WithFields(logrus.Fields{
				"request": map[string]string{
					"proto": req.Proto,
					"host":  req.URL.Host,
					"path":  req.URL.Path,
				},
				"attempt": attempt,
			}).Debug("Sending request")
		}
	}
	return &client{
		debugEnable: debugEnable,
		logger:      clientlogger,
		HTTPClient:  httpClient.StandardClient(),
	}
}

// SetContentType sets a content type of the request
func (opt SendOptions) SetContentType(t string) SendOptions {
	newOpt := opt
	if newOpt == nil {
		newOpt = make(SendOptions)
	}
	if newOpt[HeaderParam] == nil {
		newOpt[HeaderParam] = make(map[string]interface{})
	}
	newOpt[HeaderParam]["Content-Type"] = t

	return newOpt
}

// SetContentType sets query params of the request
func (opt SendOptions) SetQueryParam(params map[string]interface{}) SendOptions {
	newOpt := opt
	if newOpt == nil {
		newOpt = make(SendOptions)
	}
	newOpt[QueryParam] = params

	return newOpt
}

// Get request and returns response from target URL
func (c client) Get(targetURL string, opts SendOptions) (*Response, error) {
	return c.Send(http.MethodGet, targetURL, opts, nil)
}

// Put request and returns response from target URL
func (c client) Put(targetURL string, opts SendOptions, body []byte) (*Response, error) {
	return c.Send(http.MethodPut, targetURL, opts, body)
}

// Post request and returns response from target URL
func (c client) Post(targetURL string, opts SendOptions, body []byte) (*Response, error) {
	return c.Send(http.MethodPost, targetURL, opts, body)
}

// Patch request and returns response from target URL
func (c client) Patch(targetURL string, opts SendOptions, body []byte) (*Response, error) {
	return c.Send(http.MethodPatch, targetURL, opts, body)
}

// Delete request and returns response from target URL
func (c client) Delete(targetURL string, opts SendOptions, body []byte) (*Response, error) {
	return c.Send(http.MethodDelete, targetURL, opts, body)
}

// Send a request and returns response from target URL
func (c client) Send(method, targetURL string, opts SendOptions, body []byte) (*Response, error) {
	if c.debugEnable {
		c.logger.WithFields(logrus.Fields{
			"method": method,
			"url":    targetURL,
			"opts":   opts,
			"body":   string(body),
		}).Debug("[Send]: http request")
	}
	method = strings.ToUpper(method)
	urlSchema, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	queryString := url.Values{}
	if query := opts[QueryParam]; len(query) != 0 {
		for key, val := range query {
			queryString.Add(key, val.(string))
		}
	}
	if len(queryString) > 0 {
		urlSchema.RawQuery = queryString.Encode()
	}
	requestAPIUrl := urlSchema.String()

	var bBody io.Reader
	if body != nil {
		bBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, requestAPIUrl, bBody)
	if err != nil {
		return nil, err
	}

	if header := opts[HeaderParam]; len(header) != 0 {
		for key, val := range header {
			req.Header.Set(key, val.(string))
		}
	}
	if c.debugEnable {
		c.logger.Debugf("request %+v", req)
	}

	netResponse, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(netResponse.Body)
	defer func() {
		if err := netResponse.Body.Close(); err != nil {
			c.logger.WithError(err).Error("[Client.Send]: unable to close a response body")
		}
	}()

	if err != nil {
		return nil, err
	}

	response := &Response{
		StatusCode: netResponse.StatusCode,
		Header:     netResponse.Header,
		Body:       contents,
	}

	return response, nil
}

// Get http standard client
func (c client) GetStandardClient() *http.Client {
	return c.HTTPClient
}
