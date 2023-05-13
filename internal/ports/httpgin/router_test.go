package httpgin

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework10/internal/app"
	mocks "homework10/internal/mocks/appemocks"
	"log"
	"net/http"
	"testing"
)

type httpAppSuiteRoute struct {
	suite.Suite
	r      *gin.RouterGroup
	a      app.App
	logger *log.Logger
}

func (s *httpAppSuiteRoute) SetupTest() {
	s.r = gin.New().Group("")
	s.a = new(mocks.App)
	s.logger = log.Default()
}
func TestSuiteHttpRoute(t *testing.T) {
	u := new(httpAppSuiteRoute)
	suite.Run(t, u)

}

func (s *httpAppSuiteRoute) TestAppRouter() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		method string
		url    string
	}{
		{http.MethodGet, "/ads/:ad_id"},
		{http.MethodGet, "/ads"},
		{http.MethodPost, "/ads"},
		{http.MethodPut, "/ads/:ad_id/status"},
		{http.MethodPut, "/ads/:ad_id"},
		{http.MethodDelete, "/ads/:ad_id"},
		{http.MethodGet, "/users/:user_id"},
		{http.MethodPost, "/users"},
		{http.MethodPut, "/users/:user_id"},
		{http.MethodDelete, "/users/:user_id"},
	}

	g := gin.New()
	r := g.Group("")
	AppRouter(r, s.a, s.logger)
	routes := g.Routes()

	assert.NotEmpty(s.T(), routes)
	assert.Lenf(s.T(), routes, len(testCases), "routes count mismatch")

	for _, tc := range testCases {
		found := false
		for _, route := range routes {
			if route.Method == tc.method && route.Path == tc.url {
				found = true
				break
			}
		}
		assert.Truef(s.T(), found, "route %s %s not found", tc.method, tc.url)
	}
}
