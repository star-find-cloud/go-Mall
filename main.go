package main

import (
	"fmt"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/protobuf/pb"
	"google.golang.org/grpc"

	"github.com/star-find-cloud/star-mall/handler"
	ds "github.com/star-find-cloud/star-mall/internal/deepseek"
	"github.com/star-find-cloud/star-mall/internal/mq"
	"github.com/star-find-cloud/star-mall/pkg/database"
	"github.com/star-find-cloud/star-mall/pkg/oss"
	"github.com/star-find-cloud/star-mall/repo"
	"github.com/star-find-cloud/star-mall/routers"
	"github.com/star-find-cloud/star-mall/service"
	"github.com/star-find-cloud/star-mall/utils"
)

// @title           Star Mall API
// @version         1.0
// @description     This is the API documentation for Star Mall service.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9090
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description 输入 JWT token，格式为：Bearer <token>

func main() {
	var (
		mqProducer     mq.Producer
		mqConsumer     mq.Consumer
		db, err        = database.NewMySQL()
		ossClient      = oss.NewTencentCos()
		cache, err     = database.NewRedis()
		deepseekClient = ds.NewDeepseekClient()
	)
	if err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		log.AppLogger.Fatalf("初始化失败: %v\n", err)
		panic(err)
	}

	imageRepo := repo.NewImageRepo(db)
	var imageService service.ImageMetaDataService
	if utils.IsEnableOSS() {
		imageService = service.NewImageForOssService(imageRepo, ossClient, cache)
	} else {
		imageService = service.NewImageForDBService(imageRepo)
	}
	imageHandler := handler.NewImageHandler(imageService)

	userRepo := repo.NewUserRepo(db, cache)
	userService := service.NewUserService(userRepo, imageRepo)
	userHandler := handler.NewUserHandler(userService)

	// 初始化商家相关组件
	merchantRepo := repo.NewMerchantRepo(db, cache)
	merchantService := service.NewMerchantService(merchantRepo, imageRepo)
	merchantHandler := handler.NewMerchantHandler(merchantService, userService)

	// 初始化商品相关组件
	productRepo := repo.NewProductRepo(db, cache)
	productService := service.NewProductService(productRepo, ossClient, imageRepo)
	productHandler := handler.NewProductHandler(productService)

	// 初始化库存相关组件
	inventoryRepo := repo.NewInventoryRepo(db, cache)
	inventoryService := service.NewInventoryService(inventoryRepo, merchantRepo, productRepo)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	// 初始化公共组件
	publicRepo := repo.NewPublicRepo(db, cache)
	publicService := service.NewPublicService(publicRepo)
	publicHandler := handler.NewPublicHandler(publicService)

	// 初始化购物车相关组件
	cartRepo := repo.NewCartRepo(db, cache)
	cartService := service.NewCartService(cartRepo, productRepo)
	cartHandler := handler.NewCartHandler(cartService)

	// 初始化订单相关组件
	orderRepo := repo.NewOrderRepo(db, cache)
	orderService := service.NewOrderService(orderRepo, productRepo, userRepo, inventoryRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	// 初始化 deepseek 相关组件
	deepseekService := service.NewDeepseekService(deepseekClient, userRepo, productRepo)
	deepseekHandler := handler.NewDeepseekHandler(deepseekService)

	//grpcServer := grpc.NewServer()
	//pb.RegisterImageServiceServer(grpcServer, imageService)

	fmt.Println("配置读取完成")
	var r = routers.InitRouter(userHandler, imageHandler, merchantHandler, productHandler, inventoryHandler, publicHandler, cartHandler, orderHandler, deepseekHandler)

	fmt.Println("gin 配置完成")
	fmt.Println("正在启动服务器...")

	if err := r.Run(":8080"); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		panic(err)
	}
}
