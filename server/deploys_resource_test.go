package server

import (
	"fmt"
	"testing"

	"github.com/bmorton/deployster/server/mocks"
	"github.com/coreos/fleet/schema"
	"github.com/rcrowley/go-tigertonic/mocking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeploysResourceTestSuite struct {
	suite.Suite
	Subject         DeploysResource
	FleetClientMock *mocks.FleetClient
	Service         *DeploysterService
}

func (suite *DeploysResourceTestSuite) SetupSuite() {
	suite.Service = NewDeploysterService("0.0.0.0:3000", "v1.0", "username", "password", "mmmhm")
}

func (suite *DeploysResourceTestSuite) SetupTest() {
	suite.FleetClientMock = new(mocks.FleetClient)
	suite.Subject = DeploysResource{suite.FleetClientMock, "mmmhm"}
}

func (suite *DeploysResourceTestSuite) TestCreateWithoutPassedInstancesAndNoInstancesRunning() {
	expectedOptions := getUnitOptions(UnitTemplate{"carousel", "abc123", "mmmhm", "2006.01.02-15.04.05"})
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{}, nil)

	// Should only start 1 unit
	suite.FleetClientMock.On("CreateUnit", &schema.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@1.service", "launched").Return(nil)

	suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{Deploy{Version: "abc123", DestroyPrevious: false, Timestamp: "2006.01.02-15.04.05"}},
	)

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithoutPassedInstancesAndMultipleInstancesRunning() {
	expectedOptions := getUnitOptions(UnitTemplate{"carousel", "abc123", "mmmhm", "2006.01.02-15.04.05"})
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@2.service", []*schema.UnitOption{}},
	}, nil)

	// Should start 2 units
	suite.FleetClientMock.On("CreateUnit", &schema.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@1.service", "launched").Return(nil)
	suite.FleetClientMock.On("CreateUnit", &schema.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@2.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@2.service", "launched").Return(nil)

	suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{Deploy{Version: "abc123", DestroyPrevious: false, Timestamp: "2006.01.02-15.04.05"}},
	)

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithoutPassedInstancesAndMultipleVersionsRunning() {
	expectedOptions := getUnitOptions(UnitTemplate{"carousel", "abc123", "mmmhm", "2006.01.02-15.04.05"})
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "abbbbbb", "carousel:abbbbbb:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
	}, nil)

	// Should start 1 unit
	suite.FleetClientMock.On("CreateUnit", &schema.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@1.service", "launched").Return(nil)

	suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{Deploy{Version: "abc123", DestroyPrevious: false, Timestamp: "2006.01.02-15.04.05"}},
	)

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithoutDestroyPrevious() {
	expectedOptions := getUnitOptions(UnitTemplate{"carousel", "abc123", "mmmhm", "2006.01.02-15.04.05"})
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{}, nil)
	suite.FleetClientMock.On("CreateUnit", &schema.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@1.service", "launched").Return(nil)

	code, _, response, err := suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{Deploy{"abc123", false, "2006.01.02-15.04.05", 1}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 201, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithDestroyPreviousAndNoPreviousVersions() {
	expectedOptions := getUnitOptions(UnitTemplate{"carousel", "abc123", "mmmhm", "2006.01.02-15.04.05"})
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{}, nil)
	suite.FleetClientMock.On("CreateUnit", &schema.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@1.service", "launched").Return(nil)

	code, _, response, err := suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{Deploy{"abc123", true, "2006.01.02-15.04.05", 1}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 201, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithDestroyPreviousAndTooManyInstancesRunning() {
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@2.service", []*schema.UnitOption{}},
	}, nil)

	code, _, _, err := suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{Deploy{"abc123", true, "2006.01.02-15.04.05", 1}},
	)

	assert.Contains(suite.T(), fmt.Sprintf("%s", err), "A greater number of instances")
	assert.Equal(suite.T(), 400, code)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithDestroyPreviousAndTooManyVersionsRunning() {
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "aabbccd", "carousel:aabbccd:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
	}, nil)

	code, _, _, err := suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{Deploy{"abc123", true, "2006.01.02-15.04.05", 1}},
	)

	assert.Contains(suite.T(), fmt.Sprintf("%s", err), "Too many versions")
	assert.Equal(suite.T(), 400, code)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroyPrevious() {
	mockedStates := []*schema.UnitState{
		&schema.UnitState{"", "", "carousel:cccddd:2006.01.02-15.04.05@1.service", "", "", "running"},
	}
	suite.FleetClientMock.On("UnitStates").Return(mockedStates, nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:abc123:2006.01.02-15.04.05@1.service").Return(nil)

	suite.Subject.destroyPrevious("carousel:abc123:2006.01.02-15.04.05@1.service", "carousel:cccddd:2006.01.02-15.04.05@1.service", 0)

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroySingleInstance() {
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}}}, nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@1.service").Return(nil)

	code, _, response, err := suite.Subject.Destroy(
		mocking.URL(suite.Service.RootMux, "DELETE", "http://example.com/v1/services/carousel/deploys/efefeff"),
		mocking.Header(nil),
		&DeployRequest{Deploy{"efefeff", false, "2006.01.02-15.04.05", 1}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 204, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroyMultipleInstances() {
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@2.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "3e33333", "carousel:3e33333:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
	}, nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@1.service").Return(nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@2.service").Return(nil)

	code, _, response, err := suite.Subject.Destroy(
		mocking.URL(suite.Service.RootMux, "DELETE", "http://example.com/v1/services/carousel/deploys/efefeff"),
		mocking.Header(nil),
		&DeployRequest{Deploy{"efefeff", false, "2006.01.02-15.04.05", 1}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 204, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroyWithUnmanagedUnits() {
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@2.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "3e33333", "carousel:3e33333:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "3e33333", "vulcand.service", []*schema.UnitOption{}},
	}, nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@1.service").Return(nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@2.service").Return(nil)

	code, _, response, err := suite.Subject.Destroy(
		mocking.URL(suite.Service.RootMux, "DELETE", "http://example.com/v1/services/carousel/deploys/efefeff"),
		mocking.Header(nil),
		&DeployRequest{Deploy{"efefeff", false, "2006.01.02-15.04.05", 1}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 204, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroyMultipleInstancesWithTimestampSpecified() {
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "efefeff", "carousel:efefeff:2012.01.02-15.04.05@2.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "3e33333", "carousel:3e33333:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
	}, nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@1.service").Return(nil)

	code, _, response, err := suite.Subject.Destroy(
		mocking.URL(suite.Service.RootMux, "DELETE", "http://example.com/v1/services/carousel/deploys/efefeff?timestamp=2006.01.02-15.04.05"),
		mocking.Header(nil),
		nil,
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 204, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func TestDeploysResourceTestSuite(t *testing.T) {
	suite.Run(t, new(DeploysResourceTestSuite))
}
