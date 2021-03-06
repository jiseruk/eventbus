// Code generated by mockery v1.0.0. DO NOT EDIT.

package test

import mock "github.com/stretchr/testify/mock"
import model "github.com/wenance/wequeue-management_api/app/model"

// TopicsDaoMock is an autogenerated mock type for the TopicsDaoMock type
type TopicsDaoMock struct {
	mock.Mock
}

// CreateTopic provides a mock function with given fields: name, engine, resourceID
func (_m *TopicsDaoMock) CreateTopic(name string, engine string, resourceID string) (*model.Topic, error) {
	ret := _m.Called(name, engine, resourceID)

	var r0 *model.Topic
	if rf, ok := ret.Get(0).(func(string, string, string) *model.Topic); ok {
		r0 = rf(name, engine, resourceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Topic)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(name, engine, resourceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteTopic provides a mock function with given fields: name
func (_m *TopicsDaoMock) DeleteTopic(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetTopic provides a mock function with given fields: name
func (_m *TopicsDaoMock) GetTopic(name string) (*model.Topic, error) {
	ret := _m.Called(name)

	var r0 *model.Topic
	if rf, ok := ret.Get(0).(func(string) *model.Topic); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Topic)
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
