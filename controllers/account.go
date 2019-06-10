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

	////验证器
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
		info["AccountCompany"] = Cinfo[0].Id

		//相同电话不可以重复注册
		res,_ := models.Newaccount().Saccount(info["AccountNum"].(string))

		if res{
			//不允许 新增
			this.Data["json"] = types.Successre{Status: 400, Message: "电话以绑定帐号", Code: -1}
		}else{

			ac, acid := models.Newaccount().Iaccount(info) //新增子帐号主键

			if ac {
				//维护 主-子关系
				uainfo := models.Uafiliation{UserId: Cinfo[0].Id, AccountId: acid,UserName:info["AccountName"].(string),UserMailbox:info["AccountMailbox"].(string),UserPhone: info["AccountNum"].(string),Status:1}
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

	}
	this.ServeJSON()
}

//子帐号 -- 修改