package main

import (
	"github.com/yundun-php/ydsdk_go/ydsdk"
	"fmt"
)

var (
	app_id  string = "xxxxxxxxxxxxxxxxxxxx"
	app_secret string = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	base_api_url string = "http://apiv4.yundun.com/V4/"
)


func main() {
	fmt.Println("Hello", "world")
	//构造
	//get
	get()
}

func get()  {
	var client = ydsdk.New(app_id, app_secret,base_api_url)
	args := map[string]interface{}{

	}
	resp, err := client.Get("test",args)

	fmt.Println("结果")
	fmt.Println(err)
	fmt.Println(resp.StatusCode())
}

