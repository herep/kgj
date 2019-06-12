package models

import (
	"bytes"
	. "finance/database"
	"finance/types"
	"fmt"
	"strings"
	"time"
)

type Role struct {
	RoleID     int    `json:"role_id"`
	RoleName   string `json:"role_name"`
	RolePsIds  string `json:"role_ps_ids"`
	RolePsCa   string `json:"role_ps_ca"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
	DeleteTime int64  `json:"delete_time"`
}

func Newrole() *Role {
	return &Role{}
}

//role -- permission
func (R *Role) RoleLsit() (item map[int]map[string][]interface{}, res bool) {

	var roles []types.Roles
	//role -- permission 内容
	Db.Table("kg_role").Select("kg_role.role_id,kg_role.role_name,group_concat(kg_permission.ps_name ORDER BY kg_permission.ps_id DESC) as role_names").
		Joins("join kg_permission on FIND_IN_SET(kg_permission.ps_id,kg_role.role_ps_ids)").Group("kg_role.role_name").
		Find(&roles)

	if len(roles) != 0 {
		item = make(map[int]map[string][]interface{})
		//整合数据
		for k, v := range roles {
			//多个权限
			res := strings.Contains(v.RoleNames, ",")
			if res {
				items := make(map[string][]interface{})
				arr := strings.Split(v.RoleNames, ",")
				items["son"] = append(items["son"], arr)
				items["name"] = append(items["name"], v.RoleName)
				items["role_id"] = append(items["role_id"], v.RoleId)
				item[k] = items

			} else {
				items := make(map[string][]interface{})
				items["son"] = append(items["son"], v.RoleNames)
				items["name"] = append(items["name"], v.RoleName)
				items["role_id"] = append(items["role_id"], v.RoleId)
				item[k] = items
			}
		}
		return item, true
	}
	return item, false
}

//role -- 分配
func (R *Role) Roledistribution(item map[string]interface{}) (res bool) {

	//入库数据 整合
	var items Role
	items.RoleName = item["RoleName"].(string)
	// 权限分组名称
	items.RolePsIds = strings.Replace(strings.Trim(fmt.Sprint(item["RolePsIds"]), "[]"), " ", ",", -1) // 所属 权限 id

	// 库内权限格式  控制器-方法名  拼接
	var Cas []types.Ca
	Db.Table("kg_permission").Where("ps_id in (?)", item["RolePsIds"]).Select("concat(ps_c,'-',ps_a)as ca").Find(&Cas) // 查询分配权限名称

	var buffer bytes.Buffer
	for _, v := range Cas {
		buffer.WriteString(strings.Trim(fmt.Sprint(v), "{}"))
		buffer.WriteString(",")
	}
	s := buffer.String()
	i := len([]rune(s)) //字符串长度
	items.RolePsCa = string([]rune(s)[:i-1])
	items.CreateTime = time.Now().Unix()

	sql := Db.Table("kg_role").Create(&items)
	if sql.Error != nil {

		return false
	} else {
		return true
	}
}

//role -- name 条件去重复
func (R *Role) RoleRepeat(rolename string) (info Role) {

	Db.Table("kg_role").Where("role_name = ?", rolename).Find(&info)
	return
}

//role -- id 条件去重复
func (R *Role) RoleidRepeat(roleid float64) (info Role) {

	Db.Table("kg_role").Where("role_id = ?", roleid).Find(&info)
	return
}

//role -- id 修改
func (R *Role) RoleUp(Uinfo map[string]interface{}) (res bool) {

	//改 数据库
	var item Role
	item.RolePsIds = strings.Replace(strings.Trim(fmt.Sprint(Uinfo["RolePsIds"]), "[]"), " ", ",", -1) // 所属 权限 id

	// 库内权限格式  控制器-方法名  拼接
	var Cas []types.Ca
	Db.Table("kg_permission").Where("ps_id in (?)", Uinfo["RolePsIds"]).Select("concat(ps_c,'-',ps_a)as ca").Find(&Cas) // 查询分配权限名称

	// 控制器-方法名 role_ps_ca 拼接
	var buffer bytes.Buffer
	for _, v := range Cas {
		buffer.WriteString(strings.Trim(fmt.Sprint(v), "{}"))
		buffer.WriteString(",")
	}

	s := buffer.String()
	i := len([]rune(s)) //字符串长度
	item.RolePsCa = string([]rune(s)[:i-1])
	item.UpdateTime = time.Now().Unix()

	//sql 修改
	err := Db.Table("kg_role").Where("role_id = ?", Uinfo["RoleID"]).Update(item)
	if err.Error != nil {

		return false
	} else {
		return true
	}
}

//role -- id 删除
func (R *Role) RoleDe(role_id float64) (res bool) {

	var info Role
	err := Db.Table("kg_role").Where("role_id = ?", role_id).Delete(&info)
	if err.Error != nil {

		return false
	} else {
		return true
	}
}
