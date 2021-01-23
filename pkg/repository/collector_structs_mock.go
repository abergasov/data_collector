// Code generated by MockGen. DO NOT EDIT.
// Source: collector_structs.go

// Package repository is a generated GoMock package.
package repository

import (
	sql "database/sql"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIDatabase is a mock of IDatabase interface
type MockIDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockIDatabaseMockRecorder
}

// MockIDatabaseMockRecorder is the mock recorder for MockIDatabase
type MockIDatabaseMockRecorder struct {
	mock *MockIDatabase
}

// NewMockIDatabase creates a new mock instance
func NewMockIDatabase(ctrl *gomock.Controller) *MockIDatabase {
	mock := &MockIDatabase{ctrl: ctrl}
	mock.recorder = &MockIDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIDatabase) EXPECT() *MockIDatabaseMockRecorder {
	return m.recorder
}

// Select mocks base method
func (m *MockIDatabase) Select(dest interface{}, query string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Select", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Select indicates an expected call of Select
func (mr *MockIDatabaseMockRecorder) Select(dest, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockIDatabase)(nil).Select), varargs...)
}

// Exec mocks base method
func (m *MockIDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec
func (mr *MockIDatabaseMockRecorder) Exec(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockIDatabase)(nil).Exec), varargs...)
}
