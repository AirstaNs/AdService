// Code generated by mockery v2.25.0. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "homework10/internal/entities"

	mock "github.com/stretchr/testify/mock"

	service "homework10/internal/service"

	util "homework10/internal/util"
)

// App is an autogenerated mock type for the App type
type App struct {
	mock.Mock
}

// ChangeAdStatus provides a mock function with given fields: ctx, adID, authorID, published
func (_m *App) ChangeAdStatus(ctx context.Context, adID int64, authorID int64, published bool) (*entities.Ad, error) {
	ret := _m.Called(ctx, adID, authorID, published)

	var r0 *entities.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, bool) (*entities.Ad, error)); ok {
		return rf(ctx, adID, authorID, published)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, bool) *entities.Ad); ok {
		r0 = rf(ctx, adID, authorID, published)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Ad)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, bool) error); ok {
		r1 = rf(ctx, adID, authorID, published)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAd provides a mock function with given fields: ctx, title, text, authorID
func (_m *App) CreateAd(ctx context.Context, title string, text string, authorID int64) (*entities.Ad, error) {
	ret := _m.Called(ctx, title, text, authorID)

	var r0 *entities.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int64) (*entities.Ad, error)); ok {
		return rf(ctx, title, text, authorID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int64) *entities.Ad); ok {
		r0 = rf(ctx, title, text, authorID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Ad)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int64) error); ok {
		r1 = rf(ctx, title, text, authorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, nickname, email
func (_m *App) CreateUser(ctx context.Context, nickname string, email string) (*entities.User, error) {
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

// GetAdByID provides a mock function with given fields: ctx, adID
func (_m *App) GetAdByID(ctx context.Context, adID int64) (*entities.Ad, error) {
	ret := _m.Called(ctx, adID)

	var r0 *entities.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*entities.Ad, error)); ok {
		return rf(ctx, adID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entities.Ad); ok {
		r0 = rf(ctx, adID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Ad)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, adID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAdsByFilter provides a mock function with given fields: ctx, filters
func (_m *App) GetAdsByFilter(ctx context.Context, filters service.AdFilters) ([]entities.Ad, error) {
	ret := _m.Called(ctx, filters)

	var r0 []entities.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, service.AdFilters) ([]entities.Ad, error)); ok {
		return rf(ctx, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, service.AdFilters) []entities.Ad); ok {
		r0 = rf(ctx, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entities.Ad)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, service.AdFilters) error); ok {
		r1 = rf(ctx, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDateTimeFormat provides a mock function with given fields:
func (_m *App) GetDateTimeFormat() util.DateTimeFormatter {
	ret := _m.Called()

	var r0 util.DateTimeFormatter
	if rf, ok := ret.Get(0).(func() util.DateTimeFormatter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(util.DateTimeFormatter)
		}
	}

	return r0
}

// GetUserByID provides a mock function with given fields: ctx, userID
func (_m *App) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
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

// RemoveAd provides a mock function with given fields: ctx, adID, authorID
func (_m *App) RemoveAd(ctx context.Context, adID int64, authorID int64) error {
	ret := _m.Called(ctx, adID, authorID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, adID, authorID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveUser provides a mock function with given fields: ctx, userID
func (_m *App) RemoveUser(ctx context.Context, userID int64) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAd provides a mock function with given fields: ctx, adID, authorID, title, text
func (_m *App) UpdateAd(ctx context.Context, adID int64, authorID int64, title string, text string) (*entities.Ad, error) {
	ret := _m.Called(ctx, adID, authorID, title, text)

	var r0 *entities.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, string, string) (*entities.Ad, error)); ok {
		return rf(ctx, adID, authorID, title, text)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, string, string) *entities.Ad); ok {
		r0 = rf(ctx, adID, authorID, title, text)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Ad)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, string, string) error); ok {
		r1 = rf(ctx, adID, authorID, title, text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: ctx, UserID, Nickname, Email
func (_m *App) UpdateUser(ctx context.Context, UserID int64, Nickname string, Email string) (*entities.User, error) {
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

type mockConstructorTestingTNewApp interface {
	mock.TestingT
	Cleanup(func())
}

// NewApp creates a new instance of App. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewApp(t mockConstructorTestingTNewApp) *App {
	mock := &App{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}