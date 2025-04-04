package repository

import (
	"github.com/fauzancodes/videoverse-api/app/config"
	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/models"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/google/uuid"
)

func CreateChannel(data models.VAChannel) (models.VAChannel, error) {
	err := config.DB.Create(&data).Error

	return data, err
}

func GetChannelByID(id uuid.UUID, preloadFields []string) (response models.VAChannel, err error) {
	db := utils.BuildPreload(config.DB, preloadFields)

	err = db.Where("id = ?", id).First(&response).Error

	return
}

func GetChannels(param dto.FindParameter, preloadFields []string) (responses []models.VAChannel, total int64, totalFiltered int64, err error) {
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

func UpdateChannel(data models.VAChannel) (models.VAChannel, error) {
	err := config.DB.Save(&data).Error

	return data, err
}

func DeleteChannel(data models.VAChannel) error {
	err := config.DB.Delete(&data).Error

	return err
}

func CreateSubscription(data models.VASubscription) (models.VASubscription, error) {
	err := config.DB.Create(&data).Error

	return data, err
}

func GetSubscriptions(param dto.FindParameter, preloadFields []string) (responses []models.VASubscription, total int64, totalFiltered int64, err error) {
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

func DeleteSubscription(data models.VASubscription) error {
	err := config.DB.Delete(&data).Error

	return err
}
