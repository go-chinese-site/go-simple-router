package router

import (
	"fmt"
	"net/http"
	"path"
)

const (
	// ErrNotFound is not found error.
	ErrNotFound = "not found page"
	// ErrMethodUnsupported is method upsupported error.
	ErrMethodUnsupported = "http method unsupported"
	// RouterKey is route key format.
	RouterKey = "%s-%s"
)

type (
	// Router is a http.Handler. store all routes.
	Router struct {
		// store parent all HandlerFunc.
		globalHandlers []HandlerFunc
		// store parent path.
		basePath string
		// store all routes.
		routers map[string]*route
	}
	// route is storage http method and handle function.
	route struct {
		// http method.
		method string
		// handle function.
		handlers []HandlerFunc
	}
	// Context is storage request response information.
	Context struct {
		Request *http.Request
		Writer  http.ResponseWriter
		// handle function.
		handlers []HandlerFunc
		// The current handle function index is executed.
		index int8
	}
	// HandlerFunc is a function that can be registered to a route to handle HTTP requests.
	HandlerFunc func(*Context)
)

// New is returns an initialized Router.
func New() *Router {
	return &Router{
		routers:  make(map[string]*route),
		basePath: "/",
	}
}

// Use is add global handle function.
func (r *Router) Use(handlers ...HandlerFunc) {
	r.globalHandlers = append(r.globalHandlers, handlers...)
}

// Group is add route group.
func (r *Router) Group(partPath string, fn func(), handlers ...HandlerFunc) {
	rootBasePath := r.basePath
	rootHandlers := r.globalHandlers
	r.basePath = path.Join(r.basePath, partPath)
	r.globalHandlers = r.combineHandlers(handlers)
	fn()
	r.basePath = rootBasePath
	r.globalHandlers = rootHandlers
}

// GET is register GET method HandlerFunc to Router.
func (r *Router) GET(partPath string, handlers ...HandlerFunc) {
	path := path.Join(r.basePath, partPath)
	handlers = r.combineHandlers(handlers)
	r.addRoute(http.MethodGet, path, handlers)
}

// POST is register POST method HandlerFunc to Router.
func (r *Router) POST(partPath string, handlers ...HandlerFunc) {
	path := path.Join(r.basePath, partPath)
	handlers = r.combineHandlers(handlers)
	r.addRoute(http.MethodPost, path, handlers)
}

// Run listens on the TCP network address addr.
func (r *Router) Run(addr string) error {
	return http.ListenAndServe(addr, r)
}

// combineHandlers is merge multiple HnalderFunc slice into one HandlerFunc slice.
func (r *Router) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finallyLen := len(r.globalHandlers) + len(handlers)
	finallyHandlers := make([]HandlerFunc, finallyLen)
	copy(finallyHandlers, r.globalHandlers)
	copy(finallyHandlers[len(r.globalHandlers):], handlers)
	return finallyHandlers
}

// addRoute is add to routes.
func (r *Router) addRoute(method, path string, handlers []HandlerFunc) {
	route := &route{
		method:   method,
		handlers: handlers,
	}
	r.routers[fmt.Sprintf(RouterKey, path, method)] = route
}

// ServeHTTP is implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	httpMethod := req.Method
	path := req.URL.Path
	route, ok := r.routers[fmt.Sprintf(RouterKey, path, httpMethod)]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, ErrNotFound)
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

// Next is call the next handler function.
func (c *Context) Next() {
	c.index++
	if n := int8(len(c.handlers)); c.index < n {
		c.handlers[c.index](c)
	}
}
