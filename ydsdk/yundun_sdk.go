package ydsdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)


// YdSdk 用来构造请求，设置各项参数，和执行请求的接口
type YdSdkClient interface {
	NewUrl(api string) *YdSdk
	NewRequest(params interface{}) ([]byte, error)

	SetAPPID(appid string) *YdSdk
	SetAPPSecret(appSecret string) *YdSdk
	SetLogger(logger *log.Logger) *YdSdk
	SetMethod(method string) *YdSdk
	SetParams(params interface{}) *YdSdk
	Get(url string,params interface{}) *YdSdk
	Post(url string,params interface{}) *YdSdk
	Put(url string,params interface{}) *YdSdk
	Delete(url string,params interface{}) *YdSdk
	New(appid string, appSecret string) *YdSdk
}

// YdSdk 是请求的结构
// 一次请求具体功能由 YdSdkClient 接口实现
type YdSdk struct {
	URL     string
	ReqTime int64
	Options Options
	Logger  *log.Logger
}

const (
	SDKName = "yundun-go-sdk"
	BASE_API_URL string = "http://apiv4.yundun.cn/V4/"
)

type  Options struct {
	METHOD string
	APP_ID string
	APP_SECRET string
	// 用户id, 仅代理需要
	USER_ID string
	//客户端ip
	CLIENT_IP string
	//客户端userAgent
	CLIENT_USER_AGENT string
	//api base url
	PARAMS map[string]interface{}

	HOST string
	LOG string
	HTTP struct {
		Timeout time.Duration
	}
	Debug bool
}

// NewOptions 返回一个新的 *Options
func NewOptions(appid string, appSecret string) *Options {
	opt := &Options{
		APP_ID:  appid,
		APP_SECRET: appSecret,
		//RandomLen: 6,
		CLIENT_USER_AGENT: "yundun_go_client",

		Debug: false,

		HTTP: struct {
			Timeout time.Duration
		}{Timeout: 10 * time.Second},
	}

	return opt
}

func New(appId string,appSecret string)*YdSdk   {
	opt := NewOptions(appId, appSecret)
	c := &YdSdk{}
	c.Options = *opt

	//c.NewRandom(c.Options.RandomLen)
	c.ReqTime = time.Now().Unix()

	c.Logger = log.New(os.Stderr, "["+SDKName+"]", log.LstdFlags)
	return c
}

// NewClient 生成一个新的 client 实例
func NewClient(o *Options) *YdSdk {
	c := &YdSdk{}
	c.Options = *o

	//c.NewRandom(c.Options.RandomLen)
	c.ReqTime = time.Now().Unix()

	c.Logger = log.New(os.Stderr, "["+SDKName+"]", log.LstdFlags)
	return c
}

// SetAPPID 为实例设置 APPID
func (c *YdSdk) SetAPPID(appid string) *YdSdk {
	c.Options.APP_ID = appid
	return c
}

// SetAPPSecret为实例设置 APPSecret
func (c *YdSdk) SetAPPSecret(appSecret string) *YdSdk {
	c.Options.APP_SECRET = appSecret
	return c
}

// SetLogger 为实例设置 logger
func (c *YdSdk) SetLogger(logger *log.Logger) *YdSdk {
	c.Logger = logger
	return c
}

// SetMETHOD 为实例设置 Method
func (c *YdSdk) SetMethod(method string) *YdSdk {
	c.Options.METHOD = method
	return c
}
// SetMETHOD 为实例设置 Params
func (c *YdSdk) SetParams(params map[string]interface{}) *YdSdk {
	c.Options.PARAMS = params
	return c
}

// SetDebug 为实例设置调试模式
func (c *YdSdk) SetDebug(debug bool) *YdSdk {
	if debug {
		c.Options.Debug = debug
	}

	return c
}

// NewURL 为实例设置 URL
func (c *YdSdk) NewURL(api string) *YdSdk {
	c.URL = BASE_API_URL + api
	return c
}
type  payload struct {
	Body interface{} `json:"body"`
	Algorithm    string   `json:"algorithm"`
	Issued_at   int64    `json:"issued_at"`
	UserId   int64    `json:"user_id"`
	ClientIp string  `json:"client_ip"`
	ClientUserAgent string  `json:"client_userAgent"`
}
// NewRequest 执行实例发送请求
func (c *YdSdk) NewRequest() ([]byte, error) {
	reqest, err := http.NewRequest(c.Options.METHOD, c.URL,nil)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	//给一个key设定为响应的value.
	//reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value") //必须设定该参数,POST参数才能正常提交
	reqest.Header.Set("Content-Type", "application/json")
	reqest.Header.Set("User-Agent", c.Options.CLIENT_USER_AGENT)
	sign := SignedRequest(c.Options.METHOD,c.Options.PARAMS,c.Options.APP_SECRET)
	reqest.Header.Set("X-Auth-Sign",sign)
	reqest.Header.Set("X-Auth-App-Id", c.Options.APP_ID)
	//客户端,被Get,Head以及Post使用
	client := &http.Client{
		Timeout:c.Options.HTTP.Timeout,
	}
	resp, err := client.Do(reqest)//发送请求
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求异常")
	}
	defer resp.Body.Close()//一定要关闭resp.Body
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	if c.Options.Debug {
		c.Logger.Printf("Request Url : %s, Request Params : %s, Request Res : %s\n", c.URL,c.Options.PARAMS, string(content))
	}

	return content, err
}

func (c *YdSdk)Get(url string,params map[string]interface{}) ([]byte, error) {
	c.SetMethod("GET")
	c = c.NewURL(url)
	c = c.SetParams(params)
	rs,e := c.NewRequest()
	return  rs ,e
}
func (c *YdSdk) Post(url string,params map[string]interface{}) ([]byte, error) {
	c.SetMethod("POST")
	c = c.NewURL(url)
	c = c.SetParams(params)
	rs,e := c.NewRequest()
	return  rs ,e
}
func (c *YdSdk) Put(url string,params map[string]interface{}) ([]byte, error) {
	c.SetMethod("PUT")
	c = c.NewURL(url)
	c = c.SetParams(params)
	rs,e := c.NewRequest()
	return  rs ,e
}

func (c *YdSdk) Delete(url string,params map[string]interface{}) ([]byte, error) {
	c.SetMethod("DELETE")
	c = c.NewURL(url)
	c = c.SetParams(params)
	rs,e := c.NewRequest()
	return  rs ,e
}


func SignedRequest(method string,params map[string]interface{} ,app_secret string) (sign string) {
	Payload := payload{
		Algorithm: "HMAC-SHA256",
		Issued_at:time.Now().Unix(),
	}
	if method!="GET"{
		json_str,_ := json.Marshal(params)
		Payload.Body = string(json_str)
	}else{
		Payload.Body = params
	}
	jsons, errs := json.Marshal(Payload) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	encodeString := base64.StdEncoding.EncodeToString(jsons)
	encodedPayload := strings.Replace(encodeString, "+/", "-_", -1)
	hashedSig  := hmacSha256(encodedPayload, app_secret)
	encodedSig := strings.Replace(hashedSig, "+/", "-_", -1)
	sign = encodedSig+ "." + encodedPayload
	return  sign
}

func hmacSha256(encodedData string, appSecret string)(hashedSig string) {
	key:=[]byte(appSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(encodedData))
	hashedSig = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return  hashedSig
}
