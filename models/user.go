package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type KgUser struct {
	Id          int    `json:"id"`
	CompanyName string `json:"company_name"`
	AdminNum    string `json:"admin_num"`
	Password    string `json:"password"`
	AdminName   string `json:"admin_name"`
	PhoneNum    string `json:"phone_num"`
	Mailbox     string `json:"mailbox"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	Consent     int    `json:"consent"`
}

func NewUser() *KgUser {
	return &KgUser{}
}

//注册信息入库
func (U *KgUser) Insertv(info map[string]interface{}) (res bool) {

	//入库数据
	company_name := info["company_name"].(string)
	admin_num := info["admin_num"].(string)
	admin_name := info["admin_name"].(string)
	password := info["password"].(string)
	phone_num := info["phone_num"].(string)
	mailbox := info["mailbox"].(string)
	create_time := info["create_time"].(int64)
	consent := info["consent"].(string)


	qb, _ := orm.NewQueryBuilder("mysql")
	qb.InsertInto("kg_user", "kg_user.company_name,kg_user.admin_num,kg_user.admin_name,kg_user.password,"+
		"kg_user.phone_num,kg_user.mailbox,kg_user.create_time,kg_user.consent").Values("?", "?", "?", "?", "?", "?", "?", "?")
	sql := qb.String()

	//执行
	o := orm.NewOrm()
	if _, error := o.Raw(sql, company_name, admin_num, admin_name, password, phone_num, mailbox, create_time, consent).Exec(); error != nil {
		fmt.Println(error)
		return false
	}
	return true
}

//验证账户密码
func (U *KgUser) Checkuser(username string,) []KgUser {
	//存储信息
	var userinfo []KgUser

	//查询
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("kg_user").
		Where("admin_num = ?")

	sql := qb.String()

	//执行
	o := orm.NewOrm()
	if _, err := o.Raw(sql, username).QueryRows(&userinfo); err != nil {
		fmt.Println(err)
	}
	return userinfo
}

//验证账户密码 -- 字段名可变
func (U *KgUser) Checkfuser(username string,field string) []KgUser {
	//存储信息
	var userinfo []KgUser

	//查询
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("kg_user").
		Where(field+" = ?")

	sql := qb.String()

	//执行
	o := orm.NewOrm()
	if _, err := o.Raw(sql, username).QueryRows(&userinfo); err != nil {
		fmt.Println(err)
	}
	return userinfo
}

//验证账户密码
func (U *KgUser) IdGetInfo(userid float64) []KgUser {
	//存储信息
	var userinfo []KgUser

	//查询
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("kg_user").
		Where("id = ?")
	sql := qb.String()

	//执行
	o := orm.NewOrm()
	if _, err := o.Raw(sql, userid).QueryRows(&userinfo); err != nil {
		fmt.Println(err) 
	}
	return userinfo
}
