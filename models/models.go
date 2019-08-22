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
	ImageCount 			int			`json:"image_count" gorm:"type:int;column:image_count"`
	ChartCount 			int			`json:"chart_count" gorm:"type:int;column:chart_count"`
	PullCount 			int			`json:"pull_count" gorm:"type:int;column:pull_count"`
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
	TagsCount 		int64		`json:"tag_count" gorm:"type:int;column:tags_count"`
	VerifyStatus 	string		`json:"verify_status" gorm:"type:varchar(50);column:verify_status"`
	Cached 			bool		`json:"cached" gorm:"type:bool;column:cached"`
}
func (s Repository) TableName() string {
	return "repository"
}

type RepositoryTag struct {
	gorm.Model
	Digest 			string		`json:"digest" gorm:"type:varchar(255);column:digest"`
	TagName 		string		`json:"tag_name" gorm:"type:varchar(50);column:tag_name"`
	Size 			int64		`json:"size" gorm:"type:int;column:size"`
	Architecture 	string		`json:"architecture" gorm:"type:varchar(50);column:architecture"`
	OS 				string		`json:"os" gorm:"type:varchar(50);column:os"`
	OSVersion 		string		`json:"os_version" gorm:"type:varchar(50);column:os_version"`
	DockerVersion 	string		`json:"docker_version" gorm:"type:varchar(50);column:docker_version"`
	Author 			string		`json:"author" gorm:"type:varchar(50);column:author"`
	VerifyStatus 	string		`json:"verify_status" gorm:"type:varchar(50);column:verify_status"`
	Cached 			bool		`json:"cached" gorm:"type:bool;column:cached"`
}
func (t RepositoryTag) TableName() string {
	return "repository_tag"
}

type ChartRepo struct {
	OwnerNode 		string		`json:"owner_name" gorm:"type:varchar(50);column:owner_name"`
	Name 			string		`json:"name" gorm:"type:varchar(50);column:name"`
	VersionCount 	int64		`json:"version_count" gorm:"type:int;column:size"`
	LatestVersion 	string		`json:"latest_version" gorm:"type:varchar(50);column:latest_version"`
	Icon 			string		`json:"icon" gorm:"type:varchar(50);column:icon"`
	Home 			string		`json:"home" gorm:"type:varchar(50);column:home"`
	VerifyStatus 	string		`json:"verify_status" gorm:"type:varchar(50);column:verify_status"`
	Cached 			bool		`json:"cached" gorm:"type:bool;column:cached"`
}
func (c ChartRepo) TableName() string{
	return "chart_repo"
}

type ChartVersion struct {
	Name 		string			`json:"name" gorm:"type:varchar(50);column:name"`
	Version 	string			`json:"version" gorm:"type:varchar(50);column:version"`
	Description	string			`json:"description" gorm:"type:varchar(255);column:description"`
	AppVersion 	string			`json:"app_version" gorm:"type:varchar(50);column:app_version"`
	Url 		string			`json:"url" gorm:"type:varchar(50);column:url"`
	Digest 		string			`json:"digest" gorm:"type:varchar(50);column:digest"`
}
func (c ChartVersion) TableName() string {
	return "chart_version"
}

func initDBTables() {
	tables := []interface{}{&SeederNode{}, &Repository{}, &RepositoryTag{}, &ChartRepo{}, &ChartVersion{}}
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

	Db.AutoMigrate(&SeederNode{}, &Repository{}, &RepositoryTag{}, &ChartRepo{}, &ChartVersion{})
	glog.V(2).Infof("InitDB(%s) spend %v", DbAddr, time.Since(start))

	return nil
}