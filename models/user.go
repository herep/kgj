package models

import (
	. "finance/database"
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
	CreateTime  int64  `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	Consent     int    `json:"consent"`
}

func NewUser() *KgUser {
	return &KgUser{}
}

//注册信息入库
func (U *KgUser) Insertv(info map[string]interface{}) (res bool, id int) {

	//入库数据
	var insert KgUser
	insert.CompanyName = info["company_name"].(string)
	insert.AdminNum = info["admin_num"].(string)
	insert.AdminName = info["admin_name"].(string)
	insert.Password = info["password"].(string)
	insert.PhoneNum = info["phone_num"].(string)
	insert.Mailbox = info["mailbox"].(string)
	insert.CreateTime = info["create_time"].(int64)
	insert.Consent = info["consent"].(int)

	err := Db.Table("kg_account").Create(&insert)

	if err.Error != nil {
		return false, 0
	} else {

		return true, insert.Id
	}
}

//验证账户密码
func (U *KgUser) Checkuser(username string) []KgUser {
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

//验证账户密码
func (U *KgUser) Checkmailboxuser(mailbox string) []KgUser {
	//存储信息
	var userinfo []KgUser

	//查询
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("kg_user").
		Where("account_mailbox = ?")

	sql := qb.String()

	//执行
	o := orm.NewOrm()
	if _, err := o.Raw(sql, mailbox).QueryRows(&userinfo); err != nil {
		fmt.Println(err)
	}
	return userinfo
}

//验证账户密码
func (U *KgUser) Checknameuser(name string) []KgUser {
	//存储信息
	var userinfo []KgUser

	//查询
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("kg_user").
		Where("account_name = ?")

	sql := qb.String()

	//执行
	o := orm.NewOrm()
	if _, err := o.Raw(sql, name).QueryRows(&userinfo); err != nil {
		fmt.Println(err)
	}
	return userinfo
}

//验证账户密码 -- 字段名可变
func (U *KgUser) Checkfuser(username string, field string) []KgUser {
	//存储信息
	var userinfo []KgUser

	//查询
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("kg_user").
		Where(field + " = ?")

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

func (U *KgUser) IdGetIntInfo(userid int) []KgUser {
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
