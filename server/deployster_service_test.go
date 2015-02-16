package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DeploysterServiceTestSuite struct {
	suite.Suite
	Subject *DeploysterService
}

func (suite *DeploysterServiceTestSuite) SetupSuite() {
	suite.Subject = NewDeploysterService("0.0.0.0:3000", "v1.0", "username", "password", "mmmhm")
}

func (suite *DeploysterServiceTestSuite) TestGetVersionRequiresAuthentication() {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/v1/version", nil)
	suite.Subject.RootMux.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *DeploysterServiceTestSuite) TestPostDeploysRequiresAuthentication() {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "http://example.com/v1/services/test/deploys", nil)
	suite.Subject.RootMux.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *DeploysterServiceTestSuite) TestDeleteDeploysRequiresAuthentication() {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "http://example.com/v1/services/test/deploys/abc123", nil)
	suite.Subject.RootMux.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *DeploysterServiceTestSuite) TestGetUnitsRequiresAuthentication() {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/v1/services/test/units", nil)
	suite.Subject.RootMux.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func TestDeploysterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DeploysterServiceTestSuite))
}
