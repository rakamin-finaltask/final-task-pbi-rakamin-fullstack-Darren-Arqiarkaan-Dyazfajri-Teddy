package usecase

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"user-personalize/internal/model/dto"
	"user-personalize/internal/repository"
	"user-personalize/pkg/util/exception"
	"user-personalize/pkg/util/mapping"
	"user-personalize/pkg/util/service"
)

type AuthUC interface {
	Login(payload dto.LoginRequest) (dto.LoginResponse, error)
	Register(payload dto.UserRequest) (dto.UserResponse, error)
}

type authUCImpl struct {
	userRepository repository.UserRepository
	jwtService     service.JwtService
	validate       *validator.Validate
}

func NewAuthUC(userRepository repository.UserRepository, jwtService service.JwtService, validate *validator.Validate) AuthUC {
	return &authUCImpl{userRepository: userRepository, jwtService: jwtService, validate: validate}
}

func (a *authUCImpl) Login(payload dto.LoginRequest) (dto.LoginResponse, error) {
	err := a.validate.Struct(payload)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("validate payload failed: %w", err)
	}

	user, err := a.userRepository.GetByEmail(payload.Email)
	if err != nil {
		return dto.LoginResponse{}, exception.NotFoundErr
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("LoginUC : %w", err)
	}

	token, err := a.jwtService.GenerateToken(user.Id)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("LoginUC : %w", err)
	}

	return dto.LoginResponse{Token: *token}, nil
}

func (a *authUCImpl) Register(payload dto.UserRequest) (dto.UserResponse, error) {
	err := a.validate.Struct(payload)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("validate payload failed: %w", err)
	}

	_, err = a.userRepository.GetByEmail(payload.Email)
	if err == nil {
		return dto.UserResponse{}, exception.DuplicateErr
	}

	id := uuid.NewString()
	password, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("RegisterUC: %w", err)
	}

	entity := mapping.MapUserToEntity(payload)
	entity.Id = id
	entity.Password = string(password)

	user, err := a.userRepository.Create(entity)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("RegisterUC : %w", err)
	}

	return mapping.MapUserToResponse(user), nil
}
