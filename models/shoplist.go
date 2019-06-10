package models

import (
	. "finance/database"
	"finance/types"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type KgShoplist struct {
	Id             int    `json:"id"`
	Session        string `json:"session"`
	Status         int    `json:"status"`
	SitName        string `json:"site_name"`
	RefreshToken   string `json:"refresh_token"`
	ShopPath       string `json:"shop_path"`
	DeliveryPath   string `json:"delivery_path"`
	TaobaoUserNick string `json:"taobao_user_nick"`
	Restatus       int    `json:"restatus"`
	CompanyId      int    `json:"company_id"`
	CreateTime     string `json:"create_time"`
	UpdateTime     int64  `json:"update_time"`
}

//注册表
func init() {
	orm.RegisterModel(new(KgShoplist))
}

//表 实力对象
func NewShoplist() *KgShoplist {
	return &KgShoplist{}
}

//判断 是否重复添加
func (S *KgShoplist) Chackshop(shopname string, company_id int) []KgShoplist {

	var shopinfo []KgShoplist

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("kg_shoplist").
		Where("taobao_user_nick = ? AND company_id = ?")

	sql := qb.String()

	o := orm.NewOrm()
	if _, err := o.Raw(sql, shopname, company_id).QueryRows(&shopinfo); err != nil {
		fmt.Println(err)
	}
	return shopinfo

}

//tf,o
func (S *KgShoplist) Sinsert(info map[string]interface{}) (result bool) {

	//店铺信息
	shopname := info["shopanme"].(string)
	create_time := info["create_time"].(int64)
	company_id := info["company_id"].(int)

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.InsertInto("kg_shoplist", "kg_shoplist.company_id,kg_shoplist.taobao_user_nick,kg_shoplist.create_time").Values("?", "?", "?")
	sql := qb.String()

	//执行
	o := orm.NewOrm()

	if _, error := o.Raw(sql, company_id, shopname, create_time).Exec(); error != nil {
		fmt.Println(error)
		return false
	} else {
		return true
	}
}

//授权信息 入库
func (S *KgShoplist) Sauthorize(info types.Pcallback) (message string, result bool) {

	//入库数据
	seesion := info.Seesion
	TaobaoUserNick := info.TaobaoUserNick
	refresh_token := info.RefreshToken
	UpdateTime := time.Now().Unix()

	//执行
	o := orm.NewOrm()
	user := KgShoplist{TaobaoUserNick: TaobaoUserNick}

	if err := o.Read(&user, "taobao_user_nick"); err == nil {

		user.Session = seesion
		user.Status = 1
		user.RefreshToken = refresh_token
		user.UpdateTime = UpdateTime

		if _, err := o.Update(&user, "session", "refresh_token", "update_time"); err == nil {
			message := "授权成功"
			return message, true
		} else {
			fmt.Println(err)
			message := "服务器内部错误"
			return message, false
		}

	} else {
		fmt.Println(err)
		message := "授权店铺旺旺号与填写店铺名不一致，请重新填写" 
		return message, false
	}

}

//条件显示 列表
func (S *KgShoplist) Sshoplist(where *types.WshopList) (data []KgShoplist, count int64) {

	//查询 -- 存储数据
	var info []KgShoplist

	sql := Db.Table("kg_shoplist").Where("company_id = ?", where.CompanyId).
		Where("status = ?", where.Status).
		Where("restatus = ?", where.Restatus)

	// 状态查询
	if where.SiteName != "" {
		sql.Where("site_name = ?", where.SiteName) // 站点查询
	}

	if where.Strattime != "" {
		sql.Where("create_time >= ?", where.Strattime) // 创建时间
	}
	if where.Endtime != "" {
		sql.Where("create_time <= ?", where.Endtime) // 创建时间
	}

	if where.TaobaoUserNick != "" {
		sql.Where("taobao_user_nick = ?", where.TaobaoUserNick) // 店铺名
	}

	sql.Count(&count)

	//分页 查询
	bepage := (where.Page - 1) * where.PageSize
	sql.Limit(where.PageSize).Offset(bepage).Find(&info)
	return info, count
}

//完善条件
func (S *KgShoplist) Ushoplist(info map[string]interface{}) (res bool) {

	shop_id := info["id"]

	sql := Db.Table("kg_shoplist").Where("id = ?", shop_id).Update(info)

	if sql.Error != nil {
		fmt.Println(sql.Error)
		return false
	} else {
		return true
	}

}

//id 查询
func (S *KgShoplist) Iidinfo(id float64) []KgShoplist {

	var data []KgShoplist
	sql := Db.Table("kg_shoplist").Where("id = ?", id).Find(&data)

	if sql.Error != nil || len(data) == 0 {
		fmt.Println(sql.Error)
		return nil
	} else {
		return data
	}
}
