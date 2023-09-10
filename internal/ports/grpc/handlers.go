package grpc

import (
	"context"
	"errors"
	"github.com/AirstaNs/ValidationAds"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"homework10/internal/adapters/repository/userrepo"
	"homework10/internal/app"
	"homework10/internal/entities"
	"homework10/internal/service"
	"homework10/internal/util"
	"time"
)

var (
	errInvalidArgument = status.Error(codes.InvalidArgument, "invalid argument")
	errNotFound        = status.Error(codes.NotFound, "not found")
	errUnknown         = status.Error(codes.Unknown, "unknown error")
	errForbidden       = status.Error(codes.PermissionDenied, "permission denied")
)

type GServer struct {
	app.App
}

func (s GServer) AddAd(ctx context.Context, req *CreateAdRequest) (*AdResponse, error) {
	empty := &AdResponse{}
	id, err := s.GetUserByID(ctx, req.UserId)
	if err != nil {
		return empty, errNotFound
	}
	ad, err := s.App.CreateAd(ctx, req.Title, req.Text, id.ID)
	if err != nil {
		if isEmptyAuthorID := errors.Is(err, userrepo.ErrEmptyUser); isEmptyAuthorID {
			return empty, errNotFound
		}
		isBadTitle := errors.Is(err, ValidationAds.ErrBadTitle)
		isBadText := errors.Is(err, ValidationAds.ErrBadText)
		isBadAuthorID := errors.Is(err, ValidationAds.ErrBadAuthorID)
		if isBadTitle || isBadText || isBadAuthorID {
			return empty, errInvalidArgument
		}
		return empty, errUnknown
	}
	return AdSuccessResponse(ad), nil

}

func (s GServer) UpdateAdStatus(ctx context.Context, req *ChangeAdStatusRequest) (*AdResponse, error) {
	empty := &AdResponse{}
	app := s.App
	_, err := app.GetUserByID(ctx, req.UserId)

	if err != nil {
		return empty, errNotFound
	}
	adStatus, err := app.ChangeAdStatus(ctx, req.AdId, req.UserId, req.Published)

	if err != nil {
		isBadAuthorID := errors.Is(err, ValidationAds.ErrBadAuthorID)
		if isBadAuthorID {
			return empty, errForbidden
		}
		isBadTitle := errors.Is(err, ValidationAds.ErrBadTitle)
		isBadText := errors.Is(err, ValidationAds.ErrBadText)
		if isBadTitle || isBadText {
			return empty, errInvalidArgument
		}
		return empty, errUnknown
	}
	return AdSuccessResponse(adStatus), nil

}

func (s GServer) ModifyAd(ctx context.Context, req *UpdateAdRequest) (*AdResponse, error) {
	empty := &AdResponse{}
	app := s.App
	_, err := app.GetUserByID(ctx, req.UserId)

	if err != nil {
		return empty, errNotFound
	}
	ad, err := app.UpdateAd(ctx, req.AdId, req.UserId, req.Title, req.Text)
	if err != nil {
		isBadAuthorID := errors.Is(err, ValidationAds.ErrBadAuthorID)
		if isBadAuthorID {
			return empty, errForbidden
		}
		isBadTitle := errors.Is(err, ValidationAds.ErrBadTitle)
		isBadText := errors.Is(err, ValidationAds.ErrBadText)
		if isBadTitle || isBadText {
			return empty, errInvalidArgument
		}
		return empty, errUnknown
	}
	return AdSuccessResponse(ad), nil
}

func (s GServer) GetAd(ctx context.Context, req *GetADByIDRequest) (*AdResponse, error) {
	empty := &AdResponse{}
	ad, err := s.App.GetAdByID(ctx, req.AdId)
	if err != nil {
		isNotFound := errors.Is(err, util.ErrNotFound)
		if isNotFound {
			return empty, errNotFound
		}
		return empty, errUnknown
	}
	return AdSuccessResponse(ad), nil
}

func (s GServer) GetAds(ctx context.Context, filters *AdFilters) (*ListAdResponse, error) {
	empty := &ListAdResponse{}
	dateTime := time.Time{}

	cDate := filters.GetOptionalCreateDate()
	if cDate == nil {
		formatter := s.App.GetDateTimeFormat()
		dateTime, _ = formatter.ToTime(dateTime)
	} else {
		dateTime = cDate.AsTime().UTC()
	}

	title := filters.GetOptionalTitle()
	if title == nil {
		title = &wrapperspb.StringValue{Value: ""}
	}

	AuthorId := filters.GetOptionalAuthorId()
	if AuthorId == nil {
		AuthorId = &wrapperspb.Int64Value{Value: -1}
	}

	published := filters.GetOptionalPublished()
	if published == nil {
		published = &wrapperspb.BoolValue{Value: true}
	}

	adFilters := service.AdFilters{
		CreateDate: dateTime,
		Title:      title.GetValue(),
		AuthorID:   AuthorId.GetValue(),
		Published:  published.GetValue(),
	}

	ads, err := s.App.GetAdsByFilter(ctx, adFilters)
	if err != nil {
		return empty, errUnknown
	}

	response := AdListSuccessResponse(&ads)
	return &response, nil
}

func (s GServer) RemoveAd(ctx context.Context, req *DeleteAdRequest) (*DeleteAdResponse, error) {
	empty := &DeleteAdResponse{}
	err := s.App.RemoveAd(ctx, req.AdId, req.AuthorId)
	if err != nil {
		isBadAuthorID := errors.Is(err, ValidationAds.ErrBadAuthorID)
		if isBadAuthorID {
			return empty, errForbidden
		}
	}
	return &DeleteAdResponse{AdId: req.AdId, UserId: req.AuthorId}, nil
}

func (s GServer) ModifyUser(ctx context.Context, req *UserUpdateRequest) (*UserResponse, error) {
	empty := &UserResponse{}
	user, err := s.App.UpdateUser(ctx, req.Id, req.Nickname, req.Email)
	if err != nil {
		return empty, errUnknown
	}
	return UserSuccessResponse(user), nil
}

func (s GServer) AddUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	empty := &UserResponse{}
	user, err := s.App.CreateUser(ctx, req.Nickname, req.Email)
	if err != nil {
		return empty, errUnknown
	}
	return UserSuccessResponse(user), nil
}

func (s GServer) GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error) {
	empty := &UserResponse{}
	user, err := s.App.GetUserByID(ctx, req.Id)
	if err != nil {
		isNotFound := errors.Is(err, util.ErrNotFound)
		if isNotFound {
			return empty, errNotFound
		}
		return empty, errUnknown
	}
	return UserSuccessResponse(user), nil
}

func (s GServer) RemoveUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error) {
	_ = s.App.RemoveUser(ctx, req.Id)
	return &DeleteUserResponse{Id: req.Id}, nil
}

func (s GServer) mustEmbedUnimplementedAdServiceServer() {
}

func AdSuccessResponse(ad *entities.Ad) *AdResponse {
	return &AdResponse{
		Id:         ad.ID,
		Title:      ad.Title,
		Text:       ad.Text,
		AuthorId:   ad.AuthorID,
		Published:  ad.Published,
		CreateDate: timestamppb.New(ad.CreateDate),
		UpdateDate: timestamppb.New(ad.UpdateDate),
	}
}

func UserSuccessResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		Id:       user.ID,
		Nickname: user.Nickname,
		Email:    user.Email,
	}
}

func AdListSuccessResponse(ads *[]entities.Ad) ListAdResponse {
	adsResponse := make([]*AdResponse, 0)
	for _, a := range *ads {
		cDate := a.CreateDate
		uDate := a.UpdateDate
		ad := AdResponse{
			Id:         a.ID,
			Title:      a.Title,
			Text:       a.Text,
			AuthorId:   a.AuthorID,
			Published:  a.Published,
			CreateDate: &timestamppb.Timestamp{Seconds: cDate.Unix(), Nanos: int32(cDate.Nanosecond())},
			UpdateDate: &timestamppb.Timestamp{Seconds: uDate.Unix(), Nanos: int32(uDate.Nanosecond())},
		}
		adsResponse = append(adsResponse, &ad)
	}
	return ListAdResponse{List: adsResponse}
}
