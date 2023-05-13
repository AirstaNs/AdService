package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework10/internal/entities"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
	AuthorID   int64     `json:"author_id"`
	Published  bool      `json:"published"`
	CreateDate time.Time `json:"create_date"`
	UpdateDate time.Time `json:"update_date"`
}

type changeAdStatusRequest struct {
	UserID    int64 `json:"user_id"`
	Published bool  `json:"published"`
}

type updateAdRequest struct {
	UserID int64  `json:"user_id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

type FilterAdRequest struct {
	Published  bool      `json:"published"`
	UserID     int64     `json:"user_id"`
	CreateDate time.Time `json:"create_Date"`
}

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func AdSuccessResponse(ad *entities.Ad) gin.H {
	return gin.H{
		"data": adResponse{
			ID:         ad.ID,
			Title:      ad.Title,
			Text:       ad.Text,
			AuthorID:   ad.AuthorID,
			Published:  ad.Published,
			CreateDate: ad.CreateDate,
			UpdateDate: ad.UpdateDate,
		},
		"error": nil,
	}
}
func AdListSuccessResponse(ads *[]entities.Ad) gin.H {
	adsResponse := make([]adResponse, 0)
	for _, a := range *ads {
		ad := adResponse{
			ID:         a.ID,
			Title:      a.Title,
			Text:       a.Text,
			AuthorID:   a.AuthorID,
			Published:  a.Published,
			CreateDate: a.CreateDate,
			UpdateDate: a.UpdateDate,
		}
		adsResponse = append(adsResponse, ad)
	}
	return gin.H{
		"data":  adsResponse,
		"error": nil,
	}
}

func ErrorResponse(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func UserSuccessResponse(user *entities.User) gin.H {
	return gin.H{
		"data":  user,
		"error": nil,
	}
}

func DeleteUserSuccessResponse(userID int64) gin.H {
	return gin.H{
		"data":  gin.H{"user_id": userID},
		"error": nil,
	}
}

func DeleteAdSuccessResponse(adID int64, authorID int64) gin.H {
	return gin.H{
		"data":  gin.H{"ad_id": adID, "author_id": authorID},
		"error": nil,
	}
}
