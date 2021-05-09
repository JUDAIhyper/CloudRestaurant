package dao

import (
	"CloudRestaurant/model"
	"CloudRestaurant/tool"
	"fmt"
)

type MemberDao struct {
	*tool.Orm
}

//根据用户id查询
func (md *MemberDao) QueryMemberById(id int64) *model.Member {
	var member model.Member
	err := md.Where("id=?", id).Find(&member).Error
	if err != nil {
		return nil
	}
	return &member
}

//更新member记录，头像
func (md *MemberDao) UpdateMemberAvatar(userId int64, fileName string) int64 {
	member := model.Member{Avatar: fileName}
	err := md.Where("id=?", userId).Update(&member).Error
	if err != nil {
		panic(err.Error())
	}
	//获取插入记录的id
	var id []int64
	md.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id)
	return id[0]
}

//根据用户名密码进行查询
func (md *MemberDao) Qery(name string, password string) *model.Member {
	var member model.Member

	password = tool.EncoderSha256(password)

	err := md.Where("user_name=? and password=?", name, password).First(&member).Error
	if err != nil {
		panic(err.Error())
	}
	return &member
}

//验证手机号和验证码是否存在
func (md *MemberDao) ValidateSmsCode(phone string, code string) *model.SmsCode {
	var sms model.SmsCode
	//TODO 更改为gorm语法
	err := md.Where("phone=? and code=?", phone, code).First(&sms).Error
	if err != nil {
		panic(err.Error())
	}
	return &sms
}

//使用手机号比对数据表
func (md *MemberDao) QueryByPhone(phone string) *model.Member {
	var member model.Member
	//if _, err := md.Where("mobile=?", phone).Get(&member); err != nil {
	//	fmt.Println(err.Error())
	//}
	//return &member
	err := md.Where("mobile=?", phone).First(&member).Error
	if err != nil {
		return nil
	}
	return &member
}

//向表中插入新的会员信息
func (md *MemberDao) InsertMember(member model.Member) int64 {
	//result, err := md.InsertOne(&member)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return 0
	//}
	//return result
	err := md.Create(&member).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	//获取插入记录的id
	var id []int64
	md.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id)
	return id[0]
}

func (md *MemberDao) InsertCode(sms model.SmsCode) int64 {
	//result, err := md.InsertOne(&sms)
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//return result
	err := md.Create(&sms).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	//获取插入记录的id
	var id []int64
	md.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id)
	return id[0]
}
