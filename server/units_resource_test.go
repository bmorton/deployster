package server

import (
	"testing"

	"github.com/bmorton/deployster/server/mocks"
	"github.com/coreos/fleet/schema"
	"github.com/rcrowley/go-tigertonic/mocking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UnitsResourceTestSuite struct {
	suite.Suite
	Subject         UnitsResource
	FleetClientMock *mocks.FleetClient
	Service         *DeploysterService
}

func (suite *UnitsResourceTestSuite) SetupSuite() {
	suite.Service = NewDeploysterService("0.0.0.0:3000", "v1.0", "username", "password", "mmmhm")
}

func (suite *UnitsResourceTestSuite) SetupTest() {
	suite.FleetClientMock = new(mocks.FleetClient)
	suite.Subject = UnitsResource{suite.FleetClientMock}
}

func (suite *UnitsResourceTestSuite) TestIndexWithNoResults() {
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{}, nil)

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
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{&schema.Unit{"running", "running", "abc123", "differentapp:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}}}, nil)

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
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{&schema.Unit{"running", "running", "abc123", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}}}, nil)

	code, _, response, err := suite.Subject.Index(
		mocking.URL(suite.Service.RootMux, "GET", "http://example.com/v1/services/carousel/units"),
		mocking.Header(nil),
		nil,
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), &UnitsResponse{Units: []VersionedUnit{VersionedUnit{Service: "carousel", Instance: "1", Version: "efefeff", CurrentState: "running", DesiredState: "running", MachineID: "abc123", Timestamp: "2006.01.02-15.04.05"}}}, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *UnitsResourceTestSuite) TestIndexWithNonDeploysterManagedUnits() {
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{
		&schema.Unit{"running", "running", "abc123", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "abc123", "vulcand.service", []*schema.UnitOption{}},
	}, nil)

	code, _, response, err := suite.Subject.Index(
		mocking.URL(suite.Service.RootMux, "GET", "http://example.com/v1/services/carousel/units"),
		mocking.Header(nil),
		nil,
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), &UnitsResponse{Units: []VersionedUnit{VersionedUnit{Service: "carousel", Instance: "1", Version: "efefeff", CurrentState: "running", DesiredState: "running", MachineID: "abc123", Timestamp: "2006.01.02-15.04.05"}}}, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func TestUnitsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(UnitsResourceTestSuite))
}
