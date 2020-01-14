package ydsdk

import (
	"bytes"
	"log"
	"testing"
)

func TestNewOptions(t *testing.T)  {
	t.Log(NewOptions("a","b"))
}
func TestNew(t *testing.T)  {
	t.Log(New("a","b"))
}

func TestSetAPPID(t *testing.T)  {
	var c = New("a","b")
	t.Log(c.SetAPPID("c"))
}

func TestSetAPPSecret(t *testing.T)  {
	var c = New("a","b")
	t.Log(c.SetAPPSecret("d"))
}
func TestSetLogger(t *testing.T)  {
	var c = New("a","b")
	buffer := bytes.NewBuffer(make([]byte, 0, 64))
	logger := log.New(buffer, "prefix: ", 0)
	t.Log(c.SetLogger(logger))
}
func TestSetMethod(t *testing.T)  {
	var c = New("a","b")
	t.Log(c.SetMethod("get"))
}

func TestSetParams(t *testing.T)  {
	args := map[string]interface{}{

	}
	var c = New("a","b")
	t.Log(c.SetParams(args))
}
func TestSetDebug(t *testing.T)  {
	var c = New("a","b")
	t.Log(c.SetDebug(false))
	t.Log(c.SetDebug(true))
}

func TestNewURL(t *testing.T)  {
	var c = New("a","b")
	t.Log(c.NewURL("aaa"))
}

func TestNewRequest(t *testing.T)  {
	var c = New("a","b")
	c.SetMethod("GET")
	c = c.NewURL("aaa")
	c = c.SetDebug(false)
	args := map[string]interface{}{

	}
	c = c.SetParams(args)
	c = c.SetMethod("aaa")
	t.Log(c.NewRequest())
	t.Log(c.NewRequest())
	c = c.SetDebug(true)
	t.Log(c.NewRequest())
}
func TestGet(t *testing.T)  {
	var c = New("a","b")
	c = c.SetDebug(false)
	args := map[string]interface{}{

	}
	t.Log(c.Get("aaa",args))
	c = c.SetDebug(true)
	t.Log(c.Get("aaa",args))
}
func TestPost(t *testing.T)  {
	var c = New("a","b")
	args := map[string]interface{}{

	}
	t.Log(c.Post("aaa",args))
}

func TestPut(t *testing.T)  {
	var c = New("a","b")
	args := map[string]interface{}{

	}
	t.Log(c.Put("aaa",args))
}
func TestDelete(t *testing.T)  {
	var c = New("a","b")
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

