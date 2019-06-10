package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type KgBindcode struct {
	Id         int    `json:"id"`
	Userid     int    `json:"userid"`
	Phone      string `json:"phone"`
	Verifycode string `json:"verifycode"`
	Status     int    `json:"status"`
	Rescode    string `json:"rescode"`
	Resmessage string `json:"resmessage"`
	Createtime int64  `json:"createtime"`
	Expirytime int64  `json:"expirytime"`
}

func NewBindcode() *KgBindcode {
	return &KgBindcode{}
}

//注册---验证
func (K *KgBindcode) CheckBinde(phone string) []KgBindcode {
	//存储数据
	var info []KgBindcode

	//构建查询
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("kg_bindcode").
		Where("phone = ? ").
		OrderBy("kg_bindcode.createtime").
		Desc().Limit(1)
	sql := qb.String()

	//执行
	o := orm.NewOrm()
	if _, err := o.Raw(sql, phone).QueryRows(&info); err != nil {
		fmt.Println(err)
	}
	return info
}

//注册--短信入库
func (K *KgBindcode) InsertBcode(data map[string]interface{}) (res map[string]interface{}) {
	//message
	var message map[string]interface{}
	code := data["code"].(int)

	//接受参数
	if code == 0 {
		//0 发短信成功 入库
		phone := data["phone"].(string)
		verifycode := data["VerifyCode"].(string)
		createtime := data["CreateTime"].(int64)
		expirytime := data["ExpiryTime"].(int64)

		qb, _ := orm.NewQueryBuilder("mysql")
		// 构建查询对象
		qb.InsertInto("kg_bindcode", "kg_bindcode.phone", "kg_bindcode.verifycode", "kg_bindcode.createtime", "kg_bindcode.expirytime").
			Values("?", "?", "?", "?")
		//返回sql语句
		sql := qb.String()
		// 执行 SQL 语句
		o := orm.NewOrm()
		if _, error := o.Raw(sql, phone, verifycode, createtime, expirytime).Exec(); error != nil {
			message = map[string]interface{}{"status": 400, "message": "入库失败"}
		}
		message = map[string]interface{}{"status": 200, "message": "入库成功"}

	} else if code == 1 { 
		//1 发短信失败 入库
		phone := data["phone"].(string)
		verifycode := data["VerifyCode"].(string)
		createtime := data["CreateTime"].(int64)
		expirytime := data["ExpiryTime"].(int64)
		rescode := data["ResCode"].(float64)
		resmessage := data["ResMessage"].(string)

		qb, _ := orm.NewQueryBuilder("mysql")

		// 构建查询对象
		qb.InsertInto("kg_bindcode", "kg_bindcode.phone", "kg_bindcode.verifycode", "kg_bindcode.createtime", "kg_bindcode.expirytime", "kg_bindcode.rescode", "kg_bindcode.resmessage").
			Values("?", "?", "?", "?", "?", "?")

		//返回sql语句
		sql := qb.String()
		// 执行 SQL 语句
		o := orm.NewOrm()
		if _, error := o.Raw(sql, phone, verifycode, createtime, expirytime, rescode, resmessage).Exec(); error != nil {
			message = map[string]interface{}{"status": 400, "message": "入库失败"}
		}
		message = map[string]interface{}{"status": 200, "message": "入库成功"}
	}

	return message
}
