package models

import (
	. "finance/database"
	"finance/types"
	"strings"
)

type Role struct {
	RoleID     int    `json:"role_id"`
	RoleName   string `json:"role_name"`
	RolePsIds  string `json:"role_ps_ids"`
	RolePsCa   string `json:"role_ps_ca"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time"`
	DeleteTime int    `json:"delete_time"`
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
			if (res) {
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
