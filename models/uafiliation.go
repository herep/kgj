package models

import (
	. "finance/database"
	"finance/types"
	"fmt"
)

type Uafiliation struct {
	Id          int     `json:"id"`
	UserId      int     `json:"user_id"`
	AccountId   float64 `json:"account_id"`
	UserName    string  `json:"user_name"`
	UserMailbox string  `json:"user_mailbox"`
	UserPhone   string  `json:"user_phone"`
	Status      int     `json:"status"`
}

func Newuafiliation() *Uafiliation {
	return &Uafiliation{}
}

//判断 主子帐号
func (U *Uafiliation) Suainfo(userPhone string, field string) (info []Uafiliation, err bool) {

	sql := Db.Table("kg_uafiliation").Where(field+"= ? ", userPhone).Find(&info)

	if len(info) != 0 {
		return info, true
	} else {
		fmt.Println(sql.Error) // 查询出错
		return nil, false
	}
}

//维护 主子帐号 关系 -- 新增
func (U *Uafiliation) Iuainfo(info Uafiliation) (res bool) {

	sql := Db.Table("kg_uafiliation").Create(&info)

	if sql.Error != nil {
		return false
	} else {
		return true
	}
}

//维护 主子帐号关系 -- 修改
func (U *Uafiliation) Uuainfo(ua Uafiliation) (res bool) {

	sql := Db.Table("kg_uafiliation").Where("account_id = ? ", ua.AccountId).Update(&ua)

	if sql.Error != nil {
		return false
	} else {
		return true
	}
}

//维护 主子帐号关系 -- 删除
func (U *Uafiliation) Duainfo(ac_id float64) (res bool) {

	var ua Uafiliation
	sql := Db.Table("kg_uafiliation").Where("account_id = ? ", ac_id).Delete(&ua)

	if sql.Error != nil {
		return false
	} else {
		return true
	}
}

//查询 权限
func (U *Uafiliation) SelectRoleinfo(info Uafiliation) (roles []types.RolePsCa) {

	Db.Table("kg_account").
		Select("kg_role.role_ps_ca as role_ps_cas ").
		Joins("join kg_role on FIND_IN_SET(kg_role.role_id,kg_account.account_role)").
		Where("kg_account.account_company = ?", info.UserId).
		Find(&roles)

	return
}
