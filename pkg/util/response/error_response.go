package response

import (
	"github.com/gin-gonic/gin"
	"user-personalize/internal/model/dto"
)

func ErrorResponse(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, dto.WebResponse{
		Code:    code,
		Message: message,
	})
}
