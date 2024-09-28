package validator

import (
	"github.com/GermanBogatov/user-service/internal/common/apperror"
	"github.com/GermanBogatov/user-service/internal/handler/http/model"
	"strings"
)

func ValidateSignUpUser(user model.SignUpRequest) error {
	if strings.TrimSpace(user.Name) == "" {
		return apperror.ErrEmptyName
	}
	if strings.TrimSpace(user.Surname) == "" {
		return apperror.ErrEmptySurname
	}
	if strings.TrimSpace(user.Email) == "" {
		return apperror.ErrEmptyEmail
	}
	if strings.TrimSpace(user.Password) == "" {
		return apperror.ErrEmptyPassword
	}

	return nil
}

func ValidateSignInUser(user model.SignInRequest) error {
	if strings.TrimSpace(user.Email) == "" {
		return apperror.ErrEmptyEmail
	}
	if strings.TrimSpace(user.Password) == "" {
		return apperror.ErrEmptyPassword
	}

	return nil
}
