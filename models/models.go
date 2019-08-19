package models

import (
	"time"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/golang/glog"
)

var Db *gorm.DB

type SeederNode struct {
	gorm.Model
	NodeName 			string   	`json:"node_name" gorm:"type:varchar(50);column:node_name"`
	AdvertiseAddr		string 		`json:"advertise_addr" gorm:"type:varchar(50);column:advertise_addr"`
	BindAddr 			string		`json:"bind_addr" gorm:"type:varchar(50);column:bind_addr"`
	Role 				string		`json:"role" gorm:"type:varchar(50);column:role"`
	Status 				string		`json:"status" gorm:"type:varchar(50);column:status"`
}
func (s SeederNode) TableName() string {
	return "seeder_node"
}

type Repository struct {
	gorm.Model
	OwnerNode 		string		`json:"owner_name" gorm:"type:varchar(50);column:owner_name"`
	Name 			string		`json:"repo_name" gorm:"type:varchar(50);column:repo_name"`
	Description 	string		`json:"description" gorm:"type:varchar(255);column:description"`
	PullCount 		int64		`json:"pull_count" gorm:"type:int;column:pull_count"`
	StarCount 		int64		`json:"star_count" gorm:"type:int;column:star_count"`
	TagsCount 		int64		`json:"tag_count" gorm:"type:int;column:pull_count"`
	VerifyStatus 	string		`json:"verify_status" gorm:"type:varchar(50);column:verify_status"`
	Cached 			bool		`json:"cached" gorm:"type:bool;column:cached"`
}
func (s Repository) TableName() string {
	return "repository"
}

type RepositoryTag struct {
	gorm.Model
	Digest 			string
	TagName 		string
	Size 			int64
	Architecture 	string
	OS 				string
	OSVersion 		string
	DockerVersion 	string
	Author 			string

	VerifyStatus 	string
	Cached 			bool
}
func (t RepositoryTag) TableName() string {
	return "repository_tag"
}

type ChartRepo struct {
	Name 			string
	VersionCount 	int64
	LatestVersion 	string
	Icon 			string
	Home 			string

	VerifyStatus 	string
	Cached 			bool
}
func (c ChartRepo) TableName() string{
	return "chart_repo"
}

type ChartVersion struct {
	Name 		string
	Version 	string
	Description	string
	AppVersion 	string
	Url 		string
	Digest 		string
}
func (c ChartVersion) TableName() string {
	return "chart_version"
}

func initDBTables() {
	tables := []interface{}{&SeederNode{}, &Repository{}, &ChartRepo{}, &ChartVersion{}}
	Db.DropTable(tables...)
	Db.CreateTable(tables...)

	Db = Db.Model(&SeederNode{})
	node1 := SeederNode{
		NodeName:"center-node",
		AdvertiseAddr:"10.23.100.2:15300",
		BindAddr:"192.168.0.2:8080",
		Role:"master",
		Status:"active",
	}
	if err := Db.Create(&node1).Error; err != nil {
		panic(err)
	}

	node2 := SeederNode{
		NodeName:"edge-node",
		AdvertiseAddr:"10.23.100.3:15300",
		BindAddr:"192.168.0.2:8080",
		Role:"follower",
		Status:"active",
	}
	if err := Db.Create(&node2).Error; err != nil {
		panic(err)
	}

	node3 := SeederNode{
		NodeName:"edge-node",
		AdvertiseAddr:"10.23.100.4:15300",
		BindAddr:"192.168.0.2:8080",
		Role:"follower",
		Status:"active",
	}
	if err := Db.Create(&node3).Error; err != nil {
		panic(err)
	}
}

func InitDB(DbAddr, DbName, User, Password string, needInitDb bool) error {
	var err error
	var connectURL string

	start := time.Now()
	connectURL = fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		User, Password, DbAddr, DbName)

	Db, err = gorm.Open("mysql", connectURL)
	if err != nil {
		return err
	}

	if needInitDb {
		initDBTables()
	}

	Db.AutoMigrate(&SeederNode{}, &Repository{}, &ChartRepo{}, &ChartVersion{})
	glog.V(2).Infof("InitDB(%s) spend %v", DbAddr, time.Since(start))

	return nil
}