package repository

import (
	"github.com/falcon78/realtime_iot/pkg/models"
)

func (r *Repo) GetAllRecordsByChannelId(channelId int) (*[]models.Record, error) {
	var records []models.Record
	if err := r.DB.Find(&records, "channel_id = ?", channelId).Limit(100).Error; err != nil {
		return &records, err
	}
	return &records, nil
}

func (r *Repo) DeleteRecordsByChannelId(channelId int) error {
	return r.DB.Where("channel_id = ?", channelId).Delete(&models.Record{}).Error
}

func (r *Repo) AddRecord(accessKey string, c1, c2, c3, c4 float64) error {
	channel, err := r.GetChannelByAccessKey(accessKey)
	if err != nil {
		return err
	}

	newRecord := &models.Record{
		ChannelId:    channel.Id,
		ChannelOne:   c1,
		ChannelTwo:   c2,
		ChannelThree: c3,
		ChannelFour:  c4,
	}

	return r.DB.Create(&newRecord).Error
}
