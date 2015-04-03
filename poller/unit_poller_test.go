package poller

import (
	"testing"
	"time"

	"github.com/bmorton/deployster/schema"
	"github.com/bmorton/deployster/server/mocks"
	fleet "github.com/coreos/fleet/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MockSuccessHandler struct {
	called bool
}

func (m *MockSuccessHandler) Handle(s *fleet.UnitState) {
	m.called = true
}

type UnitPollerTestSuite struct {
	suite.Suite
	Subject         *UnitPoller
	FleetClientMock *mocks.FleetClient
	ServiceInstance *schema.ServiceInstance
}

func (suite *UnitPollerTestSuite) SetupTest() {
	suite.FleetClientMock = new(mocks.FleetClient)
	suite.ServiceInstance = &schema.ServiceInstance{Name: "railsapp", Version: "latest", Instance: "1", Timestamp: "2006.01.02-15.04.05"}
	suite.Subject = New(suite.ServiceInstance, suite.FleetClientMock)
	suite.Subject.Timeout = 100 * time.Millisecond
	suite.Subject.Delay = 0
}

func (suite *UnitPollerTestSuite) TestSuccessHandlerCalledWhenStateRunning() {
	handler := &MockSuccessHandler{}
	suite.FleetClientMock.On("UnitStates").Return(suite.expectedForState("running"), nil)

	suite.Subject.AddSuccessHandler(handler)
	suite.Subject.Watch()

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
	assert.True(suite.T(), handler.called)
}

func (suite *UnitPollerTestSuite) TestSuccessHandlerNotCalledWhenStateFailed() {
	handler := &MockSuccessHandler{}
	suite.FleetClientMock.On("UnitStates").Return(suite.expectedForState("failed"), nil)

	suite.Subject.AddSuccessHandler(handler)
	suite.Subject.Watch()

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
	assert.False(suite.T(), handler.called)
}

func (suite *UnitPollerTestSuite) TestPollsAgainWhenStateUnresolved() {
	handler := &MockSuccessHandler{}
	suite.FleetClientMock.On("UnitStates").Return(suite.expectedForState("launching"), nil).Times(1)
	suite.FleetClientMock.On("UnitStates").Return(suite.expectedForState("running"), nil).Times(1)

	suite.Subject.AddSuccessHandler(handler)
	suite.Subject.Watch()

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func TestUnitsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(UnitPollerTestSuite))
}

func (suite *UnitPollerTestSuite) expectedForState(state string) []*fleet.UnitState {
	return []*fleet.UnitState{
		&fleet.UnitState{
			Name:            suite.ServiceInstance.FleetUnitName(),
			SystemdSubState: state,
		},
	}
}
