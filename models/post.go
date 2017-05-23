package models

var (
	PostModel = &Post{}
)

type Post struct {
	Id         uint   `gorm:"primary_key, column:id"`
	CategoryId uint   `gorm:"column:categoryId"`
	Title      string `gorm:"column:title"`
	Content    string `gorm:"column:content"`
	InsertTime uint   `gorm:"column:insertTime"`
	ModifyTime uint   `gorm:"column:modifyTime"`
}

func (this Post) TableName() string {
	return "post"
}

func (this Post) ConnName() string {
	return "default"
}
