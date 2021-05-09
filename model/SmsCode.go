package model

type SmsCode struct {
	Id         int64  `gorm:"pk autoincr" json:"id"`
	Phone      string `gorm:"varchar(11)" json:"phone"`
	BizId      string `gorm:"varchar(30)" json:"biz_id"`
	Code       string `gorm:"varchar(6)" json:"code"`
	CreateTime int64  `gorm:"bigint" json:"create_time"`
}