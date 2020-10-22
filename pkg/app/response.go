package app

import (
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

// ErrorResp 错误返回值
func ErrorResp(c *gin.Context, code int, errMsg string) {
	resp(c, http.StatusOK, code, errMsg, nil)
}

// ErrorResp 错误返回值
func UnauthorizedResp(c *gin.Context, code int, errMsg string) {
	resp(c, http.StatusUnauthorized, code, errMsg, nil)
}

// SuccessResp 正确返回值
func SuccessRespByCode(c *gin.Context, code int, data interface{}) {
	resp(c, http.StatusOK, code, "", data)
}

// SuccessResp 正确返回值
func SuccessResp(c *gin.Context, data interface{}) {
	resp(c, http.StatusOK, e.SUCCESS, "", data)
}

// SuccessResp 正确返回值
func SuccessPureResp(c *gin.Context, data interface{}) {
	pureResp(c, http.StatusOK, e.SUCCESS, "", data)
}

// resp 返回
func resp(c *gin.Context, httpCode, code int, errMsg string, data interface{}) {
	resp := Response{
		Code:   code,
		Msg:    e.GetMsg(code),
		ErrMsg: errMsg,
		Data:   data,
	}
	c.Set(util.LogResponse, &resp)
	c.JSON(httpCode, resp)
}

// resp 返回
func pureResp(c *gin.Context, httpCode, code int, errMsg string, data interface{}) {
	resp := Response{
		Code:   code,
		Msg:    e.GetMsg(code),
		ErrMsg: errMsg,
		Data:   data,
	}
	c.Set(util.LogResponse, &resp)
	c.PureJSON(httpCode, resp)
}
