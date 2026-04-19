package auth_test

import (
	"context"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
)

type Method string

var (
	CreateUser           Method = "CreateUser"
	GetUserByID          Method = "GetUserByID"
	GetUserByEmail       Method = "GetUserByEmail"
	CheckTakenUserFields Method = "CheckTakenUserFields"
	SetUserLastLogin     Method = "SetUserLastLogin"
	CreateRefreshToken   Method = "CreateRefreshToken"
	GetRefreshToken      Method = "GetRefreshToken"
	DeleteRefreshToken   Method = "DeleteRefreshToken"
)

// call is a mock call for testing
type call struct {
	returns []interface{}
}

// mockStore is a mock auth store for testing
type mockStore struct {
	t     *testing.T
	calls map[Method][]call
}

func newMockStore(t *testing.T) *mockStore {
	t.Helper()

	return &mockStore{
		t:     t,
		calls: make(map[Method][]call),
	}
}

// on sets the return values for the next call of a method
func (s *mockStore) on(method Method, returns ...any) *mockStore {
	s.t.Helper()
	s.calls[method] = append(s.calls[method], call{returns: returns})
	return s
}

// next gets the return values for the next call of a method
func (s *mockStore) next(method Method) call {
	s.t.Helper()
	calls, ok := s.calls[method]
	if !ok || len(calls) == 0 {
		s.t.Fatalf("unexpected call to %s", method)
	}
	call := calls[0]
	s.calls[method] = calls[1:]
	return call
}

func (s *mockStore) CreateUser(ctx context.Context, user models.NewUser) (int64, error) {
	call := s.next(CreateUser)
	return call.returns[0].(int64), testutil.ConvertError(call.returns[1])
}

func (s *mockStore) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	call := s.next(GetUserByID)
	return call.returns[0].(models.User), testutil.ConvertError(call.returns[1])
}

func (s *mockStore) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	call := s.next(GetUserByEmail)
	return call.returns[0].(models.User), testutil.ConvertError(call.returns[1])
}

func (s *mockStore) CheckTakenUserFields(ctx context.Context, request models.NewUserRequest) (models.UserFieldsTaken, error) {
	call := s.next(CheckTakenUserFields)
	return call.returns[0].(models.UserFieldsTaken), testutil.ConvertError(call.returns[1])
}

func (s *mockStore) SetUserLastLogin(ctx context.Context, id int64) error {
	call := s.next(SetUserLastLogin)
	return testutil.ConvertError(call.returns[0])
}

func (s *mockStore) CreateRefreshToken(ctx context.Context, refreshToken models.RefreshToken) error {
	call := s.next(CreateRefreshToken)
	return testutil.ConvertError(call.returns[0])
}

func (s *mockStore) GetRefreshToken(ctx context.Context, tokenHash string) (models.RefreshToken, error) {
	call := s.next(GetRefreshToken)
	return call.returns[0].(models.RefreshToken), testutil.ConvertError(call.returns[1])
}

func (s *mockStore) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	call := s.next(DeleteRefreshToken)
	return testutil.ConvertError(call.returns[0])
}
