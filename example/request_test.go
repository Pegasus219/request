package example

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Pegasus219/request"
	"io/ioutil"
	"testing"
)

type TestResp struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

//demo:get请求（同步，获取返回信息）
func TestGetRequest(t *testing.T) {
	url := "http://127.0.0.1/get.test"
	params := map[string]interface{}{
		"first_param":  "kaka卡罗卡",
		"second_param": 2341,
	}
	data := &TestResp{}
	errRes, err := request.NewHttpRequest(url).Params(params).Method("get").Unmarshal(data)
	fmt.Println("err:", err)
	fmt.Println("err response:", errRes)
	fmt.Println("rtn first:", data.First)
	fmt.Println("rtn second:", data.Second)
}

//demo:post请求（异步，仅发送信息）
func TestPostRequest(t *testing.T) {
	url := "http://127.0.0.1/post.test"
	params := map[string]interface{}{
		"first_param":  "kaka卡罗卡",
		"second_param": 2341,
	}
	request.InitAsyncPool(10)
	try, err := request.NewHttpRequest(url).Async(true).Params(params).Success()
	fmt.Println("err:", err)
	fmt.Println("try result", try)
}

//demo:本地文件上传
func TestFileUpload(t *testing.T) {
	url := "http://127.0.0.1/upload.test"
	params := map[string]interface{}{
		"first_param":  "kaka卡罗卡",
		"second_param": []int{1024, 2088},
	}
	files := map[string]interface{}{}
	files["newFile"] = []string{
		"all_async_search.png",
		"all_async_search_view.png",
	}
	try, err := request.NewHttpRequest(url).Params(params).Files(files).Success()
	fmt.Println("err:", err)
	fmt.Println("try result", try)
}

//demo:直接请求，返回[]byte数据
func TestRtnByte(t *testing.T) {
	url := "https://www.baidu.com"
	rspByte, err := request.NewHttpRequest(url).Method("get").Response()
	fmt.Println("err:", err)
	fmt.Println("result:", string(rspByte))
}

//demo:直接在body中传入二进制数据
func TestByteBody(t *testing.T) {
	url := "http://127.0.0.1/byte.body.test"
	filename := "test01.jpg"
	fileBuf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	encodeString := base64.StdEncoding.EncodeToString(fileBuf)
	p := map[string]interface{}{
		"image": encodeString,
	}
	b, _ := json.Marshal(p)
	h := map[string]string{
		"Content-Type": "application/json",
	}
	rspByte, err := request.NewHttpRequest(url).Headers(h).Body(b).Response()
	fmt.Println("err:", err)
	fmt.Println("result:", string(rspByte))
}
