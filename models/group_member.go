package models

import (
	"time"
)

type GroupMemberModel struct {
	Id  int `json:"omitempty" gorm:"primaryKey"`
	GroupId int `json:"-" `
	AuthId int `json:"-"`
	CreatedAt time.Time `json:"-" `
	IsAdmin int `json:"-" gorm:"notnull;default:0;type:tinyint(1)"`
	IsLeave int `json:"-" gorm:"notnull;default:0;type:tinyint(1)"`
	IsDeleted int `json:"-" gorm:"notnull;default:0;type:tinyint(1)"`
	AuthForeign AuthModel `validate:"-" json:"-" gorm:"references:id;foreignKey:AuthId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	GroupForeign GroupModel `validate:"-" json:"-" gorm:"references:id;foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}


func (GroupMemberModel) TableName() string {
	return "group_member"
}