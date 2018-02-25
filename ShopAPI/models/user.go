package models

import (
	"crypto/rand"
	"github.com/astaxie/beego/orm"
	"io"
	//"errors"
	"golang.org/x/crypto/scrypt"
	"strconv"
	//"time"
	"ShopAPI/utils"
	"fmt"
)

var (
	UserList map[string]*User
)

func init() {
}

type User struct {
	Id          int    `json:"-"`
	Realname    string `json:"realnme"`
	Passwd      string `json:"passwd"`
	Salt        string `json:"salt"`
	Create_time string `json:"createtime"`
	Appid       string `json:"uid"`
	Secret      string `json:"secret"`
	Role_id     int    `json:"roleid"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	Address     string `json:"address"`
	Nickname    string `json:"nickname"`
	Phone       string `json:"phone"`
	UpdateTime  string `json:"updatetime"`
	Token       string `json:"token"`
	UrlStr      string `json:"urlstr"`
	Pyname      string `json:"pyname"`
}

const pwHashBytes = 32

func generateSalt() (salt string, err error) {
	buf := make([]byte, pwHashBytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", buf), nil
}

func generatePassHash(password string, salt string) (hash string, err error) {
	h, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, pwHashBytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h), nil
}

func init() {
	orm.RegisterModel(new(User))
}

func NewUser(r *RegisterForm, appid string, secret string) (u *User, err error) {
	salt, err := generateSalt()
	regDate := utils.GetNowTimeStr()
	if err != nil {
		return nil, err
	}
	hash, err := generatePassHash(r.Password, salt)
	if err != nil {
		return nil, err
	}
	user := User{
		Nickname:    r.Nickname,
		Passwd:      hash,
		Email:       r.Email,
		Salt:        salt,
		Create_time: regDate,
		Appid:       appid,
		Secret:      secret,
		Role_id:     1,
	}
	return &user, nil
}

func (self *User) Insert() error {
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.Insert(self)
	return err
}

func (self *User) FindUersById() (cou int64, err error) {
	qs := orm.NewOrm().QueryTable(new(User))
	count, err := qs.Filter("appid", self.Appid).Count()
	return count, err
}

func (self *User) FindUerById(appid string) (uer *User, err error) {
	qs := orm.NewOrm().QueryTable(new(User))
	erro := qs.Filter("appid", appid).One(self)
	return self, erro
}

func (self *User) CheckPassWord(inputPass string) (exit bool, err error) {
	inputHash, err := generatePassHash(inputPass, self.Salt)
	if err != nil {
		return false, err
	} else {
		return self.Passwd == inputHash, err
	}
}

func (self *User) CreateOrUpdateUserToken() bool {
	token, expire := CreateToken(self.Appid, self.Secret)
	userToken := User_token{
		Token:       token,
		Express_in:  strconv.FormatInt(expire, 10),
		Appid:       self.Appid,
		Create_time: utils.GetNowTimeStr(),
	}
	o := orm.NewOrm()
	count, err := o.QueryTable(new(User_token)).Filter("appid", self.Appid).Count()
	if err != nil {
		return false
	} else {
		if count == 0 {
			_, err := o.Insert(&userToken)
			if err != nil {
				return false
			} else {
				self.Token = token
				return true
			}
			return false
		} else {
			_, err := o.QueryTable(new(User_token)).Filter("appid", self.Appid).Update(orm.Params{
				"token":       userToken.Token,
				"express_in":  userToken.Express_in,
				"create_time": userToken.Create_time,
			})
			if err != nil {
				return false
			} else {
				self.Token = token
				return true
			}
		}
	}
}

func (self *User) GetUser(token, secret string) (resut *utils.ControllerError) {
	uid, repon := GetUserUidByToken(token, secret)
	var responerr *utils.ControllerError
	if repon == utils.Actionsuccess {
		o := orm.NewOrm()
		err := o.QueryTable(new(User)).Filter("appid", uid).One(self)
		if err != nil {
			responerr = utils.ErrDatabase
		} else {
			responerr = utils.Actionsuccess
		}
	} else {
		responerr = repon
	}
	return responerr
}

func GetAllUsers(limit, offset int, token, secret string) (userList *[]orm.Params, resut *utils.ControllerError) {
	var list []orm.Params
	var responerr *utils.ControllerError
	_, err := TokenAuth(token, secret)
	if err == nil {
		o := orm.NewOrm()
		_, err := o.QueryTable(new(User)).Limit(limit, offset).Values(&list, "realname", "email", "gender", "age", "address", "nickname", "phone", "appid")
		if err != nil {
			responerr = utils.ErrDatabase
		} else {
			responerr = utils.Actionsuccess
		}
	} else {
		responerr = utils.ErrExpired
	}
	return &list, responerr
}

func UpdateUser(form *ChangeForm) (resut *utils.ControllerError) {
	uid, repon := GetUserUidByToken(form.Token, form.Secret)
	var responerr *utils.ControllerError
	if repon == utils.Actionsuccess {
		o := orm.NewOrm()
		dict := orm.Params{}
		if len(form.Address) > 0 {
			dict["address"] = form.Address
		}
		if len(form.Gender) > 0 {
			dict["gender"] = form.Gender
		}
		if len(form.Phone) > 0 {
			dict["phone"] = form.Phone
		}
		if len(form.Realname) > 0 {
			dict["realname"] = form.Realname
		}
		if form.Roleid > 0 {
			dict["roleid"] = form.Roleid
		}
		if form.Age > 0 {
			dict["age"] = form.Age
		}
		dict["updatetime"] = utils.GetNowTimeStr()
		_, err := o.QueryTable(new(User)).Filter("appid", uid).Update(dict)
		if err != nil {
			responerr = utils.ErrDatabase
		} else {
			responerr = utils.Actionsuccess
		}
	} else {
		responerr = repon
	}
	return responerr
}

func ChangePassword(form *ChangePassWordForm) (resut *utils.ControllerError) {
	uid, repon := GetUserUidByToken(form.Token, form.Secret)
	var responerr *utils.ControllerError
	if repon == utils.Actionsuccess {
		salt, err := generateSalt()
		if err != nil {
			responerr = utils.ErrSystem
		} else {
			hashPw, err := generatePassHash(form.Passwd, salt)
			if err != nil {
				responerr = utils.ErrSystem
			} else {
				dict := orm.Params{}
				dict["salt"] = salt
				dict["passwd"] = hashPw
				dict["updatetime"] = utils.GetNowTimeStr()
				o := orm.NewOrm()
				_, err := o.QueryTable(new(User)).Filter("appid", uid).Update(dict)
				if err != nil {
					responerr = utils.ErrDatabase
				} else {
					responerr = utils.Actionsuccess
				}
			}
		}
	} else {
		responerr = repon
	}
	return responerr

}

func Logout(token, secret string) (resut *utils.ControllerError) {
	uid, repon := GetUserUidByToken(token, secret)
	var responerr *utils.ControllerError
	if repon == utils.Actionsuccess {
		o := orm.NewOrm()
		_, err := o.QueryTable(new(User_token)).Filter("appid", uid).Delete()
		if err != nil {
			responerr = utils.ErrDatabase
		} else {
			responerr = utils.Actionsuccess
		}
	} else {
		responerr = repon
	}
	return responerr
}

func UpdateImageUrl(imageUrl, token, secret string) (resut *utils.ControllerError) {
	uid, repon := GetUserUidByToken(token, secret)
	var responerr *utils.ControllerError
	if repon == utils.Actionsuccess {
		dict := orm.Params{}
		dict["urlstr"] = imageUrl
		dict["updatetime"] = utils.GetNowTimeStr()
		o := orm.NewOrm()
		_, err := o.QueryTable(new(User)).Filter("appid", uid).Update(dict)
		if err != nil {
			responerr = utils.ErrDatabase
		} else {
			responerr = utils.Actionsuccess
		}
	} else {
		responerr = repon
	}
	return responerr
}

func Login(username, password string) bool {
	return false
}

func DeleteUser(uid string) {
}
