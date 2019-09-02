package models

import (
	"github.com/golang/glog"
	"fmt"
)

const (
	RoleMaster 				= "master"
	RoleFollower 			= "follower"

	NodeStatusActive 		= "active"
	NodeStatusUnknown		= "unknown"
	NodeStatusDown 			= "down"
)

func GetNodeByName(name string) (*SeederNode, error) {
	var count int
	var node SeederNode
	db := Db.Model(&SeederNode{})

	db = db.Where("cluster_name = ?", name)
	db = db.Count(&count)
	glog.V(5).Infof("db query nodes count %d", count)
	db = db.Find(&node)

	if count < 1 {
		return nil, fmt.Errorf(fmt.Sprintf("can not fild node with name %s", name))
	} else if count > 1 {
		glog.Warningf("node %s have record %d", name, count)
	}

	return &node, nil
}

func GetNodesByRole(role string) ([]*SeederNode, error) {
	var count int
	var nodes []*SeederNode
	db := Db.Model(&SeederNode{})

	db = db.Where("role = ?", role)
	db = db.Count(&count)
	glog.V(5).Infof("db query nodes count %d", count)
	db = db.Find(&nodes)

	if count < 1 {
		return nil, fmt.Errorf(fmt.Sprintf("can not fild node with role %s", role))
	} else if count > 1 {
		glog.Warningf("have %d records for role %s", count, role)
	}

	return nodes, nil
}