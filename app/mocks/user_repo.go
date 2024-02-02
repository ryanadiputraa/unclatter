// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	user "github.com/ryanadiputraa/unclatter/app/user"
	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// SaveOrUpdate provides a mock function with given fields: c, arg
func (_m *UserRepository) SaveOrUpdate(c context.Context, arg user.User) error {
	ret := _m.Called(c, arg)

	if len(ret) == 0 {
		panic("no return value specified for SaveOrUpdate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, user.User) error); ok {
		r0 = rf(c, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
