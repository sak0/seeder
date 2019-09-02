package models

import "github.com/golang/glog"

func GetAllCharts(page, pageSize int, chartName string) ([]*ChartRepo, int, error) {
	var count int
	var charts []*ChartRepo
	db := Db.Model(&ChartRepo{})

	if chartName != "" {
		db = db.Where("name = ?", chartName)
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

func CreateChart(chart *ChartRepo) error {
	glog.V(5).Infof("create chart repo item: %v", chart)

	db := Db.Model(&ChartRepo{})
	if err := db.Create(chart).Error; err != nil {
		return err
	}
	return nil
}