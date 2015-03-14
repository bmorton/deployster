package templates

import (
	"github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"
)

// UnitTemplate is the view model that is passed to the template parser that
// renders a unit file.
type Unit struct {
	Name        string
	Version     string
	ImagePrefix string
	Timestamp   string
}

// Options renders the unit file and converts it to an array of
// UnitOption structs.
func (u *Unit) Options(template *Template) []*schema.UnitOption {
	unitTemplate, _ := template.Generate(u)
	unitFile, _ := unit.NewUnitFile(unitTemplate)

	return schema.MapUnitFileToSchemaUnitOptions(unitFile)
}
