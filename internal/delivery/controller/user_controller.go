package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"user-personalize/internal/model/dto"
	"user-personalize/internal/usecase"
	response2 "user-personalize/pkg/util/response"
)

type UserController struct {
	userUC usecase.UserUC
	rg     *gin.RouterGroup
}

func NewUserController(userUC usecase.UserUC, rg *gin.RouterGroup) *UserController {
	return &UserController{userUC: userUC, rg: rg}
}

func (u *UserController) RouteGroup() {
	u.rg.POST("/users", u.CreateUser)
	u.rg.GET("/users/:userId", u.GetUserById)
	u.rg.GET("/users", u.GetListUser)
	u.rg.PUT("/users/:userId", u.updateUser)
	u.rg.DELETE("/users/:userId", u.DeleteUser)
	u.rg.PUT("/users/updatePassword/:userId", u.updatePassword)
}

func (u *UserController) CreateUser(ctx *gin.Context) {
	userReq := dto.UserRequest{}

	err := ctx.BindJSON(&userReq)
	if err != nil {
		log.Println(err)
		response2.ErrorResponse(ctx, http.StatusBadRequest, "request body is invalid")
		return
	}

	user, err := u.userUC.CreateUser(userReq)
	if err != nil {
		log.Println(err)
		response2.ErrorResponse(ctx, http.StatusInternalServerError, "error while creating user")
		return
	}

	response2.CreatedResponse(ctx, "success create new user", user)
}

func (u *UserController) GetListUser(ctx *gin.Context) {
	user, err := u.userUC.GetAllUser()
	if err != nil {
		log.Println(err)
		response2.ErrorResponse(ctx, http.StatusInternalServerError, "error while getting user")
		return
	}

	response2.SuccessResponse(ctx, "success get all user", user)
}

func (u *UserController) GetUserById(ctx *gin.Context) {
	id := ctx.Param("userId")

	user, err := u.userUC.GetUserById(id)
	if err != nil {
		log.Println(err)
		response2.ErrorResponse(ctx, http.StatusInternalServerError, "error while getting user")
		return
	}

	response2.SuccessResponse(ctx, "success get user by id", user)
}

func (u *UserController) updateUser(ctx *gin.Context) {
	id := ctx.Param("userId")

	userReq := dto.UserUpdateRequest{}
	err := ctx.BindJSON(&userReq)
	if err != nil {
		log.Println(err)
		response2.ErrorResponse(ctx, http.StatusBadRequest, "request body is not valid")
		return
	}

	updatedUser, err := u.userUC.Update(id, userReq)
	if err != nil {
		log.Println(err)
		response2.ErrorResponse(ctx, http.StatusInternalServerError, "error while updating user")
		return
	}

	response2.SuccessResponse(ctx, "success updatedUser user", updatedUser)
}

func (u *UserController) updatePassword(ctx *gin.Context) {
	id := ctx.Param("userId")
	userReq := dto.UpdatePasswordRequest{}

	err := ctx.BindJSON(&userReq)

	if err != nil {
		log.Println(err)
		response2.ErrorResponse(ctx, http.StatusBadRequest, "request not valid")
		return
	}

	updatedUser, err := u.userUC.UpdatePassword(id, userReq)
	if err != nil {
		log.Println(err)
		response2.ErrorResponse(ctx, http.StatusInternalServerError, "error while updating user")
		return
	}

	response2.SuccessResponse(ctx, "success update password user", updatedUser)
}

func (u *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("userId")

	err := u.userUC.DeleteUser(id)
	if err != nil {
		log.Println(err)
		response2.ErrorResponse(ctx, http.StatusInternalServerError, "error while deleting user")
		return
	}

	response2.SuccessResponse(ctx, "success delete user", id)
}
