package main

import (
	"fmt"
	"github.com/star-find-cloud/star-mall/pkg/mail"
)

func main() {
	//mysqlSlaves := db.GetMySQLSlaves()
	//mysqlMaster := pkg.GetMasterDB()
	//fmt.Println(mysqlSlaves)

	err := mail.SendVerificationCode("3223590891@qq.com", mail.GenerateCode())
	if err != nil {
		fmt.Println(err)
	}

	//captData, err := model.RotateCapt.Generate()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//dotData := captData.GetData()
	//if dotData == nil {
	//	log.Fatalln(">>>>> generate err")
	//}
	//
	//dots, _ := json.Marshal(dotData)
	//fmt.Println(">>>>> ", string(dots))
	//
	//var mBase64, tBase64 string
	//mBase64, err = captData.GetMasterImage().ToBase64()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//tBase64, err = captData.GetThumbImage().ToBase64()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(">>>>> ", mBase64)
	//fmt.Println(">>>>> ", tBase64)
}
