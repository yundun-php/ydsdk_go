package main

import (
	"github.com/yundun-php/ydsdk_go/ydsdk"
	"fmt"
)

func main() {
	fmt.Println("test start")
	sdk_test()
	fmt.Println("test end")
}

func sdk_test() {
	app_id := "xxxxxxxxxxxxxxxxxxxx"
	app_secret := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	api_pre := "http://apiv4.local.com/V4/"

	sdk := ydsdk.YdSdk{
		AppId: app_id,
		AppSecret: app_secret,
		ApiPre: api_pre,
		Timeout: 30,
	}
	var api string
	var err error
	//ReqParams 有三个属性，如果用不到，不设置即可：
	//ReqParams.Query 对应的是GET请求的参数，map[string]interface{}
	//ReqParams.Data  对应的是非GET请求的参数，map[string]interface{}
	//ReqParams.Header  对应的是发起请求的header头，map[string]string
	var reqParams ydsdk.ReqParams

	// Response包括响应的全部信息，其中：
	// Response.HttpCode Http请求响应状态码，成功是200
	// Response.RespBody 请求返回的body字符串
	// Response.BizCode 是业务的状态码，1代码请求成功，非1代码请求失败
	// Response.BizMsg  是业务的状态码对应的信息
	// Response.BizData  返回的业务数据，只有BizCode为1时，才会有数据
	var resp *ydsdk.Response

	api = "test.sdk.get"
	reqParams = ydsdk.ReqParams{
		Query: map[string]interface{}{
			"page": 1,
			"pagesize": 10,
			"data": map[string]interface{}{
				"name": "name名称",
				"domain": "baidu.com",
			},
		},
	}
	resp, err = sdk.Get(api, reqParams)
	if err == nil {
		if resp.BizCode == 1 {
			fmt.Println(api, " 业务处理成功")
		} else {
			fmt.Println(api, " 业务处理错误: ")
		}
		fmt.Println(api, " http_code: ", resp.HttpCode)
		fmt.Println(api, " body: ", resp.RespBody)
		fmt.Println(api, " biz_code: ", resp.BizCode)
		fmt.Println(api, " biz_msg: ", resp.BizMsg)
		fmt.Println(api, " biz_data: ", resp.BizData)
		fmt.Println(api, " err: ", err)
	} else {
		fmt.Println(api, " request error: ", err)
	}
	fmt.Println("")


	api = "test.sdk.post"
	reqParams = ydsdk.ReqParams{
		Data: map[string]interface{}{
			"name": 1,
			"age": 10,
			"data": map[string]interface{}{
				"name": "name名称",
				"domain": "baidu.com",
			},
		},
	}
	resp, err = sdk.Post(api, reqParams)
	if err == nil {
		if resp.BizCode == 1 {
			fmt.Println(api, " 业务处理成功")
		} else {
			fmt.Println(api, " 业务处理错误: ")
		}
		fmt.Println(api, " http_code: ", resp.HttpCode)
		fmt.Println(api, " body: ", resp.RespBody)
		fmt.Println(api, " biz_code: ", resp.BizCode)
		fmt.Println(api, " biz_msg: ", resp.BizMsg)
		fmt.Println(api, " biz_data: ", resp.BizData)
		fmt.Println(api, " err: ", err)
	} else {
		fmt.Println(api, " request error: ", err)
	}
	fmt.Println("")


	api = "test.sdk.delete"
	reqParams = ydsdk.ReqParams{
		Data: map[string]interface{}{
			"id": 10,
		},
	}
	resp, err = sdk.Delete(api, reqParams)
	if err == nil {
		if resp.BizCode == 1 {
			fmt.Println(api, " 业务处理成功")
		} else {
			fmt.Println(api, " 业务处理错误: ")
		}
		fmt.Println(api, " http_code: ", resp.HttpCode)
		fmt.Println(api, " body: ", resp.RespBody)
		fmt.Println(api, " biz_code: ", resp.BizCode)
		fmt.Println(api, " biz_msg: ", resp.BizMsg)
		fmt.Println(api, " biz_data: ", resp.BizData)
		fmt.Println(api, " err: ", err)
	} else {
		fmt.Println(api, " request error: ", err)
	}
	fmt.Println("")


	api = "test.sdk.put"
	reqParams = ydsdk.ReqParams{
		Data: map[string]interface{}{
			"name": 1,
			"age": 10,
			"data": map[string]interface{}{
				"name": "name名称",
				"domain": "baidu.com",
			},
		},
	}
	resp, err = sdk.Put(api, reqParams)
	if err == nil {
		if resp.BizCode == 1 {
			fmt.Println(api, " 业务处理成功")
		} else {
			fmt.Println(api, " 业务处理错误: ")
		}
		fmt.Println(api, " http_code: ", resp.HttpCode)
		fmt.Println(api, " body: ", resp.RespBody)
		fmt.Println(api, " biz_code: ", resp.BizCode)
		fmt.Println(api, " biz_msg: ", resp.BizMsg)
		fmt.Println(api, " biz_data: ", resp.BizData)
		fmt.Println(api, " err: ", err)
	} else {
		fmt.Println(api, " request error: ", err)
	}
	fmt.Println("")
}
