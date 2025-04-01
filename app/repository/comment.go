package repository

import (
	"github.com/fauzancodes/videoverse-api/app/config"
	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/models"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/google/uuid"
)

func CreateComment(data models.VAComment) (models.VAComment, error) {
	err := config.DB.Create(&data).Error

	return data, err
}

func GetCommentByID(id uuid.UUID, preloadFields []string) (response models.VAComment, err error) {
	db := utils.BuildPreload(config.DB, preloadFields)

	err = db.Where("id = ?", id).First(&response).Error

	return
}

func GetCommentByIDs(ids []uuid.UUID, preloadFields []string) (responses []models.VAComment, err error) {
	db := utils.BuildPreload(config.DB, preloadFields)

	err = db.Where("id IN ?", ids).Find(&responses).Error

	return
}

func GetComments(param dto.FindParameter, preloadFields []string) (responses []models.VAComment, total int64, totalFiltered int64, err error) {
	err = config.DB.Model(responses).Where(param.BaseFilter, param.BaseFilterValues...).Count(&total).Error
	if err != nil {
		return
	}

	err = config.DB.Model(responses).Where(param.Filter, param.FilterValues...).Count(&totalFiltered).Error
	if err != nil {
		return
	}

	db := utils.BuildPreload(config.DB, preloadFields)

	if param.Limit == 0 {
		err = db.Offset(param.Offset).Order(param.Order).Where(param.Filter, param.FilterValues...).Find(&responses).Error
	} else {
		err = db.Limit(param.Limit).Offset(param.Offset).Order(param.Order).Where(param.Filter, param.FilterValues...).Find(&responses).Error
	}

	return
}

func UpdateComment(data models.VAComment) (models.VAComment, error) {
	err := config.DB.Save(&data).Error

	return data, err
}

func DeleteComment(data models.VAComment) error {
	err := config.DB.Delete(&data).Error

	return err
}
