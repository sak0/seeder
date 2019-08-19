package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	c.JSON(http.StatusOK, resp)
}