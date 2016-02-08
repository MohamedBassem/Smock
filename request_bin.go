// RequestBin is a package for testing the http requests initiated from your application. It takes a function and returns back all the requests that were initiated by this function.
package requestBin

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/eapache/channels"
)

type MockServer struct {
	// The mock server instance
	testServer *httptest.Server
	// The channel in which each new request received by the mock server will be passed to.
	Reqs *channels.InfiniteChannel
}

type MockRequest struct {
	Body    string
	Headers http.Header
}

// Stops the mock server and closes the requests channel
func (m *MockServer) Close() {
	m.testServer.Close()
	m.Reqs.Close()
}

// Returns the URL of the created MockServer
func (m *MockServer) URL() string {
	return m.testServer.URL
}

// NewRequestBin returns a MockServer instance that you can use on your own instead of using the higher level CaptureRequests function.
func NewRequestBin(responseStatusCode int) MockServer {

	server := MockServer{
		Reqs: channels.NewInfiniteChannel(),
	}

	server.testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(responseStatusCode)
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		server.Reqs.In() <- MockRequest{
			Body:    string(body),
			Headers: r.Header,
		}
	}))

	return server
}

// CaptureRequests takes a function that we need to test, the response status code that's returned from the mockserver and a timeout. CaptureRequests passes the URL of the mock server to the function passed, and any request that is sent to this URL is logged. The server is stopped and the function returns whenever there aren't any new requests within <timeout> seconds from the last request.
func CaptureRequests(f func(string), responseStatusCode, timeout int) []MockRequest {

	server := NewRequestBin(responseStatusCode)
	defer server.Close()

	f(server.URL())

	ret := make([]MockRequest, 0)

	for {
		select {
		case in := <-server.Reqs.Out():
			request, _ := in.(MockRequest)
			ret = append(ret, request)
		case <-time.After(time.Second * time.Duration(timeout)):
			return ret
		}
	}
}
