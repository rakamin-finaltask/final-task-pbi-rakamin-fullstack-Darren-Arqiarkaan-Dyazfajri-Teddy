package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-personalize/internal/model/dto"
)

func CreatedResponse(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusCreated, dto.WebResponse{
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	})
}
