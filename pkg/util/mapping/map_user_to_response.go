package mapping

import (
	"user-personalize/internal/model/dto"
	"user-personalize/internal/model/entity"
)

func MapUserToResponse(user entity.User) dto.UserResponse {
	return dto.UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}
