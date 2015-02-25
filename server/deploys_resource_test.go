package server

import (
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

func (suite *DeploysResourceTestSuite) TestCreateWithoutDestroyPrevious() {
	expectedOptions := getUnitOptions("carousel", "abc123", "mmmhm")
	suite.FleetClientMock.On("CreateUnit", &schema.Unit{Name: "carousel-abc123@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel-abc123@1.service", "launched").Return(nil)

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
	expectedOptions := getUnitOptions("carousel", "abc123", "mmmhm")
	suite.FleetClientMock.On("Units").Return([]*schema.Unit{}, nil)
	suite.FleetClientMock.On("CreateUnit", &schema.Unit{Name: "carousel-abc123@1.service", Options: expectedOptions}).Return(nil)
	suite.FleetClientMock.On("SetUnitTargetState", "carousel-abc123@1.service", "launched").Return(nil)

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
	mockedStates := []*schema.UnitState{
		&schema.UnitState{"", "", "carousel-abc123@1.service", "", "", "running"},
	}
	suite.FleetClientMock.On("UnitStates").Return(mockedStates, nil)
	suite.FleetClientMock.On("DestroyUnit", "carousel-cccddd@1.service").Return(nil)

	suite.Subject.destroyPrevious("carousel", "cccddd", "abc123", 0)

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *DeploysResourceTestSuite) TestDestroy() {
	suite.FleetClientMock.On("DestroyUnit", "carousel-cccddd@1.service").Return(nil)

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
