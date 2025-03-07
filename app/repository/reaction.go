package repository

import (
	"github.com/fauzancodes/videoverse-api/app/config"
	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/models"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
)

func CreateVideoLike(data models.VAVideoLike) (models.VAVideoLike, error) {
	err := config.DB.Create(&data).Error

	return data, err
}

func GetVideoLikes(param dto.FindParameter, preloadFields []string) (responses []models.VAVideoLike, total int64, totalFiltered int64, err error) {
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

func DeleteVideoLike(data models.VAVideoLike) error {
	err := config.DB.Delete(&data).Error

	return err
}

func CreateVideoDislike(data models.VAVideoDislike) (models.VAVideoDislike, error) {
	err := config.DB.Create(&data).Error

	return data, err
}

func GetVideoDislikes(param dto.FindParameter, preloadFields []string) (responses []models.VAVideoDislike, total int64, totalFiltered int64, err error) {
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

func DeleteVideoDislike(data models.VAVideoDislike) error {
	err := config.DB.Delete(&data).Error

	return err
}
