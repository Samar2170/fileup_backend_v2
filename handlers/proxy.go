package handlers

import (
	"io"
	"net"
	"net/http"
	"time"
)

// var customTransport = http.DefaultTransport
var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func HandleImageProxy(w http.ResponseWriter, r *http.Request) {
	targetUrl := r.URL.Query().Get("resource")
	proxyReq, err := http.NewRequest(r.Method, targetUrl, nil)
	if err != nil {
		http.Error(w, "Creating proxy reqyest", http.StatusInternalServerError)
	}

	proxyReq.Header.Set("Content-Type", "application/octet-stream")
	proxyReq.Header.Set("accept", "*/*")
	// proxyReq.Header.Set("AccessKey", config.StorageAccessKey)

	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(w, resp.Body)
	// resp, err := customTransport.RoundTrip(proxyReq)
	// if err != nil {
	// 	http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
	// 	return
	// }
	// defer resp.Body.Close()
	// for name, values := range resp.Header {
	// 	for _, value := range values {
	// 		w.Header().Add(name, value)
	// 	}
	// }
	// w.WriteHeader(resp.StatusCode)

	// io.Copy(w, resp.Body)
}
