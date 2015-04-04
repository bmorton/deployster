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
	timesCalled int
}

func (m *MockSuccessHandler) Handle(e *Event) {
	m.timesCalled++
}

func (m *MockSuccessHandler) wasCalled() bool {
	return m.timesCalled > 0
}

type PollerTestSuite struct {
	suite.Suite
	Subject         *Poller
	FleetClientMock *mocks.FleetClient
	Deploy          *schema.Deploy
}

func (suite *PollerTestSuite) SetupTest() {
	suite.FleetClientMock = new(mocks.FleetClient)
	suite.Deploy = &schema.Deploy{ServiceName: "railsapp", Version: "latest", InstanceCount: 1, Timestamp: "2006.01.02-15.04.05"}
	suite.Subject = New(suite.Deploy, suite.FleetClientMock)
	suite.Subject.Timeout = 100 * time.Millisecond
	suite.Subject.Delay = 0
}

func (suite *PollerTestSuite) TestSuccessHandlerCalledWhenStateRunning() {
	handler := &MockSuccessHandler{}
	suite.FleetClientMock.On("UnitStates").Return(suite.expectedForState("running"), nil)

	suite.Subject.AddSuccessHandler(handler)
	suite.Subject.Watch()

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
	assert.True(suite.T(), handler.wasCalled())
}

func (suite *PollerTestSuite) TestSuccessHandlerNotCalledWhenStateFailed() {
	handler := &MockSuccessHandler{}
	suite.FleetClientMock.On("UnitStates").Return(suite.expectedForState("failed"), nil)

	suite.Subject.AddSuccessHandler(handler)
	suite.Subject.Watch()

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
	assert.False(suite.T(), handler.wasCalled())
}

func (suite *PollerTestSuite) TestPollsAgainWhenStateUnresolved() {
	handler := &MockSuccessHandler{}
	suite.FleetClientMock.On("UnitStates").Return(suite.expectedForState("launching"), nil).Times(1)
	suite.FleetClientMock.On("UnitStates").Return(suite.expectedForState("running"), nil).Times(1)

	suite.Subject.AddSuccessHandler(handler)
	suite.Subject.Watch()

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func (suite *PollerTestSuite) TestPollsMultipleInstances() {
	suite.Deploy.InstanceCount = 3
	suite.Subject = New(suite.Deploy, suite.FleetClientMock)
	suite.Subject.Timeout = 100 * time.Millisecond
	suite.Subject.Delay = 0

	bogus := suite.Deploy.ServiceInstance("1")
	bogus.Version = "older"

	states := make(map[string]string)
	states[bogus.FleetUnitName()] = "running"
	states[suite.Deploy.ServiceInstance("1").FleetUnitName()] = "running"
	states[suite.Deploy.ServiceInstance("2").FleetUnitName()] = "running"
	states[suite.Deploy.ServiceInstance("3").FleetUnitName()] = "launching"
	suite.FleetClientMock.On("UnitStates").Return(suite.expectedForStates(states), nil).Times(1)
	states[suite.Deploy.ServiceInstance("3").FleetUnitName()] = "failed"
	suite.FleetClientMock.On("UnitStates").Return(suite.expectedForStates(states), nil).Times(1)

	handler := &MockSuccessHandler{}
	suite.Subject.AddSuccessHandler(handler)
	suite.Subject.Watch()

	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
	assert.Equal(suite.T(), 2, handler.timesCalled)
}

func (suite *PollerTestSuite) expectedForState(state string) []*fleet.UnitState {
	states := make(map[string]string)
	states[suite.Deploy.ServiceInstance("1").FleetUnitName()] = state
	return suite.expectedForStates(states)
}

func (suite *PollerTestSuite) expectedForStates(states map[string]string) []*fleet.UnitState {
	var generated []*fleet.UnitState
	for unit, state := range states {
		generated = append(generated,
			&fleet.UnitState{
				Name:            unit,
				SystemdSubState: state,
			},
		)
	}
	return generated
}

func TestPollerTestSuite(t *testing.T) {
	suite.Run(t, new(PollerTestSuite))
}
