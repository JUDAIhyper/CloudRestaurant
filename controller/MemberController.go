package controller

import (
	"CloudRestaurant/model"
	"CloudRestaurant/param"
	"CloudRestaurant/service"
	"CloudRestaurant/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type MemberController struct {
}

func (mc *MemberController) Router(engine *gin.Engine) {
	engine.GET("/api/sendcode", mc.sendSmsCode)
	engine.POST("/api/smsLogin", mc.smsLogin)
	engine.GET("/api/captcha", mc.captcha)
	//postman测试
	engine.POST("api/vertifycha", mc.vertifyCaptcha)
	//login_pwd
	engine.POST("/api/login_pwd", mc.nameLogin)
	//头像上传
	engine.POST("/api/upload/avator", mc.uploadAvator)
	//用户信息查询
	engine.GET("api/userinfo", mc.userInfo)
}

//查询用户信息
func (mc *MemberController) userInfo(context *gin.Context) {
	cookie, err := tool.CookieAuth(context)
	if err != nil {
		context.Abort()
		tool.Failed(context, "还未登录，请先登录")
		return
	}
	memberService := service.MemberService{}
	member := memberService.GetUserInfo(cookie.Value)
	if member != nil {
		//返回成功
		tool.Success(context, map[string]interface{}{
			"id":            member.Id,
			"user_name":     member.UserName,
			"mobile":        member.Mobile,
			"register_time": member.RegisterTime,
			"avatar":        member.Avatar,
			"balance":       member.Balance,
			"city":          member.City,
		})
	}
	tool.Failed(context, "获取用户信息失败")
}

//头像上传
func (mc *MemberController) uploadAvator(context *gin.Context) {
	//1.解析上传的参数: file, user_id
	userId := context.PostForm("user_id") //用户id
	fmt.Println(userId)
	file, err := context.FormFile("avatar")
	if err != nil || userId == "" {
		tool.Failed(context, "参数解析失败")
		return
	}
	//2.判断user_id对应的用户是否已经登录
	//sess := tool.GetSess(context, "user_"+userId)
	//if sess == nil {
	//	tool.Failed(context, "参数不合法")
	//	return
	//}
	tool.GetSession("user_" + userId)
	var member model.Member
	//json.Unmarshal(sess.([]byte), &member)

	//3.file 保存到本地
	fileName := "./uploadfile/" + strconv.FormatInt(time.Now().Unix(), 10) + file.Filename
	err = context.SaveUploadedFile(file, fileName)
	if err != nil {
		tool.Failed(context, "头像更新失败")
		return
	}

	//http://localhost:8080/static/.../xxx.png

	//4.将保存后的文件本地路径 保存到用户表中的头像字段
	memberService := service.MemberService{}
	path := memberService.UploadAvatar(member.Id, fileName[1:])
	if path != "" {
		tool.Success(context, "http://localhost:9000"+path)
		return
	}
	//5.返回结果
	tool.Failed(context, "上传失败")
}

//第2种登录方法：用户名密码登录
func (mc *MemberController) nameLogin(context *gin.Context) {
	//1.解析用户登录传递参数
	var loginParam param.LoginParam
	err := tool.Decode(context.Request.Body, &loginParam)
	if err != nil {
		tool.Failed(context, "参数解析失败")
		return
	}
	//2.验证验证码
	validate := tool.VertifyCaptcha(loginParam.Id, loginParam.Value)
	if !validate {
		tool.Failed(context, "验证码不正确，请重新验证")
		return
	}
	//3.登录
	ms := service.MemberService{}
	member := ms.Login(loginParam.Name, loginParam.Password)
	if member.Id != 0 {
		//将用户信息保存到session当中
		//序列化
		//sess, _ := json.Marshal(member)
		//err := tool.SetSess(context, "user_"+string(member.Id), sess)
		//if err != nil {
		//	tool.Failed(context, "登录失败")
		//	return
		//}
		sessionValue := map[string]interface{}{
			"user_name": member.UserName,
			"password":  member.Password,
		}
		tool.SetSession("user_"+string(member.Id), sessionValue)
		//调用cookie储存登录信息
		context.SetCookie("cookie_user", strconv.Itoa(int(member.Id)),
			10*60, "/", "192.168.17.1", true, true)
		tool.Success(context, member)
		return
	}
	tool.Failed(context, "登录失败")
}

//生成图形验证码
func (mc *MemberController) captcha(context *gin.Context) {
	//生成验证码并返回客户端
	tool.GenerateCaptcha(context)
}

//验证验证码是否正确
func (mc *MemberController) vertifyCaptcha(context *gin.Context) {
	var captcha tool.CaptchaResult
	err := tool.Decode(context.Request.Body, &captcha)
	if err != nil {
		tool.Failed(context, "参数解析失败")
		return
	}
	result := tool.VertifyCaptcha(captcha.Id, captcha.VerifyValue)
	if result {
		fmt.Println("验证通过")
	} else {
		fmt.Println("验证失败")
	}
}

//短信验证码服务
// http://localhost:8090/api/sendcode?phone=15811420676
func (mc *MemberController) sendSmsCode(context *gin.Context) {
	//发送验证码
	phone, exist := context.GetQuery("phone")
	if !exist {
		tool.Failed(context, "参数解析失败")
		return
	}
	//绑定服务
	ms := service.MemberService{}
	isSend := ms.SendCode(phone)
	if isSend {
		tool.Success(context, "发送成功")
		return
	}
	tool.Failed(context, "发送失败")
}

//第1种登录方法: 手机号+短信登录方法
//http://localhost:8090/api/login_sms
func (mc *MemberController) smsLogin(context *gin.Context) {
	var smsLoginParam param.SmsLoginParam
	err := tool.Decode(context.Request.Body, &smsLoginParam)
	if err != nil {
		tool.Failed(context, "参数解析失败")
		return
	}
	//完成手机+验证码登录
	us := service.MemberService{}
	member := us.SmsLogin(smsLoginParam)
	if member != nil {
		tool.Success(context, member)
		return
	}
	tool.Failed(context, "登录失败")
}
