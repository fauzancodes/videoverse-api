package repository

import (
	"github.com/fauzancodes/videoverse-api/app/config"
	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/models"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/google/uuid"
)

func CreatePlaylist(data models.VAPlaylist) (models.VAPlaylist, error) {
	err := config.DB.Create(&data).Error

	return data, err
}

func GetPlaylistByID(id uuid.UUID, preloadFields []string) (response models.VAPlaylist, err error) {
	db := utils.BuildPreload(config.DB, preloadFields)

	err = db.Where("id = ?", id).First(&response).Error

	return
}

func GetPlaylists(param dto.FindParameter, preloadFields []string) (responses []models.VAPlaylist, total int64, totalFiltered int64, err error) {
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

func UpdatePlaylist(data models.VAPlaylist) (models.VAPlaylist, error) {
	err := config.DB.Save(&data).Error

	return data, err
}

func DeletePlaylist(data models.VAPlaylist) error {
	err := config.DB.Delete(&data).Error

	return err
}

func AddVideosToPlaylist(videos []models.VAVideo, playlist models.VAPlaylist) (err error) {
	err = config.DB.Model(&playlist).Association("Videos").Append(videos)

	return
}

func ClearVideosFromPlaylist(playlist models.VAPlaylist) (err error) {
	err = config.DB.Model(&playlist).Association("Videos").Clear()

	return
}
