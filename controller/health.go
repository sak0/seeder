package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code 		string 		`json:"code"`
	Message 	string 		`json:"message"`
	Data 		interface{}	`json:"data"`
}

type PageList struct {
	Total 		int 		`json:"total"`
	DataList 	interface{}	`json:"list"`
}

// @Summary 获取服务健康状态
// @Produce  json
// @Success 200 {string} json "{"code":S200,"data":{},"msg":"ok"}"
// @Router /api/v1/health [get]
func HealthCheck(c *gin.Context) {
	resp := Response{
		Code : "S200",
		Message : "health check pass.",
	}

	c.JSON(http.StatusOK, resp)
}