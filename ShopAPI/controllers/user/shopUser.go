package user

import (
	"ShopAPI/models"
	"ShopAPI/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/mingzhehao/goutils/filetool"
	"math/rand"
	"strings"
	"time"
)

type RegisterController struct {
	beego.Controller
}

type LoginController struct {
	beego.Controller
}

type ShopUserController struct {
	beego.Controller
}
type ShopUserListController struct {
	beego.Controller
}

type ChangeUserInfoController struct {
	beego.Controller
}

type ChangePassWordController struct {
	beego.Controller
}

type LogoutController struct {
	beego.Controller
}

type ChangeHeaderImageController struct {
	beego.Controller
}

/*注册用户的控制器*/
func (self *RegisterController) Post() {
	form := &models.RegisterForm{}
	if err := self.ParseForm(form); err != nil {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	}
	beego.Debug("ParseRegsigerFrorm", &form)
	appid := utils.GetMD5(form.Nickname)
	secret := utils.GetMD5(form.Email)
	user, err := models.NewUser(form, appid, secret)
	if err != nil {
		self.Data["json"] = utils.ErrSystem
		self.ServeJSON()
		return
	}
	if len(form.Nickname) == 0 {
		self.Data["json"] = utils.ErrNilNickName
		self.ServeJSON()
		return
	}
	if len(form.Email) == 0 {
		self.Data["json"] = utils.ErrNilEmail
		self.ServeJSON()
		return
	} else {
		if m := utils.IsEmail(form.Email); m == false {
			self.Data["json"] = utils.ErrEmailFormat
			self.ServeJSON()
			return
		}
	}
	if len(form.Password) == 0 {
		self.Data["json"] = utils.ErrNilPassWord
		self.ServeJSON()
		return
	}

	cou, err := user.FindUersById()
	if err != nil {
		self.Data["json"] = utils.ErrSystem
		self.ServeJSON()
		return
	} else {
		if cou > 0 {
			self.Data["json"] = utils.ErrExitUser
			self.ServeJSON()
			return
		}
	}
	err = user.Insert()
	if err != nil {
		self.Data["json"] = utils.ErrDatabase
		self.ServeJSON()
	} else {
		self.Data["json"] = utils.Actionsuccess
		self.ServeJSON()
	}
}

/*登陆的控制器*/
func (self *LoginController) Post() {
	form := &models.LoginForm{}
	//检验参数
	if err := self.ParseForm(form); err != nil {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	}
	if len(form.Nickname) == 0 {
		self.Data["json"] = utils.ErrNilLoginNickName
		self.ServeJSON()
		return
	}
	if len(form.Password) == 0 {
		self.Data["json"] = utils.ErrNilLoginPassWord
		self.ServeJSON()
		return
	}
	//检验用户存在和密码
	appid := utils.GetMD5(form.Nickname)
	user := models.User{}
	_, err := user.FindUerById(appid)
	if err != nil {
		if err == orm.ErrNoRows {
			self.Data["json"] = utils.ErrNoUser
			self.ServeJSON()
			return
		} else {
			self.Data["json"] = utils.ErrDatabase
			self.ServeJSON()
			return
		}
	} else {
		if ok, err := user.CheckPassWord(form.Password); err != nil {
			self.Data["json"] = utils.ErrNoUserPass
			self.ServeJSON()
			return
		} else {
			if ok {
				success := user.CreateOrUpdateUserToken()
				if success {
					self.Data["json"] = utils.CommendFormat{
						Data:    &user,
						Code:    1,
						Message: "恭喜你，登陆成功！",
					}
					self.ServeJSON()
				} else {
					self.Data["json"] = utils.ErrTokenCreate
					self.ServeJSON()
				}
			} else {
				self.Data["json"] = utils.ErrPass
				self.ServeJSON()
			}
		}
	}
}

func (self *ShopUserController) Post() {
	token := self.GetString("token")
	secret := self.GetString("secret")
	if len(token) == 0 || len(secret) == 0 {
		self.Data["json"] = utils.ErrLessParms
		self.ServeJSON()
		return
	} else {
		user := &models.User{}
		repon := user.GetUser(token, secret)
		if repon == utils.Actionsuccess {
			self.Data["json"] = utils.CommendFormat{
				Data: map[string]interface{}{
					"realnme":  user.Realname,
					"uid":      user.Appid,
					"secret":   user.Secret,
					"roleid":   user.Role_id,
					"email":    user.Email,
					"gender":   strings.TrimSpace(user.Gender),
					"age":      user.Age,
					"address":  user.Address,
					"nickname": user.Nickname,
					"phone":    user.Phone,
					"urlstr":   user.UrlStr,
					"ip":       self.Ctx.Input.IP(),
				},
				Code:    9000,
				Message: "获取用户信息成功",
			}
			self.ServeJSON()
		} else {
			self.Data["json"] = repon
			self.ServeJSON()
		}
	}
}

func (self *ChangeUserInfoController) Post() {
	form := &models.ChangeForm{}
	self.ParseForm(form)
	if (len(form.Address) == 0 && form.Age == 0 && len(form.Gender) == 0 && len(form.Phone) == 0 && form.Roleid == 0) || (len(form.Token) == 0 || len(form.Secret) == 0) {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else if strings.ToLower(strings.TrimSpace(form.Gender)) != "f" && strings.ToLower(strings.TrimSpace(form.Gender)) != "m" {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		repon := models.UpdateUser(form)
		self.Data["json"] = repon
		self.ServeJSON()
	}
}

func (self *ShopUserListController) Post() {
	form := &models.AllUserForm{}
	self.ParseForm(form)
	if len(form.Token) == 0 || len(form.Secret) == 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		var limit, offset int
		if form.Limit < 1 {
			limit = -1
		} else {
			limit = form.Limit
		}
		if form.Offset < 1 {
			form.Offset = 0
		} else {
			offset = form.Offset
		}
		list, err := models.GetAllUsers(limit, offset, form.Token, form.Secret)
		if err == utils.Actionsuccess {
			self.Data["json"] = utils.CommendFormat{
				Data:    list,
				Code:    9000,
				Message: "获取信息成功",
			}
			self.ServeJSON()
		} else {
			self.Data["json"] = err
			self.ServeJSON()
		}
	}

}

func (self *ChangePassWordController) Post() {
	form := &models.ChangePassWordForm{}
	self.ParseForm(form)
	if len(form.Passwd) == 0 || len(form.Secret) == 0 || len(form.Token) == 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
	} else {
		err := models.ChangePassword(form)
		self.Data["json"] = err
		self.ServeJSON()
	}
}

func (self *LogoutController) Post() {
	token := self.GetString("token")
	secret := self.GetString("secret")
	if len(token) == 0 && len(secret) == 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		err := models.Logout(token, secret)
		self.Data["json"] = err
		self.ServeJSON()
	}
}

func (self *ChangeHeaderImageController) Post() {
	token := self.GetString("token")
	secret := self.GetString("secret")
	_, h, err := self.GetFile("headerImage")
	//defer f.Close()
	if err != nil || len(token) < 0 || len(secret) < 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		//var url string
		ext := filetool.Ext(h.Filename)
		fileExt := strings.TrimLeft(ext, ".")
		if fileExt != "jpeg" && fileExt != "png" {
			self.Data["json"] = utils.ErrFileType
			self.ServeJSON()
			return
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		fieleSaveName := fmt.Sprintf("%s_%d%d%s", fileExt, time.Now().Unix(), r.Intn(100), ext)
		imgPath := fmt.Sprintf("%s/%s", utils.LOCAL_FILE_DIR, fieleSaveName)
		filetool.InsureDir(utils.LOCAL_FILE_DIR)
		err := self.SaveToFile("headerImage", imgPath)
		if err == nil {
			err := models.UpdateImageUrl(imgPath, token, secret)
			if err == utils.Actionsuccess {
				self.Data["json"] = utils.CommendFormat{
					Data:    imgPath,
					Code:    9000,
					Message: "更新图像成功",
				}
				self.ServeJSON()
			} else {
				self.Data["json"] = err
				self.ServeJSON()
			}
			return
		} else {
			self.Data["json"] = utils.ErrWriteFile
			self.ServeJSON()
			return
		}
	}

}
