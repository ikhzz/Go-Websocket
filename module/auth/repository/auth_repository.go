package repository

import (
	"clean_arch_v2/config"
	"clean_arch_v2/models"
	"fmt"

	"gorm.io/gorm"
)

type authRepository struct {
	sqlDB	*gorm.DB
}

func NewAuthRepository() models.AuthRepository {
	return &authRepository{
		sqlDB: config.SqlConn,
	}
}

func(a authRepository) SignUp(reqParam *models.AuthModel) error {
	res := a.sqlDB.Create(reqParam)
	
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func(a authRepository) SignIn(reqParam *models.ReqSigninModel) (models.ResSigninModel, error) {
	var result models.ResSigninModel
	res := a.sqlDB.Table("auth").Where("is_deleted = ?", 0).Select("id", "password").First(&result, fmt.Sprintf(`email = '%s'`, reqParam.Email))
	if res.Error != nil {
		return result, res.Error
	}

	return result, nil
}

func(gm *authRepository) GetAllUser(response *models.MainResponse) []models.ResAllUser {
	result := make([]models.ResAllUser, 0) 	
	gm.sqlDB.Table("auth").Where("is_deleted = 0 AND id != ?", response.Payload.Id).Select("id", "username", "profile_image").Find(&result)
	return result
}