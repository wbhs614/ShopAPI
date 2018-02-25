package models

import (
	"ShopAPI/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"sort"
	"strings"
	"time"
)

type User_token struct {
	Id          int
	Token       string
	Express_in  string
	Appid       string
	Create_time string
}

type Claims struct {
	Appid string `json:"appid"`
	jwt.StandardClaims
}

func init() {
	orm.RegisterModel(new(User_token))
}

func NewToken(r *CreateTokenForm, token string, express_in string) (u *User_token) {
	regDate := time.Now().Format("2006-01-02 15:04:05")
	useToken := User_token{
		Appid:       r.Appid,
		Token:       token,
		Express_in:  express_in,
		Create_time: regDate,
	}
	return &useToken
}

func (self *User_token) insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(self)
	return err
}

func CreateToken(appid, seret string) (string, int64) {
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	claims := Claims{
		Appid: appid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    appid,
		},
	}

	//使用claim创建一个token
	c_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//signs the token with a secret
	signedToken, _ := c_token.SignedString([]byte(seret))
	return signedToken, expireToken

}

func GetUserUidByToken(token, secret string) (uid string, err *utils.ControllerError) {
	o := orm.NewOrm()
	count, errs := o.QueryTable(new(User_token)).Filter("token", token).Count()
	if err != nil {
		if errs == orm.ErrNoRows {
			return "", utils.ErrNoToken
		}
		return "", utils.ErrCheckErrToken
	} else {
		if count > 0 {
			uid, errss := TokenAuth(token, secret)
			if errss != nil {
				return "", utils.ErrExpired
			} else {
				return uid, utils.Actionsuccess
			}
		} else {
			return "", utils.ErrNoToken
		}

	}

}

func TokenAuth(signedToken, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Appid, err
	}
	return "", err
}

//wxpay计算签名的函数
func wxpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	fmt.Println("微信支付签名计算, API KEY:", key)
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}
	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}
