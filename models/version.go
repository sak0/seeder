package models

import "github.com/golang/glog"

func GetAllVersions(page, pageSize int) ([]*ChartVersion, int, error) {
	var count int
	var versions []*ChartVersion
	db := Db.Model(&ChartVersion{})

	db = db.Count(&count)
	if page > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	db = db.Find(&versions)
	glog.V(5).Infof("db query version count %d", count)

	return versions, count, nil
}

func CreateVersion(version *ChartVersion) error {
	glog.V(5).Infof("create version item: %v", version)

	db := Db.Model(&ChartVersion{})
	if err := db.Create(version).Error; err != nil {
		return err
	}
	return nil
}