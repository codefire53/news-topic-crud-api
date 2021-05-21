package models
import (
	"gorm.io/gorm"
)

//Tag ...
type Tag struct {
	gorm.Model
	Name string `gorm:"not null, unique" json:"name"`
}

//TagsList ...
type TagsList struct {
	Data []Tag `json:"data"`
}



