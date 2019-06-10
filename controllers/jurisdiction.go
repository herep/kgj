package controllers

import (
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
		data.CreateTime = time.Now().Unix();
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
			data.UpdateTime = time.Now().Unix();
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
