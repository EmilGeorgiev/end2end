package end2end_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/EmilGeorgiev/end2end"
)

type GithubResponse struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

func TestRequestWithBasicAuth(t *testing.T) {
	got := GithubResponse{}

	end2end.NewRequest(http.MethodGet, "https://api.github.com", "").
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
