package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//http异步请求池，默认上限100
var httpAsyncPool chan *HttpRequest

//初始化异步池（只能初始化一次，之后重复调用返回false）
func InitAsyncPool(poolSize int) bool {
	if httpAsyncPool != nil {
		return false
	}
	if poolSize < 1 {
		poolSize = 100
	}
	//创建异步池
	httpAsyncPool = make(chan *HttpRequest, poolSize)
	//监听消费异步请求
	go handleAsyncRequest()
	return true
}

//消费异步请求任务
func handleAsyncRequest() {
	for {
		select {
		case req := <-httpAsyncPool:
			_, err := doRequest(req)
			if err != nil {
				log.Printf("handleAsyncRequest Error:%s, Req:%#v", err.Error(), req)
			}
		}
	}
}

// 创建Http请求
func NewHttpRequest(url string) *HttpRequest {
	req := &HttpRequest{
		method:  http.MethodPost,
		url:     url,
		async:   false,
		timeout: 10 * time.Second,
	}
	return req
}

// 设置异步
func (req *HttpRequest) Async(on bool) *HttpRequest {
	req.async = on
	return req
}

// 设置Http请求方法
func (req *HttpRequest) Method(method string) *HttpRequest {
	req.method = strings.ToUpper(method)
	return req
}

// 设置请求参数
func (req *HttpRequest) Params(params map[string]interface{}) *HttpRequest {
	req.params = params
	return req
}

// 设置Headers
func (req *HttpRequest) Headers(headers map[string]string) *HttpRequest {
	req.headers = headers
	return req
}

// 设置文件上传
func (req *HttpRequest) Files(files map[string]interface{}) *HttpRequest {
	req.files = files
	return req
}

// 直接设置body
func (req *HttpRequest) Body(body []byte) *HttpRequest {
	req.body = body
	return req
}

// 设置超时
func (req *HttpRequest) Timeout(timeout time.Duration) *HttpRequest {
	req.timeout = timeout
	return req
}

//http请求（获取请求执行成功与否）
func (req *HttpRequest) Success() (bool, error) {
	rspData, err := req.Response()
	if err != nil {
		return false, err
	}
	var response HttpResponse
	err = json.Unmarshal(rspData, &response)
	if err != nil {
		return false, err
	}
	return response.Success, nil
}

// HTTP请求（填充data结构体信息）
func (req *HttpRequest) Unmarshal(data interface{}) (*HttpResponseError, error) {
	rspError := &HttpResponseError{}
	rspData, err := req.Response()
	if err != nil {
		return nil, err
	}
	var response HttpResponse
	response.Data = data
	err = json.Unmarshal(rspData, &response)
	if err != nil {
		return nil, err
	}

	if response.Success && response.Data != nil {
		data = response.Data
	} else {
		rspError = response.Error
	}

	return rspError, nil
}

//原始http请求（返回byte）
func (req *HttpRequest) Response() ([]byte, error) {
	//发布异步请求任务
	if req.async {
		//检测异步池是否已初始化
		if httpAsyncPool == nil {
			return nil, errors.New("async request pool not init")
		}
		select {
		case httpAsyncPool <- req:
			break
		case <-time.After(req.timeout):
			return nil, errors.New("async request timeout")
		}
		response := &HttpResponse{
			Success: true,
		}
		return json.Marshal(response)
	}
	//处理同步请求
	return doRequest(req)
}

//执行http请求
func doRequest(req *HttpRequest) ([]byte, error) {
	var request *http.Request
	var err error
	//构建请求，填充请求参数
	if req.files != nil {
		request, err = buildFormDataRequest(req.url, req.files, req.params)
	} else if req.body != nil {
		request, err = buildBodyRequest(req.url, req.body)
	} else {
		switch req.method {
		case http.MethodPost:
			request, err = buildPostRequest(req.url, req.params)
		case http.MethodGet:
			request, err = buildGetRequest(req.url, req.params)
		default:
			return nil, errors.New("bad request method")
		}
	}

	if err != nil {
		return nil, err
	}
	//设置默认User-Agent
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:60.0) Gecko/20100101 Firefox/60.0")
	//设置请求头部
	if len(req.headers) > 0 {
		for k, v := range req.headers {
			request.Header.Set(k, v)
		}
	}
	//发送http请求
	client := &http.Client{}
	client.Timeout = req.timeout
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("status code:%d", resp.StatusCode))
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
