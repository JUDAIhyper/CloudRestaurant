package model

// 商家服务表
type Service struct {
	//id
	Id int64 `gorm:"pk autoincr" json:"id"`
	//服务名称
	Name string `gorm:"varchar(20)" json:"name"`
	//服务描述
	Description string `gorm:"varchar(30)" json:"description"`
	//图标名称
	IconName string `gorm:"varchar(3)" json:"icon_name"`
	//图标颜色
	IconColor string `gorm:"varchar(6)" json:"icon_color"`
}
