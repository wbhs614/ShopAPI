// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"ShopAPI/controllers/goods"
	"ShopAPI/controllers/user"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSRouter("/regsiter", &user.RegisterController{}),
			beego.NSRouter("/login", &user.LoginController{}),
			beego.NSRouter("/useInfo", &user.ShopUserController{}),
			beego.NSRouter("/updatUserInfo", &user.ChangeUserInfoController{}),
			beego.NSRouter("/useListInfo", &user.ShopUserListController{}),
			beego.NSRouter("/changePassword", &user.ChangePassWordController{}),
			beego.NSRouter("/logout", &user.LogoutController{}),
			beego.NSRouter("/changeheadir", &user.ChangeHeaderImageController{}),
		),
		beego.NSNamespace("goods",
			beego.NSRouter("/addorUpdateGoods", &goods.AddorUpdateGoodsController{}),
			beego.NSRouter("/addGoodsType", &goods.AddTypeController{}),
			beego.NSRouter("/goodDetail", &goods.GoodDetailController{}),
			beego.NSRouter("/goodList", &goods.GoodListController{}),
			beego.NSRouter("/addorUpdateCarts", &goods.AddOrUpdateCartController{}),
			beego.NSRouter("/getCartList", &goods.CartListController{}),
			beego.NSRouter("/addOrder", &goods.CreateOrderController{}),
			beego.NSRouter("/getOrderDetail", &goods.OrderDetailController{}),
			beego.NSRouter("/getOrderList", &goods.OrderListController{}),
			beego.NSRouter("/getTest", &goods.TestGetController{}),
		),
	)
	beego.AddNamespace(ns)
}
