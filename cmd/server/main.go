package main

import (
	"log"

	"github.com/anandabhimanyu/employee-backend-api/internal/auth"
	"github.com/anandabhimanyu/employee-backend-api/internal/config"
	"github.com/anandabhimanyu/employee-backend-api/internal/db"
	"github.com/anandabhimanyu/employee-backend-api/internal/employee"
	"github.com/anandabhimanyu/employee-backend-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	postgres, err := db.NewPostgres(
		cfg.DBHost, cfg.DBPort,
		cfg.DBUser, cfg.DBPass, cfg.DBName,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.DB.Close()

	r := gin.Default()

	// AUTH
	authRepo := auth.NewRepository(postgres.DB)
	authService := auth.NewService(authRepo)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)
	authHandler := auth.NewHandler(authService, jwtManager)

	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// EMPLOYEES
	empRepo := employee.NewRepository(postgres.DB)
	empHandler := employee.NewHandler(empRepo)

	emp := r.Group("/api/v1/employees")
	emp.Use(middleware.JWTAuth(jwtManager))
	{
		emp.POST("", empHandler.Create)
		emp.GET("", empHandler.List)
		emp.GET("/:id", empHandler.GetByID)
		emp.PUT("/:id", empHandler.Update)
		emp.DELETE("/:id", empHandler.Delete)
	}

	log.Println("Server running on port", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
