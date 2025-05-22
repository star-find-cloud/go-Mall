package main

import (
	"fmt"
	"github.com/star-find-cloud/star-mall/handler"
	"github.com/star-find-cloud/star-mall/pkg/database"
	"github.com/star-find-cloud/star-mall/pkg/oss"
	repo2 "github.com/star-find-cloud/star-mall/repo"
	"github.com/star-find-cloud/star-mall/routers"
	service2 "github.com/star-find-cloud/star-mall/service"
)

func main() {
	fmt.Println("程序启动中...")
	db := database.GetDB()
	ossClient := oss.GetCos()

	fmt.Println("配置读取完成")
	userRepo := repo2.NewUserRepositoryImpl(db)
	userService := service2.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	fmt.Println("配置读取完成")
	imageRepo := repo2.NewImageRepositoryImpl(db)
	imageService := service2.NewOSSService(ossClient, imageRepo)
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
