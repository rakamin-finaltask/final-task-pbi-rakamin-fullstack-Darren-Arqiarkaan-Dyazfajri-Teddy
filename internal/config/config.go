package config

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
	Driver   string
}

type ApiConfig struct {
	ApiPort string
}

type Config struct {
	DbConfig  DbConfig
	ApiConfig ApiConfig
	JwtConfig JwtConfig
}

type JwtConfig struct {
	JwtSecretKey     []byte
	JwtSigningMethod *jwt.SigningMethodHMAC
	JwtExpiredTime   time.Duration
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cfg.ConfigConfiguration()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) ConfigConfiguration() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error load .env file")
	}

	// config db
	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	// config server
	c.ApiConfig = ApiConfig{
		ApiPort: os.Getenv("API_PORT"),
	}

	// config jwt
	tokenExpired, err := strconv.Atoi(os.Getenv("JWT_EXPIRED_TIME"))
	if err != nil {
		return fmt.Errorf("config :" + err.Error())
	}

	c.JwtConfig = JwtConfig{
		JwtSecretKey:     []byte(os.Getenv("JWT_SECRET")),
		JwtExpiredTime:   time.Duration(tokenExpired) * time.Hour,
		JwtSigningMethod: jwt.SigningMethodHS256,
	}

	if c.DbConfig.Dbname == "" || c.DbConfig.Password == "" || c.DbConfig.Username == "" || c.DbConfig.Host == "" || c.DbConfig.Port == "" || c.DbConfig.Driver == "" || c.ApiConfig.ApiPort == "" || c.JwtConfig.JwtSecretKey == nil || c.JwtConfig.JwtSigningMethod == nil || c.JwtConfig.JwtExpiredTime == 0 {
		return fmt.Errorf("missing required environment variables")
	}

	return nil
}
