package router

import (
	"net/http"
	"strings"
)

type MultiDomainRouter map[string]http.Handler

func (hs MultiDomainRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// remove port from given host
	host := r.Host
	_idx_port := strings.LastIndex(host, ":")
	if _idx_port != -1 {
		host = host[:_idx_port]
	}

	// lookup the hostname
	if handler := hs[host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		ShortenerHandler(w, r)
	}
}
