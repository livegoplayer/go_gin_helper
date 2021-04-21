package helper

//处理全局panic的返回值，重写gin.Recover中间件的内容
import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	mylogger "github.com/livegoplayer/go_logger/logger"
)

const (
	COMMON_STATUS = 0
	COMMON_ERROR  = 1
	AUTH_ERROR    = 2
	OTHER_ERROR   = 3
)

// 错误处理的结构体
type Error struct {
	//只是说他们两个差不多，没用到
	Resp       `json:"-"`
	StatusCode int         `json:"-"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

var (
	ServerError = NewError(http.StatusInternalServerError, COMMON_ERROR, "系统异常，请稍后重试!")
	NotFound    = NewError(http.StatusNotFound, COMMON_ERROR, http.StatusText(http.StatusNotFound))
)

func OtherError(message string) *Error {
	return NewError(http.StatusOK, COMMON_ERROR, message)
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
		Data:       &EmptyData{},
	}
}

func NewErrorWithData(statusCode, Code int, data interface{}, msg string) *Error {
	if data == nil {
		data = &EmptyData{}
	}

	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
		Data:       data,
	}
}

// 404处理
func HandleNotFound(c *gin.Context) {
	err := NotFound
	c.JSON(err.StatusCode, err)
	return
}

// 服务异常处理
func HandleServerError(c *gin.Context) {
	err := ServerError
	c.JSON(err.StatusCode, err)
	return
}

func ErrHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var Err *Error
				var realErr error
				//如果是通过本文件定义的Error，如果是调试模式，则输出所有的错误内容，否则，只输出自定义内容
				if e, ok := err.(*Error); ok {
					Err = e
				} else if e, ok := err.(error); ok {
					Err = ServerError
					realErr = OtherError(e.Error())
				} else {
					//普通panic,直接传string的那种
					Err = ServerError
					realErr = OtherError(err.(string))
				}

				//测试环境非服务器返回的错误打印下来
				if realErr != nil {
					//这种程度的error, 输出
					mylogger.Error(realErr.Error())
					if gin.IsDebugging() {
						//这里打印错误
						PrintStack(Err)
					}
				}

				c.JSON(Err.StatusCode, Err)

				c.Abort()
			}
		}()
		c.Next()
	}
}

// 打印异常栈的方法
func PrintStack(err error) {
	msg := err.Error()
	_, _ = os.Stderr.Write([]byte("error: " + msg + "\n"))
	debug.PrintStack()
}

func TimeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}

//
////// 跨域

func Cors(allowedHostList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		//获取配置文件中的host
		accessControlAllowOrigin := "*"
		//跨域允许的域名
		for _, host := range allowedHostList {

			if origin == host {
				accessControlAllowOrigin = origin
			}
		}

		AccessControlAllowCredentials := true
		if accessControlAllowOrigin == "*" {
			AccessControlAllowCredentials = false
		}

		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", accessControlAllowOrigin)                  // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			//允许设置cookie
			c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(AccessControlAllowCredentials)) //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                       // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

//可以用一个类似handler的东西作为验证逻辑
func AuthenticationMiddleware(CheckTokenHandler func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		CheckTokenHandler(c)

		c.Next()
	}
}
