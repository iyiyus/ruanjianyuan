package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type PageResult struct {
	Records interface{} `json:"records"`
	Current int         `json:"current"`
	Size    int         `json:"size"`
	Total   int64       `json:"total"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{Code: 200, Msg: "success", Data: data})
}

func OKMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Result{Code: 200, Msg: msg, Data: data})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Result{Code: 400, Msg: msg, Data: nil})
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Result{Code: 401, Msg: msg, Data: nil})
}

func Page(c *gin.Context, records interface{}, total int64, current, size int) {
	c.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "success",
		Data: PageResult{
			Records: records,
			Current: current,
			Size:    size,
			Total:   total,
		},
	})
}
