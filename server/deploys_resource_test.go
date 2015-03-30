package server

import (
	"fmt"
	"testing"

	"github.com/bmorton/deployster/schema"
	"github.com/bmorton/deployster/server/mocks"
	fleet "github.com/coreos/fleet/schema"
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
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{}, nil)

	// Should only start 1 unit
	suite.FleetClientMock.On("CreateUnit", &fleet.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@1.service", "launched").Return(nil)

	suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "abc123", DestroyPrevious: false, Timestamp: "2006.01.02-15.04.05"}},
	)

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithoutPassedInstancesAndMultipleInstancesRunning() {
	expectedOptions := getUnitOptions(UnitTemplate{"carousel", "abc123", "mmmhm", "2006.01.02-15.04.05"})
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@2.service", []*fleet.UnitOption{}},
	}, nil)

	// Should start 2 units
	suite.FleetClientMock.On("CreateUnit", &fleet.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@1.service", "launched").Return(nil)
	suite.FleetClientMock.On("CreateUnit", &fleet.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@2.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@2.service", "launched").Return(nil)

	suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "abc123", DestroyPrevious: false, Timestamp: "2006.01.02-15.04.05"}},
	)

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithoutPassedInstancesAndFailedInstances() {
	expectedOptions := getUnitOptions(UnitTemplate{"carousel", "abc123", "mmmhm", "2008.01.02-15.04.05"})
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
		&fleet.Unit{"failed", "failed", "efefeff", "carousel:efefeff:2007.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
	}, nil)

	// Should start 1 unit
	suite.FleetClientMock.On("CreateUnit", &fleet.Unit{Name: "carousel:abc123:2008.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2008.01.02-15.04.05@1.service", "launched").Return(nil)

	suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "abc123", DestroyPrevious: false, Timestamp: "2008.01.02-15.04.05"}},
	)

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithoutPassedInstancesAndMultipleVersionsRunning() {
	expectedOptions := getUnitOptions(UnitTemplate{"carousel", "abc123", "mmmhm", "2006.01.02-15.04.05"})
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "abbbbbb", "carousel:abbbbbb:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
	}, nil)

	// Should start 1 unit
	suite.FleetClientMock.On("CreateUnit", &fleet.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@1.service", "launched").Return(nil)

	suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "abc123", DestroyPrevious: false, Timestamp: "2006.01.02-15.04.05"}},
	)

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithoutDestroyPrevious() {
	expectedOptions := getUnitOptions(UnitTemplate{"carousel", "abc123", "mmmhm", "2006.01.02-15.04.05"})
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{}, nil)
	suite.FleetClientMock.On("CreateUnit", &fleet.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@1.service", "launched").Return(nil)

	code, _, response, err := suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "abc123", DestroyPrevious: false, Timestamp: "2006.01.02-15.04.05", InstanceCount: 1}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 201, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithDestroyPreviousAndNoPreviousVersions() {
	expectedOptions := getUnitOptions(UnitTemplate{"carousel", "abc123", "mmmhm", "2006.01.02-15.04.05"})
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{}, nil)
	suite.FleetClientMock.On("CreateUnit", &fleet.Unit{Name: "carousel:abc123:2006.01.02-15.04.05@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel:abc123:2006.01.02-15.04.05@1.service", "launched").Return(nil)

	code, _, response, err := suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "abc123", DestroyPrevious: true, Timestamp: "2006.01.02-15.04.05", InstanceCount: 1}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 201, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithDestroyPreviousAndTooManyInstancesRunning() {
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@2.service", []*fleet.UnitOption{}},
	}, nil)

	code, _, _, err := suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "abc123", DestroyPrevious: true, Timestamp: "2006.01.02-15.04.05", InstanceCount: 1}},
	)

	assert.Contains(suite.T(), fmt.Sprintf("%s", err), "A greater number of instances")
	assert.Equal(suite.T(), 400, code)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestCreateWithDestroyPreviousAndTooManyVersionsRunning() {
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "aabbccd", "carousel:aabbccd:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
	}, nil)

	code, _, _, err := suite.Subject.Create(
		mocking.URL(suite.Service.RootMux, "POST", "http://example.com/v1/services/carousel/deploys"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "abc123", DestroyPrevious: true, Timestamp: "2006.01.02-15.04.05", InstanceCount: 1}},
	)

	assert.Contains(suite.T(), fmt.Sprintf("%s", err), "Too many versions")
	assert.Equal(suite.T(), 400, code)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroySingleInstance() {
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}}}, nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@1.service").Return(nil)

	code, _, response, err := suite.Subject.Destroy(
		mocking.URL(suite.Service.RootMux, "DELETE", "http://example.com/v1/services/carousel/deploys/efefeff"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "efefeff", DestroyPrevious: false, Timestamp: "2006.01.02-15.04.05", InstanceCount: 1}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 204, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroyMultipleInstances() {
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@2.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "3e33333", "carousel:3e33333:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
	}, nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@1.service").Return(nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@2.service").Return(nil)

	code, _, response, err := suite.Subject.Destroy(
		mocking.URL(suite.Service.RootMux, "DELETE", "http://example.com/v1/services/carousel/deploys/efefeff"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "efefeff", DestroyPrevious: false, Timestamp: "2006.01.02-15.04.05", InstanceCount: 1}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 204, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroyWithUnmanagedUnits() {
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@2.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "3e33333", "carousel:3e33333:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "3e33333", "vulcand.service", []*fleet.UnitOption{}},
	}, nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@1.service").Return(nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel:efefeff:2006.01.02-15.04.05@2.service").Return(nil)

	code, _, response, err := suite.Subject.Destroy(
		mocking.URL(suite.Service.RootMux, "DELETE", "http://example.com/v1/services/carousel/deploys/efefeff"),
		mocking.Header(nil),
		&DeployRequest{&schema.Deploy{Version: "efefeff", DestroyPrevious: false, Timestamp: "2006.01.02-15.04.05", InstanceCount: 1}},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 204, code)
	assert.Equal(suite.T(), nil, response)
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroyMultipleInstancesWithTimestampSpecified() {
	suite.FleetClientMock.On("Units").Return([]*fleet.Unit{
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "efefeff", "carousel:efefeff:2012.01.02-15.04.05@2.service", []*fleet.UnitOption{}},
		&fleet.Unit{"running", "running", "3e33333", "carousel:3e33333:2006.01.02-15.04.05@1.service", []*fleet.UnitOption{}},
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
