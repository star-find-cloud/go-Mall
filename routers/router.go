package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/star-find-cloud/star-mall/handler"
	"github.com/star-find-cloud/star-mall/middleware"
)

func InitRouter(userHandler *handler.UserHandler,
	imageHandler *handler.ImageHandler,
	merchantHandler *handler.MerchantHandler,
	product *handler.ProductHandler,
	inventoryHandler *handler.InventoryHandler,
	public *handler.PublicHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	deepseekHandler *handler.DeepseekHandler) *gin.Engine {
	// 设置 gin 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(
		middleware.GinLogger(),
		middleware.GinRecoveryWithZap(true),
	)

	publicGroup := r.Group("/api/v1")
	{
		publicGroup.GET("/health", public.HealthCheck)
		publicGroup.POST("/sendVerifyCodeByEmail", public.SendVerifyCodeByEmail)
	}

	userGroup := r.Group("/api/v1/user")
	{
		userGroup.POST("/login", userHandler.Login)
		userGroup.PUT("/register", userHandler.Register)
		userGroup.DELETE("/delete", userHandler.Delete)
		userGroup.PATCH("/forgetPassword", userHandler.ForgetPassword)
	}
	userGroup.Use(middleware.JwtAuth())
	{
		userGroup.GET("/:id", userHandler.GetUser)
		userGroup.PATCH("/update", userHandler.Update)
		userGroup.PATCH("/update/password", userHandler.UpdatePassword)
	}

	imageGroup := r.Group("/api/v1/image")
	{
		imageGroup.GET("/getImage", imageHandler.GetImage)
	}
	imageGroup.Use(middleware.JwtAuth())
	{
		imageGroup.POST("/upload", imageHandler.UploadImage)
		//imageGroup.POST("/:owner_type/:id/images/upload", imageHandler.)

	}

	// 商家相关路由组
	merchantGroup := r.Group("/api/v1/merchants")
	{
		// 获取单个商家信息
		merchantGroup.GET("/getMerchant/:merchantID", merchantHandler.GetMerchant)
		// 搜索商家
		merchantGroup.GET("/search/:merchant_name", merchantHandler.Select)
		// 商家注册
		merchantGroup.POST("/register", merchantHandler.Register)
		// 更新商家信息
		merchantGroup.PATCH("/update", merchantHandler.Update)
		// 删除商家
		merchantGroup.DELETE("/delete/:merchantID", merchantHandler.Delete)
	}

	// 商品相关路由组
	productGroup := r.Group("/api/v1/products")
	{
		// 获取单个商品信息
		productGroup.GET("/getProduct/:id", product.GetProduct)
	}
	productGroup.Use(middleware.JwtAuth())
	{
		// 创建单个商品
		productGroup.PUT("/create", product.Create)
		// 搜索商品
		productGroup.GET("/search/:MerchantID", product.GetProductByMerchant)
		// 更新商品信息
		productGroup.PATCH("/update", product.Update)
		//
		productGroup.POST("/search", product.SearchProduct)
	}

	// 库存相关路由组
	inventoryGroup := r.Group("/api/v1/inventory")
	{
		// 创建单个库存
		inventoryGroup.PUT("/create", inventoryHandler.Create)
		// 获取单个库存信息
		inventoryGroup.GET("inventory/{id}", inventoryHandler.GetInventory)
		//// 搜索库存
		inventoryGroup.GET("/search/:MerchantID", inventoryHandler.GetInventoryByMerchant)
		// 更新库存信息
		inventoryGroup.PATCH("/update", inventoryHandler.Update)
	}

	// 购物车相关路由组
	cartGroup := r.Group("/api/v1/cart")
	{
		// 获取单个购物车信息
		cartGroup.GET("/{id}", cartHandler.GetByID)
	}
	cartGroup.Use(middleware.JwtAuth())
	{
		// 创建单个购物车
		cartGroup.PUT("/create", cartHandler.Create)
		// 获取单个用户购物车信息
		cartGroup.GET("/user", cartHandler.GetByUserID)
		// 添加商品到购物车
		cartGroup.PUT("/add", cartHandler.AddProduct)
	}

	// 订单相关路由组
	orderGroup := r.Group("/api/v1/order")
	{
		// 创建单个订单
		orderGroup.PUT("/create", orderHandler.CreateOrder)
		// 获取单个订单信息
		orderGroup.GET("/get", orderHandler.GetOrder)
		//// 获取单个用户订单信息
		//orderGroup.GET("user/{id}", orderHandler.GetByUserID)
		//// 添加商品到订单
		//orderGroup.PUT("/add", orderHandler.AddProduct)
	}

	// AI相关路由组
	aiGroup := r.Group("/api/v1/deepseek")
	aiGroup.Use(middleware.JwtAuth())
	{
		aiGroup.POST("/suggest", deepseekHandler.SuggestProductByUserTags)
		aiGroup.POST("/search", deepseekHandler.SuggestProductBySearch)
	}

	return r
}
