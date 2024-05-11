package usecase

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"user-personalize/internal/model/dto"
	"user-personalize/internal/model/entity"
	"user-personalize/internal/repository"
	mapping2 "user-personalize/pkg/util/mapping"
)

type UserUC interface {
	CreateUser(payload dto.UserRequest) (dto.UserResponse, error)
	GetAllUser() ([]dto.UserResponse, error)
	GetUserById(userId string) (dto.UserResponse, error)
	Update(id string, payload dto.UserUpdateRequest) (dto.UserResponse, error)
	UpdatePassword(id string, payload dto.UpdatePasswordRequest) (dto.UserResponse, error)
	DeleteUser(id string) error
}

type userUCImpl struct {
	userRepository repository.UserRepository
	validate       *validator.Validate
}

func NewUserUC(userRepository repository.UserRepository, validate *validator.Validate) UserUC {
	return &userUCImpl{userRepository: userRepository, validate: validate}
}

func (u *userUCImpl) CreateUser(payload dto.UserRequest) (dto.UserResponse, error) {
	err := u.validate.Struct(payload)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("CreateUserUC : %w", err)
	}

	user := mapping2.MapUserToEntity(payload)

	user.Id = uuid.NewString()

	password := []byte(payload.Password)
	passwordHashed, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("CreateUserUC : %w", err)
	}

	user.Password = string(passwordHashed)

	userCreated, err := u.userRepository.Create(user)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("CreateUserUC : %w", err)
	}

	return mapping2.MapUserToResponse(userCreated), nil
}

func (u *userUCImpl) GetAllUser() ([]dto.UserResponse, error) {
	users, err := u.userRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("GetAllUserUC : %w", err)
	}

	result := make([]dto.UserResponse, len(users))

	for i, user := range users {
		result[i] = mapping2.MapUserToResponse(user)
	}

	return result, nil
}

func (u *userUCImpl) GetUserById(userId string) (dto.UserResponse, error) {
	user, err := u.userRepository.GetById(userId)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("GetUserByIdUC : %w", err)
	}

	return mapping2.MapUserToResponse(user), nil
}

func (u *userUCImpl) Update(id string, payload dto.UserUpdateRequest) (dto.UserResponse, error) {
	err := u.validate.Struct(payload)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("UpdateUC : %w", err)
	}

	user := entity.User{
		Id:       id,
		Username: payload.Username,
		Email:    payload.Email,
	}

	updatedUser, err := u.userRepository.Update(user)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("UpdateUserUC : %w", err)
	}

	return mapping2.MapUserToResponse(updatedUser), nil
}

func (u *userUCImpl) UpdatePassword(id string, payload dto.UpdatePasswordRequest) (dto.UserResponse, error) {
	err := u.validate.Struct(payload)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("UpdateUserPasswordUC : %w", err)
	}

	user, err := u.userRepository.UpdatePassword(id, payload.Password)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("UpdateUserPasswordUC : %w", err)
	}

	return mapping2.MapUserToResponse(user), nil
}

func (u *userUCImpl) DeleteUser(id string) error {
	_, err := u.GetUserById(id)
	if err != nil {
		return err
	}

	err = u.userRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("DeleteUserUC : %w", err)
	}

	return nil
}
