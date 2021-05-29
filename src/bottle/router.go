package bottle

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	// trie roots for "GET", "POST", etc...
	roots map[string]*node
	// handlers for all routes
	handlers map[string]HandlerFunc
}

// examples for router:
// roots['GET'], roots['POST']
// handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parse(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := parse(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	realParts := parse(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(realParts, 0)

	if n != nil {
		parts := parse(n.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = realParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(realParts[i:], "/")
				break
			}
			return n, params
		}
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	handler, ok := r.handlers[key]
	if ok {
		handler(c)
	} else {
		c.Text(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
