package request

import "time"

type (
	//httpRequest对象结构体
	HttpRequest struct {
		url     string
		method  string //默认是POST方法
		async   bool
		body    []byte
		params  map[string]interface{}
		headers map[string]string
		files   map[string]interface{}
		timeout time.Duration //仅同步请求有效, 默认10秒钟
	}

	//http请求错误信息
	HttpResponseError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	//http请求返回的通用结构体
	HttpResponse struct {
		Success bool               `json:"success"`
		Data    interface{}        `json:"data"`
		Error   *HttpResponseError `json:"error"`
	}
)
