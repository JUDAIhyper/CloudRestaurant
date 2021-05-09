package model

type FoodCategory struct {
	//类别ID
	Id int64 `gorm:"pk autoincr" json:"id"`
	//食品类别标题
	Title string `gorm:"varchar(20)" json:"title"`
	//食品描述
	Description string `gorm:"varchar(30)" json:"description"`
	//食品种类图片
	ImageUrl string `gorm:"varchar(255)" json:"image_url"`
	//食品类别连接
	LinkUrl string `gorm:"varchar(255)" json:"link_url"`
	//该类别是否在服务状态
	IsInServing bool `json:"is_in_serving"`
}
