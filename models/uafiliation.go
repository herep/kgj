package models

import (
	. "finance/database"
	"fmt"
)

type Uafiliation struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	AccountId   int    `json:"account_id"`
	UserName    string `json:"user_name"`
	UserMailbox string `json:"user_mailbox"`
	UserPhone   string `json:"user_phone"`
	Status      int    `json:"status"`
}

func Newuafiliation() *Uafiliation {
	return &Uafiliation{}
}

//判断 主子帐号
func (U *Uafiliation) Suainfo(userPhone string,field string) (info []Uafiliation, err bool) {

	sql := Db.Table("kg_uafiliation").Where(field +"= ? ", userPhone).Find(&info)

	if len(info) != 0 {
		return info, true
	} else {
		fmt.Println(sql.Error) // 查询出错
		return nil, false
	}
}

//维护 主子帐号 关系
func (U *Uafiliation) Iuainfo(info Uafiliation) (res bool) {

	sql := Db.Table("kg_uafiliation").Create(&info)

	if sql.Error != nil {
		return false
	} else {
		return true
	}
}
