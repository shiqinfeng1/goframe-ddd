// Code generated by MockGen. DO NOT EDIT.
// Source: ./interfaces.go
//
// Generated by this command:
//
//	mockgen -destination=mock_client.go -package=nats -source=./interfaces.go ConnIntf,ConnMgr,SubMgr
//

// Package nats is a generated GoMock package.
package nats

import (
	context "context"
	reflect "reflect"

	nats "github.com/nats-io/nats.go"
	jetstream "github.com/nats-io/nats.go/jetstream"
	health "github.com/shiqinfeng1/goframe-ddd/pkg/health"
	pubsub "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
	gomock "go.uber.org/mock/gomock"
)

// MockConnIntf is a mock of ConnIntf interface.
type MockConnIntf struct {
	ctrl     *gomock.Controller
	recorder *MockConnIntfMockRecorder
	isgomock struct{}
}

// MockConnIntfMockRecorder is the mock recorder for MockConnIntf.
type MockConnIntfMockRecorder struct {
	mock *MockConnIntf
}

// NewMockConnIntf creates a new mock instance.
func NewMockConnIntf(ctrl *gomock.Controller) *MockConnIntf {
	mock := &MockConnIntf{ctrl: ctrl}
	mock.recorder = &MockConnIntfMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnIntf) EXPECT() *MockConnIntfMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConnIntf) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockConnIntfMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConnIntf)(nil).Close))
}

// Conn mocks base method.
func (m *MockConnIntf) Conn() *nats.Conn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Conn")
	ret0, _ := ret[0].(*nats.Conn)
	return ret0
}

// Conn indicates an expected call of Conn.
func (mr *MockConnIntfMockRecorder) Conn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Conn", reflect.TypeOf((*MockConnIntf)(nil).Conn))
}

// NewJetStream mocks base method.
func (m *MockConnIntf) NewJetStream() (jetstream.JetStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewJetStream")
	ret0, _ := ret[0].(jetstream.JetStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewJetStream indicates an expected call of NewJetStream.
func (mr *MockConnIntfMockRecorder) NewJetStream() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewJetStream", reflect.TypeOf((*MockConnIntf)(nil).NewJetStream))
}

// Status mocks base method.
func (m *MockConnIntf) Status() nats.Status {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(nats.Status)
	return ret0
}

// Status indicates an expected call of Status.
func (mr *MockConnIntfMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockConnIntf)(nil).Status))
}

// MockConnector is a mock of Connector interface.
type MockConnector struct {
	ctrl     *gomock.Controller
	recorder *MockConnectorMockRecorder
	isgomock struct{}
}

// MockConnectorMockRecorder is the mock recorder for MockConnector.
type MockConnectorMockRecorder struct {
	mock *MockConnector
}

// NewMockConnector creates a new mock instance.
func NewMockConnector(ctrl *gomock.Controller) *MockConnector {
	mock := &MockConnector{ctrl: ctrl}
	mock.recorder = &MockConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnector) EXPECT() *MockConnectorMockRecorder {
	return m.recorder
}

// Connect mocks base method.
func (m *MockConnector) Connect(arg0 string, arg1 ...nats.Option) (ConnIntf, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Connect", varargs...)
	ret0, _ := ret[0].(ConnIntf)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connect indicates an expected call of Connect.
func (mr *MockConnectorMockRecorder) Connect(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockConnector)(nil).Connect), varargs...)
}

// MockJetStreamCreator is a mock of JetStreamCreator interface.
type MockJetStreamCreator struct {
	ctrl     *gomock.Controller
	recorder *MockJetStreamCreatorMockRecorder
	isgomock struct{}
}

// MockJetStreamCreatorMockRecorder is the mock recorder for MockJetStreamCreator.
type MockJetStreamCreatorMockRecorder struct {
	mock *MockJetStreamCreator
}

// NewMockJetStreamCreator creates a new mock instance.
func NewMockJetStreamCreator(ctrl *gomock.Controller) *MockJetStreamCreator {
	mock := &MockJetStreamCreator{ctrl: ctrl}
	mock.recorder = &MockJetStreamCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJetStreamCreator) EXPECT() *MockJetStreamCreatorMockRecorder {
	return m.recorder
}

// New mocks base method.
func (m *MockJetStreamCreator) New(conn ConnIntf) (jetstream.JetStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", conn)
	ret0, _ := ret[0].(jetstream.JetStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// New indicates an expected call of New.
func (mr *MockJetStreamCreatorMockRecorder) New(conn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockJetStreamCreator)(nil).New), conn)
}

// MockConnMgr is a mock of ConnMgr interface.
type MockConnMgr struct {
	ctrl     *gomock.Controller
	recorder *MockConnMgrMockRecorder
	isgomock struct{}
}

// MockConnMgrMockRecorder is the mock recorder for MockConnMgr.
type MockConnMgrMockRecorder struct {
	mock *MockConnMgr
}

// NewMockConnMgr creates a new mock instance.
func NewMockConnMgr(ctrl *gomock.Controller) *MockConnMgr {
	mock := &MockConnMgr{ctrl: ctrl}
	mock.recorder = &MockConnMgrMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnMgr) EXPECT() *MockConnMgrMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConnMgr) Close(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close", ctx)
}

// Close indicates an expected call of Close.
func (mr *MockConnMgrMockRecorder) Close(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConnMgr)(nil).Close), ctx)
}

// Connect mocks base method.
func (m *MockConnMgr) Connect(ctx context.Context, opts ...nats.Option) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Connect", varargs...)
}

// Connect indicates an expected call of Connect.
func (mr *MockConnMgrMockRecorder) Connect(ctx any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockConnMgr)(nil).Connect), varargs...)
}

// GetJetStream mocks base method.
func (m *MockConnMgr) GetJetStream() (jetstream.JetStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJetStream")
	ret0, _ := ret[0].(jetstream.JetStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJetStream indicates an expected call of GetJetStream.
func (mr *MockConnMgrMockRecorder) GetJetStream() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJetStream", reflect.TypeOf((*MockConnMgr)(nil).GetJetStream))
}

// Health mocks base method.
func (m *MockConnMgr) Health() *health.Health {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health")
	ret0, _ := ret[0].(*health.Health)
	return ret0
}

// Health indicates an expected call of Health.
func (mr *MockConnMgrMockRecorder) Health() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockConnMgr)(nil).Health))
}

// JsPublish mocks base method.
func (m *MockConnMgr) JsPublish(ctx context.Context, subject string, message []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JsPublish", ctx, subject, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// JsPublish indicates an expected call of JsPublish.
func (mr *MockConnMgrMockRecorder) JsPublish(ctx, subject, message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JsPublish", reflect.TypeOf((*MockConnMgr)(nil).JsPublish), ctx, subject, message)
}

// Publish mocks base method.
func (m *MockConnMgr) Publish(ctx context.Context, subject string, message []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, subject, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockConnMgrMockRecorder) Publish(ctx, subject, message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockConnMgr)(nil).Publish), ctx, subject, message)
}

// isConnected mocks base method.
func (m *MockConnMgr) isConnected() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "isConnected")
	ret0, _ := ret[0].(bool)
	return ret0
}

// isConnected indicates an expected call of isConnected.
func (mr *MockConnMgrMockRecorder) isConnected() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isConnected", reflect.TypeOf((*MockConnMgr)(nil).isConnected))
}

// MockJsSubMgr is a mock of JsSubMgr interface.
type MockJsSubMgr struct {
	ctrl     *gomock.Controller
	recorder *MockJsSubMgrMockRecorder
	isgomock struct{}
}

// MockJsSubMgrMockRecorder is the mock recorder for MockJsSubMgr.
type MockJsSubMgrMockRecorder struct {
	mock *MockJsSubMgr
}

// NewMockJsSubMgr creates a new mock instance.
func NewMockJsSubMgr(ctrl *gomock.Controller) *MockJsSubMgr {
	mock := &MockJsSubMgr{ctrl: ctrl}
	mock.recorder = &MockJsSubMgrMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJsSubMgr) EXPECT() *MockJsSubMgrMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockJsSubMgr) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockJsSubMgrMockRecorder) Close(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockJsSubMgr)(nil).Close), ctx)
}

// DeleteSubscriber mocks base method.
func (m *MockJsSubMgr) DeleteSubscriber(ctx context.Context, identity []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSubscriber", ctx, identity)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSubscriber indicates an expected call of DeleteSubscriber.
func (mr *MockJsSubMgrMockRecorder) DeleteSubscriber(ctx, identity any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubscriber", reflect.TypeOf((*MockJsSubMgr)(nil).DeleteSubscriber), ctx, identity)
}

// NewSubscriber mocks base method.
func (m *MockJsSubMgr) NewSubscriber(ctx context.Context, stream streamIntf, identity []string, consumeType SubType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSubscriber", ctx, stream, identity, consumeType)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewSubscriber indicates an expected call of NewSubscriber.
func (mr *MockJsSubMgrMockRecorder) NewSubscriber(ctx, stream, identity, consumeType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSubscriber", reflect.TypeOf((*MockJsSubMgr)(nil).NewSubscriber), ctx, stream, identity, consumeType)
}

// Subscribe mocks base method.
func (m *MockJsSubMgr) Subscribe(ctx context.Context, identity []string, handler pubsub.JsSubscribeFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, identity, handler)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockJsSubMgrMockRecorder) Subscribe(ctx, identity, handler any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockJsSubMgr)(nil).Subscribe), ctx, identity, handler)
}

// MockSubMgr is a mock of SubMgr interface.
type MockSubMgr struct {
	ctrl     *gomock.Controller
	recorder *MockSubMgrMockRecorder
	isgomock struct{}
}

// MockSubMgrMockRecorder is the mock recorder for MockSubMgr.
type MockSubMgrMockRecorder struct {
	mock *MockSubMgr
}

// NewMockSubMgr creates a new mock instance.
func NewMockSubMgr(ctrl *gomock.Controller) *MockSubMgr {
	mock := &MockSubMgr{ctrl: ctrl}
	mock.recorder = &MockSubMgrMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubMgr) EXPECT() *MockSubMgrMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockSubMgr) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockSubMgrMockRecorder) Close(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSubMgr)(nil).Close), ctx)
}

// DeleteSubscriber mocks base method.
func (m *MockSubMgr) DeleteSubscriber(ctx context.Context, topicName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSubscriber", ctx, topicName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSubscriber indicates an expected call of DeleteSubscriber.
func (mr *MockSubMgrMockRecorder) DeleteSubscriber(ctx, topicName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubscriber", reflect.TypeOf((*MockSubMgr)(nil).DeleteSubscriber), ctx, topicName)
}

// NewSubscriber mocks base method.
func (m *MockSubMgr) NewSubscriber(ctx context.Context, conn *nats.Conn, topicName string, consumeType SubType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSubscriber", ctx, conn, topicName, consumeType)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewSubscriber indicates an expected call of NewSubscriber.
func (mr *MockSubMgrMockRecorder) NewSubscriber(ctx, conn, topicName, consumeType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSubscriber", reflect.TypeOf((*MockSubMgr)(nil).NewSubscriber), ctx, conn, topicName, consumeType)
}

// Subscribe mocks base method.
func (m *MockSubMgr) Subscribe(ctx context.Context, topicName string, handler pubsub.SubscribeFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, topicName, handler)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockSubMgrMockRecorder) Subscribe(ctx, topicName, handler any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockSubMgr)(nil).Subscribe), ctx, topicName, handler)
}
