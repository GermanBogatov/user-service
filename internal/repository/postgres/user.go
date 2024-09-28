package postgres

import (
	"context"
	"github.com/GermanBogatov/user-service/internal/common/apperror"
	"github.com/GermanBogatov/user-service/internal/entity"
	"github.com/GermanBogatov/user-service/pkg/postgresql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

var _ IUser = &User{}

type IUser interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	GetUserByEmailAndPassword(ctx context.Context, email, password string) (entity.User, error)
}

type User struct {
	client postgresql.Client
}

func NewUser(client postgresql.Client) IUser {
	return &User{
		client: client,
	}
}

// CreateUser - создание пользователя
func (u *User) CreateUser(ctx context.Context, user entity.User) error {
	q := `
	INSERT INTO users 
    	(id,name,surname,email,password,roles,created_date) 
    VALUES 
		($1,$2,$3,$4,$5,$6,$7);
		`

	_, err := u.client.Exec(ctx, q, user.ID, user.Name, user.Surname, user.Email, user.Password, user.Roles, user.CreatedDate)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return apperror.ErrUserIsExistWithEmail
			}
			return err
		}
	}
	return nil
}

// GetUserByEmailAndPassword - получение пользователя по емайл и паролю
func (u *User) GetUserByEmailAndPassword(ctx context.Context, email, password string) (entity.User, error) {
	q := `
		SELECT id,name,surname,email,password,roles,created_date,updated_date 
		FROM users
		WHERE email=$1 AND password=$2;	
		`

	var user entity.User
	err := u.client.QueryRow(ctx, q, email, password).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Roles, &user.CreatedDate, &user.UpdatedDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, apperror.ErrUserNotFound
		}
		return entity.User{}, err
	}

	return user, nil
}

// GetUserByID - получение пользователя по идентификатору
func (u *User) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	q := `
		SELECT id,name,surname,email,password,roles,created_date,updated_date 
		FROM users
		WHERE id=$1;	
		`

	var user entity.User
	err := u.client.QueryRow(ctx, q, id).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Roles, &user.CreatedDate, &user.UpdatedDate)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
