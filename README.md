## This is a test library that can be used for end2end testing of HTTP server. The server must except and return JSON data.

## Example

```go
func TestHappyPath(t *testing.T) {
    baseURL = "http://localhost:8080"
    userName = "AK10003"
    password = "SK10003"

    newUser := User{
        Name: "Emil",
        Address: "London",
        Age: 31,
    }
    var response api.AuditLog
    end2end.NewRequestToEndpoint(baseURL + "/v1/api/users").
        Create(newUser).
        WithBasicAuth(userName, password).
        Read(&response, http.StatusCreated).
        Call(t)

    var got []User
    end2end.NewRequestToEndpoint(baseURL + "/v1/api/users").
        Get("").
        WithBasicAuth(userName, password).
        Read(&got, http.StatusOK).
        Call(t)

    want := []User{
        {
            ID: result.ID,
            Name: "Emil",
            Address: "London",
            Age: 31,
        },
    }

    assert.Equal(t, want, got)
}
```

The library has dependency to github.com/stretchr/testify/assert