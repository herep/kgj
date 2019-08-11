package types

//搜寻条件
type GetOrderParamenter struct {
	Page     int `json:"page"`      //页数
	PageSize int `json:"page_size"` // 显示数据

	Status       string `json:"status"`        //订单状态
	SellerFlag   string `json:"seller_flag"`   //订单旗帜
	RefundStatus string `json:"refund_status"` // 订单退款状态
	HebingStatus string `json:"hebing_status"` //订单合并状态
	BuyerNick    string `json:"buyer_nick"`    //旺旺号
	Logisticcode string `json:"logisticcode"`  //运单号
	Link         string `json:"link"`          //宝贝编号
	ShipmentName string `json:"shipment_name"` //店铺名
	Oid          string `json:"oid"`           //订单号

	Time        []string  `json:"time"`         //时间条件
	Price       []float64 `json:"price"`        //价格查询
	Num         []float64 `json:"num"`          //数量查询
	NoReceiver  string    `json:"no_receiver"`  //排除城市查询
	YesReceiver string    `json:"yes_receiver"` //显示城市条件
}

type Tidsinfo struct {
	Tid string `json:"tid"`
}
