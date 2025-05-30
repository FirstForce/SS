// Code generated by MockGen. DO NOT EDIT.
// Source: mqtt-streaming-server/domain (interfaces: UserRepository,PhotoRepository,DeviceRepository)
//
// Generated by this command:
//
//	mockgen mqtt-streaming-server/domain UserRepository,PhotoRepository,DeviceRepository
//

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	domain "mqtt-streaming-server/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
	isgomock struct{}
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// FindByEmail mocks base method.
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserRepositoryMockRecorder) FindByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindByEmail), ctx, email)
}

// Save mocks base method.
func (m *MockUserRepository) Save(ctx context.Context, email, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, email, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockUserRepositoryMockRecorder) Save(ctx, email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserRepository)(nil).Save), ctx, email, password)
}

// MockPhotoRepository is a mock of PhotoRepository interface.
type MockPhotoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPhotoRepositoryMockRecorder
	isgomock struct{}
}

// MockPhotoRepositoryMockRecorder is the mock recorder for MockPhotoRepository.
type MockPhotoRepositoryMockRecorder struct {
	mock *MockPhotoRepository
}

// NewMockPhotoRepository creates a new mock instance.
func NewMockPhotoRepository(ctrl *gomock.Controller) *MockPhotoRepository {
	mock := &MockPhotoRepository{ctrl: ctrl}
	mock.recorder = &MockPhotoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPhotoRepository) EXPECT() *MockPhotoRepositoryMockRecorder {
	return m.recorder
}

// GetPhotos mocks base method.
func (m *MockPhotoRepository) GetPhotos(ctx context.Context, filters map[string]any) ([]*domain.Photo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPhotos", ctx, filters)
	ret0, _ := ret[0].([]*domain.Photo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPhotos indicates an expected call of GetPhotos.
func (mr *MockPhotoRepositoryMockRecorder) GetPhotos(ctx, filters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPhotos", reflect.TypeOf((*MockPhotoRepository)(nil).GetPhotos), ctx, filters)
}

// Save mocks base method.
func (m *MockPhotoRepository) Save(ctx context.Context, photo *domain.Photo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, photo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockPhotoRepositoryMockRecorder) Save(ctx, photo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockPhotoRepository)(nil).Save), ctx, photo)
}

// MockDeviceRepository is a mock of DeviceRepository interface.
type MockDeviceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceRepositoryMockRecorder
	isgomock struct{}
}

// MockDeviceRepositoryMockRecorder is the mock recorder for MockDeviceRepository.
type MockDeviceRepositoryMockRecorder struct {
	mock *MockDeviceRepository
}

// NewMockDeviceRepository creates a new mock instance.
func NewMockDeviceRepository(ctrl *gomock.Controller) *MockDeviceRepository {
	mock := &MockDeviceRepository{ctrl: ctrl}
	mock.recorder = &MockDeviceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeviceRepository) EXPECT() *MockDeviceRepositoryMockRecorder {
	return m.recorder
}

// GetAllDevices mocks base method.
func (m *MockDeviceRepository) GetAllDevices(ctx context.Context) ([]*domain.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDevices", ctx)
	ret0, _ := ret[0].([]*domain.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDevices indicates an expected call of GetAllDevices.
func (mr *MockDeviceRepositoryMockRecorder) GetAllDevices(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDevices", reflect.TypeOf((*MockDeviceRepository)(nil).GetAllDevices), ctx)
}

// GetByID mocks base method.
func (m *MockDeviceRepository) GetByID(ctx context.Context, id string) (*domain.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*domain.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockDeviceRepositoryMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockDeviceRepository)(nil).GetByID), ctx, id)
}

// Save mocks base method.
func (m *MockDeviceRepository) Save(ctx context.Context, device *domain.Device) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, device)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockDeviceRepositoryMockRecorder) Save(ctx, device any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockDeviceRepository)(nil).Save), ctx, device)
}

// Update mocks base method.
func (m *MockDeviceRepository) Update(ctx context.Context, id string, device *domain.Device) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, device)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockDeviceRepositoryMockRecorder) Update(ctx, id, device any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDeviceRepository)(nil).Update), ctx, id, device)
}
