package models

import "github.com/golang/glog"

func GetAllCharts(page, pageSize int, chartName, typeName string) ([]*ChartRepo, int, error) {
	var count int
	var charts []*ChartRepo
	db := Db.Model(&ChartRepo{})

	if chartName != "" {
		db = db.Where("name = ?", chartName)
	}
	if typeName != "" {
		db = db.Where("type = ?", typeName)
	}

	db = db.Count(&count)
	if page > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	db = db.Order("created desc")

	db = db.Find(&charts)
	glog.V(5).Infof("db query charts count %d", count)

	return charts, count, nil
}

func GetAllCachedCharts(page, pageSize int, chartName, typeName string, cached bool) ([]*ChartRepo, int, error) {
	var count int
	var charts []*ChartRepo
	db := Db.Model(&ChartRepo{})

	if chartName != "" {
		db = db.Where("name = ?", chartName)
	}
	if typeName != "" {
		db = db.Where("type = ?", typeName)
	}

	db = db.Where("cached = ?", cached)

	db = db.Count(&count)
	if page > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	db = db.Order("created desc")

	db = db.Find(&charts)
	glog.V(5).Infof("db query charts count %d", count)

	return charts, count, nil
}

func CreateChart(chart *ChartRepo) error {
	glog.V(5).Infof("create chart repo item: %v", chart)

	db := Db.Model(&ChartRepo{})
	if err := db.Create(chart).Error; err != nil {
		return err
	}
	return nil
}

func UpdateChartCached(chartName string) error {
	db := Db.Model(&ChartRepo{})
	db = db.Where("name = ?", chartName)
	db.Update("cached", true)
	return nil
}

func DeleteChartByName(chartName string) error {
	var count int
	var charts []*ChartRepo
	glog.V(5).Infof("delete chart item: %s", chartName)

	db := Db.Model(&ChartRepo{})
	db = db.Where("name = ?", chartName)

	db = db.Count(&count)
	if count > 1 {
		glog.Warningf("there have %d charts with name %s", count, chartName)
	}
	db = db.Find(&charts)

	db.Delete(&charts)

	return nil
}