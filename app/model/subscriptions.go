package model

import "github.com/jinzhu/gorm"

type SubscriptionsDao interface {
	CreateSubscription(name string, topic string, endpoint string, resource string) (*Subscriber, error)
	GetSubscription(name string) (*Subscriber, error)
}

func (db *DB) CreateSubscription(name string, topic string, endpoint string, resource string) (*Subscriber, error) {
	subscription := Subscriber{Name: name, Topic: topic, Endpoint: endpoint, ResourceID: resource}
	subscription.CreatedAt = Clock.Now()
	subscription.UpdatedAt = Clock.Now()

	if err := db.Save(&subscription).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (db *DB) GetSubscription(name string) (*Subscriber, error) {
	var subscription Subscriber
	if err := db.Where(&Subscriber{Name: name}).First(&subscription).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &subscription, nil
}
