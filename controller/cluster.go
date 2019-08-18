package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Node struct {
	Addr 	string		`json:"addr"`
	Role 	string		`json:"role"`
	Status 	string		`json:"status"`
}

type ClusterInfo struct {
	Nodes 	[]*Node		`json:"nodes"`
}

// @Summary 获取seed集群信息
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":S200,"data":{},"msg":"ok"}"
// @Router /api/v1/cluster [get]
func GetCluster(c *gin.Context) {
	cluster := &ClusterInfo{
		Nodes : []*Node{
			&Node{
				Addr:"192.168.0.2",
				Role:"master",
				Status:"active",
			},
		},
	}

	resp := Response{}
	resp.Message = "get cluster info success."
	resp.Data = cluster
	resp.Code = "S200"

	c.JSON(http.StatusOK, resp)
}