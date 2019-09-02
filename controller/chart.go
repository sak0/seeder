package controller

import (
	"net/http"
	"strconv"
	"fmt"
	"io/ioutil"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/gin-gonic/gin"

	"github.com/sak0/seeder/models"
	"github.com/sak0/seeder/pkg/utils"
	"github.com/sak0/seeder/pkg/transfer"
	"github.com/sak0/go-harbor"
)

// @Summary 获取Chart仓库列表
// @Accept  json
// @Produce json
// @Param page query int false "Page"
// @Param pageSize query int false "PageSize"
// @Param status query bool false "VerifyStatus"
// @Param cached query bool false "Cached"
// @Param chart query string false "chart_name"
// @Param cluster query string false "ClusterName"
// @Success 200 {object} models.ChartRepo
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/chart [get]
func GetChartRepo(c *gin.Context) {
	resp := Response{}

	clusterName := c.Query("ClusterName")
	glog.V(3).Infof("get chart for remote cluster: %s", clusterName)

	page, _ := strconv.Atoi(c.Query("Page"))
	pageSize, _ := strconv.Atoi(c.Query("PageSize"))
	chartName := c.Query("chart_name")

	if clusterName == "" {
		charts, count, err := models.GetAllCharts(page, pageSize, chartName)
		if err != nil {
			RespErr(ERRBADREQUEST, ERROR_INVALID_PARAMS, "get chart failed.", c)
			return
		}

		resp.Message = "get repos success."
		resp.Data = PageList{
			Total:count,
			DataList:charts,
		}
		resp.Code = "200"

		c.JSON(http.StatusOK, resp)
	} else {
		node, err := models.GetNodeByName(clusterName)
		if err != nil {
			glog.V(2).Infof("get node %s info failed: %v", clusterName, node)
			RespErr(ERRBADREQUEST, ERROR_INVALID_PARAMS, err.Error(), c)
			return
		}
		glog.V(2).Infof("get chart from remote edge: %s", clusterName)

		client := http.Client{
			Transport:utils.GetHTTPTransport(true),
		}

		var url string
		if pageSize > 0 && page > 0 {
			url = fmt.Sprintf("http://%s/api/v1/chart?page=%d&page_size=%d&chart_name=%s",
				node.AdvertiseAddr, page, pageSize, chartName)
		} else {
			url = fmt.Sprintf("http://%s/api/v1/chart?chart_name=%s", node.AdvertiseAddr, chartName)
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
			return
		}

		var remoteResp Response
		remoteRawResp, err := client.Do(req)
		if err != nil {
			RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
			return
		}
		data, err := ioutil.ReadAll(remoteRawResp.Body)
		if err != nil {
			RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
			return
		}
		remoteRawResp.Body.Close()

		err = json.Unmarshal(data, &remoteResp)
		if err != nil {
			RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
			return
		}

		c.JSON(http.StatusOK, remoteResp)
	}
}


// @Summary 获取指定Chart仓库的版本列表
// @Accept  json
// @Produce json
// @Param page query int false "Page"
// @Param pageSize query int false "PageSize"
// @Param status query bool false "VerifyStatus"
// @Param cached query bool false "Cached"
// @Param cluster query string false "ClusterName"
// @Success 200 {object} models.ChartVersion
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/chart/{chartName}/versions [get]
func GetChartVersion(c *gin.Context) {
	resp := Response{}
	chartName := c.Param("id")

	if chartName == "" {
		RespErr(ERRBADREQUEST, ERROR_INVALID_PARAMS, "must have chartRepo name", c)
		return
	}
	glog.V(5).Infof("ctr: get versions for chart %v", chartName)

	clusterName := c.Query("ClusterName")
	glog.V(3).Infof("get chart version for remote cluster: %s", clusterName)

	page, _ := strconv.Atoi(c.Query("Page"))
	pageSize, _ := strconv.Atoi(c.Query("PageSize"))


	if clusterName == "" {
		versions, count, err := models.GetVersionByChart(page, pageSize, chartName)
		if err != nil {
			RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, "get version failed", c)
			return
		}

		resp.Message = "get versions success."
		resp.Data = PageList{
			Offset:page,
			Size:pageSize,
			Total:count,
			DataList:versions,
		}
		resp.Code = "200"

		c.JSON(http.StatusOK, resp)
	} else {
		node, err := models.GetNodeByName(clusterName)
		if err != nil {
			glog.V(2).Infof("get node %s info failed: %v", clusterName, node)
			RespErr(ERRBADREQUEST, ERROR_INVALID_PARAMS, err.Error(), c)
			return
		}
		glog.V(2).Infof("get version from remote edge: %s", clusterName)

		client := http.Client{
			Transport:utils.GetHTTPTransport(true),
		}

		var url string
		if pageSize > 0 && page > 0 {
			url = fmt.Sprintf("http://%s/api/v1/chart/%s/versions?page=%d&page_size=%d",
				node.AdvertiseAddr, chartName, page, pageSize)
		} else {
			url = fmt.Sprintf("http://%s/api/v1/chart/%s/versions", node.AdvertiseAddr, chartName)
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
			return
		}

		var remoteResp Response
		remoteRawResp, err := client.Do(req)
		if err != nil {
			RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
			return
		}
		data, err := ioutil.ReadAll(remoteRawResp.Body)
		if err != nil {
			RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
			return
		}
		remoteRawResp.Body.Close()

		err = json.Unmarshal(data, &remoteResp)
		if err != nil {
			RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
			return
		}

		c.JSON(http.StatusOK, remoteResp)
	}
}

// @Summary 查询指定Version的文件详情，例如：README
// @Accept  json
// @Produce json
// @Param chart query string false "chart_name"
// @Param version query string false "version"
// @Param file query string false "file_name"
// @Success 202 {object} models.ChartVersion
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/versiondetail/file [get]
func GetChartVersionFiles(c *gin.Context) {
	resp := Response{}
	c.JSON(http.StatusOK, resp)
}

// @Summary 查询指定Version的参数Key-Value详情
// @Accept  json
// @Produce json
// @Param chart query string false "chart_name"
// @Param version query string false "version"
// @Success 202 {object} models.ChartVersion
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/versiondetail/params [get]
func GetChartVersionParam(c *gin.Context) {
	resp := Response{}
	chartName := c.Query("chart_name")
	version := c.Query("version")
	if chartName == "" || version == "" {
		RespErr(ERRBADREQUEST, ERROR_INVALID_PARAMS, "must have chart_name and version.", c)
		return
	}

	nodeInfo, err := models.GetNodeByName(utils.GetMyNodeName())
	if err != nil {
		RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS,
			fmt.Sprintf("can not get registry information for node %s", utils.GetMyNodeName()), c)
		return
	}
	harborCli := harbor.NewClient(nil, nodeInfo.RepoAddr, "admin", "Harbor12345")
	detail, _, errs := harborCli.ChartRepos.GetChartVersionDetail(utils.DefaultProjectName, chartName, version)
	if len(errs) > 0 {
		RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS,
			fmt.Sprintf("can not get version %s detail %v", version, errs[0]), c)
		return
	}

	//remoteNodeName := c.Query("ClusterName")
	resp.Code = "200"
	resp.Message = fmt.Sprintf("get chart version %s/%s params success", chartName, version)
	resp.Data = detail.Values
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


// @Summary 推送指定Chart仓库的指定版本到远端仓库
// @Accept  json
// @Produce json
// @Param remote query string true "remote"
// @Success 202 {object} models.ChartVersion
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/chart/{repo}/{version}/push [post]
func PushChartVersion(c *gin.Context) {
	resp := Response{}
	chartName := c.Param("id")
	version := c.Param("version")
	if chartName == "" || version == "" {
		RespErr(ERRBADREQUEST, ERROR_INVALID_PARAMS, "must have chartName and version.", c)
		return
	}

	remoteNodeName := c.Query("remote")
	if remoteNodeName == "" {
		RespErr(ERRBADREQUEST, ERROR_INVALID_PARAMS, "must have remote node name.", c)
		return
	}

	remoteNode, err := models.GetNodeByName(remoteNodeName)
	if err != nil {
		RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
		return
	}
	localNode, err := models.GetNodeByName(utils.GetMyNodeName())
	if err != nil {
		RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
		return
	}

	glog.V(2).Infof("prepare push %s:%s to remote node %v", chartName, version, remoteNode.RepoAddr)

	trans, err := transfer.NewTransfer(localNode.RepoAddr, remoteNode.RepoAddr)
	if err != nil {
		RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
		return
	}
	if err := trans.Transfer(chartName, version); err != nil {
		RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, err.Error(), c)
		return
	}

	resp.Message = "push completed"
	resp.Code = "200"
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