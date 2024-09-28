package http

import (
	"encoding/json"
	"fmt"
	"github.com/GermanBogatov/user-service/internal/common/apperror"
	"github.com/GermanBogatov/user-service/internal/common/helpers"
	"github.com/GermanBogatov/user-service/internal/common/response"
	"github.com/GermanBogatov/user-service/internal/config"
	"github.com/GermanBogatov/user-service/internal/handler/http/mapper"
	"github.com/GermanBogatov/user-service/internal/handler/http/model"
	"github.com/GermanBogatov/user-service/internal/handler/http/validator"
	"github.com/GermanBogatov/user-service/pkg/logging"
	"github.com/pkg/errors"
	"net/http"
)

// SignUp - хэндлер регистрации
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var createUser model.SignUpRequest
	defer func() {
		err := r.Body.Close()
		if err != nil {
			logging.Error("error close request body")
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&createUser); err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "json decode"))
	}

	err := validator.ValidateSignUpUser(createUser)
	if err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "validate create user"))
	}

	user := mapper.MapToEntityUser(createUser)
	user.GenerateID()
	user.SetPasswordHash(helpers.GeneratePasswordHash(createUser.Password))
	user.GenerateCreatedDate()
	// todo когда админ появится условия предусмотреть
	user.AddRoleDeveloper()

	token, refreshToken, err := h.jwtService.GenerateAccessToken(user)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	user.SetJWT(token, refreshToken)

	err = h.userService.CreateUser(ctx, user)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return response.RespondSuccessCreate(w, mapper.MapToUserWithJWTResponse(http.StatusCreated, user))

}

// SignIn - хэндлер авторизации
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var signInUser model.SignInRequest
	defer func() {
		err := r.Body.Close()
		if err != nil {
			logging.Error("error close request body")
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&signInUser); err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "json decode"))
	}

	err := validator.ValidateSignInUser(signInUser)
	if err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "validate create user"))
	}

	passwordHash := helpers.GeneratePasswordHash(signInUser.Password)
	user, err := h.userService.GetUserByEmailAndPassword(ctx, signInUser.Email, passwordHash)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	token, refreshToken, err := h.jwtService.GenerateAccessToken(user)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	user.SetJWT(token, refreshToken)

	return response.RespondSuccess(w, mapper.MapToUserWithJWTResponse(http.StatusOK, user))
}

// GetUserByID - хэндлер получения пользователя по идентификатору
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	res1 := ctx.Value("id")
	res2 := ctx.Value("roles")
	fmt.Println("res1=", res1)
	fmt.Println("res2=", res2)
	userID, err := helpers.GetUuidFromPath(r, config.ParamID)
	if err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "get uuid from header"))
	}

	user, err := h.userService.GetUserByID(ctx, userID.String())
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return response.RespondSuccess(w, mapper.MapToUserResponse(http.StatusOK, user))
}

// UpdateRefreshToken - хэндлер обновления jwt-токена
func (h *Handler) UpdateRefreshToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	refreshToken, err := helpers.GetStringFromPath(r, config.ParamID)
	if err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "get uuid from header"))
	}

	token, newRefreshToken, err := h.jwtService.UpdateRefreshToken(ctx, refreshToken)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return response.RespondSuccess(w, mapper.MapToJWTResponse(http.StatusOK, token, newRefreshToken))
}
