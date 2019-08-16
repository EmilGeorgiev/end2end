package end2end

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Config contains basic information about the application
// that will be tested and user that wil send requests.
type Config struct {
	// URL is the base URL of the server, For example: "http://localhost:8080"
	URL string
	UserName string
	Password string
}

// Requester make a call to server and can assert the returned response with expected one.
type Requester struct {
	config             Config
	httpRequest        *http.Request
	actualResponse     interface{}
	expectedResponse   interface{}
	expectedStatusCode int64
	method             string
}

// NewRequestTo create and return a new Requester.
func NewRequestTo(a Config) Requester {
	return Requester{config: a}
}

// Create accept as parameters resource path and payload that you want to send to the server.
// Create a new http.Request and return the updated Requester.
func (r Requester) Create(path string, payload interface{}) Requester {
	b, _ :=json.Marshal(payload)
	body := bytes.NewReader(b)
	req, _ := http.NewRequest(http.MethodPost, r.config.URL +path, body)
	r.httpRequest = req
	return r
}

// Update accept as parameters resource path and payload that you want to send to the server.
// Create a new http.Request and return the updated Requester.
func (r Requester) Update(path string, payload interface{}) Requester {
	b, _ :=json.Marshal(payload)
	body := bytes.NewReader(b)
	req, _ := http.NewRequest(http.MethodPost, r.config.URL +path, body)

	r.httpRequest = req
	return r
}

// Delete accept as parameters resource path to the resource that you want to delete.
// Create a new http.Request and return the updated Requester.
func (r Requester) Delete(path string) Requester {
	req, _ := http.NewRequest(http.MethodPost, r.config.URL +path, nil)

	r.httpRequest = req
	return r
}

// Get accept as parameters resource path to the resource that you want to get.
// Create a new http.Request and return the updated Requester.
func (r Requester) Get(path string) Requester {
	req, _ := http.NewRequest(http.MethodGet, r.config.URL +path, nil)

	r.httpRequest = req
	return r
}

// Read accept as parameter a type in which you will store
func (r Requester) Read(actual interface{}) Requester {
	r.actualResponse = actual
	return r
}

func (r Requester) Headers(headers map[string]string) Requester {
	if r.httpRequest == nil {
		return r
	}

	for k, v := range headers {
		r.httpRequest.Header.Add(k, v)
	}

	return r
}

// Assert accept as parameters actual and expected response from the server.
func (r Requester) Assert(actual, expected interface{}) Requester {
	r.actualResponse = actual
	r.expectedResponse = expected
	return r
}

// ExpectStatusCode set the status code that you expect to be returned from the server.
func (r Requester) ExpectStatusCode(status int64) Requester {
	r.expectedStatusCode = status
	return r
}

// Call make actual make request to the server and assert the returned response with expected one.
func (r Requester) Call(t *testing.T) {
	r.httpRequest.SetBasicAuth(r.config.UserName, r.config.Password)

	resp, err := http.DefaultClient.Do(r.httpRequest)
	if err != nil {
		t.Fatal(err)

	}
	defer resp.Body.Close()

	if r.actualResponse == nil {
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(&r.actualResponse); err != nil {
		t.Error(err)
		return
	}

	if r.expectedResponse == nil {
		return
	}

	assert.Equal(t, r.expectedResponse, r.actualResponse)

	if r.expectedStatusCode != 0 {
		assert.EqualValues(t, r.expectedStatusCode, resp.StatusCode)
	}
}


