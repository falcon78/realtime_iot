package repository

import (
	"github.com/falcon78/realtime_iot/pkg/models"
	"time"
)

type AllRecordsResponse struct {
	Min     float64          `json:"min"`
	Max     float64          `json:"max"`
	Records *[]models.Record `json:"records"`
}

func (r *Repo) GetAllRecordsByChannelId(channelId int, limit int) (*AllRecordsResponse, error) {
	var records []models.Record
	if err := r.DB.Order("timestamp desc").Limit(limit).Find(&records, "channel_id = ?", channelId).Error; err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return &AllRecordsResponse{
			Min:     0,
			Max:     0,
			Records: &records,
		}, nil
	}

	var max float64
	var min float64

	r.DB.Raw(""+
		"select case "+
		"when m1 >= m2 and m1 >= m3 and m1 >= m4 then m1 "+
		"when m2 >= m1 and m2 >= m3 and m2 >= m4 then m2 "+
		"when m3 >= m1 and m3 >= m2 and m3 >= m4 then m3 "+
		"else m4 "+
		"end max "+
		"from (select max(channel_four) as m1, max(channel_one) as m2, "+
		"max(channel_three) as m3, max(channel_two) as m4 "+
		"from Records where channel_id = ?) as r;", channelId,
	).Scan(&max)

	r.DB.Raw("select case "+
		"when m1 <= m2 and m1 <= m3 and m1 <= m4 then m1 "+
		"when m2 <= m1 and m2 <= m3 and m2 <= m4 then m2 "+
		"when m3 <= m1 and m3 <= m2 and m3 <= m4 then m3 "+
		"else m4 "+
		"end min "+
		"from (select min(channel_four) as m1, min(channel_one) as m2, "+
		"min(channel_three) as m3, min(channel_two) as m4 "+
		"from Records where channel_id = ?) as r;", channelId,
	).Scan(&min)

	max = max * 1.1
	min = min * 1.1

	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}

	return &AllRecordsResponse{
		Min:     min,
		Max:     max,
		Records: &records,
	}, nil
}

func (r *Repo) GetAllRecordsByChannelKey(key string) ([]*models.Record, error) {
	channel, err := r.GetChannelByAccessKey(key)
	if err != nil {
		return nil, err
	}

	var records []*models.Record
	err = r.DB.Find(&records, "channel_id = ?", channel.Id).Error
	if err != nil {
		return records, err
	}

	return records, nil
}

func (r *Repo) DeleteRecordsByChannelId(channelId int) error {
	return r.DB.Where("channel_id = ?", channelId).Delete(&models.Record{}).Error
}

func (r *Repo) AddRecord(accessKey string, c1, c2, c3, c4 float64, current time.Time) error {
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
		Timestamp:    current,
	}

	return r.DB.Create(&newRecord).Error
}
