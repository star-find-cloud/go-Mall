package main

import (
	"github.com/star-find-cloud/star-mall/handler"
	"github.com/star-find-cloud/star-mall/internal/repo"
	"github.com/star-find-cloud/star-mall/internal/service"
	"github.com/star-find-cloud/star-mall/pkg/database"
	"github.com/star-find-cloud/star-mall/routers"
)

func main() {
	db := database.GetDB()

	userRepo := repo.NewUserRepositoryImpl(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	var r = routers.InitRouter(userHandler)
	err := r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}
