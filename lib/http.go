package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"net/http"
	"net/url"
	"time"

	"github.com/opentracing/opentracing-go"
	"gopkg.in/tokopedia/logging.v1"
)

// Why myhttp?
// https://hackernoon.com/avoiding-memory-leak-in-golang-api-1843ef45fca8
// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
// https://awmanoj.github.io/tech/2016/12/16/keep-alive-http-requests-in-golang/

// How to use it?
// Example 1
// Let's say you need to call 2 separate domains.
// Domain 1: Typically replies back in 10 seconds and you are making few calls(so you can have just a few cached conns).
// domain1Client, err := Client{Timeout: 10 * time.Second, MaxIdleConns: 20}.New("domain1")
// Domain 2: Typically replies back in 2 seconds and you are making a lot of calls(so need to have a lot more cached conns).
// domain2Client, err := Client{Timeout: 2 * time.Second, MaxIdleConns: 100}.New("domain2")
// Example 2
// Let's say you need to call a number of domains with no skewed pattern and typically you get replies back in <= 10 sec
// defaultClient, err := Client{}.New()

const (
	DEFAULT = "default"
	GET     = "GET"
	POST    = "POST"

	// Default value for client Timeout
	DefaultTimeout = 10 * time.Second
)

// Maintains map of connection identifier and http client.
// http client is singleton and created only once for an identifier
// Each http client is thread-safe
var clientConns map[string]*Client

// This is the custom client
type Client struct {
	*http.Client

	// Timeout specifies a time limit for requests made by this Client
	Timeout time.Duration

	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	MaxIdleConns int

	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// DefaultMaxIdleConnsPerHost is used (which is 2)
	MaxIdleConnsPerHost int

	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself. Zero means no limit. Default is 0
	IdleConnTimeout time.Duration
}

type Request struct {
	RequestClient        *Client
	Url                  string
	Headers              map[string]string
	QueryParams          map[string][]string
	FormParams           map[string]string
	RequestPostBody      interface{}
	DisAllowedStatusCode map[int]bool
}

type HTTPRequest interface {
	AddQueryParams(params map[string][]string) HTTPRequest
	AddQueryParam(key, val string) HTTPRequest
	CreateQueryString() string
	AddHeaders(headers map[string]string) HTTPRequest
	AddHeader(key, val string) HTTPRequest
	AddPostBody(body interface{}) HTTPRequest
	SetDisAllowedStatusCode(codes []int) HTTPRequest
	DoGet(ctx context.Context, isRawUrl bool, dest interface{}) (result Response, err error)
	DoPost(ctx context.Context, isRawUrl bool, dest interface{}) (result Response, err error)
}

type Response struct {
	StatusCode int    // e.g. 500
	Status     string // e.g. "500 Internal Server Error"
	Body       string
}

func RequestBuilder(ctx context.Context, url string, client *Client) HTTPRequest {
	r := &Request{
		RequestClient: client,
		Url:           url,
	}

	return r
}

func (r *Request) AddQueryParams(params map[string][]string) HTTPRequest {
	if r.QueryParams == nil {
		r.QueryParams = params
		return r
	}
	for k, v := range params {
		r.QueryParams[k] = v
	}
	return r
}

func (r *Request) AddQueryParam(key, val string) HTTPRequest {
	if r.QueryParams == nil {
		r.QueryParams = make(map[string][]string)
	}
	arr, ok := r.QueryParams[key]
	if !ok {
		r.QueryParams[key] = []string{val}
	} else {
		r.QueryParams[key] = append(arr, val)
	}
	// r.queryParams[key] = url.QueryEscape(val) //feed is returning 400 when using this
	return r
}

func (r *Request) CreateQueryString() string {
	var buffer bytes.Buffer
	for k, v := range r.QueryParams {
		for _, val := range v {
			buffer.WriteString(k)
			buffer.WriteString("=")
			buffer.WriteString(val)
			buffer.WriteString("&")
		}
	}
	if buffer.Len() > 0 {
		buffer.Truncate(buffer.Len() - 1)
	}
	return buffer.String()
}

func (r *Request) AddHeaders(headers map[string]string) HTTPRequest {
	if r.Headers == nil {
		r.Headers = headers
		return r
	}

	for k, v := range headers {
		r.Headers[k] = v
	}

	return r
}

func (r *Request) AddHeader(key, val string) HTTPRequest {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = val
	return r
}

func (r *Request) AddPostBody(body interface{}) HTTPRequest {
	r.RequestPostBody = body
	return r
}

func (r *Request) SetDisAllowedStatusCode(codes []int) HTTPRequest {
	r.DisAllowedStatusCode = make(map[int]bool)
	for _, c := range codes {
		if c != http.StatusOK {
			r.DisAllowedStatusCode[c] = true
		}
	}
	return r
}

func (c *Request) DoGet(ctx context.Context, isRawUrl bool, dest interface{}) (result Response, err error) {
	return c.executeRequest(ctx, GET, isRawUrl, dest)
}

func (c *Request) DoPost(ctx context.Context, isRawUrl bool, dest interface{}) (result Response, err error) {
	return c.executeRequest(ctx, POST, isRawUrl, dest)
}

// @Summary Helper function to make Http calls
// @Description This function makes Http calls optimally and more safely.
// Caller does not need o bother about closing resources.
// Problems it solves are:
// 1. Sets context's deadline to `time.Now().Add(c.Timeout)`. Thus, no HTTP calls
// are made without a timeout.
// 2. Releases resources when deadline is crossed or If the response is received
// before deadline then context is canceled by calling deferred cancel()
// 3. Makes sure that a response body is read completely and closed
// if response code is 2xx, http response body is un-marshaled into dest
func (c *Request) executeRequest(ctx context.Context, method string, isRawUrl bool, dest interface{}) (result Response, err error) {

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span, ctx = opentracing.StartSpanFromContext(ctx, "executeRequest")
		defer span.Finish()
	}

	var newurl *url.URL
	newurl, err = url.Parse(c.Url)
	if err != nil {
		return
	}

	if !isRawUrl {
		parameters := url.Values{}
		for key, val := range c.QueryParams {
			parameters[key] = val
		}
		newurl.RawQuery = parameters.Encode()
	}

	var body io.Reader
	if method == POST && c.RequestPostBody != nil {
		buf, _ := json.Marshal(c.RequestPostBody)
		body = bytes.NewBuffer(buf)
	}

	req, err := http.NewRequest(method, newurl.String(), body)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, c.RequestClient.Timeout)
	defer cancel()

	req = req.WithContext(ctx)

	for key, val := range c.Headers {
		req.Header.Add(key, val)
	}

	//will be removed in future
	logging.Debug.Println("[executeRequest]Prepared URL:", newurl.String())

	resp, err := c.RequestClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	result.Status = resp.Status
	result.StatusCode = resp.StatusCode
	result.Body = string(responseBody)

	if _, ok := c.DisAllowedStatusCode[resp.StatusCode]; ok {
		err = errors.New(fmt.Sprintf("StatusCode Mismatch. Found:%d", resp.StatusCode))
		return
	}

	if dest != nil {
		err = json.Unmarshal(responseBody, dest)
		if err != nil {
			return
		}
	}
	return
}
