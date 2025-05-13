package main

import (
	"fmt"
	"github.com/star-find-cloud/star-mall/conf"
	appTLS "github.com/star-find-cloud/star-mall/pkg/tls"
)

func main() {
	//mysqlSlaves := db.GetMySQLSlaves()
	//mysqlMaster := pkg.GetMasterDB()
	//fmt.Println(mysqlSlaves)

	//err := mail.SendVerificationCode("3223590891@qq.com", mail.GenerateCode())
	//if err != nil {
	//	fmt.Println(err)
	//}

	var c = conf.GetConfig()
	fmt.Println(c.TLS)

	tls, err := appTLS.GetTLS()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tls)
}
