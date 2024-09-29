package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

func main() {
	http.HandleFunc("/", handleRequest)
	log.Println("Server starting on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Collect headers
	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = strings.Join(values, ", ")
	}
	// Add Host to headers if not already present
	if _, exists := headers["Host"]; !exists {
		headers["Host"] = r.Host
	}

	// Collect other request info
	info := map[string]string{
		"Method":     r.Method,
		"RequestURI": r.RequestURI,
		"RemoteAddr": r.RemoteAddr,
		"Protocol":   r.Proto,
	}

	// Add TLS info if present
	if r.TLS != nil {
		info["TLS-Version"] = versionToString(r.TLS.Version)
		info["TLS-CipherSuite"] = cipherSuiteToString(r.TLS.CipherSuite)
		info["TLS-ServerName"] = r.TLS.ServerName
	}

	// Print headers
	fmt.Fprintln(w, "Headers received:")
	fmt.Fprintln(w, "==================")
	printSorted(w, headers)

	// Print other request info
	fmt.Fprintln(w, "\nOther request information:")
	fmt.Fprintln(w, "===========================")
	printSorted(w, info)

	// Print JSON format
	fmt.Fprintln(w, "\nComplete data in JSON format:")
	fmt.Fprintln(w, "==============================")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"headers": headers,
		"info":    info,
	})
}

func printSorted(w http.ResponseWriter, data map[string]string) {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s: %s\n", k, data[k])
	}
}

func versionToString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (%d)", version)
	}
}

func cipherSuiteToString(cipherSuite uint16) string {
	switch cipherSuite {
	case tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:
		return "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:
		return "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
	// Add more cases as needed
	default:
		return fmt.Sprintf("Unknown (%d)", cipherSuite)
	}
}
