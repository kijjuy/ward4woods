package application

import (
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) AddRoute(url string, handlerFunc http.HandlerFunc) {
	r.mux.HandleFunc(url, handlerFunc)
}

func (r *Router) Serve() *http.ServeMux {
	return r.mux
}
