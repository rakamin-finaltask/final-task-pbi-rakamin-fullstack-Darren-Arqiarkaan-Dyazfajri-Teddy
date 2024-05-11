package mapping

import (
	"user-personalize/internal/model/dto"
	"user-personalize/internal/model/entity"
)

func MapUserToEntity(request dto.UserRequest) entity.User {
	return entity.User{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
	}
}
