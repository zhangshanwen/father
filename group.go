package father

import (
	"errors"
	"net/http"
	"strings"
)

var (
	RepeatedRouterError = errors.New("RepeatedRouterError")
)

type (
	Group struct {
		Path     string
		Method   string
		Father   *Group
		Children []*Group
		Handlers []HandlerFunc
	}
)

func (g *Group) New(path string) *Group {
	return &Group{
		Path:     path,
		Handlers: []HandlerFunc{},
		Father:   g,
		Children: []*Group{},
	}
}

func (g *Group) Middleware(handlers ...HandlerFunc) {
	g.Handlers = append(g.Handlers, handlers...)
}

func (g *Group) addRouter(method, path string, handlers ...HandlerFunc) {
	g.Children = append(g.Children, &Group{
		Method:   method,
		Path:     path,
		Handlers: append(g.Handlers, handlers...),
	})
}
func (g *Group) GetPath(path string) string {
	if !strings.HasSuffix(path, httpSeparator) {
		if len(g.Path) == 0 || strings.HasPrefix(g.Path, httpSeparator) {
			return path + g.Path
		}
		g.Path = httpSeparator + g.Path
		return g.GetPath(path)
	}
	return g.GetPath(path[:len(path)-1])
}

func (g *Group) Post(path string, handler HandlerFunc) {
	g.addRouter(http.MethodPost, path, handler)
}
func (g *Group) Get(path string, handler HandlerFunc) {
	g.addRouter(http.MethodGet, path, handler)
}

func (g *Group) Put(path string, handler HandlerFunc) {
	g.addRouter(http.MethodPut, path, handler)
}
func (g *Group) Delete(path string, handler HandlerFunc) {
	g.addRouter(http.MethodDelete, path, handler)
}
func (g *Group) Patch(path string, handler HandlerFunc) {
	g.addRouter(http.MethodPatch, path, handler)
}
