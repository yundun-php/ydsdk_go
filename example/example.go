package main

import (
	"github.com/yundun-php/ydsdk_go/ydsdk"
	"fmt"
)

var (
	app_id  string = "raKtReMGOfQwAQDpFqco"
	app_secret string = "e08230881e6fd6706deba5159a3913d7"
)


func main() {
	fmt.Println("Hello", "world")
	//构造
	//get
	get()
}

func get()  {
	var client = ydsdk.New(app_id, app_secret)
	args := map[string]interface{}{

	}
	resp, err := client.Get("test",args)

	fmt.Println("结果")
	fmt.Println(err)
	fmt.Println(string(resp))
}

