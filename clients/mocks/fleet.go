package mocks

import "github.com/stretchr/testify/mock"

import fleet "github.com/coreos/fleet/schema"

type Fleet struct {
	mock.Mock
}

func (m *Fleet) Units() ([]*fleet.Unit, error) {
	ret := m.Called()

	var r0 []*fleet.Unit
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*fleet.Unit)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Fleet) CreateUnit(_a0 *fleet.Unit) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *Fleet) DestroyUnit(_a0 string) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *Fleet) UnitStates() ([]*fleet.UnitState, error) {
	ret := m.Called()

	var r0 []*fleet.UnitState
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*fleet.UnitState)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Fleet) SetUnitTargetState(_a0 string, _a1 string) error {
	ret := m.Called(_a0, _a1)

	r0 := ret.Error(0)

	return r0
}
