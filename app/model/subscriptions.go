package model

import "github.com/jinzhu/gorm"

type SubscriptionsDao interface {
	CreateSubscription(name string, topic string, endpoint string, resource string, pullResource string) (*Subscriber, error)
	GetSubscription(name string) (*Subscriber, error)
	GetSubscriptionByEndpoint(endpoint string) (*Subscriber, error)
}

type SubscriberDaoImpl struct {
	Db DB
}

//func (db *DB) CreateSubscription(name string, topic string, endpoint string, resource string, pullResource string) (*Subscriber, error) {
func (s *SubscriberDaoImpl) CreateSubscription(name string, topic string, endpoint string, resource string, pullResource string) (*Subscriber, error) {
	subscription := Subscriber{Name: name, Topic: topic, Endpoint: endpoint, ResourceID: resource, DeadLetterQueue: pullResource}
	subscription.CreatedAt = Clock.Now()
	subscription.UpdatedAt = Clock.Now()

	if err := s.Db.Save(&subscription).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (s *SubscriberDaoImpl) GetSubscription(name string) (*Subscriber, error) {
	return s.getSubscriptionByField("name", name)
	/*	var subscription Subscriber
		if err := db.Where(&Subscriber{Name: name}).First(&subscription).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return nil, nil
			}
			return nil, err
		}
		return &subscription, nil*/
}

func (s *SubscriberDaoImpl) GetSubscriptionByEndpoint(endpoint string) (*Subscriber, error) {
	return s.getSubscriptionByField("endpoint", endpoint)
}

func (s *SubscriberDaoImpl) getSubscriptionByField(field string, value interface{}) (*Subscriber, error) {
	var subscription Subscriber
	if err := s.Db.Where(map[string]interface{}{field: value}).First(&subscription).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &subscription, nil

}
