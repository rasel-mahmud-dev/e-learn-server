package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main34() {
	proxyTarget := "http://localhost:3200/"

	// Parse the target URL
	targetURL, err := url.Parse(proxyTarget)
	if err != nil {
		fmt.Println("Error parsing proxy target URL:", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Handle requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Proxying request for %s to %s\n", r.URL.Path, targetURL)
		proxy.ServeHTTP(w, r)
	})

	// Start the proxy server
	fmt.Println("Proxy server listening on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting proxy server:", err)
		return
	}
}
