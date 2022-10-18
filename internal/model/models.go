package model

type Visit struct {
	UserIp string `json:"user_ip"`
	URL    string `json:"url"`
	Count  int    `json:"count"`
}
