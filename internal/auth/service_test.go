package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/auth"
	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
)

func TestService_Register(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		t.Parallel()
		tc := []struct {
			name                 string
			request              models.NewUserRequest
			userAgent            string
			remember             bool
			storeSetup           func(s *mockStore)
			expectedErr          error
			accessTokenExpected  bool
			refreshTokenExpected bool
		}{
			{
				name: "success: registered (don't remember)",
				request: models.NewUserRequest{
					Email:           "test.user@critiquefi.com",
					DisplayName:     "test.user",
					Name:            "Test User",
					Password:        "password",
					ConfirmPassword: "password",
				},
				userAgent: "test",
				remember:  false,
				storeSetup: func(s *mockStore) {
					s.on(CheckTakenUserFields, models.UserFieldsTaken{}, nil).
						on(CreateUser, int64(1), nil).
						on(GetUserByID, models.User{
							ID:           1,
							Email:        "user@critiquefi.com",
							DisplayName:  "user.name",
							Name:         "Test User",
							IsAdmin:      false,
							PasswordHash: "password",
							IsActive:     true,
						}, nil)
				},
				expectedErr:          nil,
				accessTokenExpected:  true,
				refreshTokenExpected: false,
			},
			{
				name: "success: registered (remember)",
				request: models.NewUserRequest{
					Email:           "test.user@critiquefi.com",
					DisplayName:     "test.user",
					Name:            "Test User",
					Password:        "password",
					ConfirmPassword: "password",
				},
				userAgent: "test",
				remember:  true,
				storeSetup: func(s *mockStore) {
					s.on(CheckTakenUserFields, models.UserFieldsTaken{}, nil).
						on(CreateUser, int64(1), nil).
						on(GetUserByID, models.User{
							ID:           1,
							Email:        "user@critiquefi.com",
							DisplayName:  "user.name",
							Name:         "Test User",
							IsAdmin:      false,
							PasswordHash: "password",
							IsActive:     true,
						}, nil).
						on(CreateRefreshToken, nil)
				},
				expectedErr:          nil,
				accessTokenExpected:  true,
				refreshTokenExpected: true,
			},
		}

		for _, tc := range tc {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				s := newMockStore(t)
				tc.storeSetup(s)

				svc := auth.NewService(auth.ServiceConfig{
					Store:                    s,
					AccessTokenKey:           "key",
					AccessTokenTTL:           time.Hour,
					RefreshTokenTTL:          time.Hour,
					RefreshTokenCookieName:   "rt",
					RefreshTokenCookieDomain: "critiquefi.com",
				})

				accessToken, refreshToken, err := svc.Register(context.Background(), tc.request, tc.userAgent, tc.remember)
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}

				if tc.accessTokenExpected && accessToken == "" {
					t.Error("expected access token, got nothing")
				}
				if !tc.accessTokenExpected && accessToken != "" {
					t.Errorf("expected no access token, got %v", accessToken)
				}
				if tc.refreshTokenExpected && refreshToken == "" {
					t.Error("expected refresh token, got nothing")
				}
				if !tc.refreshTokenExpected && refreshToken != "" {
					t.Errorf("expected no refresh token, got %v", refreshToken)
				}
			})
		}
	})
}
