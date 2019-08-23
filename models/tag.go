package models

import (
	"github.com/golang/glog"
)

func GetAllTags(page, pageSize int) ([]*RepositoryTag, int, error) {
	var count int
	var tags []*RepositoryTag
	db := Db.Model(&RepositoryTag{})

	db = db.Count(&count)
	if page > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	db = db.Find(&tags)
	glog.V(5).Infof("db query tag count %d", count)

	return tags, count, nil
}


func CreateTag(tag *RepositoryTag) error {
	glog.V(5).Infof("create tag item: %v", tag)

	db := Db.Model(&RepositoryTag{})
	if err := db.Create(tag).Error; err != nil {
		return err
	}
	return nil
}