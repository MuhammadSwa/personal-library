package repositories

import (
	"context"

	"github.com/muhammadswa/personal-library/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (ur *Repositories) GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	user, err := ur.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *Repositories) CreateUser(ctx context.Context, email, password, username string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	id, err := ur.db.CreateUser(ctx, database.CreateUserParams{
		Email:          email,
		HashedPassword: string(hashedPassword),
		Username:       username,
	})
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
