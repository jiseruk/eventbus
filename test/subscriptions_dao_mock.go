package test

import mock "github.com/stretchr/testify/mock"
import model "github.com/wenance/wequeue-management_api/app/model"

// SubscriptionsDaoMock is an autogenerated mock type for the SubscriptionsDaoMock type
type SubscriptionsDaoMock struct {
	mock.Mock
}

// CreateSubscription provides a mock function with given fields: name, topic, Type, resource, endpoint, deadLetterQueue, pullingQueue, visibilityTimeout
func (_m *SubscriptionsDaoMock) CreateSubscription(name string, topic string, Type string, resource string, endpoint *string, deadLetterQueue string, pullingQueue string, visibilityTimeout *int) (*model.Subscriber, error) {
	ret := _m.Called(name, topic, Type, resource, endpoint, deadLetterQueue, pullingQueue, visibilityTimeout)

	var r0 *model.Subscriber
	if rf, ok := ret.Get(0).(func(string, string, string, string, *string, string, string, *int) *model.Subscriber); ok {
		r0 = rf(name, topic, Type, resource, endpoint, deadLetterQueue, pullingQueue, visibilityTimeout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Subscriber)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, *string, string, string, *int) error); ok {
		r1 = rf(name, topic, Type, resource, endpoint, deadLetterQueue, pullingQueue, visibilityTimeout)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSubscription provides a mock function with given fields: name
func (_m *SubscriptionsDaoMock) DeleteSubscription(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTopicSubscriptions provides a mock function with given fields: topic
func (_m *SubscriptionsDaoMock) DeleteTopicSubscriptions(topic string) error {
	ret := _m.Called(topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetSubscription provides a mock function with given fields: name
func (_m *SubscriptionsDaoMock) GetSubscription(name string) (*model.Subscriber, error) {
	ret := _m.Called(name)

	var r0 *model.Subscriber
	if rf, ok := ret.Get(0).(func(string) *model.Subscriber); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Subscriber)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscriptionByEndpoint provides a mock function with given fields: endpoint
func (_m *SubscriptionsDaoMock) GetSubscriptionByEndpoint(endpoint string) (*model.Subscriber, error) {
	ret := _m.Called(endpoint)

	var r0 *model.Subscriber
	if rf, ok := ret.Get(0).(func(string) *model.Subscriber); ok {
		r0 = rf(endpoint)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Subscriber)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(endpoint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscriptionsByTopic provides a mock function with given fields: topic
func (_m *SubscriptionsDaoMock) GetSubscriptionsByTopic(topic string) ([]model.Subscriber, error) {
	ret := _m.Called(topic)

	var r0 []model.Subscriber
	if rf, ok := ret.Get(0).(func(string) []model.Subscriber); ok {
		r0 = rf(topic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Subscriber)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(topic)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
