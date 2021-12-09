package repository

import "github.com/falcon78/realtime_iot/pkg/models"

func (r *Repo) GetAllRecordsByChannelId(channel_id int) (*[]models.Record, error) {
	var records []models.Record
	if err := r.DB.Find(&records, "channel_id = ?", channel_id).Error; err != nil {
		return &records, err
	}
	return &records, nil
}
