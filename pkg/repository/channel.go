package repository

import (
	"github.com/falcon78/realtime_iot/pkg/models"
	"github.com/google/uuid"
)

func (r *Repo) GetChannels() (*[]models.Channel, error) {
	var channels []models.Channel
	if err := r.DB.Find(&channels).Error; err != nil {
		return &channels, err
	}
	return &channels, nil
}

func (r *Repo) GetChannelByName(channelName string) (*models.Channel, error) {
	var channel models.Channel
	if err := r.DB.First(&channel, "name = ?", channelName).Error; err != nil {
		return &channel, err
	}
	return &channel, nil
}

func (r *Repo) GetChannelByAccessKey(key string) (*models.Channel, error) {
	var channel models.Channel
	if err := r.DB.First(&channel, "access_key = ?", key).Error; err != nil {
		return &channel, err
	}
	return &channel, nil
}

func (r *Repo) CreateChannel(channelName string) error {
	id := uuid.New()

	newChannel := models.Channel{
		Name:      channelName,
		AccessKey: id.String(),
	}

	if err := r.DB.Create(&newChannel).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repo) DeleteChannel(id int) error {
	if err := r.DeleteRecordsByChannelId(id); err != nil {
		return err
	}
	return r.DB.Delete(&models.Channel{}, id).Error
}
