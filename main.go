package main

import (
	"fmt"
	"github.com/star-find-cloud/star-mall/handler"
	"github.com/star-find-cloud/star-mall/internal/repo"
	"github.com/star-find-cloud/star-mall/internal/service"
	"github.com/star-find-cloud/star-mall/pkg/database"
	"github.com/star-find-cloud/star-mall/pkg/oss"
	"github.com/star-find-cloud/star-mall/routers"
)

func main() {
	fmt.Println("程序启动中...")
	db := database.GetDB()
	ossClient := oss.GetCos()

	fmt.Println("配置读取完成")
	userRepo := repo.NewUserRepositoryImpl(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	fmt.Println("配置读取完成")
	imageRepo := repo.NewImageRepositoryImpl(db)
	imageService := service.NewOSSService(ossClient, imageRepo)
	imageHandler := handler.NewImageHandler(imageService)

	fmt.Println("配置读取完成")
	var r = routers.InitRouter(userHandler, imageHandler)

	fmt.Println("gin 配置完成")
	err := r.Run(":9090") // 修改此行绑定到9090端口

	fmt.Println("gin 启动")
	if err != nil {
		panic(err)
	}

	fmt.Println("程序启动成功")
}
