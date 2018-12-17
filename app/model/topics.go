package model

import (
	"github.com/jinzhu/gorm"
)

type TopicsDao interface {
	CreateTopic(name string, engine string, owner string, description string, resourceID string) (*Topic, error)
	GetTopic(name string) (*Topic, error)
	DeleteTopic(name string) error
	ListTopics() ([]Topic, error)
}

func (db *DB) CreateTopic(name string, engine string, owner string, description string, resourceID string) (*Topic, error) {
	topic := Topic{Name: name, Engine: engine, ResourceID: resourceID, Owner: owner, Description: description}
	topic.CreatedAt = Clock.Now()
	topic.UpdatedAt = Clock.Now()
	//topic.ID = uuid.New()
	//topic.ID = 1
	if err := db.Save(&topic).Error; err != nil {
		return nil, err
	}
	return &topic, nil
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
	if err := db.Delete(&Topic{Name: name}).Error; err != nil {
		return err
	}
	return nil
}

func (db *DB) ListTopics() ([]Topic, error) {
	return nil, nil
}
