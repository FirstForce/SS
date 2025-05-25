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

func TestUserController_Register(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      string
		mockSaveReturn error
		expectedStatus int
		expectedUser   *domain.User
	}{
		{
			name:           "successful registration",
			inputBody:      `{"email": "test@example.com", "password": "securepass"}`,
			mockSaveReturn: nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "user already exists",
			inputBody:      `{"email": "test@example.com", "password": "securepass"}`,
			mockSaveReturn: errors.New("user already exists"), // Simulate existing user
			expectedStatus: http.StatusConflict,
			expectedUser: &domain.User{
				Email: "test@example.com",
			},
		},
		{
			name:           "invalid JSON",
			inputBody:      `invalid-json`,
			mockSaveReturn: nil, // Save won't be called
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "repository error",
			inputBody:      `{"email": "test@example.com", "password": "securepass"}`,
			mockSaveReturn: errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_domain.NewMockUserRepository(ctrl)
			ctlr := routes.UserController{UserRepository: mockRepo}

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			if tt.expectedStatus != http.StatusBadRequest {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(tt.expectedUser, nil)
			}
			if tt.expectedStatus != http.StatusBadRequest && tt.expectedStatus != http.StatusConflict {
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.mockSaveReturn)
			}

			ctlr.Register(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestUserController_Login(t *testing.T) {
	tests := []struct {
		name             string
		inputBody        string
		mockUser         *domain.User
		mockError        error
		expectedStatus   int
		expectedContains string // optional: check part of response
	}{
		{
			name:      "successful login",
			inputBody: `{"email": "test@example.com", "password": "password123"}`,
			mockUser: &domain.User{
				Email:    "test@example.com",
				Password: "$2a$12$.OZ5oYXEsFvcaaVh/nmgt.cknGSFzKVlr.wkrzyCl5rgHuAGGkhiS",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:             "invalid JSON",
			inputBody:        `not-a-json`,
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid request body",
		},
		{
			name:             "user not found",
			inputBody:        `{"email": "missing@example.com", "password": "password123"}`,
			mockUser:         nil,
			mockError:        errors.New("user not found"),
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "Invalid email or password",
		},
		{
			name:           "repository error",
			inputBody:      `{"email": "test@example.com", "password": "password123"}`,
			mockUser:       nil,
			mockError:      errors.New("db error"),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_domain.NewMockUserRepository(ctrl)
			ctlr := routes.UserController{UserRepository: mockRepo}

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			if tt.mockUser != nil || tt.mockError != nil {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(tt.mockUser, tt.mockError)
			}

			ctlr.Login(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedContains != "" && !strings.Contains(rr.Body.String(), tt.expectedContains) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedContains, rr.Body.String())
			}
		})
	}
}

func TestUserController_GetProfile(t *testing.T) {
	tests := []struct {
		name             string
		userEmail        string
		mockUser         *domain.User
		mockError        error
		expectedStatus   int
		expectedContains string
	}{
		{
			name:      "successful profile fetch",
			userEmail: "test@example.com",
			mockUser: &domain.User{
				Email: "test@example.com",
			},
			expectedStatus:   http.StatusOK,
			expectedContains: "test@example.com",
		},
		{
			name:             "user not found",
			userEmail:        "missing@example.com",
			mockUser:         nil,
			mockError:        errors.New("user not found"),
			expectedStatus:   http.StatusNotFound,
			expectedContains: "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_domain.NewMockUserRepository(ctrl)
			ctlr := routes.UserController{UserRepository: mockRepo}

			// Build request and context
			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			ctx := context.WithValue(req.Context(), "email", tt.userEmail)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			// Set expectation
			if tt.userEmail != "" {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), tt.userEmail).Return(tt.mockUser, tt.mockError)
			}

			ctlr.GetProfile(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
			if tt.expectedContains != "" && !strings.Contains(rr.Body.String(), tt.expectedContains) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedContains, rr.Body.String())
			}
		})
	}
}

func TestUserController_GetProfile_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockUserRepository(ctrl)
	ctlr := routes.UserController{UserRepository: mockRepo}

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rr := httptest.NewRecorder()

	// No email in context
	ctlr.GetProfile(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Email not found in context") {
		t.Errorf("expected body to contain 'Email not found in context', got %q", rr.Body.String())
	}
}

func TestUserController_GetProfile_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockUserRepository(ctrl)
	ctlr := routes.UserController{UserRepository: mockRepo}

	req := httptest.NewRequest(http.MethodPost, "/profile", nil) // Using POST instead of GET
	rr := httptest.NewRecorder()

	ctlr.GetProfile(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Method not allowed") {
		t.Errorf("expected body to contain 'Method not allowed', got %q", rr.Body.String())
	}
}

func TestUserController_Register_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockUserRepository(ctrl)
	ctlr := routes.UserController{UserRepository: mockRepo}

	req := httptest.NewRequest(http.MethodGet, "/register", nil) // Using GET instead of POST
	rr := httptest.NewRecorder()

	ctlr.Register(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Method not allowed") {
		t.Errorf("expected body to contain 'Method not allowed', got %q", rr.Body.String())
	}
}

func TestUserController_Login_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockUserRepository(ctrl)
	ctlr := routes.UserController{UserRepository: mockRepo}

	req := httptest.NewRequest(http.MethodGet, "/login", nil) // Using GET instead of POST
	rr := httptest.NewRecorder()

	ctlr.Login(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Method not allowed") {
		t.Errorf("expected body to contain 'Method not allowed', got %q", rr.Body.String())
	}
}
