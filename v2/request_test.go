package v2_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/EmilGeorgiev/end2end/v2"
)

type GithubResponse struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

func TestRequestWithBasicAuth(t *testing.T) {
	got := GithubResponse{}

	v2.NewRequest(http.MethodGet, "https://api.github.com", "").
		WithBasicAuth("fooo", "barr").
		Expect(&got, http.StatusUnauthorized).
		Call(t)

	want := GithubResponse{
		Message:  "Bad credentials",
		DocumentationURL: "https://developer.github.com/v3",
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected: %#v\n", want)
		t.Errorf("got     : %#v\n", got)
	}
}

func TestRequestWithExpectNil(t *testing.T) {
	v2.NewRequest(http.MethodGet, "https://api.github.com", "").
		WithBasicAuth("fooo", "barr").
		Expect(nil, http.StatusUnauthorized).
		Call(t)
}

func TestSendRequestWithoutAssertTheResponse(t *testing.T) {
	v2.NewRequest(http.MethodGet, "https://api.github.com", "").
		WithBasicAuth("fooo", "barr").
		Call(t)
}

func TestSendRequestWithoutCheckStatusCode(t *testing.T) {
	resp, err := v2.NewRequest(http.MethodGet, "https://api.github.com", "").
		WithBasicAuth("fooo", "barr").
		Send()
	if err != nil {
		t.Error("unexpected error:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expect status code %d:", http.StatusUnauthorized)
		t.Errorf("got staus code     %d:", resp.StatusCode)
	}
}