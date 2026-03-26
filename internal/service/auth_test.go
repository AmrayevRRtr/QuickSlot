package service

import (
	"QuickSlot/internal/model"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// mock of UserRepository
type mockUserRepo struct {
	users  map[string]*model.User
	lastID int64
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*model.User)}
}

func (m *mockUserRepo) CreateUser(user *model.User) (int64, error) {
	if _, exists := m.users[user.Email]; exists {
		return 0, errors.New("email already exists")
	}
	m.lastID++
	user.ID = m.lastID
	m.users[user.Email] = user
	return m.lastID, nil
}

func (m *mockUserRepo) GetUserByEmail(email string) (*model.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func TestRegister(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo)

	id, err := svc.Register("test@mail.com", "pass1234")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 1 {
		t.Fatalf("expected id=1, got %d", id)
	}

	// check password is hashed
	user := repo.users["test@mail.com"]
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("pass1234"))
	if err != nil {
		t.Fatal("password was not hashed correctly")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo)

	_, err := svc.Register("dup@mail.com", "pass1234")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Register("dup@mail.com", "pass1234")
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestLoginSuccess(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo)

	_, err := svc.Register("login@mail.com", "mypassword")
	if err != nil {
		t.Fatal(err)
	}

	user, err := svc.Login("login@mail.com", "mypassword")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if user.Email != "login@mail.com" {
		t.Fatalf("expected email login@mail.com, got %s", user.Email)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo)

	_, _ = svc.Register("wrong@mail.com", "correctpass")

	_, err := svc.Login("wrong@mail.com", "wrongpass")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestLoginNoUser(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo)

	_, err := svc.Login("noone@mail.com", "pass")
	if err == nil {
		t.Fatal("expected error for non-existing user")
	}
}
