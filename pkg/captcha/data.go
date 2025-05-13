package captcha

import (
	"encoding/json"
	"fmt"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
)

func GetData() (string, string, error) {
	captData, err := RotateCapt.Generate()
	if err != nil {
		applog.AppLogger.Errorln(err)
		return "", "", err
	}

	dotData := captData.GetData()
	if dotData == nil {
		applog.AppLogger.Errorln(err)
		return "", "", err
	}

	dots, _ := json.Marshal(dotData)
	fmt.Println(">>>>> ", string(dots))

	var mBase64, tBase64 string
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		applog.AppLogger.Errorf("验证码主图转换失败: %v \n", err)
		return "", "", err
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	if err != nil {
		applog.AppLogger.Errorf("验证码副图转换失败: %v \n", err)
		return "", "", err
	}

	//fmt.Println(">>>>> ", mBase64)
	//fmt.Println(">>>>> ", tBase64)
	return mBase64, tBase64, nil
}
