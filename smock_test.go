package smock

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleCaptureRequests() {

	reqs := NewMockServer(MockServerConfig{RequestTimeout: 1}).CaptureRequests(func(url string) {
		http.Post(url, "text/plain", strings.NewReader("Hello"))
		http.Post(url, "text/plain", strings.NewReader("It's me"))
	})

	fmt.Println(reqs[0].Body)
	fmt.Println(reqs[1].Body)
	// Output: Hello
	// It's me

}

func TestGlobalTimeout(t *testing.T) {
	reqs := NewMockServer(MockServerConfig{GlobalTimeout: 1}).CaptureRequests(func(url string) {
		http.Post(url, "text/plain", strings.NewReader("Hello"))
		go func() {
			<-time.After(time.Second * 3)
			http.Post(url, "text/plain", strings.NewReader("It's me"))
		}()
	})

	assert.Equal(t, 1, len(reqs))
}

func TestRequestTimeout(t *testing.T) {
	reqs := NewMockServer(MockServerConfig{RequestTimeout: 2}).CaptureRequests(func(url string) {
		http.Post(url, "text/plain", strings.NewReader("Hello"))
		http.Post(url, "text/plain", strings.NewReader("Hello"))
		go func() {
			<-time.After(time.Second * 3)
			http.Post(url, "text/plain", strings.NewReader("It's me"))
		}()
	})

	assert.Equal(t, 2, len(reqs))
}

func TestMaximumRequestCount(t *testing.T) {
	reqs := NewMockServer(MockServerConfig{MaximumRequestCount: 2}).CaptureRequests(func(url string) {
		http.Post(url, "text/plain", strings.NewReader("Hello"))
		http.Post(url, "text/plain", strings.NewReader("Hello"))
		http.Post(url, "text/plain", strings.NewReader("Hello"))
	})

	assert.Equal(t, 2, len(reqs))
}

func TestNormalFlow(t *testing.T) {

	reqs := NewMockServer(MockServerConfig{ResponseStatusCode: 201, MaximumRequestCount: 2}).CaptureRequests(func(url string) {
		http.Post(url, "text/plain", strings.NewReader("Hello"))

		req, _ := http.NewRequest("DELETE", url, strings.NewReader("It's me"))
		req.Header.Add("Content-Type", "text/plain")

		resp, _ := (&http.Client{}).Do(req)
		assert.Equal(t, 201, resp.StatusCode)
	})

	assert.Equal(t, "Hello", reqs[0].Body)
	assert.Equal(t, "POST", reqs[0].Method)
	assert.Equal(t, "text/plain", reqs[0].Headers.Get("Content-Type"))

	assert.Equal(t, "It's me", reqs[1].Body)
	assert.Equal(t, "DELETE", reqs[1].Method)
	assert.Equal(t, "text/plain", reqs[1].Headers.Get("Content-Type"))

}

func TestPanicOnMissingTimeouts(t *testing.T) {
	assert.Panics(t, func() {
		NewMockServer(MockServerConfig{})
	})
}
