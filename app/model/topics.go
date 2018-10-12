package model

import (
	"github.com/jinzhu/gorm"
)

type Topics interface {
	CreateTopic(name string, engine string, resourceID string) error
	GetTopic(name string) (*Topic, error)
	DeleteTopic(name string) error
}

func (db *DB) CreateTopic(name string, engine string, resourceID string) error {
	topic := Topic{Name: name, Engine: engine, ResourceID: resourceID}
	db.NewRecord(topic)
	if err := db.Create(topic).Error; err != nil {
		return err
	}
	return nil
}

func (db *DB) GetTopic(name string) (*Topic, error) {
	var topic Topic
	if err := db.Where(&Topic{Name: name}).First(&topic).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &topic, nil
}

func (db *DB) DeleteTopic(name string) error {
	if err := db.Delete(&Topic{Name: name}).Error; err != nil{
		return err
	}
	return nil
}
