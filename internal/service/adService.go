package service

import (
	"github.com/AirstaNs/ValidationAds"
	"golang.org/x/net/context"
	"homework10/internal/adapters/repository/adrepo"
	"homework10/internal/entities"
	"homework10/internal/util"
	"strings"
	"time"
)

type adService struct {
	adRepository   adrepo.AdRepository
	dateTimeFormat util.DateTimeFormatter
}

//go:generate go run github.com/vektra/mockery/v2@v2.25.0 --name=AdService --filename=mockAdservice.go --output ../mocks/servicemocks
type AdService interface {
	CreateAd(ctx context.Context, title string, text string, authorID int64) (*entities.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, authorID int64, published bool) (*entities.Ad, error)
	UpdateAd(ctx context.Context, adID int64, authorID int64, title string, text string) (*entities.Ad, error)
	GetAdByID(ctx context.Context, adID int64) (*entities.Ad, error)
	GetAdsByFilter(ctx context.Context, filters AdFilters) ([]entities.Ad, error)
	GetDateTimeFormat() util.DateTimeFormatter
	RemoveAd(ctx context.Context, adID int64, authorID int64) error
}

type AdFilters struct {
	AuthorID   int64     `form:"user_id,query,default=-1"`
	Published  bool      `form:"published,query,default=true"`
	CreateDate time.Time `form:"create_Date,query,default=0001-01-01T00:00:00Z"`
	Title      string    `form:"title,query"`
}

func NewAdsService(adRepo adrepo.AdRepository, dateTimeFormatter util.DateTimeFormatter) AdService {
	return &adService{
		adRepository:   adRepo,
		dateTimeFormat: dateTimeFormatter,
	}
}

func (a *adService) CreateAd(ctx context.Context, title string, text string, authorID int64) (*entities.Ad, error) {
	parse, err := a.dateTimeFormat.ToTime(time.Now().UTC())
	if err != nil {
		return nil, err
	}
	ad := entities.Ad{
		Title:      title,
		Text:       text,
		AuthorID:   authorID,
		Published:  false,
		CreateDate: parse,
	}
	ad.UpdateDate = ad.CreateDate

	if err = ValidationAds.ValidateAuthorID(ad.AuthorID, authorID); err != nil {
		return &ad, err
	}

	if err = ValidationAds.ValidateTitle(title); err != nil {
		return &ad, err
	}
	if err = ValidationAds.ValidateText(text); err != nil {
		return &ad, err
	}

	id, err := a.adRepository.AddAd(ad)
	ad.ID = id

	if err != nil {
		return &ad, err
	}

	return &ad, nil
}

func (a *adService) ChangeAdStatus(ctx context.Context, adID int64, authorID int64, published bool) (*entities.Ad, error) {
	ad, err := a.adRepository.GetAdByID(adID)
	if err != nil {
		return ad, err
	}
	if err = ValidationAds.ValidateAuthorID(ad.AuthorID, authorID); err != nil {
		return ad, err
	}

	dateUpdate, err := a.dateTimeFormat.ToTime(time.Now().UTC())
	if err != nil {
		return ad, err
	}

	return a.adRepository.EditAdStatus(ad, published, dateUpdate)

}

func (a *adService) UpdateAd(ctx context.Context, adID int64, authorID int64, title string, text string) (*entities.Ad, error) {
	ad, err := a.adRepository.GetAdByID(adID)
	if err != nil {
		return ad, err
	}

	if err = ValidationAds.ValidateAuthorID(ad.AuthorID, authorID); err != nil {
		return ad, err
	}

	if err = ValidationAds.ValidateTitle(title); err != nil {
		return ad, err
	}
	if err = ValidationAds.ValidateText(text); err != nil {
		return ad, err
	}

	dateUpdate, err := a.dateTimeFormat.ToTime(time.Now().UTC())
	if err != nil {
		return ad, err
	}
	return a.adRepository.ChangeAdText(adID, title, text, dateUpdate)
}

func (a *adService) GetAdByID(ctx context.Context, adID int64) (*entities.Ad, error) {
	return a.adRepository.GetAdByID(adID)
}

// GetAdsByFilter Поиск объявлений по названию тоже организован через фильтры
func (a *adService) GetAdsByFilter(ctx context.Context, filters AdFilters) ([]entities.Ad, error) {
	var adFilters []func(ad entities.Ad) bool

	if filters.AuthorID != -1 {
		adFilters = append(adFilters, func(ad entities.Ad) bool {
			return ad.AuthorID == filters.AuthorID
		})
	}

	if !filters.CreateDate.IsZero() {
		adFilters = append(adFilters, func(ad entities.Ad) bool {
			return ad.CreateDate.Equal(filters.CreateDate)
		})
	}

	if filters.Title != "" {
		adFilters = append(adFilters, func(ad entities.Ad) bool {
			return strings.EqualFold(ad.Title, filters.Title)
		})
	}

	emptyFilters := len(adFilters) == 0
	isPublished := !filters.Published
	if emptyFilters || isPublished {
		adFilters = append(adFilters, func(ad entities.Ad) bool {
			return ad.Published == filters.Published
		})
	}

	return a.adRepository.GetAdsByFilters(adFilters)
}

func (a *adService) RemoveAd(ctx context.Context, adID int64, authorID int64) error {
	ad, err := a.adRepository.GetAdByID(adID)
	if err != nil {
		return err
	}
	if err = ValidationAds.ValidateAuthorID(ad.AuthorID, authorID); err != nil {
		return err
	}
	return a.adRepository.DeleteAd(adID)
}

func (a *adService) GetDateTimeFormat() util.DateTimeFormatter {
	return a.dateTimeFormat
}
