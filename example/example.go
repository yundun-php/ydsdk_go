package main

import (
	//"github.com/yundun-php/ydsdk_go/ydsdk"
	"git.nodevops.cn/gcode/ydsdk_go/ydsdk"
	"fmt"
)

var (
	app_id  string = "raKtReMGOfQwAQDpFqco"
	app_secret string = "e08230881e6fd6706deba5159a3913d7"
	base_api_url string = "http://apiv4.yundun.cn/V4/"
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

