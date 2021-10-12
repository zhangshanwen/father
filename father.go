package father

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

var (
	DefaultNotFound ErrResponseFunc = func(c *Context, w http.ResponseWriter) {
		c.SetStatusCode(http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
	}
	DefaultMethodNotAllowed ErrResponseFunc = func(c *Context, w http.ResponseWriter) {
		c.SetStatusCode(http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	DefaultRouterLog RouterLogFunc = func(c *Context) string {
		return fmt.Sprintf("[%v]:%v:%v -%v- const-->>%v s", c.Req.Method, c.Req.Host, c.Req.RequestURI, c.StatusCode, c.ConstTime)
	}
)

type ErrResponseFunc func(c *Context, w http.ResponseWriter)
type RouterLogFunc = func(c *Context) string
type (
	Father struct {
		Address string
		Routers []Router
		logger  Logger
	}
	Router struct {
		Method  string
		Path    string
		Handler HandlerFunc
	}
)

type HandlerFunc func(c *Context)

func NewFather() *Father {

	return &Father{
		Routers: []Router{},
		logger:  DefaultLogger,
	}
}

func (f *Father) SetDefaultLogger(logger Logger) {
	f.logger = logger
}

func (f *Father) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startAt := time.Now().Unix()
	c := Context{
		Req:        req,
		Writer:     w,
		StatusCode: http.StatusOK,
	}
	defer func() {
		c.ConstTime = time.Now().Unix() - startAt
		f.logger.Println(DefaultRouterLog(&c))
	}()
	// 遍历 路由
	path := req.URL.Path
	method := req.Method
	for i := 0; i < len(f.Routers); i++ {
		router := f.Routers[i]
		if path != router.Path {
			continue
		}
		if method == router.Method {
			router.Handler(&c)
			return
		}
		DefaultMethodNotAllowed(&c, w)
		return
	}
	DefaultNotFound(&c, w)
}

func (f *Father) Run(host string, port int) (err error) {
	if host == "" {
		host = defaultHost
	}
	if port == 0 {
		port = defaultPort
	}
	f.Address = fmt.Sprintf("%v:%v", host, port)
	//TODO 检查路由是否重复
	if err = f.checkRouters(); err != nil {
		f.logger.Fatalf("服务启动失败,err=%v", err)
	}
	return http.ListenAndServe(f.Address, f)
}
func (f *Father) checkRouters() (err error) {
	routerMap := map[string]bool{}
	for i := 0; i < len(f.Routers); i++ {
		originLen := len(routerMap)
		router := f.Routers[i]
		routerMap[router.Method+router.Path] = true
		if originLen == len(routerMap) {
			f.logger.Fatal("重复路由,[%v]:%v", router.Method, router.Path)
		}
		fn := runtime.FuncForPC(reflect.ValueOf(router.Handler).Pointer()).Name()
		f.logger.Printf("路由------>>>>>>[%v]:[%v]---->%v", router.Method, router.Path, fn)
	}

	return
}

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

/*
暂时使用 Get Post Put Delete Patch 方法
*/
