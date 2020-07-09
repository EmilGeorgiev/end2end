package v2

import (
	"encoding/json"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

// Request make a call to server.
type Request struct {
	url                string
	httpRequest        *http.Request
}

// NewRequest a new request and return it.
func NewRequest(method string, url string, body string) Request {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	return Request{httpRequest: req}
}

// WithBasicAuth sets the request's Authorization header to use HTTP
// Basic Authentication with the provided username and password.
func (r Request) WithBasicAuth(userName, password string) Request {
	r.httpRequest.SetBasicAuth(userName, password)
	return r
}

// WithBearerToken set authentication to the http.Request with bearer token.
func (r Request) WithBearerToken(token string) Request {
	r.httpRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return r
}

// Headers adds the key, value pair headers to the http.Request.
func (r Request) Headers(headers map[string]string) Request {
	if r.httpRequest == nil {
		return r
	}

	for k, v := range headers {
		r.httpRequest.Header.Add(k, v)
	}

	return r
}

// Params adds the key, value pair query to the route.
func (r Request) Params(params map[string]string) Request {
	if r.httpRequest == nil {
		return r
	}

	q := r.httpRequest.URL.Query()
	for k, v := range params {
		q.Set(k , v)
	}

	r.httpRequest.URL.RawQuery = q.Encode()
	return r
}

// Call send a request to the server.
func (r Request) Call(t *testing.T) {
	start := time.Now()
	resp, err := http.DefaultClient.Do(r.httpRequest)
	if err != nil {
		t.Fatal(err)
	}
	duration := time.Since(start)
	defer resp.Body.Close()

	Responses <- Response{
		StatusCode: resp.StatusCode,
		TimeDuration: duration.Milliseconds(),
		Endpoint: r.httpRequest.URL.Path,
	}
}

// RequestExpectant send request and assert the expected response.
type RequestExpectant struct {
	Request
	response interface{}
	responseStatusCode int
}

// Expect accept as parameter a response and status code that you expect to be returned from the server.
func (r Request) Expect(response interface{}, statusCode int) RequestExpectant {
	return RequestExpectant{
		Request: r,
		response: response,
		responseStatusCode: statusCode,
	}
}

// Call make actual make request to the server and assert the returned response with expected one.
func (r RequestExpectant) Call(t *testing.T) {
	start := time.Now()
	resp, err := http.DefaultClient.Do(r.httpRequest)
	if err != nil {
		t.Fatal(err)
	}
	duration := time.Since(start)
	defer resp.Body.Close()

	Responses <- Response{
		StatusCode: resp.StatusCode,
		TimeDuration: duration.Milliseconds(),
		Endpoint: r.httpRequest.URL.Path,
	}

	if r.response == nil && resp.StatusCode == r.responseStatusCode {
		return
	}

	if resp.StatusCode != r.responseStatusCode {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		t.Errorf("expect: %s %s return status code: %d and response type: %T", r.httpRequest.Method, r.httpRequest.URL.String(), r.responseStatusCode, r.response)
		t.Errorf("actual: %s %s return status code: %d and response     : %s", r.httpRequest.Method, r.httpRequest.URL.String(), resp.StatusCode, b)
		t.Fatal()
	}

	switch resp.Header.Get("content-type") {
	case "image/png":
		img, err := png.Decode(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		err = png.Encode(r.response.(io.Writer), img)
		if err != nil {
			t.Fatal(err)
		}
		return
	case "image/jpeg":
		img, err := jpeg.Decode(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		err = jpeg.Encode(r.response.(io.Writer), img, nil)
		if err != nil {
			t.Fatal(err)
		}
		return
	case "image/gif":
		img, err := gif.Decode(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		err = gif.Encode(r.response.(io.Writer), img, nil)
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(&r.response); err != nil {
		t.Fatal(err)
	}
}

func (r Request) Send() (*http.Response, error) {
	start := time.Now()
	resp, err := http.DefaultClient.Do(r.httpRequest)
	duration := time.Since(start)

	Responses <- Response{
		StatusCode: resp.StatusCode,
		TimeDuration: duration.Milliseconds(),
		Endpoint: r.httpRequest.URL.Path,
	}

	return resp, err
}