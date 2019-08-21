package models

import "github.com/golang/glog"

func GetAllRepos(page, pageSize int) ([]*Repository, int, error) {
	var count int
	var repos []*Repository
	db := Db.Model(&Repository{})

	db = db.Count(&count)
	if page > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	db = db.Find(&repos)
	glog.V(5).Infof("db query shop count %d", count)

	return repos, count, nil
}

func CreateRepo(repository *Repository) error {
	glog.V(5).Infof("create repo item: %v", repository)

	db := Db.Model(&Repository{})
	if err := db.Create(repository).Error; err != nil {
		return err
	}
	return nil
}