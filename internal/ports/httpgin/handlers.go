package httpgin

import (
	"errors"
	"github.com/AirstaNs/ValidationAds"
	"github.com/gin-gonic/gin"
	"homework10/internal/adapters/repository/userrepo"
	"homework10/internal/app"
	"homework10/internal/service"
	"homework10/internal/util"
	"net/http"
	"strconv"
)

var errConvert = errors.New("ad_id is not int")

func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createAdRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		id, err2 := a.GetUserByID(c, req.UserID)

		if err2 != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(err2))
			return
		}

		ad, err := a.CreateAd(c, req.Title, req.Text, id.ID)
		if err != nil {
			if isEmptyAuthorID := errors.Is(err, userrepo.ErrEmptyUser); isEmptyAuthorID {
				c.JSON(http.StatusNotFound, ErrorResponse(err))
				return
			}
			isBadTitle := errors.Is(err, ValidationAds.ErrBadTitle)
			isBadText := errors.Is(err, ValidationAds.ErrBadText)
			isBadAuthorID := errors.Is(err, ValidationAds.ErrBadAuthorID)
			if isBadTitle || isBadText || isBadAuthorID {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		}

		c.JSON(http.StatusCreated, AdSuccessResponse(ad))

	}
}

func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeAdStatusRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		strId := c.Param("ad_id")
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(errConvert))
			return
		}
		_, err2 := a.GetUserByID(c, req.UserID)

		if err2 != nil {
			c.JSON(http.StatusForbidden, ErrorResponse(err2))
			return
		}

		ad, err := a.ChangeAdStatus(c, id, req.UserID, req.Published)
		if err != nil {
			isBadAuthorID := errors.Is(err, ValidationAds.ErrBadAuthorID)
			if isBadAuthorID {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
			isBadTitle := errors.Is(err, ValidationAds.ErrBadTitle)
			isBadText := errors.Is(err, ValidationAds.ErrBadText)

			if isBadTitle || isBadText {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
				return
			}

			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updateAdRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		strId := c.Param("ad_id")
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(errConvert))
			return
		}
		gAd, err2 := a.GetAdByID(c, id)

		if err2 != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(err2))
			return
		}

		ad, err := a.UpdateAd(c, gAd.ID, req.UserID, req.Title, req.Text)
		if err != nil {
			isBadAuthorID := errors.Is(err, ValidationAds.ErrBadAuthorID)
			if isBadAuthorID {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
			isBadTitle := errors.Is(err, ValidationAds.ErrBadTitle)
			isBadText := errors.Is(err, ValidationAds.ErrBadText)

			if isBadTitle || isBadText {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
				return
			}

			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func getAdByID(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adId := c.Param("ad_id")
		id, err := strconv.ParseInt(adId, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(errConvert))
			return
		}

		ad, err := a.GetAdByID(c, id)
		if err != nil {
			isNotFound := errors.Is(err, util.ErrNotFound)
			if isNotFound {
				c.JSON(http.StatusNotFound, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func getAdsByFilter(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters service.AdFilters
		if err := c.ShouldBindQuery(&filters); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		ads, err := a.GetAdsByFilter(c, filters)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdListSuccessResponse(&ads))
	}
}

func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		strId := c.Param("ad_id")
		strUserID := c.Query("user_id")
		id, err := strconv.ParseInt(strId, 10, 64)
		uID, err1 := strconv.ParseInt(strUserID, 10, 64)

		if err != nil || err1 != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(errConvert))
			return
		}
		err = a.RemoveAd(c, id, uID)
		if err != nil {
			isBadAuthorID := errors.Is(err, ValidationAds.ErrBadAuthorID)
			if isBadAuthorID {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
		}
		c.JSON(http.StatusOK, DeleteAdSuccessResponse(id, uID))
	}
}

/*

 */

func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateUserRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		strUserId := c.Param("user_id")
		userId, err := strconv.ParseInt(strUserId, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(errConvert))
			return
		}
		user, err := a.UpdateUser(c, userId, req.Nickname, req.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createUserRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		user, err := a.CreateUser(c, req.Nickname, req.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusCreated, UserSuccessResponse(user))
	}
}

func getUserByID(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		strUserId := c.Param("user_id")
		userId, err := strconv.ParseInt(strUserId, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(errConvert))
			return
		}
		user, err := a.GetUserByID(c, userId)
		if err != nil {
			isNotFound := errors.Is(err, userrepo.ErrEmptyUser)
			if isNotFound {
				c.JSON(http.StatusNotFound, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

func deleteUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		strUserId := c.Param("user_id")
		userId, err := strconv.ParseInt(strUserId, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(errConvert))
			return
		}
		_ = a.RemoveUser(c, userId)
		c.JSON(http.StatusOK, DeleteUserSuccessResponse(userId))
	}
}
