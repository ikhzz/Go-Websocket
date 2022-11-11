package models

import (
	"time"
	
)

type AuthUsecase interface {
	SignUp(response *MainResponse, reqParam *AuthModel) *MainResponse
	SignIn(response *MainResponse, reqParam *ReqSigninModel) *MainResponse
	GetAllUser(response *MainResponse) *MainResponse
}

type AuthRepository interface {
	SignUp(reqParam *AuthModel) error
	SignIn(reqParam *ReqSigninModel) (ResSigninModel, error)
	GetAllUser(response *MainResponse) []ResAllUser
}

func (AuthModel) TableName() string {
	return "auth"
}

type Tabler interface {
	TableName() string
}

type (
	AuthModel struct {
		Id  int `json:"omitempty" gorm:"primaryKey"`
		Email string `json:"email" gorm:"unique;notnull"`
		Username string `json:"username" gorm:"notnull;type:varchar(255)"`
		Password string `json:"password" gorm:"notnull;size:255;type:varchar(255)"`
		ProfileImage string `json:"profile_image" gorm:"type:varchar(255)"`
		CreatedAt time.Time `json:"-"`
		IsOnline	int	`json:"is_online" gorm:"notnull;default:0;type:tinyint(1)"`
		IsDeleted int `json:"-" gorm:"notnull;default:0;type:tinyint(1)"`
	}

	ReqSigninModel struct {
		Email	string	`json:"email"` 
		Password string `json:"password"`
	}

	ResSigninModel struct {
		Id	int
		Password string
	}

	ResAllUser struct {
		Id						int			`json:"id"`
		Username			string	`json:"username"`
		ProfileImage	string	`json:"profile_image"`
		IsHistory			int			`json:"is_history"`
	}
)