package models

import (
	. "finance/database"
	"finance/types"
	"github.com/astaxie/beego/validation"
	"regexp"
	"strings"
)

type Shipments struct {
	Id                int     `json:"id"`
	ShipmentId        int     `json:"shipment_id"`
	ShipmentName      string  `json:"shipment_name"`
	UserId            int     `json:"user_id"`
	BuyerNick         string  `json:"buyer_nick"`
	PicPath           string  `json:"pic_path"`
	Payment           float64 `json:"payment"`
	Shop              string  `json:"shop"`
	Link              string  `json:"link"`
	Price             float64 `json:"price"`
	TotalFee          int     `json:"total_fee"`
	Num               int     `json:"num"`
	Tid               string  `json:"tid"`
	Oid               string  `json:"oid"`
	Zid               int     `json:"zid"`
	PostFee           int     `json:"post_fee"`
	ReceiverName      string  `json:"receiver_name"`
	ReceiverState     string  `json:"receiver_state"`
	ReceiverCity      string  `json:"receiver_city"`
	ReceiverDistrict  string  `json:"receiver_district"`
	ReceiverTown      string  `json:"receiver_town"`
	ReceiverZip       string  `json:"receiver_zip"`
	ReceiverAddress   string  `json:"receiver_address"`
	ReceiverMobile    string  `json:"receiver_mobile"`
	ReceiverPhone     string  `json:"receiver_phone"`
	Title             string  `json:"title"`
	Created           string  `json:"created"`
	PayTime           string  `json:"pay_time"`
	SellerFlag        string  `json:"seller_flag"`
	SellerMemo        string  `json:"seller_memo"`
	BuyerMessage      string  `json:"buyer_message"`
	ItemMealName      string  `json:"item_meal_name"`
	RefundStatus      string  `json:"refund_status"`
	OuterIid          string  `json:"outer_iid"`
	SkuId             string  `json:"sku_id"`
	SnapshotUrl       string  `json:"snapshot_url"`
	SkuPropertiesName string  `json:"sku_properties_name"`
	Status            string  `json:"status"`
	Kuaidi            string  `json:"kuaidi"`
	Logisticcode      string  `json:"logisticcode"`
	LogisticcodeTime  string  `json:"logisticcode_time"`
	ReviewStatus      int     `json:"review_status"`
	ReviewName        string  `json:"review_name"`
	ReviewValue       string  `json:"review_value"`
	ReviewTime        string  `json:"review_time"`
	Audit             int     `json:"audit"`
	ReviewUid         int     `json:"review_uid"`
	HebingStatus      int     `json:"hebing_status"`
}

func NewShipments() *Shipments {
	return &Shipments{}
}

//订单搜索
func (S *Shipments) GetShipmentsInfo(paramenters types.GetOrderParamenter, userinfo Uafiliation) (result_create bool, orderinfo []Shipments) {
	sql := Olddb.Table("ap_shipments").Where("user_id = ?", userinfo.UserId)

	//-------> oid条件
	if paramenters.Oid != "" {
		sql = sql.Where("oid = ?", paramenters.Oid)
	}

	//------->旺旺号查询
	if paramenters.BuyerNick != "" {
		sql = sql.Where("buyer_nick = ?", paramenters.BuyerNick)
	}

	//------->订单号查询
	if paramenters.ShipmentName != "" {
		sql = sql.Where("shipment_name = ?", paramenters.ShipmentName)
	}

	//------->宝贝编号查询
	if paramenters.Link != "" {
		sql = sql.Where("link = ?", paramenters.Link)
	}

	//------->订单旗帜查询
	if paramenters.SellerFlag != "" {
		sql = sql.Where("seller_flag = ?", paramenters.SellerFlag)
	} else {
		sql = sql.Where("seller_flag in (?)", []string{"0", "1", "4", "5"})
	}

	//------->退货状态查询
	if paramenters.RefundStatus == "" {
		sql = sql.Where("refund_status in (?)", []string{"", "卖家拒绝退款", "NO_REFUND"})
	} else {
		sql = sql.Where("refund_status in (?)", []string{"买家已经申请退款，等待卖家同意", "卖家已经同意退款，等待买家退货", "买家已经退货，等待卖家确认收货"})
	}

	//------->合并订单查询
	if paramenters.HebingStatus != "" {
		sql = sql.Where("hebing_status = ?", "1")
	} else {
		sql = sql.Where("hebing_status = ?", "0")
	}

	//------->订单状态查询
	if paramenters.Status != "" {
		sql = sql.Where("status = ?", paramenters.Status).Where("review_name = ?", userinfo.UserName)
	}

	//------->运单号查询
	if paramenters.Logisticcode != "" {
		sql = sql.Where("logisticcode = ?", paramenters.Logisticcode)
	}

	//------->时间查询
	if len(paramenters.Time) > 0 {
		if paramenters.Time[0] == "" {
			paramenters.Time[0] = "2017-01-01 00:00:00"
		}
		if paramenters.Time[1] == "" {
			paramenters.Time[1] = "2030-01-01 00:00:00"
		}
		sql = sql.Where("pay_time BETWEEN ? AND ?", paramenters.Time[0], paramenters.Time[1])
	}

	//------->数量查询
	if len(paramenters.Num) > 0 {
		if paramenters.Num[1] != 0 {
			minnum := paramenters.Num[0]
			maxnum := paramenters.Num[1]
			sql = sql.Where("num BETWEEN ? AND ?", minnum, maxnum)
		} else if paramenters.Num[0] != 0 {
			sql = sql.Where("num = ?", paramenters.Num[0])
		}
	}

	//------->排除城市查询
	if paramenters.NoReceiver != "" {
		//省 省/市 省/市/区
		v := validation.Validation{}
		v.Match(paramenters.NoReceiver, regexp.MustCompile("/,/i"), "no_receiver")

		if v.HasErrors() {
			//省
			sql = sql.Where("receiver_state != ?", paramenters)
		} else {
			//切割
			no_receiverinfo := strings.Split(paramenters.NoReceiver, ",")
			if no_receiverinfo[2] != "" {
				//省/市
				sql = sql.Where("receiver_state != ?", no_receiverinfo[0])
				sql = sql.Where("receiver_city != ?", no_receiverinfo[1])
			} else {
				//省/市/区
				sql = sql.Where("receiver_state != ?", no_receiverinfo[0])
				sql = sql.Where("receiver_city != ?", no_receiverinfo[1])
				sql = sql.Where("receiver_district != ?", no_receiverinfo[2])
			}
		}
	}

	//------->选中城市查询
	if paramenters.YesReceiver != "" {
		//省 省/市 省/市/区
		v := validation.Validation{}
		v.Match(paramenters.YesReceiver, regexp.MustCompile("/,/i"), "yes_receiver")

		if v.HasErrors() {
			//省
			sql = sql.Where("receiver_state == ?", paramenters)
		} else {
			//切割
			no_receiverinfo := strings.Split(paramenters.NoReceiver, ",")
			if no_receiverinfo[2] != "" {
				//省/市
				sql = sql.Where("receiver_state == ?", no_receiverinfo[0])
				sql = sql.Where("receiver_city == ?", no_receiverinfo[1])
			} else {
				//省/市/区
				sql = sql.Where("receiver_state == ?", no_receiverinfo[0])
				sql = sql.Where("receiver_city == ?", no_receiverinfo[1])
				sql = sql.Where("receiver_district == ?", no_receiverinfo[2])
			}
		}
	}

	//------->查询店铺信息
	shipids := NewShoplist().UseridGetSelectinfo(userinfo.UserId)
	if len(shipids) != 0 { //判断是否存在店铺
		shipid := []int{}
		for _, val := range shipids {
			shipid = append(shipid, val.Id)
		}
		sql = sql.Where("shipment_id in (?)", shipid)
	}

	//查看sql
	//fmt.Println(sql.LogMode(true))
	//分页 查询
	bepage := (paramenters.Page - 1) * paramenters.PageSize
	var tids []types.Tidsinfo
	sql.Select("tid").
		Limit(paramenters.PageSize).Offset(bepage).
		Order("pay_time desc").
		Group("tid").Find(&tids)

	if len(tids) != 0 {
		result_create = true
		var Atids []string
		for _, tid := range tids {
			Atids = append(Atids, tid.Tid)
		}
		Olddb.Table("ap_shipments").Where("tid in (?)", Atids).Find(&orderinfo)
	} else {
		result_create = false
	}
	return result_create, orderinfo
}
