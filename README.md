## This is a test library that can be used for end2end testing of HTTP server. The server must except and return JSON data.

## Example

```
myServer = end2end.Config{URL: "http://localhost:8080"}

var resp User
end2end.NewRequestTo(myServer)
    .Create("/users", User{Name: "Ivan"})
    .Read(&resp)
    .ExpectStatusCode(http.StatusCreated)
    .Call(t)

var actual User
expected := User{Name: "Ivan", GUID: resp.GUID}
end2end.NewRequestTo(myServer)
    .Get("/users/" + resp.GUID)
    .Assert(&actual, &expected)
    .ExpectedStatusCode(http.StatusOK)
    .Call(t)
```

The library has dependency to github.com/stretchr/testify/assert