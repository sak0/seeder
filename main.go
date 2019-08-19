package main

import (
	"flag"
	"strconv"

	"github.com/golang/glog"
	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-gin-prometheus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "./docs"

	"github.com/sak0/seeder/pkg/utils"
	"github.com/sak0/seeder/controller"
	"fmt"
)

const (
	WhoIAm 		= "seeder"
	PortIUse 	= 15000
	healthURL   = "health"
)

// @title Swagger Example API
// @version 1.0
// @description Server for image/chart repo consistent.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host seeder.cloudminds.com
// @BasePath /v1
func main() {
	flag.Parse()

	if err := utils.ServiceRegister(WhoIAm, PortIUse, healthURL); err != nil {
		glog.V(2).Infof("service register failed: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	r.Use(gin.Recovery())

	url := ginSwagger.URL(fmt.Sprintf("http://127.0.0.1:%d/swagger/doc.json", PortIUse)) // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	v1 := r.Group("/api/v1")
	{
		v1.GET(healthURL, controller.HealthCheck)
		v1.GET("cluster", controller.GetCluster)

		//accounts := v1.Group("/image")
		//{
		//	accounts.GET(":id", c.ShowAccount)
		//	accounts.GET("", c.ListAccounts)
		//	accounts.POST("", c.AddAccount)
		//	accounts.DELETE(":id", c.DeleteAccount)
		//	accounts.PATCH(":id", c.UpdateAccount)
		//	accounts.POST(":id/images", c.UploadAccountImage)
		//}
	}

	glog.Fatal(r.Run("0.0.0.0:" + strconv.Itoa(PortIUse)))
}