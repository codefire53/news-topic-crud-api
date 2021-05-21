package models
import (
	"gorm.io/gorm"
)

//Tag ...
type Tag struct {
	gorm.Model
	Name string `gorm:"not null" json:"name"`
	News []News `gorm:"many2many:news_tag;not null" json:"news"`
}


