package repository

import (
	"clean_arch_v2/config"
	"clean_arch_v2/models"

	"gorm.io/gorm"
)

type generalRepository struct {
	sqlDB	*gorm.DB
}

func NewGeneralMessageRepository() models.GeneralMRepository {
	return &generalRepository{
		sqlDB: config.SqlConn,
	}
}

func(gm *generalRepository) UpdateUserActiveStatus(id int, status int) error {
	gm.sqlDB.Table("auth").Where("is_deleted = 0 AND id != ?", id).Update("is_online", status)

	return nil
}

func(gm *generalRepository) GetUserHistory(idFrom int, idTo int) *[]models.ChatModel {
	result := make([]models.ChatModel, 0) 	
	gm.sqlDB.Table("general_message_list").Where("(id_from = ? AND id_to = ?) OR (id_from = ? AND id_to = ?)", idFrom, idTo, idTo, idFrom).Select("send_by", "message", "created_at").Find(&result)
	return &result
}

func(gm *generalRepository) AddChat(reqParam *models.GeneralMessageListModel) error {
	
	res := gm.sqlDB.Create(reqParam)
	
	if res.Error != nil {
		return res.Error
	}
	
	return nil
}



