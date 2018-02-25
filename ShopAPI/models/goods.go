package models

import (
	"ShopAPI/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type GoodType struct {
	Id       int      `json:"-"`
	TypeId   string   `json:"typeid"`
	TypeName string   `json:"typename"`
	Goodses  []*Goods `json:"goodses" orm:"reverse(many)"`
}

type Goods struct {
	Id         int        `json:"-"`
	Name       string     `json:"name"`
	GoodsId    string     `json:"goodsid"`
	Look       int        `json:"look"`
	Price      float64    `json:"price"`
	Discount   int        `json:"discount"`
	GoodType   *GoodType  `json:"goodtype" orm:"rel(fk)"`
	ImageUrl   string     `json:"imageurl"`
	ShopTel    string     `json:"tel"`
	Detail     string     `json:"detail"`
	StoryName  string     `json:"storyname"`
	Remark     string     `json:"remark"`
	CreateDate string     `json:"createDate"`
	UpdateDate string     `json:"updateDate"`
	StoreInfo  *GoodStore `json:"storeinfo" orm:"rel(one)"`
}

type GoodStore struct {
	Id       int    `json:"-"`
	Amount   int    `json:"amount"`
	Sells    int    `json:"sells"`
	Waste    int    `json:"waste"`
	Left     int    `json:"left"`
	GoodInfo *Goods `json:"goodinfo" orm:"reverse(one)"`
}

type GoodsCart struct {
	Id        int    `json:"-"`
	CartGoods *Goods `json:"goods" orm:"rel(fk)"`
	Amount    int    `json:"amount"`
	Uid       string `json:"-"`
}

type GoodsBuy struct {
	Id        int    `json:"-"`
	SingGoods *Goods `orm:"rel(fk)"`
	Amout     int64  `json:"amount"`
	//OneOrder  *GoodOrder `json:"-" orm:"rel(fk)"`
	OrderId string `json:"orderid"`
}

type GoodOrder struct {
	Id          int     `json:"-"`
	OrderId     string  `json:"orderid"`
	OrderName   string  `json:"ordername"`
	PostId      string  `json:"postid"`
	PostCom     string  `json:"postcom"`
	Poststatus  int     `json:"poststatus"`
	Createtime  string  `json:"createtime"`
	Updatetime  string  `json:"updatetime"`
	OrderStatus int     `json:"orderstatus"`
	AmountPrice float64 `json:"payprice"`
	GoodsAmount int64   `json:"goodsnum"`
	Postto      string  `json:"postto"`
	Buyer       string  `json:"buyer"`
	Buyertel    string  `json:"buyerTel"`
	Imageurl    string  `json:"imageurl"`
	Uid         string  `json:"-"`
	//GoodsBuys   []*GoodsBuy `json:"goodsbuy" orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(GoodsBuy))
	orm.RegisterModel(new(GoodType))
	orm.RegisterModel(new(Goods))
	orm.RegisterModel(new(GoodStore))
	orm.RegisterModel(new(GoodsCart))
	orm.RegisterModel(new(GoodOrder))

}

func AddGoodsType(form *TypeForm) *utils.ControllerError {
	_, repon := GetUserUidByToken(form.Token, form.Secret)
	var returnp *utils.ControllerError
	if repon == utils.Actionsuccess {
		goodType := GoodType{
			TypeId:   utils.GetFormatCode(),
			TypeName: form.TypeName,
		}
		o := orm.NewOrm()
		_, err := o.Insert(&goodType)
		if err != nil {
			returnp = utils.ErrDatabase
		} else {
			returnp = utils.Actionsuccess
		}
	} else {
		returnp = repon
	}
	return returnp
}

func AddorUpdateGoods(form *GoodsForm) *utils.ControllerError {
	_, repon := GetUserUidByToken(form.Token, form.Secret)
	var returnp *utils.ControllerError
	o := orm.NewOrm()
	if repon == utils.Actionsuccess {
		//执行插入操作
		if len(form.GoodsId) == 0 {
			goodtype := &GoodType{}
			err := o.QueryTable(new(GoodType)).Filter("TypeId", form.GoodTypeid).One(goodtype)
			if err != nil {
				returnp = utils.ErrCheckGoodsType
			} else {
				if form.Add < form.Waste {
					returnp = utils.ErrWastErrType
					return returnp
				}
				store := &GoodStore{
					Amount: form.Add,
					Left:   form.Add - form.Waste,
					Waste:  form.Waste,
					Sells:  0,
				}
				goods := &Goods{
					Name:       form.Name,
					GoodsId:    utils.GetFormatCode(),
					Price:      form.Price,
					Look:       10,
					Discount:   form.Discount,
					GoodType:   goodtype,
					ShopTel:    form.ShopTel,
					Detail:     form.Detail,
					StoryName:  form.StoryName,
					Remark:     form.Remark,
					CreateDate: utils.GetNowTimeStr(),
					StoreInfo:  store,
				}
				_, erro := o.Insert(store)
				if erro != nil {
					returnp = utils.ErrDatabase
				} else {
					_, err := o.Insert(goods)
					if err == nil {
						returnp = utils.Actionsuccess
					} else {
						returnp = utils.ErrDatabase
					}
				}
			}
		} else {
			//执行更新操作
			goods := &Goods{}
			err := o.QueryTable(new(Goods)).Filter("GoodsId", form.GoodsId).RelatedSel().One(goods)
			if err != nil {
				returnp = utils.ErrDatabase
			} else {
				if len(form.GoodTypeid) > 0 {
					goodtype := &GoodType{}
					err := o.QueryTable(new(GoodType)).Filter("TypeId", form.GoodTypeid).One(goodtype)
					if err != nil {
						returnp = utils.ErrCheckGoodsType
						return returnp
					} else {
						goods.GoodType = goodtype
					}
				}
				if len(form.Name) > 0 {
					goods.Name = form.Name
				}
				if form.Look != 0 {
					goods.Look = form.Look
				}
				if form.Price != 0 {
					goods.Price = form.Price
				}
				if form.Discount != 0 {
					goods.Discount = form.Discount
				}
				if len(form.ShopTel) > 0 {
					goods.ShopTel = form.ShopTel
				}
				if len(form.Detail) > 0 {
					goods.Detail = form.Detail
				}
				if len(form.StoryName) > 0 {
					goods.StoryName = form.StoryName
				}
				if len(form.Remark) > 0 {
					goods.Remark = form.Remark
				}
				store := goods.StoreInfo
				if form.Add != 0 || form.Waste != 0 {
					if (store.Amount + form.Add) < form.Waste {
						returnp = utils.ErrWastErrType
						return returnp
					}
					store.Waste = form.Waste
					store.Left = store.Left + form.Add - form.Waste
					store.Amount = store.Amount + form.Add - form.Waste
					_, err := o.Update(store)
					if err != nil {
						returnp = utils.ErrUpdateStore
						return returnp
					}
				}
				goods.UpdateDate = utils.GetNowTimeStr()
				_, errd := o.Update(goods)
				if errd != nil {
					returnp = utils.ErrDatabase
				} else {
					returnp = utils.Actionsuccess
				}
			}
		}
	} else {
		returnp = repon
	}
	return returnp
}

func (self *Goods) GetGoodDatil(form *GoodDetailForm) *utils.ControllerError {
	o := orm.NewOrm()
	var returnp *utils.ControllerError
	err := o.QueryTable(new(Goods)).Filter("GoodsId", form.GoodId).RelatedSel().Limit(100, 0).One(self)
	if err != nil {
		if err == orm.ErrNoRows {
			returnp = utils.ErrGoodNoExit
		} else {
			returnp = utils.ErrDatabase
		}
	} else {
		returnp = utils.Actionsuccess
	}
	return returnp
}

func GetGoodsList1(form *GoodsListForm) (err *utils.ControllerError, list *[]orm.Params) {
	var maps []orm.Params
	var returnp *utils.ControllerError
	o := orm.NewOrm()
	qs := o.QueryTable(new(Goods)).RelatedSel()
	var limit, offset int
	var errd error
	if form.Limit < 1 {
		limit = 20
	}
	if form.Offset < 1 {
		offset = 0
	}
	if len(form.GoodTypId) == 0 {
		_, errd = qs.Limit(limit, offset).OrderBy("Look").Values(&maps, "name", "GoodType__TypeId", "GoodType__TypeName", "Look", "GoodsId", "Price", "Discount", "ImageUrl", "StoreInfo__Amount", "StoreInfo__Left")
	} else {
		_, errd = qs.Filter("GoodType__TypeId", form.GoodTypId).Limit(limit, offset).OrderBy("Look").Values(&maps, "name", "GoodType__TypeId", "GoodType__TypeName", "Look", "GoodsId", "Price", "Discount", "ImageUrl", "StoreInfo__Amount", "StoreInfo__Left")
	}
	if errd != nil {
		returnp = utils.ErrDatabase
	} else {
		returnp = utils.Actionsuccess
	}
	return returnp, &maps
}

func GetGoodsList(form *GoodsListForm) (err *utils.ControllerError, list *[]orm.Params) {
	var maps []orm.Params
	var returnp *utils.ControllerError
	var errd error
	limit, offset := form.Limit, form.Offset
	if limit < 1 {
		limit = 20
	}
	if offset < 1 {
		offset = 0
	}
	o := orm.NewOrm()
	if len(form.GoodTypId) > 0 {
		_, errd = o.Raw("SELECT T0.`name` `name`, T1.`type_id` `typeid`, T1.`type_name` `typename`, T0.`look` `look`, T0.`goods_id` `goodsid`, T0.`price` `price`, T0.`discount` `discount`, T0.`image_url` `imageurl`, T2.`amount` `amount`, T2.`left` `left` FROM `goods` T0 INNER JOIN `good_type` T1 ON T1.`id` = T0.`good_type_id` INNER JOIN `good_store` T2 ON T2.`id` = T0.`store_info_id` WHERE T1.`type_id` = ? ORDER BY T0.`look` ASC LIMIT ? OFFSET ?", form.GoodTypId, limit, offset).Values(&maps)
	} else {
		_, errd = o.Raw("SELECT T0.`name` `name`, T1.`type_id` `typeid`, T1.`type_name` `typeName`, T0.`look` `look`, T0.`goods_id` `goodsid`, T0.`price` `price`, T0.`discount` `discount`, T0.`image_url` `imageurl`, T2.`amount` `amount`, T2.`left` `left` FROM `goods` T0 INNER JOIN `good_type` T1 ON T1.`id` = T0.`good_type_id` INNER JOIN `good_store` T2 ON T2.`id` = T0.`store_info_id` ORDER BY T0.`look` ASC LIMIT ? OFFSET ?", limit, offset).Values(&maps)
	}
	if errd != nil {
		returnp = utils.ErrDatabase
	} else {
		returnp = utils.Actionsuccess
	}
	return returnp, &maps
}

func AddOrUpdateCarts(form *GoodCartForm) *utils.ControllerError {
	var returnp *utils.ControllerError
	o := orm.NewOrm()
	uid, repon := GetUserUidByToken(form.Token, form.Secret)
	if repon == utils.Actionsuccess {
		cartGood := &GoodsCart{}
		err := o.QueryTable(new(GoodsCart)).Filter("Uid", uid).Filter("CartGoods__GoodsId", form.GoodsId).One(cartGood)
		if err != nil {
			if err == orm.ErrNoRows {
				goods := &Goods{}
				goodsForm := &GoodDetailForm{
					GoodId: form.GoodsId,
				}
				gooderr := goods.GetGoodDatil(goodsForm)
				if gooderr == utils.Actionsuccess {
					cartGood.Amount = form.Amount
					cartGood.CartGoods = goods
					cartGood.Uid = uid
					_, carterr := o.Insert(cartGood)
					if carterr != nil {
						returnp = utils.ErrDatabase
					} else {
						returnp = utils.Actionsuccess
					}
				} else {
					returnp = gooderr
				}
			} else {
				returnp = utils.ErrDatabase
			}
		} else {
			cartGood.Amount = form.Amount
			_, errup := o.Update(cartGood)
			if errup != nil {
				returnp = utils.ErrDatabase
			} else {
				returnp = utils.Actionsuccess
			}
		}
	} else {
		returnp = repon
	}
	return returnp
}

func GetCartList(form *GoodCartListForm) (err *utils.ControllerError, list *[]orm.Params) {
	var returnp *utils.ControllerError
	o := orm.NewOrm()
	uid, repon := GetUserUidByToken(form.Token, form.Secret)
	var maps []orm.Params
	limit, offset := form.Limit, form.Offset
	if limit < 1 {
		limit = 20
	}
	if offset < 1 {
		offset = 0
	}
	if repon == utils.Actionsuccess {
		_, err := o.Raw("SELECT t1.amount cartamout, t2.`name`,t2.goods_id goodsid,t2.detail,t2.discount,t2.create_date createdate,t2.image_url imageurl ,t2.look,t2.price,t2.remark,t2.shop_tel shoptel,t2.story_name storyname,t2.update_date updatedate,t3.type_name typename,t4.amount storeamount,t4.`left` goodleft FROM goods_cart t1 INNER JOIN goods t2 ON t1.cart_goods_id=t2.id INNER JOIN good_type t3 on t2.good_type_id=t3.id INNER JOIN good_store t4 ON t4.id=t2.store_info_id WHERE t1.uid=? ORDER BY t2.create_date DESC limit ? OFFSET ?", uid, limit, offset).Values(&maps)
		if err != nil {
			returnp = utils.ErrDatabase
		} else {
			returnp = utils.Actionsuccess
		}
	} else {
		returnp = repon
	}
	return returnp, &maps
}

func AddOrder(form *GoodOrderForm) *utils.ControllerError {
	var returnp *utils.ControllerError
	o := orm.NewOrm()
	uid, repon := GetUserUidByToken(form.Token, form.Secret)
	if repon == utils.Actionsuccess {
		var goodbuys []GoodBuyForm
		var jsonBody = []byte(form.BuyGoods)
		err := json.Unmarshal(jsonBody, &goodbuys)
		if err != nil {
			errt := utils.ErrGoodIds
			errt.MoreInfo = form.BuyGoods
			returnp = errt
		} else {
			order := &GoodOrder{}
			var amout int64
			var orderName string
			var amountPrice float64
			order.Uid = uid
			orderid := utils.GetOrderId()
			order.OrderId = orderid
			order.OrderStatus = order.OrderStatus
			order.PostId = form.Postid
			order.Poststatus = order.Poststatus
			order.Buyer = form.Buyer
			order.Buyertel = form.BuyerTel
			order.Postto = form.Postto
			order.PostCom = form.PostComCode
			order.Createtime = utils.GetNowTimeStr()
			for i, buy := range goodbuys {
				singlegoods := &Goods{}
				if len(buy.GoodId) > 0 {
					err = o.QueryTable(new(Goods)).Filter("GoodsId", buy.GoodId).One(singlegoods)
					if err != nil {
						returnp = utils.ErrDatabase
					} else {
						amout = amout + buy.Amount
						goodbuy := &GoodsBuy{
							SingGoods: singlegoods,
							Amout:     buy.Amount,
							OrderId:   orderid,
						}
						fmt.Println(goodbuy)
						_, err := o.Insert(goodbuy)
						if err == nil {

						}
						if i == 0 {
							orderName = singlegoods.Name + "*" + strconv.FormatInt(buy.Amount, 10)
						} else {
							orderName = orderName + "," + singlegoods.Name + "*" + strconv.FormatInt(buy.Amount, 10)
						}
						order.OrderName = orderName
						order.GoodsAmount = amout
						amoutS := strconv.FormatInt(buy.Amount, 10)
						amoutF, err := strconv.ParseFloat(amoutS, 64)
						if err == nil {
							amountPrice = amoutF*singlegoods.Price + amountPrice
							order.AmountPrice = amountPrice
						}
					}
				}
			}
			_, err := o.Insert(order)
			if err != nil {
				returnp = utils.ErrDatabase
			} else {
				returnp = utils.Actionsuccess
			}

		}
	} else {
		returnp = repon
	}
	return returnp
}

func GetOrderDetail(form *GoodOrderDetailForm) (*utils.ControllerError, interface{}) {
	var returnp *utils.ControllerError
	result := make(map[string]interface{})
	o := orm.NewOrm()
	_, repon := GetUserUidByToken(form.Token, form.Secret)
	if repon == utils.Actionsuccess {
		var maps []orm.Params
		num, err := o.Raw("SELECT T0.order_name ordername,T0.post_id postid,T0.post_com postcom,T0.postto,T0.poststatus,T0.createtime,T0.updatetime,T0.amount_price amountp,T0.buyer,T0.buyertel,T0.imageurl  from good_order T0 where T0.order_id=?", form.OrderId).Values(&maps)
		if err == nil && num > 0 {
			result["oerderInfo"] = &maps[0]
			returnp = utils.Actionsuccess
		} else {
			returnp = utils.ErrDatabase
		}
		var list []orm.Params
		_, listrr := o.Raw("SELECT t0.name,t0.goods_id,t0.price,t0.image_url imageurl,t1.type_id goodstypeid,t1.type_name typename,(SELECT amout from goods_buy t2 WHERE order_id=? and t0.id=t2.sing_goods_id) amout FROM goods t0 INNER JOIN good_type t1 ON t0.good_type_id=t1.id where t0.id IN (SELECT sing_goods_id FROM goods_buy  where order_id=?) ", form.OrderId, form.OrderId).Values(&list)
		if listrr == nil {
			result["goodlist"] = list
			returnp = utils.Actionsuccess
		} else {
			returnp = utils.ErrDatabase
		}

	} else {
		returnp = repon
	}
	return returnp, result
}

func GetOrderList(form *GoodsOrderListForm) (*utils.ControllerError, interface{}) {
	var returnp *utils.ControllerError
	var list []orm.Params
	o := orm.NewOrm()
	uid, repon := GetUserUidByToken(form.Token, form.Secret)
	if repon == utils.Actionsuccess {
		limit, offset := form.Limit, form.Offset
		if limit < 1 {
			limit = 20
		}
		if offset < 1 {
			offset = 0
		}
		_, listrr := o.Raw("SELECT t0.order_name ordername,t0.order_id orderids,t0.imageurl,t0.amount_price amoutprice,t0.goods_amount,t0.order_status orderstatus FROM good_order t0 WHERE uid=? LIMIT ? OFFSET ? ", uid, limit, offset).Values(&list)
		if listrr != nil {
			returnp = utils.ErrInputData
		} else {
			returnp = utils.Actionsuccess
		}

	} else {
		returnp = repon
	}
	return returnp, list
}
