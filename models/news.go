package models
import (
	"gorm.io/gorm"
)

//News ...
type News struct {
	gorm.Model
	Title string `gorm:"not null" json:"title"`
	Thumbnail string `gorm:"not null" json:"thumbnail"`
	Summary string `gorm:"not null" json:"summary"`
	Content string `gorm:"not null" json:"content"`
	Tags []Tag `gorm:"many2many:news_tag;not null" json:"tags"`
	Topic string `gorm:"not null" json:"topic"`
	Status string `gorm:"not null" json:"status"`
}

func (News) TableName() string {
	return "news"
}

//NewsList ...
type NewsList struct {
	Data []News `json:"data"`
}


