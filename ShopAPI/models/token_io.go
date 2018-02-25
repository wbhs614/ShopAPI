package models

//RegiserForm definition
type CreateTokenForm struct {
	Appid  string `form:"appid"`
	Secret string `form:"secret"`
}
