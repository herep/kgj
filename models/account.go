package models

import (
	"crypto/md5"
	. "finance/database"
	"fmt"
	"io"
	"time"
)

type Account struct {
	Id             int    `json:"id"`
	AccountName    string `json:"account_name"`
	AccountNum     string `json:"account_num"`
	AccountPas     string `json:"account_pas"`
	AccountMailbox string `json:"account_mailbox"`
	AccountPhone   string `json:"account_phone"`
	AccountStatus  string `json:"account_status"`
	AccountRole    int    `json:"account_role"`
	AccountCompany int    `json:"account_company"`
	CreateTime     int64  `json:"create_time"`
	UpdateTime     int64  `json:"update_time"`
}

func Newaccount() *Account {
	return &Account{}
}

//子帐号 入库
func (A *Account) Iaccount(data map[string]interface{}) (res bool, id int) {

	// md5 加密
	w := md5.New()
	io.WriteString(w, data["AccountPas"].(string))
	password := fmt.Sprintf("%x", w.Sum(nil))

	//拼接入库信息
	insert := Account{
		AccountName:    data["AccountName"].(string),
		AccountNum:     data["AccountNum"].(string),
		AccountPas:     password,
		AccountMailbox: data["AccountMailbox"].(string),
		AccountPhone:   data["AccountPhone"].(string),
		AccountStatus:  "1", // 默认启动状态
		AccountRole:    data["AccountRole"].(int),
		AccountCompany: data["AccountCompany"].(int),
		CreateTime:     time.Now().Unix(),
	}

	err := Db.Table("kg_account").Create(&insert)

	if err.Error != nil {
		return false, 0
	} else {

		return true, insert.Id
	}
}

//子帐号 防止重复入库
func (A *Account) Saccount(account_phone string) (res bool, acc Account) {

	Db.Table("kg_account").Where("account_num = ?", account_phone).Find(&acc)

	if acc.Id != 0 {
		return true, acc //id 不为0 存在不允许 新增
	} else {
		return false, acc //id 为0 不存在 允许 新增
	}
}

//子帐号 登录
func (A *Account) Checkacinfo(ac_name string, field string) (data Account, res bool) {

	err := Db.Table("kg_account").Where(field+"= ?", ac_name).Find(&data)

	if data.Id != 0 {

		return data, true
	} else {
		fmt.Println(err)
		return data, false
	}
}

//查询子帐号 -- 对应权限
//func (A *Account) IdGetInfo(userid string) (info Account) {
//
//	//err := Db.Table("kg_account").Select("").Where("id = ?", userid).Joins("left join kg_role on kg_role.id = kg_account.account_role").Row();
//}
