package delivery

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"user-personalize/internal/config"
	"user-personalize/internal/delivery/controller"
	"user-personalize/internal/delivery/middleware"
	"user-personalize/internal/repository"
	"user-personalize/internal/usecase"
	"user-personalize/pkg/util/service"
)

type Server struct {
	UserUC     usecase.UserUC
	AuthUC     usecase.AuthUC
	PhotoUC    usecase.PhotosUC
	Middleware middleware.Middleware
	Host       string
	Engine     *gin.Engine
}

func (s *Server) ServerRun() {
	s.Engine.Use(s.Middleware.ValidateUser)
	s.InitRoute()
	err := s.Engine.Run(s.Host)
	if err != nil {
		panic(err)
	}
}

func (s *Server) InitRoute() {
	rg := s.Engine.Group("")
	controller.NewUserController(s.UserUC, rg).RouteGroup()
	controller.NewAuthController(s.AuthUC, rg).RouteGroup()
	controller.NewPhotosController(s.PhotoUC, rg).RouteGroup()
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbConfig.Host, cfg.DbConfig.Port, cfg.DbConfig.Username, cfg.DbConfig.Password, cfg.DbConfig.Dbname)

	db, err := sql.Open(cfg.DbConfig.Driver, dsn)
	if err != nil {
		panic(err)
	}

	validate := validator.New()

	// repository
	userRepository := repository.NewUserRepository(db)
	photosRepository := repository.NewPhotosRepository(db)

	// UC
	jwtService := service.NewJwtService(cfg.JwtConfig)

	userUC := usecase.NewUserUC(userRepository, validate)
	authUC := usecase.NewAuthUC(userRepository, jwtService, validate)
	photosUC := usecase.NewPhotosUC(photosRepository)

	newMiddleware := middleware.NewMiddleware(jwtService)

	engine := gin.Default()

	host := fmt.Sprintf(":%s", cfg.ApiConfig.ApiPort)

	return &Server{
		Host:       host,
		Engine:     engine,
		UserUC:     userUC,
		AuthUC:     authUC,
		Middleware: newMiddleware,
		PhotoUC:    photosUC,
	}
}
