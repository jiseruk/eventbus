package test

import mock "github.com/stretchr/testify/mock"
import model "github.com/wenance/wequeue-management_api/app/model"

// SubscriptionsDaoMock is an autogenerated mock type for the SubscriptionsDaoMock type
type SubscriptionsDaoMock struct {
	mock.Mock
}

// CreateSubscription provides a mock function with given fields: name, topic, endpoint, resource
func (_m *SubscriptionsDaoMock) CreateSubscription(name string, topic string, endpoint string, resource string) (*model.Subscriber, error) {
	ret := _m.Called(name, topic, endpoint, resource)

	var r0 *model.Subscriber
	if rf, ok := ret.Get(0).(func(string, string, string, string) *model.Subscriber); ok {
		r0 = rf(name, topic, endpoint, resource)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Subscriber)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string) error); ok {
		r1 = rf(name, topic, endpoint, resource)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
