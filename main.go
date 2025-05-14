package main

import (
	"fmt"
	"github.com/star-find-cloud/star-mall/conf"
)

func main() {
	//r := gin.Default()
	//r.GET("/", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "欢迎使用star mall",
	//	})
	//})
	//
	//authGroup := r.Group("/login")
	//{
	//	authGroup.POST("/login/id", loginByID)
	//	authGroup.POST("/login/email", loginByEmail)
	//}
	//
	//err := r.Run(":8080")
	//if err != nil {
	//	fmt.Errorf("app error: %v", err)
	//	return
	//}
	c := conf.GetConfig()
	fmt.Println(c)
}
