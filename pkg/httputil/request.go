package httputil

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"

	"eim/pkg/log"
)

var client *resty.Client

func init() {
	client = resty.New()
	client.SetAllowGetMethodPayload(true)
	client.SetLogger(log.Default())
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (its *HTTPError) Error() string {
	return its.Message
}

func NewClient() *resty.Client {
	return resty.New()
}

func SetLogger(l *log.Logger) {
	client.SetLogger(l)
}

// DoRequest 发送请求
// @param T 返回结果类型
// @param ctx 上下文
// @param reqUrl 请求地址
// @param method 请求方法	GET/POST/PUT/DELETE/PATCH
// @param header 请求头
// @param data 请求数据 nil表示不传数据，可以是url.Values、map[string]string、struct等，url.Values表示表单数据需要设置Content-Type: application/x-www-form-urlencoded
// @param debug 是否开启debug模式
// @return result 返回结果
// @return err 错误信息	nil表示成功，err.(*HTTPError)表示状态非200～300，但是不代表请求失败，比如401等错误，可用于特定场景判断
func DoRequest[T any](ctx context.Context, reqUrl string, method string, header http.Header, data interface{}, debug bool) (result T, err error) {
	hMap := map[string]string{}
	for k, v := range header {
		if v != nil {
			hMap[k] = v[0]
		}
	}
	req := client.R().EnableTrace().SetContext(ctx).SetHeaders(hMap)

	if data != nil {
		if formData, isOk := data.(url.Values); isOk {
			req.SetFormDataFromValues(formData)
		} else {
			req.SetBody(data)
		}
	}

	resp, err := req.SetDebug(debug).Execute(method, reqUrl)
	if err != nil {
		return result, err
	}
	body := resp.Body()
	if resp.IsSuccess() {
		if len(body) > 0 {
			var ri interface{} = result
			switch ri.(type) {
			case string:
				ri = string(body)
				return ri.(T), nil
			default:
				err = json.Unmarshal(body, &result)
				if err != nil {
					err = errors.New(fmt.Sprintf("%q", body))
				}
			}
		}
		return result, err
	}
	return result, &HTTPError{
		Code:    resp.StatusCode(),
		Message: string(body),
	}
}

func IsWebSocketRequest(r *http.Request) bool {
	contains := func(key, val string) bool {
		vv := strings.Split(r.Header.Get(key), ",")
		for _, v := range vv {
			if val == strings.ToLower(strings.TrimSpace(v)) {
				return true
			}
		}
		return false
	}
	if !contains("Connection", "upgrade") {
		return false
	}
	if !contains("Upgrade", "websocket") {
		return false
	}
	return true
}
