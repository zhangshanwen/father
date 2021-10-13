package father

import "net/http"

type (
	Router struct {
		Method   string
		Path     string
		Handler  HandlerFunc
		Handlers []HandlerFunc
	}
)

func (f *Father) addRouter(method, path string, handler HandlerFunc) {
	f.Routers = append(f.Routers, Router{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

func (f *Father) Post(path string, handler HandlerFunc) {
	f.addRouter(http.MethodPost, path, handler)
}
func (f *Father) Get(path string, handler HandlerFunc) {
	f.addRouter(http.MethodGet, path, handler)
}

func (f *Father) Put(path string, handler HandlerFunc) {
	f.addRouter(http.MethodPut, path, handler)
}
func (f *Father) Delete(path string, handler HandlerFunc) {
	f.addRouter(http.MethodDelete, path, handler)
}
func (f *Father) Patch(path string, handler HandlerFunc) {
	f.addRouter(http.MethodPatch, path, handler)
}
