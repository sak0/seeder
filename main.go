package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-gin-prometheus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "./docs"
	_ "github.com/go-sql-driver/mysql"

	"github.com/sak0/seeder/pkg/utils"
	"github.com/sak0/seeder/controller"
	"github.com/sak0/seeder/models"
)

const (
	WhoIAm 		= "seeder"
	PortIUse 	= 15000
	healthURL   = "health"
)

var (
	dbAddr 			string
	dbName			string
	dbUser			string
	dbPassword		string
	initDb			bool
)

func init() {
	flag.StringVar(&dbAddr, "db-addr", "172.16.24.103:3306", "database connection url.")
	flag.StringVar(&dbName, "db-name", "seeder", "database name to use.")
	flag.StringVar(&dbUser, "db-user", "root", "database login name.")
	flag.StringVar(&dbPassword, "db-password", "password", "database login password.")
	flag.BoolVar(&initDb, "init-db", true, "if need init database.")
	flag.Parse()
}

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
	if err := models.InitDB(dbAddr, dbName, dbUser, dbPassword, initDb); err != nil {
		glog.Fatalf("init db failed: %v", err)
		return
	}

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

		repository := v1.Group("/repository")
		{
			repository.GET("", controller.GetRepository)
			repository.GET(":id/tags", controller.GetRepositoryTags)
		}

		chart := v1.Group("/chart")
		{
			chart.GET("", controller.GetChartRepo)
			chart.GET(":id/charts", controller.GetChartVersion)
		}
	}

	glog.Fatal(r.Run("0.0.0.0:" + strconv.Itoa(PortIUse)))
}