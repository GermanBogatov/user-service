package service

import (
	"context"
	"github.com/GermanBogatov/user-service/internal/entity"
	"github.com/GermanBogatov/user-service/internal/repository/postgres"
)

var _ IUser = &User{}

type IUser interface {
	CreateUser(ctx context.Context, user entity.User) error
}

type User struct {
	client postgres.IUser
}

func NewUser(client postgres.IUser) IUser {
	return &User{
		client: client,
	}
}

func (u *User) CreateUser(ctx context.Context, user entity.User) error {
	return nil
}
