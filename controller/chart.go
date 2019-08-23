package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"github.com/sak0/seeder/models"
	"github.com/golang/glog"
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

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	charts, count, err := models.GetAllCharts(page, pageSize)
	if err != nil {
		resp.Message = "get chart failed."
		resp.Data = err
		resp.Code = "S400"
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Message = "get repos success."
	resp.Data = PageList{
		Total:count,
		DataList:charts,
	}
	resp.Code = "S200"

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
	chartName := c.Param("id")

	if chartName == "" {
		resp.Message = "must have chartRepo name"
		resp.Code = "S400"
		c.JSON(http.StatusOK, resp)
	}
	glog.V(5).Infof("ctr: get versions for chart %v", chartName)

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	versions, count, err := models.GetVersionByChart(page, pageSize, chartName)
	if err != nil {
		resp.Message = "get versions failed."
		resp.Data = err
		resp.Code = "S400"
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Message = "get versions success."
	resp.Data = PageList{
		Total:count,
		DataList:versions,
	}
	resp.Code = "S200"

	c.JSON(http.StatusOK, resp)
}

// @Summary 下载更新指定Chart仓库的指定版本到本地仓库
// @Accept  json
// @Produce json
// @Param version body models.ChartVersion true "Download the version to local"
// @Success 202 {object} models.ChartVersion
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/chart/{repo}/{version}/download [post]
func DownloadChartVersion(c *gin.Context) {
	resp := Response{}
	c.JSON(http.StatusOK, resp)
}

// @Summary 删除本地指定Chart仓库的指定版本
// @Accept  json
// @Produce json
// @Success 202 {object} models.ChartVersion
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/chart/{repo}/{version} [delete]
func DeleteChartVersion(c *gin.Context) {
	resp := Response{}
	c.JSON(http.StatusOK, resp)
}