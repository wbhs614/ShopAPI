package goods

import (
	"ShopAPI/models"
	"ShopAPI/utils"
	// "fmt"
	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/orm"
	// "github.com/mingzhehao/goutils/filetool"
	// "math/rand"
	// "strings"
	//    "time"
)

type AddTypeController struct {
	beego.Controller
}

type AddorUpdateGoodsController struct {
	beego.Controller
}

type GoodDetailController struct {
	beego.Controller
}

type GoodListController struct {
	beego.Controller
}

type AddOrUpdateCartController struct {
	beego.Controller
}

type CartListController struct {
	beego.Controller
}

type CreateOrderController struct {
	beego.Controller
}

type OrderDetailController struct {
	beego.Controller
}

type OrderListController struct {
	beego.Controller
}

type TestGetController struct {
	beego.Controller
}

func (self *AddTypeController) Post() {
	form := &models.TypeForm{}
	self.ParseForm(form)
	if len(form.Token) == 0 || len(form.TypeName) == 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		err := models.AddGoodsType(form)
		self.Data["json"] = err
		self.ServeJSON()
	}
}

func (self *AddorUpdateGoodsController) Post() {
	form := &models.GoodsForm{}
	self.ParseForm(form)
	if len(form.Name) == 0 || len(form.Secret) == 0 || len(form.Token) == 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		if len(form.GoodsId) == 0 {
			if len(form.GoodTypeid) == 0 {
				self.Data["json"] = utils.ErrInputData
				self.ServeJSON()
				return
			}
		}
		err := models.AddorUpdateGoods(form)
		self.Data["json"] = err
		self.ServeJSON()
		return
	}
}

func (self *GoodDetailController) Post() {
	form := &models.GoodDetailForm{}
	self.ParseForm(form)
	good := &models.Goods{}
	if len(form.GoodId) == 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		err := good.GetGoodDatil(form)
		if err == utils.Actionsuccess {
			self.Data["json"] = utils.CommendFormat{
				Code:    9000,
				Message: "获取信息成功",
				Data: map[string]interface{}{
					"goodsname":  good.Name,
					"goodid":     good.GoodsId,
					"price":      good.Price,
					"look":       good.Look,
					"dicount":    good.Discount,
					"shoptel":    good.ShopTel,
					"shopname":   good.StoryName,
					"remak":      good.Remark,
					"deatil":     good.Detail,
					"updateDate": good.UpdateDate,
					"createDate": good.CreateDate,
					"typeid":     good.GoodType.TypeId,
					"typename":   good.GoodType.TypeName,
					"sells":      good.StoreInfo.Sells,
					"left":       good.StoreInfo.Left,
					"imageurl":   good.ImageUrl,
				},
			}
			self.ServeJSON()
		} else {
			self.Data["json"] = err
			self.ServeJSON()
		}
	}
}

func (self *GoodListController) Post() {
	form := &models.GoodsListForm{}
	self.ParseForm(form)
	err, list := models.GetGoodsList(form)
	if err == utils.Actionsuccess {
		self.Data["json"] = utils.CommendFormat{
			Code:    9000,
			Message: "获取商品列表成功",
			Data:    list,
		}
		self.ServeJSON()
	} else {
		self.Data["json"] = err
		self.ServeJSON()
	}
}

func (self *AddOrUpdateCartController) Post() {
	form := &models.GoodCartForm{}
	self.ParseForm(form)
	if len(form.GoodsId) == 0 || len(form.Secret) == 0 || len(form.Token) == 0 || form.Amount == 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		err := models.AddOrUpdateCarts(form)
		self.Data["json"] = err
		self.ServeJSON()
	}
}

func (self *CartListController) Post() {
	form := &models.GoodCartListForm{}
	self.ParseForm(form)
	if len(form.Secret) == 0 || len(form.Token) == 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		err, list := models.GetCartList(form)
		if err == utils.Actionsuccess {
			self.Data["json"] = utils.CommendFormat{
				Code:    9000,
				Message: "获取购物车列表成功",
				Data:    list,
			}
			self.ServeJSON()
		} else {
			self.Data["json"] = err
			self.ServeJSON()
		}
	}

}

func (self *CreateOrderController) Post() {
	form := &models.GoodOrderForm{}
	self.ParseForm(form)
	if len(form.BuyGoods) == 0 || len(form.Token) == 0 || len(form.Secret) == 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		err := models.AddOrder(form)
		self.Data["json"] = err
		self.ServeJSON()
	}
}

func (self *OrderDetailController) Post() {
	form := &models.GoodOrderDetailForm{}
	self.ParseForm(form)
	if len(form.Token) == 0 || len(form.Secret) == 0 || len(form.OrderId) == 0 {
		self.Data["json"] = utils.ErrInputData
		self.ServeJSON()
		return
	} else {
		err, data := models.GetOrderDetail(form)
		if err == utils.Actionsuccess {
			self.Data["json"] = utils.CommendFormat{
				Code:    9000,
				Message: "获取订单详情成功",
				Data:    data,
			}
		} else {
			self.Data["json"] = err
		}
		self.ServeJSON()
	}
}

func (self *OrderListController) Post() {
	form := &models.GoodsOrderListForm{}
	self.ParseForm(form)
	if len(form.Token) == 0 || len(form.Secret) == 0 {
		self.Data["json"] = form
		self.ServeJSON()
		return
	} else {
		err, data := models.GetOrderList(form)
		if err == utils.Actionsuccess {
			self.Data["json"] = utils.CommendFormat{
				Code:    9000,
				Message: "获取订单成功",
				Data:    data,
			}
		} else {
			self.Data["json"] = err
		}
		self.ServeJSON()
	}
}

func (self *TestGetController) Get() {

	token := self.GetString("token")
	limit, limiterr := self.GetInt("limit")
	offset, offseterr := self.GetInt("offset")
	if limiterr != nil {
		limit = 20
	}
	if offseterr != nil {
		offset = 1
	}
	self.Data["json"] = utils.CommendFormat{
		Code:    9000,
		Message: "获取信息成功",
		Data: map[string]interface{}{
			"token":      token,
			"limit":      limit,
			"offset":     offset,
			"body":       self.Ctx.Request.Body,
			"inputToken": self.Input().Get("token"),
		},
	}
	self.ServeJSON()
}

func (self *TestGetController) Post() {

	token := self.GetString("token")
	limit, limiterr := self.GetInt("limit")
	offset, offseterr := self.GetInt("offset")
	if limiterr != nil {
		limit = 20
	}
	if offseterr != nil {
		offset = 1
	}
	self.Data["json"] = utils.CommendFormat{
		Code:    9000,
		Message: "获取信息成功",
		Data: map[string]interface{}{
			"token":     token,
			"limit":     limit,
			"offset":    offset,
			"body":      self.Ctx.Request.Body,
			"domain":    self.Ctx.Input.Domain(),
			"ip":        self.Ctx.Input.IP(),
			"method":    self.Ctx.Input.Method(),
			"isAjax":    self.Ctx.Input.IsAjax(),
			"inputData": self.Ctx.Input.GetData("sdnifn"),
			"param":     self.Ctx.Input.Param("token"),
		},
	}
	self.ServeJSON()
}
