package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-gin-prometheus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "github.com/sak0/seeder/docs"
	_ "github.com/go-sql-driver/mysql"

	"github.com/sak0/seeder/pkg/utils"
	"github.com/sak0/seeder/controller"
	"github.com/sak0/seeder/models"
	"github.com/sak0/seeder/pkg/repoer"
	"github.com/sak0/seeder/pkg/cluster"
	"github.com/sak0/seeder/pkg/keeper"
)

const (
	WhoIAm 		= "seeder"
	PortIUse 	= 15000
	healthURL   = "health"
	baseURL		= "/api/v1/"
)

var (
	myName 			string
	dbAddr 			string
	dbName			string
	dbUser			string
	dbPassword		string
	initDb			bool
	role 			string
	master 			string
	repoAddr		string
	advAddr 		string

	useNat 			bool
)

func init() {
	flag.StringVar(&advAddr, "advertise-addr", "10.12.103.89", "addr for advertise.")
	flag.StringVar(&repoAddr, "repo-addr", "http://172.16.24.103", "addr for repo.")
	flag.StringVar(&myName, "node-name", "edge-node-pc", "seeder node name.")
	flag.StringVar(&dbAddr, "db-addr", "172.16.24.103:3306", "database connection url.")
	flag.StringVar(&dbName, "db-name", "seeder", "database name to use.")
	flag.StringVar(&dbUser, "db-user", "root", "database login name.")
	flag.StringVar(&dbPassword, "db-password", "password", "database login password.")
	flag.StringVar(&role, "role", "follower", "seeder role.")
	flag.StringVar(&master,"master-addr", "", "master addr")
	flag.BoolVar(&initDb, "init-db", true, "if need init database.")
	flag.BoolVar(&useNat, "use-nat", false, "if use nat access.")
	flag.Parse()

	utils.SetNodeName(myName)
}

// @title Seeder API
// @version 0.1
// @description Server for image/chart repo consistent.
// @termsOfService

// @contact.name HaoZhi.Cui
// @contact.url http://github.com/sak0
// @contact.email 61755280@qq.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 172.16.24.200:15000
// @BasePath /
func main() {
	myIp, err := utils.GetMyIpAddr()
	if err != nil {
		panic(err)
	}
	if !useNat {
		advAddr = myIp
	}

	if err := models.InitDB(dbAddr, dbName, dbUser, dbPassword, initDb); err != nil {
		glog.Fatalf("init db failed: %v", err)
		return
	}

	if err := utils.ServiceRegister(WhoIAm, PortIUse, baseURL + healthURL); err != nil {
		glog.V(2).Infof("service register failed: %v", err)
	}

	done := make(chan interface{})
	defer close(done)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(1 * time.Minute):
				utils.DoResourceMonitor()
			}
		}
	}()

	// ***run state machine***
	// repoWatcher: watch local repo/charts
	// clusterSyncer: aggregate all nodes info
	// localKeeper: sync cluster info to local database
	{
		repoWatcher, err := repoer.NewRepoWatcher(myName, role, repoAddr, advAddr + ":" + strconv.Itoa(PortIUse), myIp, done)
		if err != nil {
			glog.Fatalf("watch repo %s failed: %v", repoAddr, err)
			return
		}
		go repoWatcher.Run()

		clusterSync := cluster.NewClusterSyncer(role, master, myName, repoAddr, "gossip", done)
		go clusterSync.Run()
		clusterSync.RegisterReporter(repoWatcher)

		localKeeper := keeper.NewLocalKeeper(role, master, myName, done)
		go localKeeper.Run()
		localKeeper.RegisterReporter(clusterSync)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	r.Use(gin.Recovery())

	r.Use(controller.RequestIdMiddleware())

	url := ginSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", myIp, PortIUse)) // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	v1 := r.Group(baseURL)
	{
		v1.GET(healthURL, controller.HealthCheck)
		v1.GET("cluster", controller.GetCluster)

		repository := v1.Group("/repository")
		{
			repository.GET("", controller.GetRepository)
			repository.GET(":id/tags", controller.GetRepositoryTags)
			repository.POST(":id/:tag/download", controller.UpdateRepositoryTag)
			repository.DELETE(":id/:tag", controller.DeleteRepositoryTag)
		}

		chart := v1.Group("/chart")
		{
			chart.GET("", controller.GetChartRepo)
			chart.GET(":id/versions", controller.GetChartVersion)
			chart.GET(":id/:version/param", controller.GetChartVersionParam)
			chart.POST(":id/:version/download", controller.DownloadChartVersion)
			chart.DELETE(":id/:version", controller.DeleteChartVersion)
			chart.POST(":id/:version/push", controller.PushChartVersion)
		}
	}

	glog.Fatal(r.Run("0.0.0.0:" + strconv.Itoa(PortIUse)))
}