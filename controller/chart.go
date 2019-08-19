package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 获取Chart仓库列表
// @Accept  json
// @Produce json
// @Param page query int false "Page"
// @Param pageSize query int false "PageSize"
// @Param status query bool false "VerifyStatus"
// @Param cached query bool false "Cached"
// @Success 200 {object} models.ChartRepo
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/chart [get]
func GetChartRepo(c *gin.Context) {
	resp := Response{}
	c.JSON(http.StatusOK, resp)
}


// @Summary 获取指定Chart仓库的版本列表
// @Accept  json
// @Produce json
// @Param page query int false "Page"
// @Param pageSize query int false "PageSize"
// @Param status query bool false "VerifyStatus"
// @Param cached query bool false "Cached"
// @Success 200 {object} models.ChartVersion
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/chart/{repo}/charts [get]
func GetChartVersion(c *gin.Context) {
	resp := Response{}
	c.JSON(http.StatusOK, resp)
}