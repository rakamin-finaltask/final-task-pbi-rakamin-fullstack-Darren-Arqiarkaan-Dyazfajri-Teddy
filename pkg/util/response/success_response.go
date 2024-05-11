package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-personalize/internal/model/dto"
)

func SuccessResponse(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}
