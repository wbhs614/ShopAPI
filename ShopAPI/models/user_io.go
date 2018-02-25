package models

// RegisterForm definiton.
type RegisterForm struct {
	Email    string `form:"email"    valid:"Required"`
	Password string `form:"password" valid:"Required"`
	Nickname string `form:"nickname"`
}

// LoginForm definiton.
type LoginForm struct {
	Nickname string `form:"nickname"`
	Password string `form:"password" valid:"Required"`
}

type ChangeForm struct {
	Address  string `form:"address"`
	Age      int    `form:"age"`
	Gender   string `form:"gender"`
	Phone    string `form:"phone"`
	Realname string `form:"realname"`
	Roleid   int    `form:"roleid"`
	Token    string `form:"token"`
	Secret   string `form:"secret"`
}

type AllUserForm struct {
	Token  string `form:"token"`
	Secret string `form:"secret"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type ChangePassWordForm struct {
	Token  string `form:"token"`
	Secret string `form:"secret"`
	Passwd string `form:"passwd"`
}

type ChangeuserroleForm struct {
	Id       string `form:"id"`
	Nickname string `form:"nickname"`
}
