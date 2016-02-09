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
	config     MockServerConfig
	// The channel in which each new request received by the mock server will be passed to.
	Reqs *channels.InfiniteChannel
}

type MockRequest struct {
	Body    string
	Headers http.Header
	Method  string
}

type MockServerConfig struct {
	// The response code returned by the mock server to all incoming requests
	ResponseStatusCode int
	// The number of seconds to stop the mock server after
	GlobalTimeout int
	// The number of seconds to stop the mock server after if the server didn't recieve any request within this period
	RequestTimeout int
}

// Stops the mock server and closes the requests channel
func (m *MockServer) Start() {
	m.testServer.Start()
	m.Reqs = channels.NewInfiniteChannel()
}

// Stops the mock server and closes the requests channel
func (m *MockServer) Close() {
	m.testServer.Close()
	m.Reqs.Close()
	m.Reqs = nil
}

// Returns the URL of the created MockServer
func (m *MockServer) URL() string {
	return m.testServer.URL
}

func mergeDefaultConfigs(config *MockServerConfig) {
	if config.ResponseStatusCode == 0 {
		config.ResponseStatusCode = 200
	}
}

// NewRequestBin returns a MockServer instance that you can use on your own or use the CaptureRequests function. Either the GlobalTimeout or the Requesttimeout must be set or the function will panic.
// If to be used on your own mockServer.Start() should be called first before starting to use the server and mockServer.Close() should be called to stop the server.
func NewRequestBin(config MockServerConfig) *MockServer {

	mergeDefaultConfigs(&config)
	if config.GlobalTimeout == 0 && config.RequestTimeout == 0 {
		panic("Either the GlobalTimeout or the RequestTimeout must be set or both")
	}

	server := MockServer{
		config: config,
	}

	server.testServer = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(config.ResponseStatusCode)
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		server.Reqs.In() <- MockRequest{
			Body:    string(body),
			Headers: r.Header,
			Method:  r.Method,
		}
	}))

	return &server
}

// CaptureRequests takes a function to be tested. It starts the mock server, runs the function passing the URL of the mock server to it, and any request that is sent to this URL is logged. Depending on the MockServerConfig, the server can be stopped by two ways. The first one is the global timeout which stops the server after a fixed number of seconds. The second one is a timeout between requests, if the server didn't recieve any request within <config.RequestTimeout> seconds it will exit.
func (s *MockServer) CaptureRequests(f func(string)) []MockRequest {

	s.Start()
	defer s.Close()

	f(s.URL())

	gtChan := time.After(time.Second * time.Duration(s.config.GlobalTimeout))

	ret := make([]MockRequest, 0)

MainLoop:
	for {
		if len(gtChan) > 0 {
			break MainLoop
		}
		select {
		case in := <-s.Reqs.Out():
			request, _ := in.(MockRequest)
			ret = append(ret, request)
		case <-time.After(time.Second * time.Duration(s.config.RequestTimeout)):
			break MainLoop
		case <-gtChan:
			break MainLoop
		}
	}

	return ret
}
