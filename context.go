package father

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	StatusCode int
	ConstTime  int64
	index      int
	handlers   []HandlerFunc
}

// BindJson 解析请求json参数
func (c *Context) BindJson(obj interface{}) (err error) {
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(obj)
}

func (c *Context) SetStatusCode(code int) {
	c.StatusCode = code
}

// Json 返回json 参数
func (c *Context) Json(obj interface{}) (err error) {
	var data []byte
	data, err = json.Marshal(obj)
	if err != nil {
		return
	}
	c.Writer.Header().Set(contentType, jsonType)
	_, err = c.Writer.Write(data)
	return
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}

}
