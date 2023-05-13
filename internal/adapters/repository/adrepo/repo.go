package adrepo

import (
	"homework10/internal/entities"
	"homework10/internal/util"
	"sync"
	"time"
)

type AdRepository interface {
	AddAd(ad entities.Ad) (int64, error)
	EditAdStatus(ad *entities.Ad, published bool, updateTime time.Time) (*entities.Ad, error)
	ChangeAdText(adID int64, title, text string, updateTime time.Time) (*entities.Ad, error)
	GetAdByID(adID int64) (*entities.Ad, error)
	GetAdsByFilters(filters []func(ad entities.Ad) bool) ([]entities.Ad, error)
	DeleteAd(adID int64) error
}

type mapRepository struct {
	rep    map[int64]entities.Ad
	mutex  sync.Mutex
	rMutex sync.RWMutex
	util.UID
}

func (m *mapRepository) AddAd(ad entities.Ad) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	const notValidID = -1
	id, err := m.UID.GenerateID()
	if err != nil {
		return notValidID, err
	}

	ad.ID = id
	m.rep[id] = ad
	return ad.ID, nil
}

func (m *mapRepository) EditAdStatus(ad *entities.Ad, published bool, updateTime time.Time) (*entities.Ad, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ad.Published = published
	ad.UpdateDate = updateTime

	m.rep[ad.ID] = *ad

	return ad, nil
}

func (m *mapRepository) ChangeAdText(adID int64, title, text string, updateTime time.Time) (*entities.Ad, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ad, err := m.GetAdByID(adID)
	if err != nil {
		return ad, err
	}
	ad.Title = title
	ad.Text = text
	ad.UpdateDate = updateTime

	m.rep[adID] = *ad
	return ad, nil
}

func (m *mapRepository) GetAdByID(adID int64) (*entities.Ad, error) {
	m.rMutex.RLock()
	defer m.rMutex.RUnlock()

	empty := &entities.Ad{}
	ad := m.rep[adID]
	if ad == (*empty) {
		return &ad, util.ErrNotFound
	}
	return &ad, nil
}

func (m *mapRepository) GetAdsByFilters(filters []func(ad entities.Ad) bool) ([]entities.Ad, error) {
	m.rMutex.RLock()
	defer m.rMutex.RUnlock()

	adsResult := make([]entities.Ad, 0)
adLoop:
	for _, ad := range m.rep {
		isValid := true
		for _, f := range filters {
			b := f(ad)
			if !isValid {
				continue adLoop
			}
			isValid = b
		}
		if isValid {
			adsResult = append(adsResult, ad)
		}
	}
	return adsResult, nil
}

func (m *mapRepository) DeleteAd(adID int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.rep, adID)
	return nil
}

func New() AdRepository {
	return &mapRepository{
		rep: make(map[int64]entities.Ad),
		UID: util.UID{Id: -1}}
}
