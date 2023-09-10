package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"homework10/internal/adapters/repository/adrepo"
	"homework10/internal/adapters/repository/userrepo"
	"homework10/internal/app"
	"homework10/internal/util"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"homework10/internal/ports/httpgin"
)

type adData struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
	AuthorID   int64     `json:"author_id"`
	Published  bool      `json:"published"`
	CreateDate time.Time `json:"create_date"`
	UpdateDate time.Time `json:"update_date"`
}

type adResponse struct {
	Data adData `json:"data"`
}

type adsResponse struct {
	Data []adData `json:"data"`
}

type userData struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}
type userResponse struct {
	Data userData `json:"data"`
}

type userDeleteResponse struct {
	UserId int64 `json:"user_id"`
}

type adDeleteResponse struct {
	AdId     int64 `json:"ad_id"`
	AuthorId int64 `json:"author_id"`
}

var (
	ErrBadRequest = fmt.Errorf("bad request")
	ErrForbidden  = fmt.Errorf("forbidden")
	ErrorNotFound = fmt.Errorf("not found")
)

type testClient struct {
	client  *http.Client
	baseURL string
}

type queryParam map[string]string

func getTestClient() *testClient {
	logger := log.New(io.Discard, "", 0)
	gin.DefaultWriter = io.Discard
	repo := adrepo.New()
	uRep := userrepo.New()
	formatter := util.NewDateTimeFormatter(time.RFC3339)
	newApp := app.NewApp(repo, uRep, formatter)
	server := httpgin.NewHTTPServer(":18080", newApp, logger, "*cert", "*key")
	httpServer := server.(*httpgin.HttpServer)
	testServer := httptest.NewServer(httpServer.App.Handler)

	return &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}

func (tc *testClient) getResponse(req *http.Request, out any) error {
	resp, err := tc.client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusAccepted {
		if resp.StatusCode == http.StatusBadRequest {
			return ErrBadRequest
		}
		if resp.StatusCode == http.StatusForbidden {
			return ErrForbidden
		}
		if resp.StatusCode == http.StatusNotFound {
			return ErrorNotFound
		}
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response: %w", err)
	}

	err = json.Unmarshal(respBody, out)
	if err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}

	return nil
}

func (tc *testClient) createAd(userID int64, title string, text string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) changeAdStatus(userID int64, adID int64, published bool) (adResponse, error) {
	body := map[string]any{
		"user_id":   userID,
		"published": published,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d/status", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) updateAd(userID int64, adID int64, title string, text string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

// фильтры передаются через query параметры.
func (tc *testClient) listAdsFilters(queryParam queryParam) (adsResponse, error) {
	params := parseQueryParams(queryParam)
	url := tc.baseURL + "/api/v1/ads" + params
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAds() (adsResponse, error) {
	var empty queryParam
	return tc.listAdsFilters(empty)
}

func parseQueryParams(queryParam queryParam) string {
	if len(queryParam) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteRune('?')

	for key, param := range queryParam {
		sb.WriteString(key)
		sb.WriteRune('=')
		sb.WriteString(param)
		sb.WriteRune('&')
	}
	return sb.String()[:len(sb.String())-1]
}

func (tc *testClient) getAdByID(adID int64) (adResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adID), nil)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) createUser(Nickname string, Email string) (userResponse, error) {
	body := map[string]any{
		"nickname": Nickname,
		"email":    Email,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/users", bytes.NewReader(data))
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) updateUser(userID int64, Nickname string, Email string) (userResponse, error) {
	body := map[string]any{
		"nickname": Nickname,
		"email":    Email,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", userID), bytes.NewReader(data))
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) getUserByID(userID int64) (userResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", userID), nil)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) deleteUser(userID int64) (userDeleteResponse, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", userID), nil)
	if err != nil {
		return userDeleteResponse{}, fmt.Errorf("unable to create request: %w", err)
	}
	var response userDeleteResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userDeleteResponse{}, err
	}

	return response, nil
}

func (tc *testClient) deleteAd(queryParam queryParam, adID int64) (adDeleteResponse, error) {
	params := parseQueryParams(queryParam)
	url := fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adID) + params
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return adDeleteResponse{}, fmt.Errorf("unable to create request: %w", err)
	}
	var response adDeleteResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adDeleteResponse{}, err
	}

	return response, nil
}
