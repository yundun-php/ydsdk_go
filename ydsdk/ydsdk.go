package ydsdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
	"strconv"
)

var SDK_VERSION = "1.0.3"

// YdSdk 是请求的结构
type YdSdk struct {
	AppId string
	AppSecret string
	ApiPre string
	UserId int
	clientIp string
	userAgent string
	Timeout int
	Debug bool
	isSetDefault bool
}

type Response struct {
	Url string
	Api string
	Method string
	Query string
	Data string
	ReqHeaders map[string]string
	Response *http.Response
	RespBody string
	HttpCode  int
	RespData map[string]interface{}
	BizCode int
	BizMsg string
	BizData interface{}
}

type ReqParams struct {
	Query map[string]interface{}
	Data map[string]interface{}
	Headers map[string]string
}

// 将深度字典全部字符串化
func MapStrval(data map[string]interface{}) (map[string]interface{}, error) {
	query := map[string]interface{}{}

	for key, row := range data {
		dtype := reflect.TypeOf(row).String()
		switch dtype {
		case "map[string]interface {}":
			subQuery, err := MapStrval(row.(map[string]interface{}))
			if err != nil {
				return query, err
			}
			query[key] = subQuery
		case "[]interface {}":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]interface{}) {
				subDtype := reflect.TypeOf(subRow).String()
				switch subDtype {
				case "map[string]interface {}":
					subQuery1, err := MapStrval(subRow.(map[string]interface{}))
					if err != nil {
						return query, err
					}
					subQuery[strconv.Itoa(subKey)] = subQuery1
				case "int":
					subQuery[strconv.Itoa(subKey)] = strconv.Itoa(row.(int))
				case "string":
					subQuery[strconv.Itoa(subKey)] = row.(string)
				default:
				}
				
			}
			query[key] = subQuery
		case "[]float32":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]float32) {
				subQuery[strconv.Itoa(subKey)] = strconv.FormatFloat(float64(subRow), 'f', 2, 64)
			}
			query[key] = subQuery
		case "[]float64":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]float64) {
				subQuery[strconv.Itoa(subKey)] = strconv.FormatFloat(subRow, 'f', 2, 64)
			}
			query[key] = subQuery
		case "[]bool":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]bool) {
				if subRow {
					subQuery[strconv.Itoa(subKey)] = "1"
				} else {
					subQuery[strconv.Itoa(subKey)] = "0"
				}
			}
			query[key] = subQuery
		case "int":
			query[key] = strconv.Itoa(row.(int))
		case "[]int":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]int) {
				subQuery[strconv.Itoa(subKey)] = strconv.Itoa(subRow)
			}
			query[key] = subQuery
		case "uint":
			query[key] = strconv.Itoa(int(row.(uint)))
		case "[]uint":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]uint) {
				subQuery[strconv.Itoa(subKey)] = strconv.Itoa(int(subRow))
			}
			query[key] = subQuery
		case "int8":
			query[key] = strconv.Itoa(int(row.(int8)))
		case "[]int8":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]int8) {
				subQuery[strconv.Itoa(subKey)] = strconv.Itoa(int(subRow))
			}
			query[key] = subQuery
		case "uint8":
			query[key] = strconv.Itoa(int(row.(uint8)))
		case "[]uint8": 	// []byte
			query[key] = string(row.([]byte))
		case "int16":
			query[key] = strconv.Itoa(int(row.(int16)))
		case "[]int16":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]int16) {
				subQuery[strconv.Itoa(subKey)] = strconv.Itoa(int(subRow))
			}
			query[key] = subQuery
		case "uint16":
			query[key] = strconv.Itoa(int(row.(uint16)))
		case "[]uint16":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]uint16) {
				subQuery[strconv.Itoa(subKey)] = strconv.Itoa(int(subRow))
			}
			query[key] = subQuery
		case "int32":
			query[key] = strconv.Itoa(int(row.(int32)))
		case "[]int32":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]int32) {
				subQuery[strconv.Itoa(subKey)] = strconv.Itoa(int(subRow))
			}
			query[key] = subQuery
		case "uint32":
			query[key] = strconv.Itoa(int(row.(uint32)))
		case "[]uint32":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]uint32) {
				subQuery[strconv.Itoa(subKey)] = strconv.Itoa(int(subRow))
			}
			query[key] = subQuery
		case "int64":
			query[key] = strconv.Itoa(int(row.(int64)))
		case "[]int64":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]int64) {
				subQuery[strconv.Itoa(subKey)] = strconv.Itoa(int(subRow))
			}
			query[key] = subQuery
		case "uint64":
			query[key] = strconv.Itoa(int(row.(uint64)))
		case "[]uint64":
			subQuery := map[string]interface{}{}
			for subKey, subRow := range row.([]int64) {
				subQuery[strconv.Itoa(subKey)] = strconv.Itoa(int(subRow))
			}
			query[key] = subQuery
		case "string":
			query[key] = row.(string)
		default:
			return query, fmt.Errorf("不支持的数据类型：%s, %s", dtype, key)
		}
	}

	return query, nil
}

// 将深度字典转换为字符串一维键值结构，用于转换为URL      MapStrval已经做了字符串化，这里不需要处理太多类型
func Map2StrKV(data map[string]interface{}) (map[string]string, error) {
	query := map[string]string{}

	for key, row := range data {
		dtype := reflect.TypeOf(row).String()
		switch dtype {
		case "map[string]interface {}":
			subQuery, err := Map2StrKV(row.(map[string]interface{}))
			if err != nil {
				return query, err
			}
			for subKey, subRow := range subQuery {
				query[key + "[" + subKey + "]"] = subRow
			}
		case "[]interface {}":
			for subKey, subRow := range row.([]interface{}) {
				subDtype := reflect.TypeOf(subRow).String()
				switch subDtype {
				case "map[string]interface {}":
					subQuery, err := Map2StrKV(subRow.(map[string]interface{}))
					if err != nil {
						return query, err
					}
					for subKey1, subRow1 := range subQuery {
						query[key + "[" + strconv.Itoa(subKey) + "]["+subKey1+"]"] = subRow1
					}
				case "int":
					query[key + "[" + strconv.Itoa(subKey) + "]"] = strconv.Itoa(row.(int))
				case "string":
					query[key + "[" + strconv.Itoa(subKey) + "]"] = row.(string)
				default:
				}
				
			}
		case "int":
			query[key] = strconv.Itoa(row.(int))
		case "string":
			query[key] = row.(string)
		default:
			return query, fmt.Errorf("不支持的数据类型：%s, %s", dtype, key)
		}
	}

	return query, nil
}

func UrlEncode(data map[string]interface{}) (string, error) {
	params, err := Map2StrKV(data)
	if err != nil {
		return "", err
	}
	querys := []string{}
	for k, v := range params {
		querys = append(querys, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	return strings.Join(querys, "&"), nil
}

//hmac 加密
func hmacSha256(encodedData string, appSecret string) (hashedSig string) {
	key:=[]byte(appSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(encodedData))
	hashedSig = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return  hashedSig
}

//处理json字符串中的非ascii字符
func jsonToUnicode(raw []byte) []byte {
	//fmt.Println("byte 0", []byte(toStr))
	toStr := strconv.QuoteToASCII(string(raw))	//全字符串非ascii转ascii
	//fmt.Println("byte 1", []byte(toStr))
	//fmt.Printf("0 %s\n", toStr)
	toStr = strings.Trim(toStr, `"`)		//转换后，前后会多加上"，需要去掉
	//fmt.Println("byte 2", []byte(toStr))
	//fmt.Printf("1 %s\n", toStr)
	toStr = strings.Replace(toStr, `\"`, `"`, -1)	//转换后，json字符串中的"，会转为\"，需要再次赚回来
	//fmt.Println("byte 3", []byte(toStr))
	//fmt.Printf("2 %s\n", toStr)
	return []byte(toStr)
}

//生成body 里的sign
func (sdk *YdSdk) Sign(method string, signData map[string]interface{}) (sign string) {
	bf := &bytes.Buffer{}
	jsonEncoder := json.NewEncoder(bf)
	//不转义html
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(signData)
	b := bf.Bytes()
	// json编码后，go会自动追加\n，去掉 https://github.com/golang/go/issues/7767
	b = bytes.TrimSpace(b)
	b = jsonToUnicode(b)
	encodeString := base64.StdEncoding.EncodeToString(b)
	tmpencodeString := strings.ReplaceAll(encodeString, "+", "-")
	encodedPayload := strings.ReplaceAll(tmpencodeString, "/", "_")
	//fmt.Println("encodedPayload: ", encodedPayload)
	hashedSig  := hmacSha256(encodedPayload, sdk.AppSecret)
	tmphashedSig := strings.ReplaceAll(hashedSig, "+", "-")
	encodedSig := strings.ReplaceAll(tmphashedSig, "/", "_")
	sign = encodedSig
	return  sign
}

func (sdk *YdSdk) payload(method string, reqParams *ReqParams) {
	issuedAt := int(time.Now().Unix())
	if method == "GET" {
		reqParams.Query["user_id"] = strconv.Itoa(sdk.UserId)
		reqParams.Query["client_ip"] = ""    //当项目内代理转发调用时，此参数用作将外部的IP传递给内部的系统，这里默认空
		reqParams.Query["client_userAgent"] = sdk.userAgent
		reqParams.Query["algorithm"] = "HMAC-SHA256"
		reqParams.Query["issued_at"] = issuedAt
		reqParams.Query, _ = MapStrval(reqParams.Query)
	} else {
		reqParams.Data["user_id"]          = strconv.Itoa(sdk.UserId)
		reqParams.Data["client_ip"]        = ""    //当项目内代理转发调用时，此参数用作将外部的IP传递给内部的系统，这里默认空
		reqParams.Data["client_userAgent"] = sdk.userAgent
		reqParams.Data["algorithm"]        = "HMAC-SHA256"
		reqParams.Data["issued_at"]        = issuedAt
	}

	signData := map[string]interface{}{}
	if method == "GET" {
		signData = reqParams.Query
	} else {
		signData = reqParams.Data
	}
	signStr := sdk.Sign(method, signData)

        reqParams.Headers["X-Auth-Sign"]        = signStr
        reqParams.Headers["X-Auth-App-Id"]      = sdk.AppId
        reqParams.Headers["X-Auth-Sdk-Version"] = SDK_VERSION
        reqParams.Headers["Content-Type"]       = "application/json; charset=utf-8"
        reqParams.Headers["User-Agent"]         = sdk.userAgent
	return
}

func (sdk *YdSdk) initDefault() bool {
	sdk.clientIp = ""
	sdk.userAgent = "YdSdk " + SDK_VERSION + "; " + runtime.Version() + "; arch/" + runtime.GOARCH + "; os/" + runtime.GOOS
	return true
}

// Request 执行实例发送请求
func (sdk *YdSdk) Request(uri, method string, reqParams ReqParams) (*Response, error) {
	//初始化数据默认值
	sdk.initDefault()
	response := Response{
		Api: uri,
	}
	if reqParams.Query == nil {
		reqParams.Query = map[string]interface{}{}
	}
	if reqParams.Data == nil {
		reqParams.Data = map[string]interface{}{}
	}
	if reqParams.Headers == nil {
		reqParams.Headers = map[string]string{}
	}

	method = strings.ToUpper(method)
	response.Method = method
	url := strings.TrimRight(sdk.ApiPre, "/") + "/" + strings.TrimLeft(uri, "/")
	response.Url = url

	var err  error
	sdk.payload(method, &reqParams)

	query, err := UrlEncode(reqParams.Query)
	if err != nil {
		return &response, err
	}
	response.Query = query
	if query != "" {
		url = url + "?" + query
	}

	jsonByte, err := json.Marshal(reqParams.Data)
	if err != nil {
		return &response, err
	}
	response.Data = string(jsonByte)
	body := strings.NewReader(string(jsonByte))
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &response, err
	}
	response.ReqHeaders = reqParams.Headers
	for k, v := range reqParams.Headers {
		req.Header.Set(k, v)
	}

	//客户端,被Get,Head以及Post使用
	client := &http.Client{
		Timeout: time.Duration(sdk.Timeout) * time.Second,
	}
	resp, err := client.Do(req)//发送请求
	if err != nil {
		return &response, err
	}
	response.Response = resp
	response.HttpCode = resp.StatusCode

	if resp.StatusCode == 200 {
        	rawByte, err := ioutil.ReadAll(resp.Body)
        	if err == nil {
			response.RespBody = string(rawByte)
			respData := map[string]interface{}{}
			err = json.Unmarshal(rawByte, &respData)
			if err == nil {
				if bizStatus, ok := respData["status"].(map[string]interface{}); ok {
					response.BizCode = int(bizStatus["code"].(float64))
					response.BizMsg = bizStatus["message"].(string)
					response.BizData = respData["data"]
				} else {
					response.BizCode = 0
					response.BizMsg = "the json format of response body has not status"
					response.BizData = map[string]interface{}{}
					err = fmt.Errorf("the json format of response body has not status")
				}
			} else {
				response.BizCode = 0
				response.BizMsg = "json parse response body error: " + err.Error()
				response.BizData = map[string]interface{}{}
				err = fmt.Errorf("json parse response body error: %s", err.Error())
			}
        	} else {
			response.BizCode = 0
			response.BizMsg = "response body read error: " + err.Error()
			response.BizData = map[string]interface{}{}
			err = fmt.Errorf("response body read error: %s", err.Error())
		}
	} else {
		response.BizCode = 0
		response.BizMsg = "response code is " + strconv.Itoa(resp.StatusCode)
		response.BizData = map[string]interface{}{}
		err = fmt.Errorf("response code is %s", strconv.Itoa(resp.StatusCode))
	}

	return &response, err
}

// GET 请求
func (sdk *YdSdk) Get(api string, reqParams ReqParams) (*Response, error) {
	return sdk.Request(api, "GET", reqParams)
}

// POST 请求
func (sdk *YdSdk) Post(api string, reqParams ReqParams) (*Response, error) {
	return sdk.Request(api, "POST", reqParams)
}

//PUT 请求
func (sdk *YdSdk) Put(api string, reqParams ReqParams) (*Response, error) {
	return sdk.Request(api, "PUT", reqParams)
}

//DELETE 请求
func (sdk *YdSdk) Delete(api string, reqParams ReqParams) (*Response, error) {
	return sdk.Request(api, "DELETE", reqParams)
}
