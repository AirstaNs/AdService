// Code generated by mockery v2.25.0. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "homework10/internal/entities"

	mock "github.com/stretchr/testify/mock"

	service "homework10/internal/service"

	util "homework10/internal/util"
)

// AdService is an autogenerated mock type for the AdService type
type AdService struct {
	mock.Mock
}

// ChangeAdStatus provides a mock function with given fields: ctx, adID, authorID, published
func (_m *AdService) ChangeAdStatus(ctx context.Context, adID int64, authorID int64, published bool) (*entities.Ad, error) {
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
func (_m *AdService) CreateAd(ctx context.Context, title string, text string, authorID int64) (*entities.Ad, error) {
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

// GetAdByID provides a mock function with given fields: ctx, adID
func (_m *AdService) GetAdByID(ctx context.Context, adID int64) (*entities.Ad, error) {
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
func (_m *AdService) GetAdsByFilter(ctx context.Context, filters service.AdFilters) ([]entities.Ad, error) {
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
func (_m *AdService) GetDateTimeFormat() util.DateTimeFormatter {
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

// RemoveAd provides a mock function with given fields: ctx, adID, authorID
func (_m *AdService) RemoveAd(ctx context.Context, adID int64, authorID int64) error {
	ret := _m.Called(ctx, adID, authorID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, adID, authorID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAd provides a mock function with given fields: ctx, adID, authorID, title, text
func (_m *AdService) UpdateAd(ctx context.Context, adID int64, authorID int64, title string, text string) (*entities.Ad, error) {
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

type mockConstructorTestingTNewAdService interface {
	mock.TestingT
	Cleanup(func())
}

// NewAdService creates a new instance of AdService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAdService(t mockConstructorTestingTNewAdService) *AdService {
	mock := &AdService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}