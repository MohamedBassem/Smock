# RequestBin

[![Build Status](https://travis-ci.org/MohamedBassem/RequestBin.svg?branch=master)](https://travis-ci.org/MohamedBassem/RequestBin)
[![Coverage Status](https://coveralls.io/repos/github/MohamedBassem/RequestBin/badge.svg?branch=master)](https://coveralls.io/github/MohamedBassem/RequestBin?branch=master)
[![GoDoc](https://godoc.org/github.com/MohamedBassem/RequestBin?status.svg)](https://godoc.org/github.com/MohamedBassem/RequestBin)


Package requestBin is a package for testing the outgoing http requests initiated from a function. The package creates a mock server and passes the URL to the function to be tested and then collects all the requests that the server received.

### Example

```go
func ExampleCaptureRequests() {
	reqs := NewRequestBin(MockServerConfig{RequestTimeout: 1}).CaptureRequests(func(url string) {
		http.Post(url, "text/plain", strings.NewReader("Hello"))
		http.Post(url, "text/plain", strings.NewReader("It's me"))
	})

	fmt.Println(reqs[0].Body) // Hello
	fmt.Println(reqs[1].Body) // It's me
}
```

### API
The API is very simple and you can check it on [GoDoc](https://godoc.org/github.com/MohamedBassem/RequestBin). Reading the tests is also a good way for understanding the API.

### Contribution
Your contributions and ideas are welcomed through issues and pull requests.
