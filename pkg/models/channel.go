package models

type Channel struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	AccessKey string `json:"accessKey"`
}
