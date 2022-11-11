package models

import (
	"time"
)

type GroupMessageListModel struct {
	Id  int `json:"omitempty" gorm:"primaryKey"`
	IdFrom int `json:"id_from"`
	GroupId int `json:"group_id"`
	Message string `json:"message" gorm:"type:varchar(255)"`
	IsForwarded int `json:"is_forwarded" gorm:"notnull;default:0;type:tinyint(1)"`
	IdReply int `json:"id_reply" gorm:"notnull;default:0;type:int(11)"`
	IsRead int `json:"is_read" gorm:"notnull;default:0;type:tinyint(1)"`
	CreatedAt time.Time `json:"-" `
	IsDeleted int `json:"-" gorm:"notnull;default:0;type:tinyint(1)"`
	IdFromForeign AuthModel `validate:"-" json:"-" gorm:"references:id;foreignKey:IdFrom;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	GroupIdForeign GroupModel `validate:"-" json:"-" gorm:"references:id;foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}


func (GroupMessageListModel) TableName() string {
	return "group_message_list"
}