package models

type Session struct {
	Id        int    `json:"id"`
	Sessionid string `json:"sessionid"`
	Userid    int    `json:"userid"`
}
