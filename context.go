package father

import (
	"encoding/json"
	"net/http"
)

// Context TODO 扩展sync.pool，以提高性能
type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
}

func (c *Context) isMethodRouterParams() bool {
	return c.Req.Method == http.MethodGet || c.Req.Method == http.MethodDelete
}

// Bind 解析请求参数
func (c *Context) Bind(v interface{}) (err error) {
	// TODO 解析参数
	// 1.分析请求方法，若为GET,DELETE 等方法则为路由传参
	if c.isMethodRouterParams() {
		// TODO 解析路由传参参数
		return
	}
	// 2.分析request里面的content-type

	return
}

// Json 返回json 参数
func (c *Context) Json(v interface{}) (err error) {
	var data []byte
	data, err = json.Marshal(v)
	if err != nil {
		return
	}
	c.Writer.Header().Set(contentType, jsonType)
	_, err = c.Writer.Write(data)
	return
}
