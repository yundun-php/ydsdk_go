
# YUNDUN API SDK Go语言实现

[![GoDoc](https://godoc.org/github.com/yundun-php/ydsdk_go/ydsdk?status.svg)](https://godoc.org/github.com/yundun-php/ydsdk_go/ydsdk)
[![license](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/yundun-php/ydsdk_go/master/LICENSE)

## 说明
> +	接口基地址(base_api_url)： 'https://apiv4.yundun.com/V4/';
> +	接口遵循RESTful,默认请求体json,接口默认返回json
> +	app_id, app_secret，base_api_url 联系技术客服，先注册一个云盾的账号，用于申请绑定api身份

> * 签名
>    * 每次请求都签名，保证传输过程数据不被篡改
>    * 客户端：sha256签名算法，将参数base64编码+app_secret用sha256签名，每次请求带上签名
>   * 服务端：拿到参数用相同的算法签名，对比签名是否正确

## 安装

```
go get github.com/yundun-php/ydsdk_go
```

## 用法

```Go
import "github.com/yundun-php/ydsdk_go/ydsdk"

var client = ydsdk.New(app_id, app_secret,base_api_url)
	args := map[string]interface{}{

	}
resp, err := client.Get("test",args)

```

更多示例可在 [Example](https://github.com/yundun-php/ydsdk_go/blob/master/example/example.go) 或 [godoc](https://godoc.org/github.com/yundun-php/ydsdk_go#pkg-examples) 查看

## Documentation

[完整文档](https://godoc.org/github.com/yundun-php/ydsdk_go)

## License

This project is under the MIT License.

