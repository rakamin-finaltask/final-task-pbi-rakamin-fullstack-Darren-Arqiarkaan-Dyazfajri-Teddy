package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"user-personalize/internal/model/dto"
	"user-personalize/internal/model/entity"
	"user-personalize/internal/usecase"
	"user-personalize/pkg/util/exception"
	"user-personalize/pkg/util/response"
)

type PhotosController struct {
	photoUC usecase.PhotosUC
	rg      *gin.RouterGroup
}

func NewPhotosController(photoUC usecase.PhotosUC, rg *gin.RouterGroup) *PhotosController {
	return &PhotosController{photoUC: photoUC, rg: rg}
}

func (p *PhotosController) RouteGroup() {
	p.rg.POST("/photos", p.UploadPhotos)
	p.rg.PUT("/photos/:photoId", p.UpdatePhotos)
	p.rg.DELETE("/photos/:photoId", p.DeletePhotos)
	p.rg.GET("photos", p.GetPhotos)
}

func (p *PhotosController) UploadPhotos(ctx *gin.Context) {
	value, exists := ctx.Get("claims")
	if !exists {
		log.Println("claims not exists")
		response.ErrorResponse(ctx, http.StatusForbidden, "you cannot access this resource")
		return
	}

	userId := value.(*dto.CustomClaims).UserId

	file, err := ctx.FormFile("photo-profile")
	if err != nil {
		log.Println(err)
		response.ErrorResponse(ctx, http.StatusBadRequest, "input is not valid")
		return
	}

	title := ctx.Request.FormValue("title")
	caption := ctx.Request.FormValue("caption")

	id := uuid.NewString()
	name := strings.Replace(id, "-", "", -1)
	ext := filepath.Ext(file.Filename)
	imageName := fmt.Sprintf("%s%s", name, ext)
	imageUrl := fmt.Sprintf("./images/%s", imageName)

	request := entity.Photos{
		Title:    title,
		Caption:  caption,
		PhotoUrl: imageUrl,
		UserId:   userId,
	}

	_, err = p.photoUC.GetPhotosByUserId(userId)
	if err == nil {
		log.Println("photo already exists")
		response.ErrorResponse(ctx, http.StatusConflict, "photo already exists")
		return
	}

	photos, err := p.photoUC.SavePhotos(request)
	if err != nil {
		log.Println(err)
		response.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	err = ctx.SaveUploadedFile(file, imageUrl)
	if err != nil {
		log.Println(err)
		response.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}
	response.SuccessResponse(ctx, "success upload photos", photos)
}

func (p *PhotosController) UpdatePhotos(ctx *gin.Context) {

	value, exists := ctx.Get("claims")
	if !exists {
		log.Println("claims not exists")
		response.ErrorResponse(ctx, http.StatusForbidden, "you cannot access this resource")
		return
	}

	userId := value.(*dto.CustomClaims).UserId

	photoId := ctx.Param("photoId")
	title := ctx.Request.FormValue("title")
	caption := ctx.Request.FormValue("caption")

	file, err := ctx.FormFile("photo-profile")
	if err != nil {
		log.Println(err)
		response.ErrorResponse(ctx, http.StatusBadRequest, "input is not valid")
		return
	}

	id := uuid.NewString()
	name := strings.Replace(id, "-", "", -1)
	ext := filepath.Ext(file.Filename)
	imageName := fmt.Sprintf("%s%s", name, ext)

	imageUrl := fmt.Sprintf("./images/%s", imageName)
	request := entity.Photos{
		Id:       photoId,
		Title:    title,
		Caption:  caption,
		PhotoUrl: imageUrl,
		UserId:   userId,
	}

	photos, err := p.photoUC.UpdatePhotos(request)
	if err != nil {
		log.Println(err)

		if errors.Is(err, exception.NotFoundErr) {
			response.ErrorResponse(ctx, http.StatusNotFound, "photo not found")
			return
		}

		response.ErrorResponse(ctx, http.StatusInternalServerError, "failed update photo")
		return
	}

	err = ctx.SaveUploadedFile(file, imageUrl)
	if err != nil {
		log.Println(err)
		response.ErrorResponse(ctx, http.StatusInternalServerError, "failed update photo")
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: "success update photos",
		Data:    photos,
	})
}

func (p *PhotosController) DeletePhotos(ctx *gin.Context) {
	value, exists := ctx.Get("claims")
	if !exists {
		log.Println("claims does not exist")
		response.ErrorResponse(ctx, http.StatusForbidden, "you cannot access this resource")
		return
	}

	userId := value.(*dto.CustomClaims).UserId

	photoId := ctx.Param("photoId")

	request := entity.Photos{
		Id:     photoId,
		UserId: userId,
	}

	err := p.photoUC.DeletePhotos(request)
	if err != nil {
		log.Println(err)
		if errors.Is(err, exception.NotFoundErr) {
			response.ErrorResponse(ctx, http.StatusNotFound, "photo not found")
			return
		}
		response.ErrorResponse(ctx, http.StatusInternalServerError, "failed delete photo")
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: "success delete photos",
	})
}

func (p *PhotosController) GetPhotos(ctx *gin.Context) {
	value, exists := ctx.Get("claims")
	if !exists {
		log.Println("claims does not exist")
		response.ErrorResponse(ctx, http.StatusForbidden, "you cannot access this resource")
		return
	}

	userId := value.(*dto.CustomClaims).UserId

	photos, err := p.photoUC.GetPhotosByUserId(userId)
	if err != nil {
		log.Println(err)
		if errors.Is(err, exception.NotFoundErr) {
			response.ErrorResponse(ctx, http.StatusNotFound, "photos not found")
			return
		}

		response.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: "success get photos",
		Data:    photos,
	})
}
