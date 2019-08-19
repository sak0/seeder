package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/golang/glog"
)

// @Summary 获取单个镜像仓库的tag列表
// @Accept  json
// @Produce json
// @Param page query int false "Page"
// @Param pageSize query int false "PageSize"
// @Param status query bool false "VerifyStatus"
// @Param cached query bool false "Cached"
// @Success 200 {object} models.RepositoryTag
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/repository/{repo}/tags [get]
func GetRepositoryTags(c *gin.Context) {
	repo := c.Param("id")
	glog.V(5).Infof("ctr: get tags for repo %v", repo)

	resp := Response{}
	c.JSON(http.StatusOK, resp)
}

// @Summary 删除本地指定镜像仓库的指定tag
// @Accept  json
// @Produce json
// @Success 202 {object} models.RepositoryTag
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/repository/{repo}/{tag} [delete]
func DeleteRepositoryTag(c *gin.Context) {
	resp := Response{}
	c.JSON(http.StatusOK, resp)
}