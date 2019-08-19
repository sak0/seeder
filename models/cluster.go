package models

import "github.com/golang/glog"

func GetSeederNodes(page, pageSize int) ([]*SeederNode, int, error) {
	var count int
	var nodes []*SeederNode
	db := Db.Model(&SeederNode{})

	db = db.Count(&count)
	if page > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	db = db.Find(&nodes)
	glog.V(5).Infof("db query node count %d", count)

	return nodes, count, nil
}