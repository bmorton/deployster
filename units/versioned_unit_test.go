package units

import (
	"testing"

	"github.com/coreos/fleet/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VersionedUnitTestSuite struct {
	suite.Suite
}

func (suite *VersionedUnitTestSuite) TestFindServiceUnitsIsExtracted() {
	units := []*schema.Unit{
		&schema.Unit{"running", "running", "m4ch1n3-1d", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
	}

	found := FindServiceUnits("carousel", "", units)
	expected := VersionedUnit{
		Service:      "carousel",
		Instance:     "1",
		Version:      "efefeff",
		Timestamp:    "2006.01.02-15.04.05",
		CurrentState: "running",
		DesiredState: "running",
		MachineID:    "m4ch1n3-1d",
	}
	assert.Equal(suite.T(), expected, found[0])
}

func (suite *VersionedUnitTestSuite) TestFindServiceUnitsByNameOnly() {
	units := []*schema.Unit{
		&schema.Unit{"running", "running", "m4ch1n3-1d", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "m4ch1n3-1d", "notcarousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
	}

	found := FindServiceUnits("carousel", "", units)
	assert.Len(suite.T(), found, 1)
	assert.Equal(suite.T(), "carousel", found[0].Service)
}

func (suite *VersionedUnitTestSuite) TestFindServiceUnitsByNameAndVersion() {
	units := []*schema.Unit{
		&schema.Unit{"running", "running", "m4ch1n3-1d", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "m4ch1n3-1d", "carousel:abababa:2007.01.02-15.04.05@1.service", []*schema.UnitOption{}},
	}

	found := FindServiceUnits("carousel", "efefeff", units)
	assert.Len(suite.T(), found, 1)
	assert.Equal(suite.T(), "efefeff", found[0].Version)
}

func (suite *VersionedUnitTestSuite) TestFindTimestampedServiceVersionsWithOneVersion() {
	units := []*schema.Unit{
		&schema.Unit{"running", "running", "m4ch1n3-1d", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
	}

	found := FindTimestampedServiceVersions("carousel", units)
	expected := []string{"efefeff:2006.01.02-15.04.05"}
	assert.Equal(suite.T(), expected, found)
}

func (suite *VersionedUnitTestSuite) TestFindTimestampedServiceVersionsWithOneVersionAtTwoTimestamps() {
	units := []*schema.Unit{
		&schema.Unit{"running", "running", "m4ch1n3-1d", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "m4ch1n3-1d", "carousel:efefeff:2007.01.02-15.04.05@1.service", []*schema.UnitOption{}},
	}

	found := FindTimestampedServiceVersions("carousel", units)
	expected := []string{"efefeff:2006.01.02-15.04.05", "efefeff:2007.01.02-15.04.05"}
	assert.Contains(suite.T(), found, expected[0], expected[1])
}

func (suite *VersionedUnitTestSuite) TestFindTimestampedServiceVersionsWithMultipleVersions() {
	units := []*schema.Unit{
		&schema.Unit{"running", "running", "m4ch1n3-1d", "carousel:efefeff:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
		&schema.Unit{"running", "running", "m4ch1n3-1d", "carousel:abababb:2006.01.02-15.04.05@1.service", []*schema.UnitOption{}},
	}

	found := FindTimestampedServiceVersions("carousel", units)
	expected := []string{"efefeff:2006.01.02-15.04.05", "abababb:2006.01.02-15.04.05"}
	assert.Contains(suite.T(), found, expected[0], expected[1])
}

func TestVersionedUnitTestSuite(t *testing.T) {
	suite.Run(t, new(VersionedUnitTestSuite))
}
