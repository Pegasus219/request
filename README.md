### 一个封装好的http请求服务

##### 支持多类型参数，以及异步请求发送

#### 1、可支持的请求参数类型
| method | params type |
| ----- | ---- |
| POST | string、[]string、int、[]int、float64、[]float64 或 二进制body |
| GET | string、int、float64 |

#### 2、提供执行http请求的方法
    1、Response() ([]byte, error)
    # 最基础的方法，直接取到请求返回的原信息
    # 例如：rspByte, err := request.NewHttpRequest(url).Params(params).Response()
    
    2、Success() (bool, error)
    # 发送请求，只关心是否执行成功
    # 例如：try, err := request.NewHttpRequest(url).Params(params).Success()
    
    3、Unmarshal(data interface{}) (*HttpResponseError, error)
    # 获取接口返回的data数据
    # 例如：errRes, err := request.NewHttpRequest(url).Params(params).Unmarshal(&data)
    
#### 3、关于请求接口的返回数据格式
    请注意，Success和Unmarshal使用时均对被调用接口返回的数据格式有一定要求，否则建议使用Response()方法。
    要求结构如下：
    //http请求返回的通用结构体
	HttpResponse struct {
		Success bool               `json:"success"`
		Data    interface{}        `json:"data"`
		Error   struct {
		    Code    int    `json:"code"`
		    Message string `json:"message"`
	    } `json:"error"`
	}

#### 4、异步请求   
    # 执行异步请求前，需先开启异步请求队列
    # request.InitAsyncPool(100)
    # try, err := request.NewHttpRequest(url).Async(true).Params(params).Success()
