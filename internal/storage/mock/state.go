// Code generated by MockGen. DO NOT EDIT.
// Source: state.go
//
// Generated by this command:
//
//	mockgen -source=state.go -destination=mock/state.go -package=mock -typed
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/dipdup-io/celestia-indexer/internal/storage"
	storage0 "github.com/dipdup-net/indexer-sdk/pkg/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIState is a mock of IState interface.
type MockIState struct {
	ctrl     *gomock.Controller
	recorder *MockIStateMockRecorder
}

// MockIStateMockRecorder is the mock recorder for MockIState.
type MockIStateMockRecorder struct {
	mock *MockIState
}

// NewMockIState creates a new mock instance.
func NewMockIState(ctrl *gomock.Controller) *MockIState {
	mock := &MockIState{ctrl: ctrl}
	mock.recorder = &MockIStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIState) EXPECT() *MockIStateMockRecorder {
	return m.recorder
}

// ByName mocks base method.
func (m *MockIState) ByName(ctx context.Context, name string) (storage.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByName", ctx, name)
	ret0, _ := ret[0].(storage.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByName indicates an expected call of ByName.
func (mr *MockIStateMockRecorder) ByName(ctx, name any) *IStateByNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByName", reflect.TypeOf((*MockIState)(nil).ByName), ctx, name)
	return &IStateByNameCall{Call: call}
}

// IStateByNameCall wrap *gomock.Call
type IStateByNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStateByNameCall) Return(arg0 storage.State, arg1 error) *IStateByNameCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStateByNameCall) Do(f func(context.Context, string) (storage.State, error)) *IStateByNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStateByNameCall) DoAndReturn(f func(context.Context, string) (storage.State, error)) *IStateByNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIState) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIStateMockRecorder) CursorList(ctx, id, limit, order, cmp any) *IStateCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIState)(nil).CursorList), ctx, id, limit, order, cmp)
	return &IStateCursorListCall{Call: call}
}

// IStateCursorListCall wrap *gomock.Call
type IStateCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStateCursorListCall) Return(arg0 []*storage.State, arg1 error) *IStateCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStateCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.State, error)) *IStateCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStateCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.State, error)) *IStateCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIState) GetByID(ctx context.Context, id uint64) (*storage.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIStateMockRecorder) GetByID(ctx, id any) *IStateGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIState)(nil).GetByID), ctx, id)
	return &IStateGetByIDCall{Call: call}
}

// IStateGetByIDCall wrap *gomock.Call
type IStateGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStateGetByIDCall) Return(arg0 *storage.State, arg1 error) *IStateGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStateGetByIDCall) Do(f func(context.Context, uint64) (*storage.State, error)) *IStateGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStateGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.State, error)) *IStateGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIState) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIStateMockRecorder) IsNoRows(err any) *IStateIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIState)(nil).IsNoRows), err)
	return &IStateIsNoRowsCall{Call: call}
}

// IStateIsNoRowsCall wrap *gomock.Call
type IStateIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStateIsNoRowsCall) Return(arg0 bool) *IStateIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStateIsNoRowsCall) Do(f func(error) bool) *IStateIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStateIsNoRowsCall) DoAndReturn(f func(error) bool) *IStateIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIState) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIStateMockRecorder) LastID(ctx any) *IStateLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIState)(nil).LastID), ctx)
	return &IStateLastIDCall{Call: call}
}

// IStateLastIDCall wrap *gomock.Call
type IStateLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStateLastIDCall) Return(arg0 uint64, arg1 error) *IStateLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStateLastIDCall) Do(f func(context.Context) (uint64, error)) *IStateLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStateLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *IStateLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIState) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIStateMockRecorder) List(ctx, limit, offset, order any) *IStateListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIState)(nil).List), ctx, limit, offset, order)
	return &IStateListCall{Call: call}
}

// IStateListCall wrap *gomock.Call
type IStateListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStateListCall) Return(arg0 []*storage.State, arg1 error) *IStateListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStateListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.State, error)) *IStateListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStateListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.State, error)) *IStateListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIState) Save(ctx context.Context, m *storage.State) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIStateMockRecorder) Save(ctx, m any) *IStateSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIState)(nil).Save), ctx, m)
	return &IStateSaveCall{Call: call}
}

// IStateSaveCall wrap *gomock.Call
type IStateSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStateSaveCall) Return(arg0 error) *IStateSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStateSaveCall) Do(f func(context.Context, *storage.State) error) *IStateSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStateSaveCall) DoAndReturn(f func(context.Context, *storage.State) error) *IStateSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIState) Update(ctx context.Context, m *storage.State) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIStateMockRecorder) Update(ctx, m any) *IStateUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIState)(nil).Update), ctx, m)
	return &IStateUpdateCall{Call: call}
}

// IStateUpdateCall wrap *gomock.Call
type IStateUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStateUpdateCall) Return(arg0 error) *IStateUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStateUpdateCall) Do(f func(context.Context, *storage.State) error) *IStateUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStateUpdateCall) DoAndReturn(f func(context.Context, *storage.State) error) *IStateUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}