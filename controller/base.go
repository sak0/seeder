package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"strconv"
	"github.com/golang/glog"
	"net/http"
)

const (
	ERRBADREQUEST					= "400"
	ERRINTERNALERR					= "500"

	MYERRCODE 						= "03"

	ERROR_INVALID_PARAMS		 	= "0001"
	ERROR_EXIST_TAG 				= "0002"
	ERROR_NOT_EXIST_TAG 			= "0003"
	ERROR_NOT_EXIST_ARTICLE 		= "0004"
	ERROR_AUTH_CHECK_TOKEN_FAIL 	= "0005"
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT 	= "0006"
	ERROR_AUTH_TOKEN 				= "0007"
	ERROR_AUTH 						= "0008"

	ERROR_UNKNOWN 					= "0000"
)

//var MsgFlags map[int]string = {
//
//}

func RespErr(codeCatg, statusCode string, msg string, c *gin.Context) {
	var resp Response
	resp.Code = codeCatg + MYERRCODE + statusCode
	resp.Message = msg

	c.JSON(http.StatusOK, resp)
}

func mustErrorCode(mergedCode string) int {
	errCode, err := strconv.Atoi(mergedCode)
	if err != nil {
		glog.V(2).Infof("error code %s: %v", mergedCode, err)
		return mustErrorCode(ERRBADREQUEST + MYERRCODE + ERROR_UNKNOWN)
	}

	return errCode
}

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}