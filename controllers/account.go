package controllers

import (
	"encoding/json"
	"finance/comm"
	"finance/models"
	"finance/types"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type AccountController struct {
	beego.Controller
}

//子帐号 -- 新增
func (this *AccountController) Iaccount() {

	//新增用户信息
	var info map[string]interface{}
	var err error
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &info)

	if err != nil {
		fmt.Println(err)
	}

	//验证器
	v := validation.Validation{}
	v.Required(info["AccountPas"], "account_pas").Message("密码不可以为空")
	v.Required(info["AccountName"], "account_name").Message("子帐号名称不可以为空！")
	v.Required(info["AccountNum"], "account_num").Message("登录帐号不可以为空！")
	v.Phone(info["AccountNum"], "account_num").Message("登录帐号格式不正确")
	v.Email(info["AccountMailbox"], "account_mailbox").Message("邮箱格式不正确")
	v.Required(info["AccountMailbox"], "account_mailbox").Message("邮箱不可以为空")
	v.Required(info["AccountPhone"], "account_phone").Message("联系电话不可以为空")
	v.Phone(info["AccountPhone"], "account_phone").Message("联系电话格式不正确")
	v.Required(info["AccountRole"], "account_role").Message("请选择用户角色")

	//验证信息 拼接
	if v.HasErrors() {
		var message string
		for _, err := range v.Errors {
			message = message + err.Message + " "
		}
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	} else {

		//此时登录 主帐号
		Cinfo := comm.GetTokeninfo(this.Ctx)
		if Cinfo.AccountId != 0 {
			info["AccountCompany"] = Cinfo.AccountId
			//相同电话不可以重复注册
			res, _ := models.Newaccount().Saccount(info["AccountNum"].(string))

			if res {
				//不允许 新增
				this.Data["json"] = types.Successre{Status: 400, Message: "电话已绑定帐号", Code: -1}
			} else {

				ac, acid := models.Newaccount().Iaccount(info) //新增子帐号主键

				if ac {
					//维护 主-子关系
					uainfo := models.Uafiliation{UserId: Cinfo.AccountId, AccountId: acid, UserName: info["AccountName"].(string), UserMailbox: info["AccountMailbox"].(string), UserPhone: info["AccountNum"].(string), Status: 1}
					ua := models.Newuafiliation().Iuainfo(uainfo)
					if ua {
						this.Data["json"] = types.Successre{Status: 200, Message: "子帐号新增成功", Code: 1}
					} else {
						this.Data["json"] = types.Successre{Status: 400, Message: "主子帐号维护失败，请联系管理员", Code: -1}
					}

				} else {
					this.Data["json"] = types.Successre{Status: 400, Message: "新增失败，请联系管理员", Code: -1}
				}
			}
		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "非主帐号不可以创建子帐号", Code: -1}
		}

	}
	this.ServeJSON()
}

//子帐号 -- 修改
func (this *AccountController) Uaccount() {

	//获取 参数
	var info models.Account
	var err error
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &info)

	if err != nil {
		fmt.Println(err)
	}

	//验证器
	v := validation.Validation{}
	v.Required(info.Id, "user_id").Message("请选择修改子帐号")

	//验证器 错误信息
	if v.HasErrors() {

		var message string
		for _, err := range v.Errors {
			message = err.Message + ""
		}
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	} else {

		//是否存在 信息
		_, result := models.Newaccount().IdGetInfo(info.Id)
		if result {
			//修改数据
			errs, data := models.Newaccount().Uinfo(info)
			if errs {
				//维护关系表 -- 整合数据
				var ua models.Uafiliation
				ua.UserMailbox = data.AccountMailbox
				ua.UserName = data.AccountName
				ua.UserPhone = data.AccountNum
				ua.AccountId = data.Id
				result := models.Newuafiliation().Uuainfo(ua)

				if result {
					this.Data["json"] = types.Successre{Status: 200, Message: "修改成功", Code: 1}
				} else {
					this.Data["json"] = types.Successre{Status: 400, Message: "主子表维护失败，服务器内部错误", Code: -1}
				}
			} else {
				this.Data["json"] = types.Successre{Status: 400, Message: "修改出错，服务器内部错误", Code: -1}
			}

		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "修改子账户不存在", Code: -1}
		}
	}

	this.ServeJSON()
}

//子帐号 删除
func (this *AccountController) Daccount() {

	//接受 参数
	var info models.Account
	var err error
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &info)
	if err != nil {
		fmt.Println(err)
	}

	//验证器
	v := validation.Validation{}
	v.Required(info.Id, "id").Message("请选择删除信息")

	if v.HasErrors() {

		var message string
		for _, err := range v.Errors {
			message = err.Message + ""
		}
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	} else {

		//是否存在 信息
		_, result := models.Newaccount().IdGetInfo(info.Id)
		if result {
			//修改数据
			errs, infos := models.Newaccount().Dinfo(info)
			if errs {
				result := models.Newuafiliation().Duainfo(infos.Id)
				if result {
					this.Data["json"] = types.Successre{Status: 200, Message: "删除出错成功", Code: 1}
				} else {
					this.Data["json"] = types.Successre{Status: 400, Message: "主子表维护失败，服务器内部错误", Code: -1}
				}

			} else {
				this.Data["json"] = types.Successre{Status: 400, Message: "删除出错，服务器内部错误", Code: -1}
			}

		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "删除出错子账户不存在", Code: -1}
		}
	}
	this.ServeJSON()
}
