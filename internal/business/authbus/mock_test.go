package authbus_test

import (
	"context"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
)

var (
	CreateUser           testutil.Method = "CreateUser"
	GetUserByID          testutil.Method = "GetUserByID"
	GetUserByEmail       testutil.Method = "GetUserByEmail"
	CheckTakenUserFields testutil.Method = "CheckTakenUserFields"
	SetUserLastLogin     testutil.Method = "SetUserLastLogin"
	CreateRefreshToken   testutil.Method = "CreateRefreshToken"
	GetRefreshToken      testutil.Method = "GetRefreshToken"
	DeleteRefreshToken   testutil.Method = "DeleteRefreshToken"
)

// mockStore is a mock authbus store for testing
type mockStore struct {
	testutil.Mock
}

func newMockStore(t *testing.T) *mockStore {
	t.Helper()

	return &mockStore{
		Mock: testutil.NewMock(t),
	}
}

func (s *mockStore) CreateUser(ctx context.Context, user models.NewUser) (int64, error) {
	call := s.Next(CreateUser)
	return call.Returns[0].(int64), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	call := s.Next(GetUserByID)
	return call.Returns[0].(models.User), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	call := s.Next(GetUserByEmail)
	return call.Returns[0].(models.User), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) CheckTakenUserFields(ctx context.Context, request models.NewUserRequest) (models.UserFieldsTaken, error) {
	call := s.Next(CheckTakenUserFields)
	return call.Returns[0].(models.UserFieldsTaken), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) SetUserLastLogin(ctx context.Context, id int64) error {
	call := s.Next(SetUserLastLogin)
	return testutil.ConvertError(call.Returns[0])
}

func (s *mockStore) CreateRefreshToken(ctx context.Context, refreshToken models.RefreshToken) error {
	call := s.Next(CreateRefreshToken)
	return testutil.ConvertError(call.Returns[0])
}

func (s *mockStore) GetRefreshToken(ctx context.Context, tokenHash string) (models.RefreshToken, error) {
	call := s.Next(GetRefreshToken)
	return call.Returns[0].(models.RefreshToken), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	call := s.Next(DeleteRefreshToken)
	return testutil.ConvertError(call.Returns[0])
}
