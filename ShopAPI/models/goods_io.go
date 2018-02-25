package models

type TypeForm struct {
	TypeName string `form:"typename"    valid:"Required"`
	Token    string `form:"token" valid:"Required"`
	Secret   string `form:"secret" valid:"Required"`
}

type GoodsForm struct {
	Token      string  `form:"token" valid:"Required"`
	Secret     string  `form:"secret" valid:"Required"`
	Name       string  `form:"name"`
	GoodsId    string  `form:"goodsid"`
	Look       int     `form:"look"`
	Price      float64 `form:"price"`
	Discount   int     `form:"discount"`
	GoodTypeid string  `form:"goodtypeid"`
	ImageUrl   string  `form:"imageurl"`
	ShopTel    string  `form:"shoptel"`
	Detail     string  `form:"detail"`
	StoryName  string  `form:"storyname"`
	Remark     string  `form:"remark"`
	CreateDate string  `form:"createDate"`
	UpdateDate string  `form:"updateDate"`
	Add        int     `form:"add"`
	Waste      int     `form:"waste"`
}

type GoodDetailForm struct {
	GoodId string `form:"goodid"`
}

type GoodListOutForm struct {
	Goodid   string `json:"gooodid"`
	Typeid   string `json:"typeid"`
	TypeName string `json:"typename"`
	Discount string `json:"discount"`
	ImageUrl string `json:"imageUrl"`
	Look     string `json:"look"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Amount   string `json:"amount"`
	Left     string `json:"left"`
}

type GoodsListForm struct {
	GoodTypId string `form:"typeid"`
	Limit     int    `form:"limit"`
	Offset    int    `form:"offset"`
}

type GoodCartForm struct {
	Token   string `form:"token"`
	GoodsId string `form:"goodsid"`
	Amount  int    `form:"amount"`
	Secret  string `form:"secret"`
}

type GoodCartListForm struct {
	Token  string `form:"token"`
	Secret string `form:"secret"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type GoodBuyForm struct {
	GoodId string `json:"goodid"`
	Amount int64  `json:"amount"`
}

type GoodOrderForm struct {
	Token       string `form:"token"`
	Secret      string `form:"secret"`
	BuyGoods    string `form:"buygoods"`
	Postto      string `form:"postto"`
	PostStatus  int    `form:"poststatus"`
	PostComCode string `form:"postcode"`
	Postid      string `form:"postid"`
	Buyer       string `form:"buyer"`
	BuyerTel    string `form:"buyertel"`
	OrderStatus int    `form:"orderstatus"`
}

type GoodsOrderListForm struct {
	Token  string `form:"token"`
	Secret string `form:"secret"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type GoodOrderDetailForm struct {
	Token   string `form:"token"`
	Secret  string `form:"secret"`
	OrderId string `form:"orderid"`
}

type UnitfyOderReq struct {
	Appid            string `xml:"appid"`
	Body             string `xml:"body"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Notify_url       string `xml:"notify_url"`
	Trade_type       string `xml:"trade_type"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
	Total_fee        int    `xml:"total_fee"`
	Out_trade_no     string `xml:"out_trade_no"`
	Sign             string `xml:"sign"`
}
