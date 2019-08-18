package main

import (
		"flag"

	"github.com/golang/glog"
	"github.com/gin-gonic/gin"
		"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	//_ "github.com/swaggo/gin-swagger/example/basic/docs"
	_ "./docs"

	"github.com/sak0/seeder/pkg/utils"
	"github.com/mcuadros/go-gin-prometheus"
	"github.com/sak0/seeder/controller"
	"strconv"
)

const (
	WhoIAm 		= "seeder"
	PortIUse 	= 15000
	healthURL   = "health"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
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
	r.GET(healthURL, controller.HealthCheck)
	r.GET("api/v1/cluster", controller.GetCluster)

	url := ginSwagger.URL("http://127.0.0.1:15000/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	glog.Fatal(r.Run("0.0.0.0:" + strconv.Itoa(PortIUse)))
}