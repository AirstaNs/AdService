// Code generated by mockery v2.25.0. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "homework10/internal/entities"

	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, nickname, email
func (_m *UserService) CreateUser(ctx context.Context, nickname string, email string) (*entities.User, error) {
	ret := _m.Called(ctx, nickname, email)

	var r0 *entities.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*entities.User, error)); ok {
		return rf(ctx, nickname, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *entities.User); ok {
		r0 = rf(ctx, nickname, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, nickname, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: ctx, userID
func (_m *UserService) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 *entities.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*entities.User, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entities.User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveUser provides a mock function with given fields: ctx, userID
func (_m *UserService) RemoveUser(ctx context.Context, userID int64) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUser provides a mock function with given fields: ctx, UserID, Nickname, Email
func (_m *UserService) UpdateUser(ctx context.Context, UserID int64, Nickname string, Email string) (*entities.User, error) {
	ret := _m.Called(ctx, UserID, Nickname, Email)

	var r0 *entities.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, string) (*entities.User, error)); ok {
		return rf(ctx, UserID, Nickname, Email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, string) *entities.User); ok {
		r0 = rf(ctx, UserID, Nickname, Email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, string, string) error); ok {
		r1 = rf(ctx, UserID, Nickname, Email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserService interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserService(t mockConstructorTestingTNewUserService) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}