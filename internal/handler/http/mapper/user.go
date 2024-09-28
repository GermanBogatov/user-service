package mapper

import (
	"github.com/GermanBogatov/user-service/internal/common/response"
	"github.com/GermanBogatov/user-service/internal/config"
	"github.com/GermanBogatov/user-service/internal/entity"
	"github.com/GermanBogatov/user-service/internal/handler/http/model"
)

func MapToEntityUser(user model.SignUpRequest) entity.User {
	return entity.User{
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
	}
}

func MapToUserWithJWTResponse(code int, user entity.User) response.ViewResponse {
	var (
		jwt         *model.JWT
		updatedDate *string
	)
	if user.UpdatedDate != nil {
		updateTime := user.UpdatedDate.Format(config.IsoTimeLayout)
		updatedDate = &updateTime
	}

	if user.JWT != nil {
		jwt = &model.JWT{
			Token:        user.JWT.Token,
			RefreshToken: user.JWT.RefreshToken,
		}
	}

	return response.ViewResponse{
		Code: code,
		Result: model.SignUpResponse{
			ID:          user.ID,
			Name:        user.Name,
			Surname:     user.Surname,
			Email:       user.Email,
			CreatedDate: user.CreatedDate.Format(config.IsoTimeLayout),
			UpdatedDate: updatedDate,
			JWT:         jwt,
		},
	}
}

func MapToUserResponse(code int, user entity.User) response.ViewResponse {
	var (
		updatedDate *string
	)

	if user.UpdatedDate != nil {
		updateTime := user.UpdatedDate.Format(config.IsoTimeLayout)
		updatedDate = &updateTime
	}

	return response.ViewResponse{
		Code: code,
		Result: model.SignUpResponse{
			ID:          user.ID,
			Name:        user.Name,
			Surname:     user.Surname,
			Email:       user.Email,
			CreatedDate: user.CreatedDate.Format(config.IsoTimeLayout),
			UpdatedDate: updatedDate,
		},
	}
}

func MapToJWTResponse(code int, token, refreshToken string) response.ViewResponse {
	return response.ViewResponse{
		Code: code,
		Result: model.JWT{
			Token:        token,
			RefreshToken: refreshToken,
		},
	}
}
