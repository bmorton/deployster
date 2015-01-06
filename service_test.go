package main

import (
	"github.com/bmorton/deployster/support"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type DeploysterServiceTestSuite struct {
	suite.Suite
	Subject *DeploysterService
}

func (suite *DeploysterServiceTestSuite) SetupSuite() {
	suite.Subject = NewDeploysterService("0.0.0.0:3000", "v1.0")
}

func (suite *DeploysterServiceTestSuite) TestGetVersionRequiresAuthentication() {
	w := &support.TestResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/v1/version", nil)
	suite.Subject.RootMux.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.StatusCode)
}

func (suite *DeploysterServiceTestSuite) TestPostDeploysRequiresAuthentication() {
	w := &support.TestResponseWriter{}
	r, _ := http.NewRequest("POST", "http://example.com/v1/services/test/deploys", nil)
	suite.Subject.RootMux.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.StatusCode)
}

func (suite *DeploysterServiceTestSuite) TestDeleteDeploysRequiresAuthentication() {
	w := &support.TestResponseWriter{}
	r, _ := http.NewRequest("DELETE", "http://example.com/v1/services/test/deploys/abc123", nil)
	suite.Subject.RootMux.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.StatusCode)
}

func (suite *DeploysterServiceTestSuite) TestGetUnitsRequiresAuthentication() {
	w := &support.TestResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/v1/services/test/units", nil)
	suite.Subject.RootMux.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.StatusCode)
}

func TestDeploysterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DeploysterServiceTestSuite))
}
