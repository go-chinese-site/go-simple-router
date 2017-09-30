package router

import (
	"fmt"
	"net/http"
	"path"
)

const (
	rootPath = "/"

	// ErrNotFound is not found error.
	ErrNotFound = "not found page"

	// ErrMethodUnsupported is method upsupported error.
	ErrMethodUnsupported = "http method unsupported"
)

// Router is responsible for managing multiple routes.
type Router struct {
	globalHandlers []HandlerFunc
	basePath       string
	routers        map[string]*route
}

// Signle route
type route struct {
	method   string
	handlers []HandlerFunc
}

// Context stores information about request and response.
type Context struct {
	// Request represents http request.
	Request *http.Request

	// Writer represents http response.
	Writer   http.ResponseWriter
	handlers []HandlerFunc
	index    int8
}

// HandlerFunc represents the function of Context.
type HandlerFunc func(*Context)

// New represents creating a router.
func New() *Router {
	return &Router{
		routers:  make(map[string]*route),
		basePath: rootPath,
	}
}

// Use represents the method of Router.
func (r *Router) Use(handlers ...HandlerFunc) {
	r.globalHandlers = append(r.globalHandlers, handlers...)
}

// Group represents the method of Router.
func (r *Router) Group(partPath string, fn func(), handlers ...HandlerFunc) {
	rootBasePath := r.basePath
	rootHandlers := r.globalHandlers
	r.basePath = path.Join(r.basePath, partPath)
	r.globalHandlers = r.combinHandlers(handlers)
	fn()
	r.basePath = rootBasePath
	r.globalHandlers = rootHandlers
}

// GET represents the method of Router.
func (r *Router) GET(partPath string, handlers ...HandlerFunc) {
	path := path.Join(r.basePath, partPath)
	handlers = r.combinHandlers(handlers)
	r.addRoute(http.MethodGet, path, handlers)
}

// POST represents the method of Router.
func (r *Router) POST(partPath string, handlers ...HandlerFunc) {
	path := path.Join(r.basePath, partPath)
	handlers = r.combinHandlers(handlers)
	r.addRoute(http.MethodPost, path, handlers)
}

// Run start the ListenAndServe of http.
func (r *Router) Run(addr string) error {
	return http.ListenAndServe(addr, r)
}

func (r *Router) combinHandlers(handlers []HandlerFunc) []HandlerFunc {
	finallyLen := len(r.globalHandlers) + len(handlers)
	finallyHandlers := make([]HandlerFunc, finallyLen)
	copy(finallyHandlers, r.globalHandlers)
	copy(finallyHandlers[len(r.globalHandlers):], handlers)
	return finallyHandlers
}

func (r *Router) addRoute(method, path string, handlers []HandlerFunc) {
	route := &route{
		method:   method,
		handlers: handlers,
	}
	r.routers[path] = route
}

// ServeHTTP implemented the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	httpmethod := req.Method
	path := req.URL.Path
	route, ok := r.routers[path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, ErrNotFound)
		return
	}
	if route.method != httpmethod {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, ErrMethodUnsupported)
		return
	}
	c := &Context{
		Request:  req,
		Writer:   w,
		handlers: route.handlers,
		index:    -1,
	}
	c.Next()
}

// Next call the next method.
func (c *Context) Next() {
	c.index++
	if n := int8(len(c.handlers)); c.index < n {
		c.handlers[c.index](c)
	}
}
