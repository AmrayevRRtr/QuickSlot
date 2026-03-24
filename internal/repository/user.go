package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
)

type UserRepository interface {
	CreateUser(user *model.User) (int64, error)
	GetUserByEmail(email string) (*model.User, error)
}
type userRepository struct {
	db *mysql.Dialect
}

func NewUserRepository(db *mysql.Dialect) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) (int64, error) {
	result, err := r.db.DB.Exec(`INSERT INTO users (email, password, role) VALUES (?, ?, ?)`,
		user.Email, user.Password, user.Role)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.DB.Get(&user, "SELECT id, email, password, role FROM users WHERE email=?", email)
	return &user, err
}
