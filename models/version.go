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

	db = db.Order("created desc")

	db = db.Find(&versions)
	glog.V(5).Infof("db query version count %d", count)

	return versions, count, nil
}

func GetUnCachedVersions(page, pageSize int, cacheStatus bool) ([]*ChartVersion, int, error) {
	var count int
	var versions []*ChartVersion
	db := Db.Model(&ChartVersion{})

	db = db.Where("cached = ?", cacheStatus)
	db = db.Count(&count)
	if page > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	db = db.Order("created desc")

	db = db.Find(&versions)
	glog.V(5).Infof("db query version count %d with cache %v", count, cacheStatus)

	return versions, count, nil
}

func UpdateVersionCached(chartName, version string) error {
	db := Db.Model(&ChartVersion{})
	db = db.Where("name = ?", chartName).Where("version = ?", version)
	db.Update("cached", true)
	return nil
}

func GetVersionByChart(page, pageSize int, chartName string)([]*ChartVersion, int, error) {
	var count int
	var versions []*ChartVersion
	db := Db.Model(&ChartVersion{})

	db = db.Where("name = ?", chartName)
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