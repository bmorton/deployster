package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExtractableUnitTestSuite struct {
	suite.Suite
}

func (suite *ExtractableUnitTestSuite) TestIsManaged() {
	subject := ExtractableUnit{Name: "railsapp:cf2e8ac:2006.01.02-15.04.05@1.service"}
	assert.True(suite.T(), subject.IsManaged())
}

func (suite *ExtractableUnitTestSuite) TestIsNotManaged() {
	subject := ExtractableUnit{Name: "vulcand.service"}
	assert.False(suite.T(), subject.IsManaged())
}

func (suite *ExtractableUnitTestSuite) TestExtractsBaseName() {
	subject := ExtractableUnit{Name: "railsapp:cf2e8ac:2006.01.02-15.04.05@1.service"}
	assert.Equal(suite.T(), "railsapp", subject.ExtractBaseName())
}

func (suite *ExtractableUnitTestSuite) TestExtractsVersion() {
	subject := ExtractableUnit{Name: "railsapp:cf2e8ac:2006.01.02-15.04.05@1.service"}
	assert.Equal(suite.T(), "cf2e8ac", subject.ExtractVersion())
}

func (suite *ExtractableUnitTestSuite) TestExtractsInstance() {
	subject := ExtractableUnit{Name: "railsapp:cf2e8ac:2006.01.02-15.04.05@1.service"}
	assert.Equal(suite.T(), "1", subject.ExtractInstance())
}

func (suite *ExtractableUnitTestSuite) TestExtractsTimestamp() {
	subject := ExtractableUnit{Name: "railsapp:cf2e8ac:2006.01.02-15.04.05@1.service"}
	assert.Equal(suite.T(), "2006.01.02-15.04.05", subject.ExtractTimestamp())
}

func TestExtractableUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ExtractableUnitTestSuite))
}
