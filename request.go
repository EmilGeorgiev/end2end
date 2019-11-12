package end2end

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"gitlab.mailjet.tech/core/go/web"
)

// Requester make a call to server and can assert the returned response with expected one.
type Requester struct {
	url         string
	httpRequest *http.Request
	response    interface{}
	responseStatusCode int
}

// NewRequestToEndpoint create and return a new Requester.
func NewRequestToEndpoint(url string) Requester {
	return Requester{url: url}
}

// Create accept as parameters resource path and payload that you want to send to the server.
// Create a new http.Request and return the updated Requester.
func (r Requester) Create(jsonPayload string) Requester {
	req, err := http.NewRequest(http.MethodPost, r.url, strings.NewReader(jsonPayload))
	if err != nil {
		log.Fatal(err)
	}

	r.httpRequest = req
	
	return r
}

// Update accept as parameters resource path and payload that you want to send to the server.
// Create a new http.Request and return the updated Requester.
func (r Requester) Update(path string, payload interface{}) Requester {
	b, err :=json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(b)
	r.httpRequest, err = http.NewRequest(http.MethodPut, r.url, body)
	if err != nil {
		log.Fatal(err)
	}

	return r
}

// Delete accept as parameters resource path to the resource that you want to delete.
// Create a new http.Request and return the updated Requester.
func (r Requester) Delete(path string) Requester {
	req, err := http.NewRequest(http.MethodDelete, r.url, nil)
	if err != nil {
		log.Fatal(err)
	}

	r.httpRequest = req
	return r
}

// Get accept as parameters resource path to the resource that you want to get.
// Create a new http.Request and return the updated Requester.
func (r Requester) Get(filters string) Requester {
	req, err := http.NewRequest(http.MethodGet, r.url + filters, nil)
	if err != nil {
		log.Fatal(err)
	}

	r.httpRequest = req
	return r
}

func (r Requester) WithBasicAuth(userName, password string) Requester {
	r.httpRequest.SetBasicAuth(userName, password)
	return r
}

// Read accept as parameter a type in which you will store
func (r Requester) Read(response interface{}, statusCode int) Requester {
	r.response = response
	r.responseStatusCode = statusCode
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
	//r.actualResponse = actual
	//r.expectedResponse = expected
	return r
}

// ExpectStatusCode set the status code that you expect to be returned from the server.
func (r Requester) ExpectStatusCode(status int64) Requester {
	//r.expectedStatusCode = status
	return r
}

// Call make actual make request to the server and assert the returned response with expected one.
func (r Requester) Call(t *testing.T) {
	resp, err := http.DefaultClient.Do(r.httpRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if r.response == nil {
		return
	}

	if resp.StatusCode != r.responseStatusCode {
		if _, ok := r.response.(*web.Error); !ok {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Errorf("expect: %s %s return status code: %d and response type: %T", r.httpRequest.Method, r.httpRequest.URL.String(), r.responseStatusCode, r.response)
			t.Errorf("actual: %s %s return status code: %d and response     : %s", r.httpRequest.Method, r.httpRequest.URL.String(), resp.StatusCode, b)
			t.Fatal()
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&r.response); err != nil {
		t.Fatal(err)
	}
}


