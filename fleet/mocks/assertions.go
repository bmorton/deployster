package mocks

import "github.com/bmorton/deployster/fleet"

// Make sure our mocked client complies to fleet's interface
var _ fleet.Client = &Client{}
