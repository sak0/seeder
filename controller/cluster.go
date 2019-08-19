package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Node struct {
	NodeName 			string   	`json:"node_name"`
	AdvertiseAddr		string 		`json:"advertise_addr"`
	BindAddr 			string		`json:"bind_addr"`
	Role 				string		`json:"role"`
	Status 				string		`json:"status"`
}

type ClusterInfo struct {
	Nodes 	[]*Node		`json:"nodes"`
}

// @Summary 获取cloud-edge整体集群信息
// @Accept  json
// @Produce json
// @Success 200 {object} controller.ClusterInfo
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/cluster [get]
func GetCluster(c *gin.Context) {
	cluster := &ClusterInfo{
		Nodes : []*Node{
			{
				NodeName:"center-node",
				AdvertiseAddr:"10.23.100.2:15300",
				BindAddr:"192.168.0.2:8080",
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