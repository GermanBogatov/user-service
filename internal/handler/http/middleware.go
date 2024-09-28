package http

import (
	"context"
	"fmt"
	"github.com/GermanBogatov/user-service/internal/common/apperror"
	"github.com/GermanBogatov/user-service/internal/common/metrics"
	"github.com/GermanBogatov/user-service/internal/common/response"
	"github.com/GermanBogatov/user-service/internal/config"
	"github.com/GermanBogatov/user-service/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func appMiddleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		routeContext := chi.RouteContext(r.Context())
		pattern := routeContext.RoutePattern()
		defer metrics.ObserveRequestDurationSeconds(method, pattern)()
		if routeContext.RoutePatterns[0] != authV1+"/*" {

			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			if len(authHeader) != 2 {
				response.RespondError(w, r, apperror.UnauthorizedError(apperror.ErrMalformedToken))
				return
			}

			accessToken := authHeader[1]
			key := []byte(config.JWTSecret)

			token, err := jwt.ParseWithClaims(accessToken, &entity.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, apperror.ErrInvalidSigningMethod
				}
				return key, nil
			})

			if err != nil {
				response.RespondError(w, r, apperror.UnauthorizedError(errors.Wrap(err, apperror.ErrMalformedToken.Error())))
				return
			}

			if !token.Valid {
				response.RespondError(w, r, apperror.UnauthorizedError(apperror.ErrTokenIsInspired))
				return
			}

			claims, ok := token.Claims.(*entity.UserClaims)
			if !ok {
				response.RespondError(w, r, apperror.UnauthorizedError(err))
				return
			}

			setCtxValue(r, config.ParamID, claims.Id)
			setCtxValue(r, config.ParamRoles, claims.Roles)
			//todo прокинуть ID
			fmt.Println(claims.Id)
		}

		err := h(w, r)
		if err != nil {
			metrics.IncRequestTotal(metrics.FailStatus, method, pattern)
			response.RespondError(w, r, err)
			return
		}

		metrics.IncRequestTotal(metrics.OkStatus, method, pattern)
	}
}

func setCtxValue(r *http.Request, key, value any) {
	ctx := r.Context()
	req := r.WithContext(context.WithValue(ctx, key, value))
	*r = *req
}
