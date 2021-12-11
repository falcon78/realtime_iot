package models

type Record struct {
	Id           int     `json:"id"`
	ChannelId    int     `json:"channelId"`
	ChannelOne   float64 `json:"channelOne"`
	ChannelTwo   float64 `json:"channelTwo"`
	ChannelThree float64 `json:"channelThree"`
	ChannelFour  float64 `json:"channelFour"`
}
