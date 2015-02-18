package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ExtractableUnitTestSuite struct {
	suite.Suite
}

func (suite *ExtractableUnitTestSuite) TestExtractsBaseName() {
	subject := ExtractableUnit{Name: "railsapp-efe1abc@1.service"}
	assert.Equal(suite.T(), "railsapp", subject.ExtractBaseName())
}

func (suite *ExtractableUnitTestSuite) TestExtractsBaseNameSupportsNamesWithHyphens() {
	subject := ExtractableUnit{Name: "hello-world-efe1abc@1.service"}
	assert.Equal(suite.T(), "hello-world", subject.ExtractBaseName())
}

func (suite *ExtractableUnitTestSuite) TestExtractsVersion() {
	subject := ExtractableUnit{Name: "railsapp-efe1abc@1.service"}
	assert.Equal(suite.T(), "efe1abc", subject.ExtractVersion())
}

func (suite *ExtractableUnitTestSuite) TestExtractsInstance() {
	subject := ExtractableUnit{Name: "railsapp-efe1abc@1.service"}
	assert.Equal(suite.T(), "1", subject.ExtractInstance())
}

func TestExtractableUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ExtractableUnitTestSuite))
}
