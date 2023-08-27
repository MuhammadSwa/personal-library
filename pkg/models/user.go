package models

type User struct {
	Id          int    `josn:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	AccessToken string `json:"access_token"`
}
