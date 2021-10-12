package father

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestNewFather(t *testing.T) {
	f := NewFather()
	f.Get("/hello", func(w http.ResponseWriter) {
		w.Header().Set("content-type", "application/json")
		js, _ := json.Marshal(map[string]interface{}{"code": 1, "msg": "hello"})
		w.Write(js)
	})
	f.Get("/nimei", func(w http.ResponseWriter) {
		js, _ := json.Marshal(map[string]interface{}{"code": 1, "msg": "nimei"})
		w.Write(js)
	})
	t.Log(f.Run("127.0.0.1", 10086))
}
