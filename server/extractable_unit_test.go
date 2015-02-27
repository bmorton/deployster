package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExtractableUnitTestSuite struct {
	suite.Suite
}

func (suite *ExtractableUnitTestSuite) TestExtractsBaseName() {
	subject := ExtractableUnit{Name: "railsapp:cf2e8ac_2013-06-05T14:10:43Z@1.service"}
	assert.Equal(suite.T(), "railsapp", subject.ExtractBaseName())
}

func (suite *ExtractableUnitTestSuite) TestExtractsVersion() {
	subject := ExtractableUnit{Name: "railsapp:cf2e8ac_2013-06-05T14:10:43Z@1.service"}
	assert.Equal(suite.T(), "cf2e8ac", subject.ExtractVersion())
}

func (suite *ExtractableUnitTestSuite) TestExtractsInstance() {
	subject := ExtractableUnit{Name: "railsapp:cf2e8ac_2013-06-05T14:10:43Z@1.service"}
	assert.Equal(suite.T(), "1", subject.ExtractInstance())
}

func (suite *ExtractableUnitTestSuite) TestExtractsTimestamp() {
	subject := ExtractableUnit{Name: "railsapp:cf2e8ac_2013-06-05T14:10:43Z@1.service"}
	assert.Equal(suite.T(), "2013-06-05T14:10:43Z", subject.ExtractTimestamp())
}

func TestExtractableUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ExtractableUnitTestSuite))
}
