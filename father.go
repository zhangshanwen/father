package father

import (
	"fmt"
	"net/http"
	"time"
)

const (
	ConstTypeUnix ConstTimeType = iota
	ConstTypeUnixMilli
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
		return fmt.Sprintf("[%v]:%v:%v -%v- cost-->>%v%v", c.Req.Method, c.Req.Host, c.Req.RequestURI, c.StatusCode, c.ConstTime, c.ConstTimeUint)
	}
)

type (
	ErrResponseFunc func(c *Context, w http.ResponseWriter)
	RouterLogFunc   = func(c *Context) string
	ConstTimeType   int
)
type (
	Father struct {
		Address       string
		Routers       []Router
		logger        Logger
		groups        []*Group
		constTimeType ConstTimeType
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
func (f *Father) getNow() int64 {
	now := time.Now()
	switch f.constTimeType {
	case ConstTypeUnixMilli:
		return now.UnixMilli()
	default:
		return now.Unix()
	}
}
func (f *Father) getTimeUint() string {
	switch f.constTimeType {
	case ConstTypeUnixMilli:
		return "ms"
	default:
		return "s"
	}
}
func (f *Father) SetConstTimeType(constTimeType ConstTimeType) {
	f.constTimeType = constTimeType
}

func (f *Father) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startAt := f.getNow()
	c := Context{
		Req:           req,
		Writer:        w,
		StatusCode:    http.StatusOK,
		index:         -1,
		ConstTimeUint: f.getTimeUint(),
	}
	defer func() {
		c.ConstTime = f.getNow() - startAt
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
			c.handlers = router.Handlers
			c.Next()
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
	if err = f.checkRouters(); err != nil {
		f.logger.Fatalf("服务启动失败,err=%v", err)
	}
	return http.ListenAndServe(f.Address, f)
}

/*
暂时使用 Get Post Put Delete Patch 方法
*/

func (f *Father) NewGroup(path string) *Group {
	g := Group{}
	gg := g.New(path)
	f.groups = append(f.groups, gg)
	return gg
}
func (f *Father) checkRouters() (err error) {
	routerMap := map[string]bool{}
	return f.initRouters(httpSeparator, &routerMap, f.groups, &f.Routers)
}

func (f *Father) initRouters(path string, routerMap *map[string]bool, groups []*Group, routers *[]Router) (err error) {
	for i := 0; i < len(groups); i++ {
		g := groups[i]
		if len(g.Children) == 0 {
			routerPath := g.GetPath(path)
			oldLength := len(*routerMap)
			r := *routerMap
			r[g.Method+routerPath] = true
			if oldLength == len(*routerMap) {
				f.logger.Fatalf("重复路由------>>>>>>[%v]%v", g.Method, routerPath)
				return RepeatedRouterError
			}
			f.logger.Printf("路由------>>>>>>[%v]%v\n", g.Method, routerPath)
			*routers = append(*routers, Router{
				Method:   g.Method,
				Path:     g.GetPath(path),
				Handlers: g.Handlers,
			})
			return
		}
		if err = f.initRouters(path+g.Path, routerMap, g.Children, routers); err != nil {
			return
		}
	}
	return
}
