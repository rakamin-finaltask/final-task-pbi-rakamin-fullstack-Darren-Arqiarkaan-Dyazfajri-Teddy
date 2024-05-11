package service

import (
	"github.com/dgrijalva/jwt-go"
	"user-personalize/internal/config"
	"user-personalize/internal/model/dto"
)

type JwtService interface {
	GenerateToken(id string) (*string, error)
	ValidateToken(token string) (*dto.CustomClaims, error)
}

type jwtServiceImpl struct {
	cfg config.JwtConfig
}

func NewJwtService(cfg config.JwtConfig) JwtService {
	return &jwtServiceImpl{cfg: cfg}
}

func (j *jwtServiceImpl) GenerateToken(id string) (*string, error) {
	claims := &dto.CustomClaims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(j.cfg.JwtExpiredTime),
		},
	}

	token := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
	resultToken, err := token.SignedString(j.cfg.JwtSecretKey)
	if err != nil {
		return nil, err
	}

	return &resultToken, nil
}

func (j *jwtServiceImpl) ValidateToken(token string) (*dto.CustomClaims, error) {
	claims := &dto.CustomClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.cfg.JwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
