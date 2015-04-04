package handlers

import (
	"testing"

	"github.com/bmorton/deployster/clients/mocks"
	"github.com/bmorton/deployster/poller"
	"github.com/bmorton/deployster/schema"
	"github.com/stretchr/testify/suite"
)

type DestroyerTestSuite struct {
	suite.Suite
	Subject   *Destroyer
	FleetMock *mocks.Fleet
}

func (suite *DestroyerTestSuite) SetupTest() {
	suite.FleetMock = new(mocks.Fleet)
	suite.Subject = &Destroyer{
		PreviousVersion: &schema.Deploy{ServiceName: "railsapp", Version: "old", Timestamp: "2006.01.02-15.04.05"},
		Client:          suite.FleetMock,
	}

}

func (suite *DestroyerTestSuite) TestDestroysUnit() {
	suite.FleetMock.On("DestroyUnit", "railsapp:old:2006.01.02-15.04.05@1.service").Return(nil).Times(1)

	suite.Subject.Handle(&poller.Event{ServiceInstance: &schema.ServiceInstance{Instance: "1"}})
	suite.FleetMock.Mock.AssertExpectations(suite.T())
}

func TestDestroyerTestSuite(t *testing.T) {
	suite.Run(t, new(DestroyerTestSuite))
}
