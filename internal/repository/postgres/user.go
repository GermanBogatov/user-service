package postgres

import (
	"context"
	"github.com/GermanBogatov/user-service/internal/entity"
	"github.com/GermanBogatov/user-service/pkg/postgresql"
)

var _ IUser = &User{}

type IUser interface {
	CreateUser(ctx context.Context, user entity.User) error
}

type User struct {
	client postgresql.Client
}

func NewUser(client postgresql.Client) IUser {
	return &User{
		client: client,
	}
}

func (u *User) CreateUser(ctx context.Context, user entity.User) error {
	return nil
}
