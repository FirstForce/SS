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

func TestPhotoController_GetPhotos(t *testing.T) {
	tests := []struct {
		name             string
		userEmail        string
		mockPhotos       []*domain.Photo
		mockError        error
		expectedStatus   int
		expectedContains string
	}{
		{
			name:             "no photos",
			userEmail:        "empty@example.com",
			mockPhotos:       []*domain.Photo{},
			expectedStatus:   http.StatusOK,
			expectedContains: "[]",
		},
		{
			name:             "repository error",
			userEmail:        "error@example.com",
			mockError:        errors.New("db error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "Failed to fetch photos",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_domain.NewMockPhotoRepository(ctrl)
			ctlr := routes.PhotoController{PhotoRepository: mockRepo}

			req := httptest.NewRequest(http.MethodGet, "/photos", nil)
			ctx := context.WithValue(req.Context(), "email", tt.userEmail)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			if tt.mockPhotos != nil || tt.mockError != nil {
				mockRepo.EXPECT().
					GetPhotos(ctx, gomock.Any()).
					Return(tt.mockPhotos, tt.mockError)
			}

			ctlr.GetPhotos(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
			if tt.expectedContains != "" && !strings.Contains(rr.Body.String(), tt.expectedContains) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedContains, rr.Body.String())
			}
		})
	}
}
