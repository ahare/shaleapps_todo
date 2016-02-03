package server

import (
	"net/http"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	method  string
	handler http.Handler
}

type router struct {
	routes []*route
}

func (r *router) HandleFunc(routePattern, method string,
	handler func(http.ResponseWriter, *http.Request)) {
	pattern := regexp.MustCompile(routePattern)
	r.routes = append(r.routes, &route{pattern, method, http.HandlerFunc(handler)})
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.pattern.MatchString(req.URL.Path) && req.Method == route.method {
			route.handler.ServeHTTP(w, req)
			return
		}
	}
	http.NotFound(w, req)
}
