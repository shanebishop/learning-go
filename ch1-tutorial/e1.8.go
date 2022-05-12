// Modify fetch to add the prefix http:// to each argument
// URL if it is missing. You might want to use
// strings.HasPrefix.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		// NOTE: This won't cover non-HTTP schemes, such as FTP. I.e., if a user
		// used ftp://some/path, then this could would fetch http://ftp://some/path,
		// and fail.
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = fmt.Sprintf("http://%s", url)
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}

		fmt.Printf("%s", b)
	}
}
