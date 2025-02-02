// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/vmware-tanzu/tanzu-framework/tkg/avi (interfaces: MiniServiceEngineGroupClient)

// Package avi is a generated GoMock package.
package avi

import (
	reflect "reflect"

	models "github.com/avinetworks/sdk/go/models"
	session "github.com/avinetworks/sdk/go/session"
	gomock "github.com/golang/mock/gomock"
)

// MockMiniServiceEngineGroupClient is a mock of MiniServiceEngineGroupClient interface
type MockMiniServiceEngineGroupClient struct {
	ctrl     *gomock.Controller
	recorder *MockMiniServiceEngineGroupClientMockRecorder
}

// MockMiniServiceEngineGroupClientMockRecorder is the mock recorder for MockMiniServiceEngineGroupClient
type MockMiniServiceEngineGroupClientMockRecorder struct {
	mock *MockMiniServiceEngineGroupClient
}

// NewMockMiniServiceEngineGroupClient creates a new mock instance
func NewMockMiniServiceEngineGroupClient(ctrl *gomock.Controller) *MockMiniServiceEngineGroupClient {
	mock := &MockMiniServiceEngineGroupClient{ctrl: ctrl}
	mock.recorder = &MockMiniServiceEngineGroupClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMiniServiceEngineGroupClient) EXPECT() *MockMiniServiceEngineGroupClientMockRecorder {
	return m.recorder
}

// GetAll mocks base method
func (m *MockMiniServiceEngineGroupClient) GetAll(arg0 ...session.ApiOptionsParams) ([]*models.ServiceEngineGroup, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAll", varargs...)
	ret0, _ := ret[0].([]*models.ServiceEngineGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockMiniServiceEngineGroupClientMockRecorder) GetAll(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMiniServiceEngineGroupClient)(nil).GetAll), arg0...)
}

// GetByName mocks base method
func (m *MockMiniServiceEngineGroupClient) GetByName(name string, options ...session.ApiOptionsParams) (*models.ServiceEngineGroup, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{name}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByName", varargs...)
	ret0, _ := ret[0].(*models.ServiceEngineGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName
func (mr *MockMiniServiceEngineGroupClientMockRecorder) GetByName(name interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockMiniServiceEngineGroupClient)(nil).GetByName), varargs...)
}
