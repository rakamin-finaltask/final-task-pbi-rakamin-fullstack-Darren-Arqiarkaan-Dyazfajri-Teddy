package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"user-personalize/internal/model/dto"
	"user-personalize/internal/usecase"
	"user-personalize/pkg/util/exception"
	"user-personalize/pkg/util/response"
)

type AuthController struct {
	authUC usecase.AuthUC
	rg     *gin.RouterGroup
}

func NewAuthController(authUC usecase.AuthUC, rg *gin.RouterGroup) *AuthController {
	return &AuthController{authUC: authUC, rg: rg}
}

func (a *AuthController) RouteGroup() {
	a.rg.POST("/users/login", a.Login)
	a.rg.POST("/users/register", a.Register)
}

func (a *AuthController) Login(ctx *gin.Context) {
	var request dto.LoginRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		log.Println(err)
		response.ErrorResponse(ctx, http.StatusBadRequest, "request invalid")
		return
	}

	loginRes, err := a.authUC.Login(request)
	if err != nil {
		log.Println(err)
		if errors.Is(err, exception.NotFoundErr) {
			response.ErrorResponse(ctx, http.StatusNotFound, "user not found")
			return
		}
		response.ErrorResponse(ctx, http.StatusUnauthorized, "login failed")
		return
	}

	webResponse := dto.WebResponse{
		Code:    http.StatusAccepted,
		Message: "login success",
		Data:    loginRes,
	}

	ctx.JSON(http.StatusAccepted, webResponse)
}

func (a *AuthController) Register(ctx *gin.Context) {
	var req dto.UserRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Println(err)
		response.ErrorResponse(ctx, http.StatusBadRequest, "request invalid")
		return
	}

	register, err := a.authUC.Register(req)
	if err != nil {
		log.Println(err)

		if errors.Is(err, exception.DuplicateErr) {
			response.ErrorResponse(ctx, http.StatusConflict, "user already exists")
			return
		}
		response.ErrorResponse(ctx, http.StatusInternalServerError, "error while register user")
		return
	}

	webResponse := dto.WebResponse{
		Code:    http.StatusCreated,
		Message: "user created",
		Data:    register,
	}

	ctx.JSON(http.StatusCreated, webResponse)
}
