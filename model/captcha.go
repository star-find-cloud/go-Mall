package model

import (
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/rotate"
)

// 旋转验证码对象
var RotateCapt rotate.Captcha

// 初始化对象
func init() {
	builder := rotate.NewBuilder(rotate.WithRangeAnglePos([]option.RangeVal{
		{Min: 50, Max: 350},
	}))

	// 背景图
	imgs, err := images.GetImages()
	if err != nil {
		log.AppLogger.Errorf("验证码图片获取失败: %v", err)
	}

	// 设置验证码图片
	builder.SetResources(rotate.WithImages(imgs))
	RotateCapt = builder.Make()
}
