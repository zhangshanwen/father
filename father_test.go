package father

import (
	"testing"
)

func TestFather(t *testing.T) {
	f := NewFather()
	f.SetConstTimeType(ConstTypeUnixMilli)
	hello := f.NewGroup("hello")
	hello.Middleware(func(c *Context) {
		f.logger.Println("我是一个中间件哦")
		a := struct {
			Name string `json:"name"`
		}{}
		a.Name = "我是中间件"
		_ = c.Json(&a)
	})
	hello.Post("", func(c *Context) {
		a := struct {
			Name string `json:"name"`
		}{}
		_ = c.BindJson(&a)
		_ = c.Json(&a)
	})
	nimei := f.NewGroup("nimei")
	nimei.Get("", func(c *Context) {
		js := map[string]interface{}{"code": 1, "msg": "nimei"}
		_ = c.Json(js)
	})
	t.Log(f.Run("127.0.0.1", 10086))
}
