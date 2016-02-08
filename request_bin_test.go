package requestBin

import (
	"fmt"
	"net/http"
	"strings"
)

func ExampleCaptureRequests() {

	reqs := CaptureRequests(func(url string) {
		http.Post(url, "text/plain", strings.NewReader("Hello World!"))
	}, 200, 3)

	fmt.Println(reqs[0].Body)
	// Output: Hello World!

}
