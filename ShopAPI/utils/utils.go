package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetNowTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

type ControllerError struct {
	Status   int    `json:"status"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	DevInfo  string `json:"dev_info"`
	MoreInfo string `json:"more_info"`
}

var (
	Err404        = &ControllerError{404, 404, "page not found", "page not found", ""}
	ErrLessParms  = &ControllerError{400, 10000, "缺少参数", "请按照文档传递参数", ""}
	ErrInputData  = &ControllerError{400, 10001, "数据输入错误", "客户端参数错误", ""}
	ErrDatabase   = &ControllerError{500, 10002, "服务器错误", "数据库操作错误", ""}
	ErrOpenFile   = &ControllerError{500, 10009, "服务器错误", "打开文件出错", ""}
	ErrWriteFile  = &ControllerError{500, 10010, "服务器错误", "写文件出错", ""}
	ErrSystem     = &ControllerError{500, 10011, "服务器错误", "操作系统错误", ""}
	ErrExpired    = &ControllerError{400, 10012, "登录已过期", "验证token过期", ""}
	ErrPermission = &ControllerError{400, 10013, "没有权限", "没有操作权限", ""}
	Actionsuccess = &ControllerError{200, 90000, "操作成功", "操作成功", ""}

	//user对应的错误码
	ErrDupUser             = &ControllerError{400, 10003, "用户信息已存在", "数据库记录重复", ""}
	ErrNoUser              = &ControllerError{400, 10004, "用户信息不存在", "数据库记录不存在", ""}
	ErrPass                = &ControllerError{400, 10005, "密码不正确", "密码不正确", ""}
	ErrNoUserPass          = &ControllerError{400, 10006, "用户信息不存在或密码不正确", "数据库记录不存在或密码不正确", ""}
	ErrNoUserChange        = &ControllerError{400, 10007, "用户信息不存在或数据未改变", "数据库记录不存在或数据未改变", ""}
	ErrInvalidUser         = &ControllerError{400, 10008, "用户信息不正确", "Session信息不正确", ""}
	ErrNilEmail            = &ControllerError{400, 10014, "请输入注册的邮箱", "用户输入参数错误", "邮箱为必填项"}
	ErrNilNickName         = &ControllerError{400, 10015, "请输入注册的用户昵称", "用户输入参数错误", "昵称必填项"}
	ErrNilPassWord         = &ControllerError{400, 10016, "请输入注册的用户密码", "用户输入参数错误", "密码为必填项"}
	ErrExitUser            = &ControllerError{400, 10017, "用户昵称已经存在", "昵称名重复", "需要重新输入昵称"}
	ErrNilLoginNickName    = &ControllerError{400, 10018, "请输入用户昵称", "用户输入参数错误", "昵称为必填项"}
	ErrNilLoginPassWord    = &ControllerError{400, 10019, "请输入用户密码", "用户输入参数错误", "密码为必填项"}
	ErrTokenCreate         = &ControllerError{400, 10020, "生成token失败", "生成token失败", ""}
	ErrEmailFormat         = &ControllerError{400, 10021, "邮箱的格式不正确", "用户输入参数错误", "邮箱为必填项"}
	ErrNickNameErrorFormat = &ControllerError{400, 10022, "邮箱的格式不正确", "用户输入参数错误", "邮箱为必填项"}
	ErrNoToken             = &ControllerError{400, 10023, "非法token", "校验token失败", ""}
	ErrCheckErrToken       = &ControllerError{400, 10024, "查询token失败", "校验token失败", ""}
	ErrNoFile              = &ControllerError{400, 10025, "请上传文件", "没有上传文件", ""}
	ErrFileType            = &ControllerError{400, 10026, "上传文件类型错误", "文件类型错误", ""}
	ErrFileWrite           = &ControllerError{500, 10027, "上传文件类型错误", "文件类型错误", ""}
	//商城对应的错误码
	ErrCheckGoodsType = &ControllerError{400, 10028, "查询商品类型失败", "请检查商品类型", ""}
	ErrUpdateStore    = &ControllerError{400, 10029, "更新库存表失败", "更新库存失败", ""}
	ErrWastErrType    = &ControllerError{400, 10030, "耗损不能大于库存总量", "请重新输入耗损量", ""}
	ErrGoodNoExit     = &ControllerError{400, 10031, "商品信息不存在", "商品信息不存在", ""}
	ErrGoodIds        = &ControllerError{400, 10032, "goodbuy的格式不正确", "请参考文档传入正确的goodBuy的格式", ""}
)

const (
	Select_all_user = "查找全部用户"
	Shop_secret     = "ShopAPIToken"
)

type CommendFormat struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func GetMD5(encode string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(encode))
	cipherStr := md5Ctx.Sum(nil)
	mdstr := string(base64Encode(cipherStr))
	return reverseString(mdstr)
}

func GetFormatCode() string {
	const shortForm = "2006-01-01 15:04:05"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	var valid = regexp.MustCompile("[0-9]")
	timeStr := strings.Join(valid.FindAllString(str, -1), "")
	r := rand.New(rand.NewSource(t.UnixNano()))
	codeStr := fmt.Sprintf("%s%d", timeStr, r.Intn(100))
	return codeStr
}

func GetOrderId() string {
	now := time.Now()
	//timeStamp := now.Unix()
	timeStamp := now.UnixNano() / 1000000
	r := rand.New(rand.NewSource(now.UnixNano()))
	stampStr := strconv.FormatInt(timeStamp, 10)
	orderId := fmt.Sprintf("%s%d", stampStr, r.Intn(100))
	return orderId
}

func reverseString(s string) string {
	runes := []rune(s)
	if len(s) > 10 {
		runes[7], runes[10] = runes[10], runes[7]
		runes[1], runes[2] = runes[2], runes[1]
	}
	return string(runes)
}

/*校验输入的字符串是否为邮箱格式*/
func IsEmail(email string) bool {
	reg := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	if m := reg.MatchString(email); m {
		return true
	} else {
		return false
	}
}

/*校验手机号是否满足要求（只校验位数）*/
func IsPhone(phone string) bool {
	reg := regexp.MustCompile(`^1[0-9]{10}$`)
	if m := reg.MatchString(phone); m {
		return true
	} else {
		return false
	}
}

func IsContainChinese() {
}

func IsDigital() {

}
