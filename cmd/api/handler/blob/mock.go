// Code generated by MockGen. DO NOT EDIT.
// Source: receiver.go
//
// Generated by this command:
//
//	mockgen -source=receiver.go -destination=mock.go -package=blob -typed
//
// Package blob is a generated GoMock package.
package blob

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockReceiver is a mock of Receiver interface.
type MockReceiver struct {
	ctrl     *gomock.Controller
	recorder *MockReceiverMockRecorder
}

// MockReceiverMockRecorder is the mock recorder for MockReceiver.
type MockReceiverMockRecorder struct {
	mock *MockReceiver
}

// NewMockReceiver creates a new mock instance.
func NewMockReceiver(ctrl *gomock.Controller) *MockReceiver {
	mock := &MockReceiver{ctrl: ctrl}
	mock.recorder = &MockReceiverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReceiver) EXPECT() *MockReceiverMockRecorder {
	return m.recorder
}

// Blobs mocks base method.
func (m *MockReceiver) Blobs(ctx context.Context, height uint64, hash ...string) ([]Blob, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, height}
	for _, a := range hash {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Blobs", varargs...)
	ret0, _ := ret[0].([]Blob)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Blobs indicates an expected call of Blobs.
func (mr *MockReceiverMockRecorder) Blobs(ctx, height any, hash ...any) *ReceiverBlobsCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, height}, hash...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Blobs", reflect.TypeOf((*MockReceiver)(nil).Blobs), varargs...)
	return &ReceiverBlobsCall{Call: call}
}

// ReceiverBlobsCall wrap *gomock.Call
type ReceiverBlobsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ReceiverBlobsCall) Return(arg0 []Blob, arg1 error) *ReceiverBlobsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ReceiverBlobsCall) Do(f func(context.Context, uint64, ...string) ([]Blob, error)) *ReceiverBlobsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ReceiverBlobsCall) DoAndReturn(f func(context.Context, uint64, ...string) ([]Blob, error)) *ReceiverBlobsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
