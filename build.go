package request

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

//构建post请求
func buildPostRequest(queryUrl string, params map[string]interface{}) (*http.Request, error) {
	postParams := url.Values{}
	for k, v := range params {
		switch v.(type) {
		case string:
			postParams.Set(k, v.(string))
		case int:
			postParams.Set(k, strconv.Itoa(v.(int)))
		case float64:
			postParams.Set(k, strconv.FormatFloat(v.(float64), 'g', -1, 64))
		case []string:
			for _, sv := range v.([]string) {
				postParams.Add(k, sv)
			}
		case []int:
			for _, sv := range v.([]int) {
				postParams.Add(k, strconv.Itoa(sv))
			}
		case []float64:
			for _, sv := range v.([]float64) {
				postParams.Add(k, strconv.FormatFloat(sv, 'g', -1, 64))
			}
		default:
			return nil, errors.New("bad params struct")
		}
	}
	body := ioutil.NopCloser(strings.NewReader(postParams.Encode())) //把form数据编下码
	request, err := http.NewRequest(http.MethodPost, queryUrl, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request, nil
}

//构建get请求
func buildGetRequest(queryUrl string, params map[string]interface{}) (*http.Request, error) {
	paramsArr := make([]string, 0)
	for k, v := range params {
		switch v.(type) {
		case string:
			paramsArr = append(paramsArr, k+"="+v.(string))
		case int:
			paramsArr = append(paramsArr, k+"="+strconv.Itoa(v.(int)))
		case float64:
			paramsArr = append(paramsArr, k+"="+strconv.FormatFloat(v.(float64), 'g', -1, 64))
		default:
			return nil, errors.New("bad params struct")
		}
	}
	paramsStr := strings.Join(paramsArr, "&")
	queryUrl += "?" + paramsStr
	request, err := http.NewRequest(http.MethodGet, queryUrl, nil)
	return request, err
}

//直接用body构建请求
func buildBodyRequest(queryUrl string, body []byte) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodPost, queryUrl, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	return request, nil
}

//构建文件上传请求
func buildFormDataRequest(queryUrl string, files, params map[string]interface{}) (*http.Request, error) {
	//创建一个模拟的form中的一个选项,这个form项现在是空的
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//设置post传参
	for k, v := range params {
		switch v.(type) {
		case string:
			_ = bodyWriter.WriteField(k, v.(string))
		case []string:
			for _, sv := range v.([]string) {
				_ = bodyWriter.WriteField(k, sv)
			}
		case int:
			_ = bodyWriter.WriteField(k, strconv.Itoa(v.(int)))
		case []int:
			for _, sv := range v.([]int) {
				_ = bodyWriter.WriteField(k, strconv.Itoa(sv))
			}
		case float64:
			_ = bodyWriter.WriteField(k, strconv.FormatFloat(v.(float64), 'g', -1, 64))
		case []float64:
			for _, sv := range v.([]float64) {
				_ = bodyWriter.WriteField(k, strconv.FormatFloat(sv, 'g', -1, 64))
			}
		default:
			return nil, errors.New("bad params struct")
		}
	}
	//设置传参文件
	for k, v := range files {
		switch v.(type) {
		case string:
			filename := v.(string)
			err := bodyWriteWithFile(bodyBuf, bodyWriter, k, filename)
			if err != nil {
				return nil, err
			}
		case []string:
			for _, filename := range v.([]string) {
				err := bodyWriteWithFile(bodyBuf, bodyWriter, k, filename)
				if err != nil {
					return nil, err
				}
			}
		case []byte:
			err := bodyWriteWithByte(bodyBuf, bodyWriter, k, v.([]byte))
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("bad file params")
		}
	}

	//获取请求头部的content-type, multipart/form-data; boundary=...
	contentType := bodyWriter.FormDataContentType()
	//这个很关键,必须在此处关闭,不能使用defer关闭,不然会导致错误
	bodyWriter.Close()

	requestReader := io.MultiReader(bodyBuf)
	request, err := http.NewRequest(http.MethodPost, queryUrl, requestReader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return request, nil
}

//按文件地址上传
func bodyWriteWithFile(bodyBuf *bytes.Buffer, bodyWriter *multipart.Writer, field, filename string) error {
	_, err := bodyWriter.CreateFormFile(field, filepath.Base(filename))
	if err != nil {
		return err
	}
	fileBuf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	bodyBuf.Write(fileBuf)
	return nil
}

//按文件内容上传
func bodyWriteWithByte(bodyBuf *bytes.Buffer, bodyWriter *multipart.Writer, field string, fileBuf []byte) error {
	_, err := bodyWriter.CreateFormFile(field, field)
	if err != nil {
		return err
	}
	bodyBuf.Write(fileBuf)
	return nil
}
