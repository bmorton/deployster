package handlers

import (
	"testing"

	"github.com/bmorton/deployster/poller"
	"github.com/bmorton/deployster/schema"
	"github.com/bmorton/deployster/server/mocks"
	"github.com/stretchr/testify/suite"
)

type DestroyerTestSuite struct {
	suite.Suite
	Subject         *Destroyer
	FleetClientMock *mocks.FleetClient
}

func (suite *DestroyerTestSuite) SetupTest() {
	suite.FleetClientMock = new(mocks.FleetClient)
	suite.Subject = &Destroyer{
		PreviousVersion: &schema.Deploy{ServiceName: "railsapp", Version: "old", Timestamp: "2006.01.02-15.04.05"},
		Client:          suite.FleetClientMock,
	}

}

func (suite *DestroyerTestSuite) TestDestroysUnit() {
	suite.FleetClientMock.On("DestroyUnit", "railsapp:old:2006.01.02-15.04.05@1.service").Return(nil).Times(1)

	suite.Subject.Handle(&poller.Event{ServiceInstance: &schema.ServiceInstance{Instance: "1"}})
	suite.FleetClientMock.Mock.AssertExpectations(suite.T())
}

func TestDestroyerTestSuite(t *testing.T) {
	suite.Run(t, new(DestroyerTestSuite))
}
