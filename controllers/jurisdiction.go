package controllers

import (
	"bytes"
	"encoding/json"
	"finance/models"
	"finance/types"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"time"
)

type JurisdictionController struct {
	beego.Controller
}

//权限 写入
func (this *JurisdictionController) Ijuinfo() {

	var data models.Permission
	//接受数据
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)

	if err != nil {
		fmt.Println(err)
	}

	//验证器
	v := validation.Validation{}
	v.Required(data.PsName, "ps_name").Message("权限名不可以为空！")
	v.Required(data.PsPid, "ps_pid").Message("请选择父级权限！")
	v.Required(data.PsC, "ps_c").Message("控制器不可以为空！")
	v.Required(data.PsA, "ps_a").Message("方法名不可以为空！")

	//验证消息 提醒
	if v.HasErrors() {

		var message string
		for _, err := range v.Errors {
			message = message + err.Message + ""
		}

		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	} else {

		//入库信息整合
		data.CreateTime = time.Now().Unix()
		if data.PsPid == 0 {
			data.PsLevel = 0
		} else {
			//查询父级等级 新增等级 父级+1
			info, res := models.NewPermission().Sfpslevel(data.PsPid)
			if res {
				data.PsLevel = info.PsLevel + 1
			}
		}

		//入库
		ress, _ := models.NewPermission().Ipsinfo(data)
		if ress {
			this.Data["json"] = types.Successre{Status: 200, Message: "新增成功", Code: 1}
		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "新增失败，请联系管理员", Code: -1}
		}
	}
	this.ServeJSON()
}

//权限 修改
func (this *JurisdictionController) Ujurinfo() {

	//接受数据
	var data models.Permission
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)

	if err != nil {
		fmt.Println(err)
	}

	//验证器
	v := validation.Validation{}
	v.Required(data.PsID, "ps_id").Message("请选则需要修改的权限")

	//错误信息
	if v.HasErrors() {

		var message string
		for _, err := range v.Errors {
			message = err.Message + ""
		}
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	} else {

		_, resulte := models.NewPermission().Sidinfo(data.PsID)

		if resulte {
			//入库信息整合
			data.UpdateTime = time.Now().Unix()
			if data.PsPid == 0 {
				data.PsLevel = 0
			} else {
				//查询父级等级 新增等级 父级+1
				info, res := models.NewPermission().Sfpslevel(data.PsPid)
				if res {
					data.PsLevel = info.PsLevel + 1
				}
			}

			//修改操作
			res := models.NewPermission().Upsinfo(data)

			if res {
				this.Data["json"] = types.Successre{Status: 200, Message: "修改成功", Code: 1}
			} else {
				this.Data["json"] = types.Successre{Status: 400, Message: "服务器内部错误", Code: -1}
			}

		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "信息不存在", Code: -1}
		}
	}
	this.ServeJSON()
}

//权限删除
func (this *JurisdictionController) Djurinfo() {

	//结收值
	var data models.Permission
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)

	if err != nil {
		fmt.Println(err)
	}

	//验证器
	v := validation.Validation{}
	v.Required(data.PsID, "ps_id").Message("请选择需要删除权限！")

	//验证消息
	if v.HasErrors() {
		var message string
		for _, err := range v.Errors {
			message = err.Message + ""
		}

		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	} else {

		_, reslute := models.NewPermission().Sidinfo(data.PsID)

		if reslute {
			res := models.NewPermission().Dpsinfo(data)

			if res {
				this.Data["json"] = types.Successre{Status: 200, Message: "删除成功", Code: 1}
			} else {
				this.Data["json"] = types.Successre{Status: 400, Message: "服务器内部错误", Code: -1}
			}
		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "信息不存在", Code: -1}
		}

	}
	this.ServeJSON()
}

//权限列表
func (this *JurisdictionController) Julist() {

	//查询数据
	infos, err := models.NewPermission().Sperlist()

	if err {
		//整合数据
		items := map[int]models.Permission{}

		for _, v := range infos {
			items[v.PsID] = v
		}

		pes := make(map[int][]interface{})

		//父级 包含 子权限
		for _, val := range items {

			if val.PsLevel == 0 {
				pes[val.PsID] = append(pes[val.PsID], val)
			} else {
				pes[val.PsPid] = append(pes[val.PsPid], val)
			}
		}
		this.Data["json"] = types.SuccessreInfo{Status: 200, Message: "权限显示", Data: pes, Code: 1}

	} else {
		this.Data["json"] = types.Successre{Status: 400, Message: "不存在信息！", Code: -1}
	}
	this.ServeJSON()
}

//分配权限 列表
func (this *JurisdictionController) Rolelist() {

	info, res := models.Newrole().RoleLsit()
	if res {
		this.Data["json"] = types.SuccessreInfo{Status: 200, Message: "返回成功", Data: info, Code: 1}
	} else {
		this.Data["json"] = types.Successre{Status: 400, Message: "ch", Code: -1}
	}

	this.ServeJSON()
}

//分配权限 -- 新增
func (this *JurisdictionController) RoleInsert() {

	//接受参数
	var info map[string]interface{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &info)
	if err != nil {
		fmt.Println(nil)
	}

	//验证器
	v := validation.Validation{}
	v.Required(info["RoleName"], "role_name").Message("分组名称不可以为空！")
	v.Required(info["RolePsIds"], "role_ps_ids").Message("分组权限不可以为空")

	if v.HasErrors() {
		var message string
		for _, err := range v.Errors {
			message = err.Message + ""
		}
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	} else {

		//权限名称不可以 重复
		roleinfo := models.Newrole().RoleRepeat(info["RoleName"].(string))
		if roleinfo.RoleID == 0 {
			//权限 分配 入库
			res := models.Newrole().Roledistribution(info)
			if res {
				this.Data["json"] = types.Successre{Status: 200, Message: "权限分配成功", Code: 1}
			} else {
				this.Data["json"] = types.Successre{Status: 400, Message: "权限分配失败，请联系管理员", Code: -1}
			}
		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "权限分组已存在", Code: -1}
		}
	}
	this.ServeJSON()
}

//分配权限 -- 修改
func (this *JurisdictionController) RoleUpdate() {

	// 接受修改参数
	var info map[string]interface{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &info)
	if err != nil {
		fmt.Println(err)
	}

	//验证器
	v := validation.Validation{}
	v.Required(info["RoleID"], "role_id").Message("请选择需要修改分组！")
	v.Required(info["RoleName"], "role_id").Message("请谨慎，权限分组名可能会清空！")
	v.Required(info["RolePsIds"], "role_ps_ids").Message("请谨慎，权限可能会清空！")

	//错误信息整合
	if v.HasErrors() {
		//获取提示错误信息
		var buff bytes.Buffer
		for _, err := range v.Errors {
			buff.WriteString(err.Message)
			buff.WriteString(" ")
		}
		message := buff.String()
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	} else {

		//空信息 判断
		resinfo := models.Newrole().RoleidRepeat(info["RoleID"].(float64))
		if resinfo.RoleID != 0 {
			//进行修改操作
			reslute := models.Newrole().RoleUp(info)
			if reslute {
				this.Data["json"] = types.Successre{Status: 200, Message: "分组修改成功", Code: 1}
			} else {
				this.Data["json"] = types.Successre{Status: 400, Message: "修改失败，联系管理员", Code: -1}
			}
		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "修改权限分组不存在", Code: -1}
		}
	}

	this.ServeJSON()
}

//分配权限 -- 删除
func (this *JurisdictionController) RoleDelect() {

	// 接受修改参数
	var info map[string]interface{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &info)
	if err != nil {
		fmt.Println(err)
	}

	//验证器
	v := validation.Validation{}
	v.Required(info["RoleID"], "role_id").Message("请选择需要删除分组！")

	if v.HasErrors() {
		var buff bytes.Buffer
		for _, err := range v.Errors {
			buff.WriteString(err.Message)
			buff.WriteString(" ")
		}
		message := buff.String()
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	} else {

		// 删除权限分组 存在
		res := models.Newrole().RoleidRepeat(info["RoleID"].(float64))
		if res.RoleID != 0 {
			result := models.Newrole().RoleDe(info["RoleID"].(float64))
			if result {

				this.Data["json"] = types.Successre{Status: 200, Message: "权限分组删除成功", Code: 1}
			} else {
				this.Data["json"] = types.Successre{Status: 400, Message: "权限删除失败，请联系管理员", Code: -1}
			}
		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "权限分组不存在", Code: -1}
		}
	}
	this.ServeJSON()
}
