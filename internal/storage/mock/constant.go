// Code generated by MockGen. DO NOT EDIT.
// Source: constant.go
//
// Generated by this command:
//
//	mockgen -source=constant.go -destination=mock/constant.go -package=mock -typed
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/dipdup-io/celestia-indexer/internal/storage"
	types "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	gomock "go.uber.org/mock/gomock"
)

// MockIConstant is a mock of IConstant interface.
type MockIConstant struct {
	ctrl     *gomock.Controller
	recorder *MockIConstantMockRecorder
}

// MockIConstantMockRecorder is the mock recorder for MockIConstant.
type MockIConstantMockRecorder struct {
	mock *MockIConstant
}

// NewMockIConstant creates a new mock instance.
func NewMockIConstant(ctrl *gomock.Controller) *MockIConstant {
	mock := &MockIConstant{ctrl: ctrl}
	mock.recorder = &MockIConstantMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIConstant) EXPECT() *MockIConstantMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockIConstant) All(ctx context.Context) ([]storage.Constant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", ctx)
	ret0, _ := ret[0].([]storage.Constant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockIConstantMockRecorder) All(ctx any) *IConstantAllCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockIConstant)(nil).All), ctx)
	return &IConstantAllCall{Call: call}
}

// IConstantAllCall wrap *gomock.Call
type IConstantAllCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IConstantAllCall) Return(arg0 []storage.Constant, arg1 error) *IConstantAllCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IConstantAllCall) Do(f func(context.Context) ([]storage.Constant, error)) *IConstantAllCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IConstantAllCall) DoAndReturn(f func(context.Context) ([]storage.Constant, error)) *IConstantAllCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByModule mocks base method.
func (m *MockIConstant) ByModule(ctx context.Context, module types.ModuleName) ([]storage.Constant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByModule", ctx, module)
	ret0, _ := ret[0].([]storage.Constant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByModule indicates an expected call of ByModule.
func (mr *MockIConstantMockRecorder) ByModule(ctx, module any) *IConstantByModuleCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByModule", reflect.TypeOf((*MockIConstant)(nil).ByModule), ctx, module)
	return &IConstantByModuleCall{Call: call}
}

// IConstantByModuleCall wrap *gomock.Call
type IConstantByModuleCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IConstantByModuleCall) Return(arg0 []storage.Constant, arg1 error) *IConstantByModuleCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IConstantByModuleCall) Do(f func(context.Context, types.ModuleName) ([]storage.Constant, error)) *IConstantByModuleCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IConstantByModuleCall) DoAndReturn(f func(context.Context, types.ModuleName) ([]storage.Constant, error)) *IConstantByModuleCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Get mocks base method.
func (m *MockIConstant) Get(ctx context.Context, module types.ModuleName, name string) (storage.Constant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, module, name)
	ret0, _ := ret[0].(storage.Constant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIConstantMockRecorder) Get(ctx, module, name any) *IConstantGetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIConstant)(nil).Get), ctx, module, name)
	return &IConstantGetCall{Call: call}
}

// IConstantGetCall wrap *gomock.Call
type IConstantGetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IConstantGetCall) Return(arg0 storage.Constant, arg1 error) *IConstantGetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IConstantGetCall) Do(f func(context.Context, types.ModuleName, string) (storage.Constant, error)) *IConstantGetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IConstantGetCall) DoAndReturn(f func(context.Context, types.ModuleName, string) (storage.Constant, error)) *IConstantGetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
