// RequestBin is a package for testing the http requests initiated from your application. It takes a function and returns back all the requests that were initiated by this function.

package requestBin

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

type mockServer struct {
	// The mock server instance
	testServer *httptest.Server

	// The channel in which each new request received by the mock server will be passed to.
	Reqs chan mockRequest
}

type mockRequest struct {
	Body    string
	Headers http.Header
}

// Stops the mock server and closes the requests channel
func (m *mockServer) Close() {
	m.testServer.Close()
	close(m.Reqs)
}

// Returns the URL of the created mockServer
func (m *mockServer) URL() string {
	return m.testServer.URL
}

// NewRequestBin returns a mockServer instance that you can use on your own instead of using the higher level CaptureRequests function.
func NewRequestBin(responseStatusCode int) mockServer {
	reqs := make(chan mockRequest)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(responseStatusCode)
		body, _ := ioutil.ReadAll(r.Body)
		reqs <- mockRequest{
			Body:    string(body),
			Headers: r.Header,
		}
	}))

	return mockServer{
		Reqs:       reqs,
		testServer: testServer,
	}
}

// CaptureRequests takes a function that we need to test, the response status code that's returned from the mockserver and a timeout. CaptureRequests passes the URL of the mock server to the function passed, and any request that is sent to this URL is logged. The server is stopped and the function returns whenever there isn't any new requests within <timeout> seconds from the last request.
func CaptureRequests(f func(string), responseStatusCode, timeout int) []mockRequest {

	server := NewRequestBin(responseStatusCode)
	defer server.Close()

	f(server.URL())

	ret := make([]mockRequest, 0)

	for {
		select {
		case request := <-server.Reqs:
			ret = append(ret, request)
		case <-time.After(time.Second * time.Duration(timeout)):
			return ret
		}
	}
}
