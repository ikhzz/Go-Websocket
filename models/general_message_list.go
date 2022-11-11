package models

import (
	"time"

	"github.com/gorilla/websocket"
)

type GeneralMUsecase interface {
	HandleInit(currentConn *WebSocketConnection, conn *MainResponse)
}

type GeneralMRepository interface {
	
	GetUserHistory(idFrom int, idTo int) *[]ChatModel
	AddChat(reqParam *GeneralMessageListModel) error
	UpdateUserActiveStatus(id int, status int) error
}

type (
	GeneralMessageListModel struct {
		Id  int `json:"omitempty" gorm:"primaryKey"`
		IdFrom int `json:"id_from"`
		IdTo int `json:"id_to"`
		Message string `json:"message" gorm:"type:varchar(255)"`
		IsForwarded int `json:"is_forwarded" gorm:"notnull;default:0;type:tinyint(1)"`
		IdReply int `json:"id_reply" gorm:"notnull;default:0;type:int(11)"`
		IsRead int `json:"is_read" gorm:"notnull;default:0;type:tinyint(1)"`
		CreatedAt time.Time `json:"-" `
		IsDeletedFrom int `json:"-" gorm:"notnull;default:0;type:tinyint(1)"`
		IsDeletedTo int `json:"-" gorm:"notnull;default:0;type:tinyint(1)"`
		SendBy int `json:"send_by"`
		IdFromForeign AuthModel `validate:"-" json:"-" gorm:"references:id;foreignKey:IdFrom;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		IdToForeign AuthModel `validate:"-" json:"-" gorm:"references:id;foreignKey:IdTo;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}

	WebSocketConnection struct {
		*websocket.Conn
		Id	int
	}
	// payload bind model
	WebsocketPayloadModel struct {
		Message string
		Id			int
	}
	// payload send model to current user
	WebsocketResponseModel struct {
		Data []ChatHistory
		Type string
	}
	// single chat model
	ChatModel struct {
		SendBy int
		Message string
		CreatedAt time.Time
		FormatedCreatedAt string
	}
	// list chat each id
	ChatHistory struct {
		Id	int
		Chat []ChatModel
	}
)


func (GeneralMessageListModel) TableName() string {
	return "general_message_list"
}