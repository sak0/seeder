package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/sak0/seeder/models"
	"strconv"
)

// @Summary 获取镜像仓库列表
// @Accept  json
// @Produce json
// @Param page query int false "Page"
// @Param pageSize query int false "PageSize"
// @Param status query bool false "VerifyStatus"
// @Param cached query bool false "Cached"
// @Success 200 {object} models.Repository
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/repository [get]
func GetRepository(c *gin.Context) {
	resp := Response{}

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	repos, count, err := models.GetAllRepos(page, pageSize)
	if err != nil {
		RespErr(ERRINTERNALERR, ERROR_INVALID_PARAMS, "get repo failed.", c)
		return
	}

	resp.Message = "get repos success."
	resp.Data = PageList{
		Offset:page,
		Size:pageSize,
		Total:count,
		DataList:repos,
	}
	resp.Code = "200"

	c.JSON(http.StatusOK, resp)
}

// @Summary 下载更新指定镜像仓库的指定tag到本地仓库
// @Accept  json
// @Produce json
// @Param tag body models.RepositoryTag true "Download the tag to local"
// @Success 202 {object} models.RepositoryTag
// @Failure 500 {string} string "Internal Error"
// @Router /api/v1/repository/{repo}/{tag}/download [post]
func UpdateRepositoryTag(c *gin.Context) {
	resp := Response{}
	c.JSON(http.StatusOK, resp)
}