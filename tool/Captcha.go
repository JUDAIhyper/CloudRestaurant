package tool

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"image/color"
)

type CaptchaResult struct {
	Id          string `json:"id"`
	CaptchaType string
	VerifyValue string `json:"code"`
}

// 设置自带的store
var store = base64Captcha.DefaultMemStore

// 生成图形化验证码
func GenerateCaptcha(ctx *gin.Context) {
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverString

	// 配置验证码信息
	captchaConfig := base64Captcha.DriverString{
		Height:          30,
		Width:           60,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	// 自定义配置，如果不需要自定义配置，则上面的结构体和下面这行代码不用写
	driverString = captchaConfig
	driver = driverString.ConvertFonts()

	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := captcha.Generate()
	if err != nil {
		Failed(ctx, err.Error())
	}
	captchaResult := CaptchaResult{
		Id:         id,
		CaptchaType: b64s,
	}

	Success(ctx, gin.H{
		"captcha_result": captchaResult,
	})
}

func VertifyCaptcha(id string, value string) bool {
	if store.Verify(id, value, true) {
		return true
	}
	return false
}
