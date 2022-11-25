package router

import (
	"net/http"
)

type MultiDomainRouter map[string]http.Handler

func (hs MultiDomainRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		ShortenerHandler(w, r)
	}
}
