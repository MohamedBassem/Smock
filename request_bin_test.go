package requestBin

import (
	"fmt"
	"net/http"
	"strings"
)

func ExampleCaptureRequests() {

	reqs := CaptureRequests(func(url string) {
		http.Post(url, "text/plain", strings.NewReader("Hello"))
		http.Post(url, "text/plain", strings.NewReader("It's me"))
	}, 200, 3)

	fmt.Println(reqs[0].Body)
	fmt.Println(reqs[1].Body)
	// Output: Hello
	// It's me

}
