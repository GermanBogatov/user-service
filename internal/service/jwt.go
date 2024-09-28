package service

import (
	"context"
	"github.com/GermanBogatov/user-service/internal/common/apperror"
	"github.com/GermanBogatov/user-service/internal/entity"
	"github.com/GermanBogatov/user-service/internal/repository/cache"
	"github.com/GermanBogatov/user-service/internal/repository/postgres"
	"github.com/GermanBogatov/user-service/pkg/logging"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

var _ IJWT = &JWT{}

type JWT struct {
	userRepo postgres.IUser
	cache    cache.ICache
	secret   string
	jwtTTL   time.Duration
}

func NewJWT(userRepo postgres.IUser, cache cache.ICache, secret string, jwtTTL int) IJWT {
	return &JWT{
		userRepo: userRepo,
		cache:    cache,
		secret:   secret,
		jwtTTL:   time.Duration(jwtTTL) * time.Second,
	}
}

type IJWT interface {
	UpdateRefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	GenerateAccessToken(user entity.User) (string, string, error)
}

// UpdateRefreshToken - обновление рефреш токена
func (j *JWT) UpdateRefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	userID, err := j.cache.Get(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", "", apperror.ErrRefreshTokenNotFound
		}
		return "", "", err
	}

	var (
		user    entity.User
		errUser error
	)
	user, errUser = j.cache.GetUser(ctx, userID)
	if errUser != nil {
		if errors.Is(errUser, redis.Nil) {
			user, errUser = j.userRepo.GetUserByID(ctx, userID)
			if errUser != nil {
				return "", "", errUser
			}
		} else {
			return "", "", errUser
		}
	}

	go func() {
		ctxDel, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
		defer cancel()
		errDelete := j.cache.Delete(ctxDel, refreshToken)
		if errDelete != nil {
			logging.Errorf("error deleting refresh token [%s]: %v", refreshToken, errDelete)
		}
	}()

	return j.GenerateAccessToken(user)

}

// GenerateAccessToken - генерация рефреш токена
func (j *JWT) GenerateAccessToken(user entity.User) (string, string, error) {
	key := []byte(j.secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &entity.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        user.ID,
			Audience:  "users",
			ExpiresAt: time.Now().Add(j.jwtTTL).Unix(),
		},
		Email: user.Email,
		Roles: user.Roles,
	})

	accessToken, err := token.SignedString(key)
	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.New().String()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
		defer cancel()

		errSet := j.cache.SetRefreshToken(ctx, refreshToken, user.ID)
		if errSet != nil {
			logging.Errorf("error set refresh token [%s]: %v", refreshToken, errSet)
		}

		errSet = j.cache.SetUser(ctx, user.ID, user)
		if errSet != nil {
			logging.Errorf("error set refresh token [%s]: %v", refreshToken, errSet)
		}
	}()

	return accessToken, refreshToken, err
}
