// Code generated by MockGen. DO NOT EDIT.
// Source: repo.go

// Package mock_repo is a generated GoMock package.
package repo

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLinkRepo is a mock of LinkRepo interface
type MockLinkRepo struct {
	ctrl     *gomock.Controller
	recorder *MockLinkRepoMockRecorder
}

// MockLinkRepoMockRecorder is the mock recorder for MockLinkRepo
type MockLinkRepoMockRecorder struct {
	mock *MockLinkRepo
}

// NewMockLinkRepo creates a new mock instance
func NewMockLinkRepo(ctrl *gomock.Controller) *MockLinkRepo {
	mock := &MockLinkRepo{ctrl: ctrl}
	mock.recorder = &MockLinkRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLinkRepo) EXPECT() *MockLinkRepoMockRecorder {
	return m.recorder
}

// SetLink mocks base method
func (m *MockLinkRepo) SetLink(ctx context.Context, url, code string, isCustom bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLink", ctx, url, code, isCustom)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLink indicates an expected call of SetLink
func (mr *MockLinkRepoMockRecorder) SetLink(ctx, url, code, isCustom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLink", reflect.TypeOf((*MockLinkRepo)(nil).SetLink), ctx, url, code, isCustom)
}

// GetLongLinkByCode mocks base method
func (m *MockLinkRepo) GetLongLinkByCode(ctx context.Context, code string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLongLinkByCode", ctx, code)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLongLinkByCode indicates an expected call of GetLongLinkByCode
func (mr *MockLinkRepoMockRecorder) GetLongLinkByCode(ctx, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLongLinkByCode", reflect.TypeOf((*MockLinkRepo)(nil).GetLongLinkByCode), ctx, code)
}

// GetCodeByLongLink mocks base method
func (m *MockLinkRepo) GetCodeByLongLink(ctx context.Context, url string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCodeByLongLink", ctx, url)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCodeByLongLink indicates an expected call of GetCodeByLongLink
func (mr *MockLinkRepoMockRecorder) GetCodeByLongLink(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCodeByLongLink", reflect.TypeOf((*MockLinkRepo)(nil).GetCodeByLongLink), ctx, url)
}

// GetNextSeq mocks base method
func (m *MockLinkRepo) GetNextSeq(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextSeq", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNextSeq indicates an expected call of GetNextSeq
func (mr *MockLinkRepoMockRecorder) GetNextSeq(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextSeq", reflect.TypeOf((*MockLinkRepo)(nil).GetNextSeq), ctx)
}

// IsCodeExists mocks base method
func (m *MockLinkRepo) IsCodeExists(ctx context.Context, code string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCodeExists", ctx, code)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsCodeExists indicates an expected call of IsCodeExists
func (mr *MockLinkRepoMockRecorder) IsCodeExists(ctx, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCodeExists", reflect.TypeOf((*MockLinkRepo)(nil).IsCodeExists), ctx, code)
}
