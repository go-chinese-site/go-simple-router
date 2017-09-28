package router

import (
	"fmt"
	"net/http"
	"path"
)

type (
	//Router 路由器
	Router struct {
		//上级处理函数集合
		globalHandlers []HandlerFunc
		//上级路径
		basePath string
		//路由集合
		routers map[string]*route
	}
	//route 单个路由
	route struct {
		//请求方法
		method string
		//处理函数
		handlers []HandlerFunc
	}
	//Context 存储请求响应信息
	Context struct {
		Request *http.Request
		Writer  http.ResponseWriter
		//处理函数
		handlers []HandlerFunc
		//执行处理函数索引
		index int8
	}
	//HandlerFunc 处理函数
	HandlerFunc func(*Context)
)

//New 创建路由器
func New() *Router {
	return &Router{
		routers:  make(map[string]*route),
		basePath: "/",
	}
}

//Use 添加全局处理函数
func (r *Router) Use(handlers ...HandlerFunc) {
	r.globalHandlers = append(r.globalHandlers, handlers...)
}

//Group 添加分组
func (r *Router) Group(partPath string, fn func(), handlers ...HandlerFunc) {
	rootBasePath := r.basePath
	rootHandlers := r.globalHandlers
	r.basePath = path.Join(r.basePath, partPath)
	r.globalHandlers = r.combinHandlers(handlers)
	fn()
	r.basePath = rootBasePath
	r.globalHandlers = rootHandlers
}

//GET 添加GET方法
func (r *Router) GET(partPath string, handlers ...HandlerFunc) {
	path := path.Join(r.basePath, partPath)
	handlers = r.combinHandlers(handlers)
	r.addRoute("GET", path, handlers)
}

//POST 添加POST方法
func (r *Router) POST(partPath string, handlers ...HandlerFunc) {
	path := path.Join(r.basePath, partPath)
	handlers = r.combinHandlers(handlers)
	r.addRoute("POST", path, handlers)
}

//Run 运行
func (r *Router) Run(addr string) error {
	return http.ListenAndServe(addr, r)
}

//combinHandlers 合并处理函数
func (r *Router) combinHandlers(handlers []HandlerFunc) []HandlerFunc {
	finallyLen := len(r.globalHandlers) + len(handlers)
	finallyHandlers := make([]HandlerFunc, finallyLen)
	copy(finallyHandlers, r.globalHandlers)
	copy(finallyHandlers[len(r.globalHandlers):], handlers)
	return finallyHandlers
}

//addRoute 添加到路由中
func (r *Router) addRoute(method, path string, handlers []HandlerFunc) {
	route := &route{
		method:   method,
		handlers: handlers,
	}
	r.routers[path] = route
}

//ServeHTTP 实现http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	httpmethod := req.Method
	path := req.URL.Path
	route, ok := r.routers[path]
	if !ok {
		fmt.Fprintf(w, "not found page")
		return
	}
	if route.method != httpmethod {
		fmt.Fprintf(w, "http method unsupported")
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

//Next 调用下一个处理函数
func (c *Context) Next() {
	c.index++
	if n := int8(len(c.handlers)); c.index < n {
		c.handlers[c.index](c)
	}
}
