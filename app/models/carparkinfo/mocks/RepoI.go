// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import carparkinfo "github.com/dare-rider/carpark/app/models/carparkinfo"
import mock "github.com/stretchr/testify/mock"
import sqlx "github.com/jmoiron/sqlx"

// RepoI is an autogenerated mock type for the RepoI type
type RepoI struct {
	mock.Mock
}

// FindAllByCarparkIDs provides a mock function with given fields: cpIds
func (_m *RepoI) FindAllByCarparkIDs(cpIds []int) ([]carparkinfo.Model, error) {
	ret := _m.Called(cpIds)

	var r0 []carparkinfo.Model
	if rf, ok := ret.Get(0).(func([]int) []carparkinfo.Model); ok {
		r0 = rf(cpIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]carparkinfo.Model)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]int) error); ok {
		r1 = rf(cpIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertOrUpdateByCarParkNo provides a mock function with given fields: mod, tx
func (_m *RepoI) InsertOrUpdateByCarParkNo(mod *carparkinfo.Model, tx ...*sqlx.Tx) error {
	_va := make([]interface{}, len(tx))
	for _i := range tx {
		_va[_i] = tx[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, mod)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*carparkinfo.Model, ...*sqlx.Tx) error); ok {
		r0 = rf(mod, tx...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
