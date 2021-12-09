package repository

import "github.com/falcon78/realtime_iot/pkg/models"

func (r *Repo) GetChannels() (*[]models.Channel, error) {
	var channels []models.Channel
	if err := r.DB.Find(&channels).Error; err != nil {
		return &channels, err
	}
	return &channels, nil
}
