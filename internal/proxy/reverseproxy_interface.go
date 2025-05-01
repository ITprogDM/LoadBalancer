package proxy

import "net/http"

type ReverseProxy interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
