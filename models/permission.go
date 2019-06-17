package models

import (
	. "finance/database"
)

type Permission struct {
	PsID       int    `json:"ps_id"`
	PsName     string `json:"ps_name"`
	PsPid      int    `json:"ps_pid"`
	PsC        string `json:"ps_c"`
	PsA        string `json:"ps_a"`
	PsLevel    int64  `json:"ps_level"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
	DeleteTime int64  `json:"delete_time"`
}

func NewPermission() *Permission {
	return &Permission{}
}

//查询 父级权限等级
func (P *Permission) Sfpslevel(pslevel int) (level Permission, res bool) {

	Db.Table("kg_permission").Where("ps_id = ?", pslevel).Find(&level)

	if level.PsID != 0 {
		return level, true
	} else {
		return level, false
	}
}

//权限入库操作
func (P *Permission) Ipsinfo(info Permission) (res bool, id int) {

	err := Db.Table("kg_permission").Create(&info)

	if err.Error != nil {
		return false, 0
	} else {

		return true, info.PsID
	}
}

//权限修改操作
func (P *Permission) Upsinfo(info Permission) (res bool) {

	err := Db.Table("kg_permission").Where("ps_id = ?", info.PsID).Update(info)

	if err.Error != nil {

		return false
	} else {
		return true
	}
}

//权限删除操作
func (P *Permission) Dpsinfo(info Permission) (res bool) {

	err := Db.Table("kg_permission").Where("ps_id = ?", info.PsID).Delete(&info)

	if err.Error != nil {

		return false
	} else {
		return true
	}
}

//权限 id查询 数据
func (P *Permission) Sidinfo(id int) (info Permission, res bool) {

	err := Db.Table("kg_permission").Where("ps_id = ?", id).Find(&info)

	if err.Error != nil {
		return info, false
	} else {
		return info, true
	}
}

//权限列表
func (P *Permission) Sperlist() (infos []Permission, res bool) {

	err := Db.Table("kg_permission").Find(&infos)

	if err.Error != nil {
		return infos, false
	} else {
		return infos, true
	}
}
