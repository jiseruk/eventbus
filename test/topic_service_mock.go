package test

import client "github.com/wenance/wequeue-management_api/app/client"
import errors "github.com/wenance/wequeue-management_api/app/errors"
import mock "github.com/stretchr/testify/mock"
import model "github.com/wenance/wequeue-management_api/app/model"

// TopicServiceMock is an autogenerated mock type for the TopicServiceMock type
type TopicServiceMock struct {
	mock.Mock
}

// CreateTopic provides a mock function with given fields: name, owner, description, engine
func (_m *TopicServiceMock) CreateTopic(name string, owner string, description string, engine client.EngineService) (*model.Topic, *errors.APIError) {
	ret := _m.Called(name, owner, description, engine)

	var r0 *model.Topic
	if rf, ok := ret.Get(0).(func(string, string, string, client.EngineService) *model.Topic); ok {
		r0 = rf(name, owner, description, engine)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Topic)
		}
	}

	var r1 *errors.APIError
	if rf, ok := ret.Get(1).(func(string, string, string, client.EngineService) *errors.APIError); ok {
		r1 = rf(name, owner, description, engine)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.APIError)
		}
	}

	return r0, r1
}

// GetTopic provides a mock function with given fields: name, adminToken
func (_m *TopicServiceMock) GetTopic(name string, adminToken ...string) (*model.Topic, *errors.APIError) {
	_va := make([]interface{}, len(adminToken))
	for _i := range adminToken {
		_va[_i] = adminToken[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *model.Topic
	if rf, ok := ret.Get(0).(func(string, ...string) *model.Topic); ok {
		r0 = rf(name, adminToken...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Topic)
		}
	}

	var r1 *errors.APIError
	if rf, ok := ret.Get(1).(func(string, ...string) *errors.APIError); ok {
		r1 = rf(name, adminToken...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.APIError)
		}
	}

	return r0, r1
}

// ListTopics provides a mock function with given fields:
func (_m *TopicServiceMock) ListTopics() ([]model.Topic, *errors.APIError) {
	ret := _m.Called()

	var r0 []model.Topic
	if rf, ok := ret.Get(0).(func() []model.Topic); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Topic)
		}
	}

	var r1 *errors.APIError
	if rf, ok := ret.Get(1).(func() *errors.APIError); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.APIError)
		}
	}

	return r0, r1
}
