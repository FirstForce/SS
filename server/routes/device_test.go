package routes_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"mqtt-streaming-server/domain"
	mock_domain "mqtt-streaming-server/mocks"
	"mqtt-streaming-server/routes"
)

func TestDeviceController_GetDevices(t *testing.T) {
	tests := []struct {
		name             string
		userEmail        string
		userRole         string
		mockDevices      []*domain.Device
		mockError        error
		expectedStatus   int
		expectedContains string
	}{
		{
			name:      "successful device fetch",
			userEmail: "user@example.com",
			userRole:  "admin",
			mockDevices: []*domain.Device{
				{ID: "dev-1", DeviceName: "iPhone"},
				{ID: "dev-2", DeviceName: "MacBook"},
			},
			expectedStatus:   http.StatusOK,
			expectedContains: "iPhone",
		},
		{
			name:             "no devices",
			userEmail:        "empty@example.com",
			userRole:         "admin",
			mockDevices:      []*domain.Device{},
			expectedStatus:   http.StatusOK,
			expectedContains: "[]",
		},
		{
			name:             "repository error",
			userEmail:        "fail@example.com",
			userRole:         "admin",
			mockError:        errors.New("db error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "Failed to fetch devices",
		},
		{
			name:             "unauthorized access",
			userEmail:        "unauthorized@example.com",
			userRole:         "user",
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_domain.NewMockDeviceRepository(ctrl)
			ctlr := routes.DeviceController{DeviceRepository: mockRepo}

			req := httptest.NewRequest(http.MethodGet, "/devices", nil)
			ctx := context.WithValue(req.Context(), "email", tt.userEmail)
			ctx = context.WithValue(ctx, "role", tt.userRole)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			if tt.mockDevices != nil || tt.mockError != nil {
				mockRepo.EXPECT().
					GetAllDevices(gomock.Any()).
					Return(tt.mockDevices, tt.mockError)
			}

			ctlr.GetDevices(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
			if tt.expectedContains != "" && !strings.Contains(rr.Body.String(), tt.expectedContains) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedContains, rr.Body.String())
			}
		})
	}
}

func TestDeviceController_SwitchDeviceMode(t *testing.T) {
	tests := []struct {
		name             string
		userEmail        string
		userRole         string
		deviceID         string
		mockError        error
		expectedStatus   int
		expectedContains string
	}{
		{
			name:             "unauthorized access",
			userEmail:        "unauthorized@example.com",
			userRole:         "user",
			deviceID:         "dev-1",
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_domain.NewMockDeviceRepository(ctrl)
			ctlr := routes.DeviceController{DeviceRepository: mockRepo}

			url := "/devices/" + tt.deviceID + "/switch"
			req := httptest.NewRequest(http.MethodPost, url, nil)
			ctx := context.WithValue(req.Context(), "email", tt.userEmail)
			ctx = context.WithValue(ctx, "role", tt.userRole)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			ctlr.SwitchDeviceMode(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
			if tt.expectedContains != "" && !strings.Contains(rr.Body.String(), tt.expectedContains) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedContains, rr.Body.String())
			}
		})
	}
}
