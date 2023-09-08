package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type debugTransport struct {
	transport http.RoundTripper
}

func NewDebugTransport() *debugTransport {
	return &debugTransport{
		transport: http.DefaultTransport,
	}
}

func (d *debugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Println("[Logging:Proxy]", r.URL)
	if os.Getenv("DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "\n=====> %s\n", r.URL)
		dump, _ := httputil.DumpRequest(r, true)
		fmt.Fprintf(os.Stderr, string(dump[:len(dump)]))
	}
	resp, err := d.transport.RoundTrip(r)
	if os.Getenv("DEBUG") != "" && r.Method == http.MethodGet {
		fmt.Fprintln(os.Stderr, "---RESONSE---")
		dump, _ := httputil.DumpResponse(resp, true)
		fmt.Fprintf(os.Stderr, string(dump[:len(dump)]))
	}
	return resp, err
}

type LogginRT struct {
	transport http.Transport
}

func (rt LogginRT) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Println("[Logging:Proxy]", req.URL)
	//log.Printf("FULL req: %#v", req)
	return http.DefaultTransport.RoundTrip(req)
}

func main() {
	tf := flag.String("t", "https://444.hu", "target url to reverseproxy")
	flag.Parse()
	target, err := url.Parse(*tf)
	if err != nil {
		panic(err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)
			// pr.Out.Header.Set("Referer", "https://444.hu")
		},
		// Transport: LogginRT{},
		Transport: NewDebugTransport(),
	}

	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
