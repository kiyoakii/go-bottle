package bottle

import (
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/:arg/ext/", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParse(t *testing.T) {
	ok := reflect.DeepEqual(parse("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parse("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parse("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parse failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, paras := r.getRoute("GET", "/hello/jin")

	if n == nil || n.pattern != "/hello/:name"{
		t.Fatal("test path /hello/:name failed")
	}

	if paras["name"] != "jin" {
		t.Fatalf("parameters parsing failed, expected \"jin\", get %s", paras["name"])
	}
}
