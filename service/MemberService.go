package service

import (
	"CloudRestaurant/dao"
	"CloudRestaurant/model"
	"CloudRestaurant/param"
	"CloudRestaurant/tool"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type MemberService struct {
}
//根据用户ID查询信息
func (ms *MemberService) GetUserInfo(userId string) *model.Member {
	id, err := strconv.Atoi(userId)
	if err != nil {
		return nil
	}
	memberDao := dao.MemberDao{tool.DbEngine}
	return memberDao.QueryMemberById(int64(id))
}

//上传头像方法
func (mc *MemberService) UploadAvatar(userId int64, fileName string) string {
	memberDao := dao.MemberDao{tool.DbEngine}
	result := memberDao.UpdateMemberAvatar(userId, fileName)
	if result == 0 {
		return ""
	}
	return fileName
}

//用户名密码登录方法
func (ms *MemberService) Login(name string, password string) *model.Member {
	//1.使用用户名+密码 查询用户信息 如果存在则直接返回
	md := dao.MemberDao{tool.DbEngine}
	member := md.Qery(name, password)
	if member.Id != 0 {
		return member
	}
	//2.用户信息不存在，作为新用户保存到数据库中
	user := model.Member{}
	user.UserName = name
	//给密码进行hash加密
	user.Password = tool.EncoderSha256(password)
	user.RegisterTime = time.Now().Unix()

	result := md.InsertMember(user)
	user.Id = result
	return &user
}

//手机号+验证码登录方法
func (ms *MemberService) SmsLogin(loginparam param.SmsLoginParam) *model.Member {
	//1.获取到手机号和验证码

	//2.验证手机号+验证码是否正确(与数据表比对)
	md := dao.MemberDao{tool.DbEngine}
	sms := md.ValidateSmsCode(loginparam.Phone, loginparam.Code)
	if sms.Id == 0 {
		return nil
	}
	//3.根据手机号member表中查询记录
	member := md.QueryByPhone(loginparam.Phone)
	if member.Id != 0 {
		return member
	}
	//4.新创建一个member记录，并保存
	user := model.Member{}
	user.UserName = loginparam.Phone
	user.Mobile = loginparam.Phone
	user.RegisterTime = time.Now().Unix()
	user.Id = md.InsertMember(user)

	return &user
}

//发送验证码方法
func (ms *MemberService) SendCode(phone string) bool {
	//1.产生一个验证码
	code := fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))

	//2.调用阿里云的sdk 完成发送
	config := tool.GetConfig().Sms
	client, err := dysmsapi.NewClientWithAccessKey(config.RegionId, config.AppKey, config.AppSecret)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = config.SignName
	request.TemplateCode = config.TemplateCode
	request.PhoneNumbers = phone

	par, err := json.Marshal(map[string]interface{}{
		"code": code,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	request.TemplateParam = string(par)
	response, err := client.SendSms(request)
	fmt.Println(response)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	//3.接收返回结果，并判断发送状态
	//判断验证码是否发送成功
	if response.Code == "OK" {
		//将验证码保存到数据库中
		smsCode := model.SmsCode{Phone: phone, Code: code, BizId: response.BizId, CreateTime: time.Now().Unix()}
		memberDao := dao.MemberDao{tool.DbEngine}
		result := memberDao.InsertCode(smsCode)
		return result > 0
	}
	return false
}
