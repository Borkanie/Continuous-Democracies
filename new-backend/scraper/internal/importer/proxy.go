package importer

import (
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

// socksTransport returns an HTTP transport that routes through the given SOCKS5 proxy URL.
func socksTransport(proxyURL string) http.RoundTripper {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return http.DefaultTransport
	}

	dialer, err := proxy.FromURL(parsedURL, proxy.Direct)
	if err != nil {
		return http.DefaultTransport
	}

	return &http.Transport{Dial: dialer.Dial}
}
