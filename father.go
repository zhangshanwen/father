package father

import (
	"fmt"
	"net/http"
)

var (
	DefaultNotFound NotFoundFunc = func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusNotFound)
	}
)

type NotFoundFunc func(w http.ResponseWriter)

type Father struct {
	Address string
	FuncMap map[string]HandlerFunc
}

type HandlerFunc func(http.ResponseWriter)

func NewFather() *Father {

	return &Father{
		FuncMap: make(map[string]HandlerFunc),
	}
}

func (f *Father) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			DefaultNotFound(w)
		}
	}()
	handler := f.FuncMap[req.Method+req.RequestURI]
	handler(w)
}

func (f *Father) Run(host string, port int) error {
	if host == "" {
		host = defaultHost
	}
	if port == 0 {
		port = defaultPort
	}
	f.Address = fmt.Sprintf("%v:%v", host, port)
	return http.ListenAndServe(f.Address, f)
}

func (f *Father) Post(path string, handler HandlerFunc) {
	f.FuncMap[http.MethodPost+path] = handler
}
func (f *Father) Get(path string, handler HandlerFunc) {
	f.FuncMap[http.MethodGet+path] = handler
}

func (f *Father) Put(path string, handler HandlerFunc) {
	f.FuncMap[http.MethodPut+path] = handler
}
func (f *Father) Delete(path string, handler HandlerFunc) {
	f.FuncMap[http.MethodDelete+path] = handler
}
func (f *Father) Head(path string, handler HandlerFunc) {
	f.FuncMap[http.MethodHead+path] = handler
}
func (f *Father) Connect(path string, handler HandlerFunc) {
	f.FuncMap[http.MethodConnect+path] = handler
}
func (f *Father) Options(path string, handler HandlerFunc) {
	f.FuncMap[http.MethodOptions+path] = handler
}
func (f *Father) Patch(path string, handler HandlerFunc) {
	f.FuncMap[http.MethodPatch+path] = handler
}
func (f *Father) Trace(path string, handler HandlerFunc) {
	f.FuncMap[http.MethodTrace+path] = handler
}
