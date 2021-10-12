package father

import (
	"testing"
)

func TestNewFather(t *testing.T) {
	f := NewFather()
	f.Post("/hello", func(c *Context) {
		a := struct {
			Name string `json:"name"`
		}{}
		_ = c.BindJson(&a)
		_ = c.Json(&a)
	})
	f.Get("/nimei", func(c *Context) {
		js := map[string]interface{}{"code": 1, "msg": "nimei"}
		_ = c.Json(js)
	})
	t.Log(f.Run("127.0.0.1", 10086))
}
