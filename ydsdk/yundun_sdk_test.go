package ydsdk

import(
	"log"
	"bytes"
	"testing"
)

func TestNewOptions(t *testing.T)  {
	t.Log(NewOptions("a","b","c"))
}
func TestNew(t *testing.T)  {
	t.Log(New("a","b","http://apiv4.yundun.com/V4/"))
}

func TestSetAPPID(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	t.Log(c.SetAPPID("c"))
}

func TestSetAPPSecret(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	t.Log(c.SetAPPSecret("d"))
}
func TestSetLogger(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	buffer := bytes.NewBuffer(make([]byte, 0, 64))
	logger := log.New(buffer, "prefix: ", 0)
	t.Log(c.SetLogger(logger))
}
func TestSetMethod(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	t.Log(c.SetMethod("get"))
}

func TestSetParams(t *testing.T)  {
	args := map[string]interface{}{

	}
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	t.Log(c.SetParams(args))
}
func TestSetDebug(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	t.Log(c.SetDebug(false))
	t.Log(c.SetDebug(true))
}

func TestNewRequest(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	c.SetMethod("GET")
	c = c.SetDebug(false)
	args := map[string]interface{}{

	}
	c = c.SetParams(args)
	c = c.SetMethod("aaa")
	t.Log(c.NewRequest("GET","AA",args))
	t.Log(c.NewRequest("GET","testaaa",args))
	t.Log(c.NewRequest("AA","AA",args))
	c = c.SetDebug(true)
	t.Log(c.NewRequest("AA","AA",args))
	t.Log(c.NewRequest("","AA",args))
}
func TestGet(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	c = c.SetDebug(false)
	args := map[string]interface{}{

	}
	t.Log(c.Get("aaa",args))
	c = c.SetDebug(true)
	t.Log(c.Get("aaa",args))
}
func TestPost(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	args := map[string]interface{}{

	}
	t.Log(c.Post("aaa",args))
}

func TestPut(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	args := map[string]interface{}{

	}
	t.Log(c.Put("aaa",args))
}
func TestDelete(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	args := map[string]interface{}{

	}
	t.Log(c.Delete("aaa",args))
}

func TestSignedRequest(t *testing.T)  {
	args := map[string]interface{}{
		"a": map[string]interface{}{
			"a":"a^4",
		},
	}
	t.Log(SignedRequest("GET",args,"aaa"))
	t.Log(SignedRequest("POST",args,"aaa"))
}

func TestSetHeaders(t *testing.T) {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	jsonHeader := map[string]string{
		"Content-Type":"application/json",
	}
	t.Log(c.SetHeaders(jsonHeader))
}

func TestResponse_Response(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	r,_:=c.NewRequest("GET","test",nil)
	t.Log(r.Response())
}
func TestResponse_Json(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	r,_:=c.NewRequest("GET","test",nil)
	jsondata :=make(map[string]interface{})
	t.Log(r.Json(jsondata))
}
func TestResponse_Body(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	r,_:=c.NewRequest("GET","test",nil)
	t.Log(r.Body())
}
func TestResponse_Content(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	r,_:=c.NewRequest("GET","test",nil)
	t.Log(r.Content())
}
func TestResponse_StatusCode(t *testing.T)  {
	var c = New("a","b","http://apiv4.yundun.com/V4/")
	r,_:=c.NewRequest("GET","test",nil)
	t.Log(r.StatusCode())
}
