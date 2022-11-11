package models

import (
	"time"
)

type GroupModel struct {
	Id  int `json:"omitempty" gorm:"primaryKey"`
	GroupName string `json:"group_name"`
	Desc string `json:"description"`
	ProfileImage string `json:"group_image" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"-" `
	IsDeleted int `json:"-" gorm:"notnull;default:0;type:tinyint(1)"`
}


func (GroupModel) TableName() string {
	return "group"
}