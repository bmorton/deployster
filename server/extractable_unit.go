package server

import (
	"strings"

	"github.com/bmorton/deployster/fleet"
)

type ExtractableUnit fleet.Unit

func (eu *ExtractableUnit) ExtractBaseName() string {
	s := strings.Split(eu.Name, "-")
	return s[0]
}

func (eu *ExtractableUnit) ExtractVersion() string {
	s := strings.Split(eu.Name, "-")
	end := strings.Index(s[1], "@")
	return s[1][:end]
}

func (eu *ExtractableUnit) ExtractInstance() string {
	s := strings.Split(eu.Name, "@")
	end := strings.Index(s[1], ".")
	return s[1][:end]
}
