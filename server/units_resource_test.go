package server

import (
	"github.com/bmorton/deployster/fleet"
	"github.com/bmorton/deployster/fleet/mocks"
	"github.com/rcrowley/go-tigertonic/mocking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UnitsResourceTestSuite struct {
	suite.Suite
	Subject         UnitsResource
	FleetClientMock *mocks.Client
	Service         *DeploysterService
}

func (suite *UnitsResourceTestSuite) SetupSuite() {
	suite.Service = NewDeploysterService("0.0.0.0:3000", "v1.0", "username", "password", "mmmhm")
}

func (suite *UnitsResourceTestSuite) SetupTest() {
	suite.FleetClientMock = new(mocks.Client)
	suite.Subject = UnitsResource{suite.FleetClientMock}
}

func (suite *UnitsResourceTestSuite) TestIndexWithNoResults() {
	suite.FleetClientMock.On("Units").Return([]fleet.Unit{}, nil)

	code, _, response, err := suite.Subject.Index(
		mocking.URL(suite.Service.RootMux, "GET", "http://example.com/v1/services/carousel/units"),
		mocking.Header(nil),
		nil,
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), &UnitsResponse{Units: []VersionedUnit{}}, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *UnitsResourceTestSuite) TestIndexWithNoMatchingResultsForService() {
	suite.FleetClientMock.On("Units").Return([]fleet.Unit{fleet.Unit{"running", "running", "abc123", "differentapp-efefeff@1.service", []fleet.UnitOption{}}}, nil)

	code, _, response, err := suite.Subject.Index(
		mocking.URL(suite.Service.RootMux, "GET", "http://example.com/v1/services/carousel/units"),
		mocking.Header(nil),
		nil,
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), &UnitsResponse{Units: []VersionedUnit{}}, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *UnitsResourceTestSuite) TestIndexWithMatchingResultsForService() {
	suite.FleetClientMock.On("Units").Return([]fleet.Unit{fleet.Unit{"running", "running", "abc123", "carousel-efefeff@1.service", []fleet.UnitOption{}}}, nil)

	code, _, response, err := suite.Subject.Index(
		mocking.URL(suite.Service.RootMux, "GET", "http://example.com/v1/services/carousel/units"),
		mocking.Header(nil),
		nil,
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), &UnitsResponse{Units: []VersionedUnit{VersionedUnit{Service: "carousel", Instance: "1", Version: "efefeff", CurrentState: "running", DesiredState: "running", MachineID: "abc123"}}}, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func TestUnitsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(UnitsResourceTestSuite))
}
