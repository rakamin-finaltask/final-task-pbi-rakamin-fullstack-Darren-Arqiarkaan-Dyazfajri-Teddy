package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"user-personalize/pkg/util/response"
	"user-personalize/pkg/util/service"
)

type Middleware interface {
	ValidateUser(ctx *gin.Context)
}

type middlewareImpl struct {
	jwtService service.JwtService
}

func (m *middlewareImpl) ValidateUser(ctx *gin.Context) {
	if !(ctx.FullPath() == "/users/login" || ctx.FullPath() == "/users/register") {
		fullToken := ctx.GetHeader("Authorization")

		if fullToken == "" {
			response.ErrorResponse(ctx, http.StatusUnauthorized, "token not found")
			ctx.Abort()
			return
		}

		tokenSplit := strings.Split(fullToken, " ")
		if len(tokenSplit) != 2 {
			response.ErrorResponse(ctx, http.StatusUnauthorized, "token invalid")
			ctx.Abort()
			return
		}

		if tokenSplit[0] != "Bearer" {
			response.ErrorResponse(ctx, http.StatusUnauthorized, "token type invalid")
			ctx.Abort()
			return
		}

		claims, err := m.jwtService.ValidateToken(tokenSplit[1])
		if err != nil {
			response.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

func NewMiddleware(jwtService service.JwtService) Middleware {
	return &middlewareImpl{jwtService: jwtService}
}
