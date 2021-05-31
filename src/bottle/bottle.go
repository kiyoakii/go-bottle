package bottle

import (
	"net/http"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix string
	middlewares []HandlerFunc
	parent *RouterGroup			// for group nesting
	engine *Engine				// all groups share the same engine
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func New() *Engine {
	// mutual assignment
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup {
		prefix: group.prefix + prefix,
		parent: group,
		engine: group.engine,
	}
	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

//func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
//	engine.router.addRoute(method, pattern, handler)
//}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	//log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(w, req)
	engine.router.handle(ctx)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

