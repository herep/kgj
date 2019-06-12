package models

import (
	"crypto/md5"
	. "finance/database"
	"fmt"
	"io"
	"time"
)

type Account struct {
	Id             int     `json:"id"`
	AccountName    string  `json:"account_name"`
	AccountNum     string  `json:"account_num"`
	AccountPas     string  `json:"account_pas"`
	AccountMailbox string  `json:"account_mailbox"`
	AccountPhone   string  `json:"account_phone"`
	AccountStatus  string  `json:"account_status"`
	AccountRole    float64 `json:"account_role"`
	AccountCompany int     `json:"account_company"`
	CreateTime     int64   `json:"create_time"`
	UpdateTime     int64   `json:"update_time"`
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
		AccountRole:    data["AccountRole"].(float64),
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

//查询子帐号
func (A *Account) IdGetInfo(userid int) (info Account, res bool) {

	err := Db.Table("kg_account").Where("id = ?", userid).Find(&info)

	if info.Id != 0 {

		return info, true
	} else {
		fmt.Println(err)
		return info, false
	}
}

//修改子账户
func (A *Account) Uinfo(info Account) (res bool, infos Account) {

	//整合 修改数据
	info.UpdateTime = time.Now().Unix()
	err := Db.Table("kg_account").Where("id = ? ", info.Id).Update(info)

	if err.Error != nil {
		return false, info
	} else {
		return true, info
	}
}

//删除子账户
func (A *Account) Dinfo(info Account) (res bool, infos Account) {

	err := Db.Table("kg_account").Where("id = ? ", info.Id).Delete(&info)

	if err.Error != nil {
		return false, infos
	} else {
		return true, infos
	}
}
