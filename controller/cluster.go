package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/sak0/seeder/models"
	"strconv"
)

type ClusterInfo struct {
	ClusterName 	string						`json:"cluster_name"`
	Nodes 			[]*models.SeederNode		`json:"nodes"`
}

// @Summary 获取edge-cloud整体集群信息
// @Accept  json
// @Produce json
// @Param page query int false "Page"
// @Param pageSize query int false "PageSize"
// @Success 200 {object} models.SeederNode
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/cluster [get]
func GetCluster(c *gin.Context) {
	resp := Response{}

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	nodes, count, err := models.GetSeederNodes(page, pageSize)
	if err != nil {
		resp.Message = "get cluster info failed."
		resp.Data = err
		resp.Code = "S400"
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Message = "get cluster info success."
	resp.Data = PageList{
		Offset:page,
		Size:pageSize,
		Total:count,
		DataList:nodes,
	}
	resp.Code = "S200"

	c.JSON(http.StatusOK, resp)
}