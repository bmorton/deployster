package server

import (
	"github.com/bmorton/deployster/fleet"
	"github.com/bmorton/deployster/fleet/mocks"
	"github.com/rcrowley/go-tigertonic/mocking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type DeploysResourceTestSuite struct {
	suite.Suite
	Subject         DeploysResource
	FleetClientMock *mocks.Client
	Service         *DeploysterService
}

func (suite *DeploysResourceTestSuite) SetupSuite() {
	suite.Service = NewDeploysterService("0.0.0.0:3000", "v1.0", "username", "password", "mmmhm")
}

func (suite *DeploysResourceTestSuite) SetupTest() {
	suite.FleetClientMock = new(mocks.Client)
	suite.Subject = DeploysResource{suite.FleetClientMock, "mmmhm"}
}

func (suite *DeploysResourceTestSuite) TestCreateWithoutDestroyPrevious() {
	suite.FleetClientMock.On("StartUnit", "carousel-abc123@1.service", mock.AnythingOfType("[]fleet.UnitOption")).Return(&http.Response{}, nil)

	code, _, response, err := suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{Deploy{"abc123", false}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 201, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithDestroyPreviousAndNoPreviousVersions() {
	suite.FleetClientMock.On("Units").Return([]fleet.Unit{}, nil)
	suite.FleetClientMock.On("StartUnit", "carousel-abc123@1.service", mock.AnythingOfType("[]fleet.UnitOption")).Return(&http.Response{}, nil)

	code, _, response, err := suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{Deploy{"abc123", true}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 201, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroyPrevious() {
	suite.FleetClientMock.On("UnitState", "carousel-abc123@1.service").Return(fleet.UnitState{"", "", "", "", "", "running"}, nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel-cccddd@1.service").Return(&http.Response{}, nil)

	suite.Subject.destroyPrevious("carousel", "cccddd", "abc123")

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroy() {
	suite.FleetClientMock.On("DestroyUnit", "carousel-cccddd@1.service").Return(&http.Response{}, nil)

	code, _, response, err := suite.Subject.Destroy(
		mocking.URL(suite.Service.RootMux, "DELETE", "http://example.com/v1/services/carousel/deploys/cccddd"),
		mocking.Header(nil),
		&DeployRequest{Deploy{"abc123", false}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 204, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func TestDeploysResourceTestSuite(t *testing.T) {
	suite.Run(t, new(DeploysResourceTestSuite))
}
