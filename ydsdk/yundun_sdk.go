package ydsdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)


// YdSdk 用来构造请求，设置各项参数，和执行请求的接口
type YdSdkClient interface {
	//NewUrl(api string) *YdSdk
	NewRequest(params interface{}) ([]byte, error)

	//SetAPPID(appid string) *YdSdk
	//SetAPPSecret(appSecret string) *YdSdk
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
	Method string
	headers           map[string]string
	ReqTime int64
	Options Options
	Logger  *log.Logger
}

const (
	SDK_VERSION = "1.0"
	SDK_NAME= "ydsdk-go"
)

type  Options struct {
	METHOD string
	APP_ID string
	APP_SECRET string
	BASE_API_URL string
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
func NewOptions(appid string, appSecret string,baseApiUrl string) *Options {
	opt := &Options{
		APP_ID:  appid,
		APP_SECRET: appSecret,
		BASE_API_URL: baseApiUrl,
		//RandomLen: 6,
		CLIENT_USER_AGENT: "yundun_go_client",

		Debug: false,

		HTTP: struct {
			Timeout time.Duration
		}{Timeout: 10 * time.Second},
	}

	return opt
}
// New 生成一个新的 client 实例
func New(appId string,appSecret string,baseApiUrl string)*YdSdk   {
	opt := NewOptions(appId, appSecret,baseApiUrl)
	c := &YdSdk{}
	c.Options = *opt

	//c.NewRandom(c.Options.RandomLen)
	c.ReqTime = time.Now().Unix()

	c.Logger = log.New(os.Stderr, "["+SDK_NAME+"]", log.LstdFlags)
	return c
}

// SetAPPID 为实例设置 APPID
/**
func (c *YdSdk) SetAPPID(appid string) *YdSdk {
	c.Options.APP_ID = appid
	return c
}
**/
// SetAPPSecret为实例设置 APPSecret
/**
func (c *YdSdk) SetAPPSecret(appSecret string) *YdSdk {
	c.Options.APP_SECRET = appSecret
	return c
}
 **/

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
/**
func (c *YdSdk) NewURL(api string) *YdSdk {
	c.URL = BASE_API_URL + api
	return c
}
 */
type  payload struct {
	Body interface{} `json:"body"`
	Algorithm    string   `json:"algorithm"`
	Issued_at   int64    `json:"issued_at"`
	UserId   int64    `json:"user_id"`
	ClientIp string  `json:"client_ip"`
	ClientUserAgent string  `json:"client_userAgent"`
}
// NewRequest 执行实例发送请求
func (c *YdSdk) NewRequest(method, url string, data  map[string]interface{}) (*Response, error) {
	// Build Response
	response := &Response{}


	if method == "" || url == "" {
		return nil, errors.New("parameter method and url is required")
	}
	 url = c.Options.BASE_API_URL + url

	var (
		err  error
		//req  *http.Request
		body io.Reader
	)

	method = strings.ToUpper(method)

	urlList := strings.Split(url, "?")
	if len(urlList) >= 2 {
		for _, val := range strings.Split(urlList[1], "&") {
			v := strings.Split(val, "=")
			if len(v) == 2 {
				data[v[0]] = v[1]
			}
		}
	}
	data["algorithm"] = "HMAC-SHA256"
	data["issued_at"] = time.Now().Unix()
	data["client_ip"] = get_external()
	data["client_userAgent"] = SDK_NAME+" "+SDK_VERSION+" "+runtime.Version()+" "+runtime.GOOS+" "+runtime.GOARCH
	if err != nil {
		fmt.Println(err)
	}
	c.Method = method
	c.JSON()
	sign := SignedRequest(method,data,c.Options.APP_SECRET)
	fmt.Println(sign)
	if method == "GET"  {
		url, err = buildUrl(url, data)
		if err != nil {
			return nil, err
		}
	}

	body, err = c.buildBody(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method,url ,body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	headers := map[string]string{
		//"Content-Type": "application/json",
		"User-Agent": c.Options.CLIENT_USER_AGENT,
		"X-Auth-App-Id":c.Options.APP_ID,
		"X-Auth-Sign":sign,
	}
	c.SetHeaders(headers)
	c.initHeaders(req)

	//客户端,被Get,Head以及Post使用
	client := &http.Client{
		Timeout:c.Options.HTTP.Timeout,
	}
	resp, err := client.Do(req)//发送请求
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求异常")
		fmt.Println(resp.StatusCode)
	}
	if c.Options.Debug {
		c.Logger.Printf("Request Url : %s, Request Params : %s, Request StatusCode : %d\n", url,body,resp.StatusCode)
	}

	response.resp = resp

	return response, nil
}

func (c *YdSdk) initHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
}
// GET 请求
func (c *YdSdk) Get(url string, data map[string]interface{}) (*Response, error) {
	return c.NewRequest(http.MethodGet, url, data)
}
// POST 请求
func (c *YdSdk) Post(url string, data  map[string]interface{}) (*Response, error) {
	return c.NewRequest(http.MethodPost, url, data)
}
//PUT 请求
func (c *YdSdk) Put(url string, data  map[string]interface{}) (*Response, error) {
	return c.NewRequest(http.MethodPut, url, data)
}
//DELETE 请求
func (c *YdSdk) Delete(url string, data  map[string]interface{}) (*Response, error) {
	return c.NewRequest(http.MethodPut, url, data)
}
//生成body 里的sign
func SignedRequest(method string,params map[string]interface{} ,app_secret string) (sign string) {
	paramsData :=make(map[string]interface{},0)
	if  method=="GET" {
		for k, v := range params {
			vv := ""
			if reflect.TypeOf(v).String() == "string" {
				vv = v.(string)
			} else {
				b, err := json.Marshal(v)
				if err != nil {
					fmt.Println("数据异常",err)
					continue
				}
				vv = string(b)
			}
			paramsData[k]=vv
		}
	}else{
		paramsData = params
	}

	bf := &bytes.Buffer{}
	jsonEncoder := json.NewEncoder(bf)
	//不转义html
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(paramsData)
	b := bf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		// 去掉 go std 给加的 \n
		// 正常的 json 末尾是不会有 \n 的
		// @see https://github.com/golang/go/issues/7767
		b = b[:len(b)-1]
	}
	encodeString := base64.StdEncoding.EncodeToString(b)
	tmpencodeString := strings.ReplaceAll(encodeString, "+", "-")
	encodedPayload := strings.ReplaceAll(tmpencodeString, "/", "_")
	hashedSig  := hmacSha256(encodedPayload, app_secret)
	tmphashedSig := strings.ReplaceAll(hashedSig, "+", "-")
	encodedSig := strings.ReplaceAll(tmphashedSig, "/", "_")
	//sign = encodedSig+ "." + encodedPayload
	sign = encodedSig
	return  sign
}
//hmac 加密
func hmacSha256(encodedData string, appSecret string)(hashedSig string) {
	key:=[]byte(appSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(encodedData))
	hashedSig = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return  hashedSig
}

// Check application/json
func (c *YdSdk) isJson() bool {
	if len(c.headers) > 0 {
		for _, v := range c.headers {
			if strings.Contains(strings.ToLower(v), "application/json") {
				return true
			}
		}
	}
	return false
}
//设置 header 传递格式为json
func (c *YdSdk) JSON() *YdSdk {
	jsonHeader := map[string]string{
		"Content-Type":"application/json",
	}
	c.SetHeaders(jsonHeader)
	return c
}

// Set headers
func (c *YdSdk) SetHeaders(headersMap map[string]string) *YdSdk {
	if len(c.headers)==0{
		c.headers=make(map[string]string)
	}
	if headersMap != nil || len(headersMap) > 0 {
		for k, v := range headersMap {
			c.headers[k] = v
		}
	}
	return c
}

// Build query data
func (c *YdSdk) buildBody(d ...interface{}) (io.Reader, error) {
	// GET and DELETE request dose not send body
	if c.Method == "GET"  {
		return nil, nil
	}

	if len(d) == 0 || d[0] == nil {
		return strings.NewReader(""), nil
	}
	t := reflect.TypeOf(d[0]).String()
	if t != "string" && !strings.Contains(t, "map[string]interface") {
		return strings.NewReader(""), errors.New("incorrect parameter format.")
	}
	if t == "string" {
		return strings.NewReader(d[0].(string)), nil
	}

	if c.isJson() {
		if b, err := json.Marshal(d[0]); err != nil {
			return nil, err
		} else {
			return bytes.NewReader(b), nil
		}
	}

	data := make([]string, 0)
	for k, v := range d[0].(map[string]interface{}) {
		if s, ok := v.(string); ok {
			data = append(data, fmt.Sprintf("%s=%v", k, s))
			continue
		}
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		data = append(data, fmt.Sprintf("%s=%s", k, string(b)))
	}

	return strings.NewReader(strings.Join(data, "&")), nil
}

// Build GET request url
func buildUrl(url string, data ...interface{}) (string, error) {
	query, err := parseQuery(url)
	if err != nil {
		return url, err
	}

	if len(data) > 0 && data[0] != nil {
		t := reflect.TypeOf(data[0]).String()
		switch t {
		case "map[string]interface {}":
			for k, v := range data[0].(map[string]interface{}) {
				vv := ""
				if reflect.TypeOf(v).String() == "string" {
					vv = v.(string)
				} else {
					b, err := json.Marshal(v)
					if err != nil {
						return url, err
					}
					vv = string(b)
				}
				query = append(query, fmt.Sprintf("%s=%s", k, vv))
			}
		case "string":
			param := data[0].(string)
			if param != "" {
				query = append(query, param)
			}
		default:
			return url, errors.New("incorrect parameter format.")
		}

	}

	list := strings.Split(url, "?")

	if len(query) > 0 {
		return fmt.Sprintf("%s?%s", list[0], strings.Join(query, "&")), nil
	}

	return list[0], nil
}

// Parse query for GET request
func parseQuery(url string) ([]string, error) {
	urlList := strings.Split(url, "?")
	if len(urlList) < 2 {
		return make([]string, 0), nil
	}
	query := make([]string, 0)
	for _, val := range strings.Split(urlList[1], "&") {
		v := strings.Split(val, "=")
		if len(v) < 2 {
			return make([]string, 0), errors.New("query parameter error")
		}
		query = append(query, fmt.Sprintf("%s=%s", v[0], v[1]))
	}
	return query, nil
}



type Response struct {
	resp *http.Response
	body []byte
}
//请求返回结构体
func (r *Response) Response() *http.Response {
	return r.resp
}
//返回状态码
func (r *Response) StatusCode() int {
	if r.resp == nil {
		return 0
	}
	return r.resp.StatusCode
}

//返回body
func (r *Response) Body() ([]byte, error) {
	defer r.resp.Body.Close()

	if len(r.body) > 0 {
		return r.body, nil
	}

	if r.resp == nil || r.resp.Body == nil {
		return nil, errors.New("response or body is nil")
	}

	b, err := ioutil.ReadAll(r.resp.Body)
	if err != nil {
		return nil, err
	}
	r.body = b

	return b, nil
}
//返回字符串
func (r *Response) Content() (string, error) {
	b, err := r.Body()
	if err != nil {
		return "", nil
	}
	return string(b), nil
}
//返回如果是json 就decode 并赋值与参数
func (r *Response) Json(v interface{}) error {
	b, err := r.Body()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	return nil
}
func get_external() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return string(content)
}