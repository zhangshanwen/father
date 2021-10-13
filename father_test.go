package father

import (
	"testing"
)

func TestFather(t *testing.T) {
	f := NewFather()
	hello := f.NewGroup("hello")
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
