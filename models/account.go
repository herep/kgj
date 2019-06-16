package models

import (
	"crypto/md5"
	. "finance/database"
	"finance/types"
	"fmt"
	"github.com/astaxie/beego/validation"
	"io"
	"strings"
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
	AccountRole    string `json:"account_role"`
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

	//role 权限入库
	roleids := strings.Replace(strings.Trim(fmt.Sprint(data["AccountRole"]), "[]"), " ", ",", -1)
	//拼接入库信息
	insert := Account{
		AccountName:    data["AccountName"].(string),
		AccountNum:     data["AccountNum"].(string),
		AccountPas:     password,
		AccountMailbox: data["AccountMailbox"].(string),
		AccountPhone:   data["AccountPhone"].(string),
		AccountStatus:  "1", // 默认启动状态
		AccountRole:    roleids,
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
func (A *Account) Saccount(account_num string, code int) (res bool, acc Account) {

	if code == 1 {
		Db.Table("kg_account").Where("account_num = ?", account_num).Find(&acc) // 电话
	} else if code == 2 {
		Db.Table("kg_account").Where("account_name = ?", account_num).Find(&acc) // 用户名
	} else if code == 3 {
		Db.Table("kg_account").Where("account_mailbox = ?", account_num).Find(&acc) // 邮箱
	}

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

	err := Db.Table("kg_account").Where("id = ? ", info.Id).Delete(&infos)

	if err.Error != nil {
		return false, infos
	} else {
		return true, infos
	}
}

//子账户 where 列表
func (A *Account) Wlist(infow map[string]interface{}) (count int64, info []types.ReturnAccount) {

	sql := Db.Table("kg_account").
		Joins("join kg_role on FIND_IN_SET(kg_role.role_id,kg_account.account_role)").
		Joins("join kg_user on kg_account.account_company = kg_user.id").
		Select("GROUP_CONCAT(kg_role.role_name ORDER BY kg_role.role_id) as role_name,kg_user.company_name,"+
			"kg_account.account_name,kg_account.account_num,"+
			"kg_account.account_mailbox,kg_account.id,kg_account.account_phone,"+
			"kg_account.account_status,kg_account.update_time,kg_account.create_time").
		Where("account_status = ? ", infow["AccountStatus"]).Group("kg_account.account_name")

	if infow["AccountNum"] != nil {
		//正则判断 筛选条件
		v := validation.Validation{}
		//条件 -- 希望通过 用户名 -- 电话 -- 邮箱
		errphone := v.Phone(infow["AccountNum"], "account_num").Message("unphone")
		errmalibox := v.Phone(infow["AccountMailbox"], "account_mailbox").Message("unmalibox")

		if errmalibox.Ok {
			sql = sql.Where("account_mailbox = ?", infow["AccountNum"])
		} else if errphone.Ok {
			sql = sql.Where("account_num = ?", infow["AccountNum"])
		} else {
			sql = sql.Where("account_name = ?", infow["AccountNum"])
		}
	}
	sql.Count(&count)

	//分页 查询
	bepage := (infow["Page"].(float64) - 1) * infow["PageSize"].(float64)
	sql.Limit(infow["PageSize"].(float64)).
		Offset(bepage).
		Find(&info)

	return count, info
}
