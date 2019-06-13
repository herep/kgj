package controllers

import (
	"encoding/json"
	"finance/comm"
	"finance/models"
	"finance/types"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"time"
)

type StroreController struct {
	beego.Controller
}

/*
--店铺入库
*/
func (this *StroreController) Ishopname() {
	//接受值
	res := make(map[string]interface{})
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &res); err != nil {
		fmt.Println(err)
	}

	shopname := res["shopname"].(string)

	v := validation.Validation{}
	v.Required(shopname, "shopname").Message("授权前输入店铺名")

	if v.HasErrors() {
		//验证信息
		var message string
		for _, err := range v.Errors {
			message = message + "," + err.Message
		}
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	} else {

		//token值
		userinfo := comm.GetTokeninfo(this.Ctx)
		info := map[string]interface{}{"shopanme": shopname, "create_time": time.Now().Unix(), "company_id": userinfo.UserId}

		//判断是否 重复 入库
		repeat := models.NewShoplist().Chackshop(shopname, userinfo.UserId)
		if len(repeat) == 0 {
			//入库操作
			res := models.NewShoplist().Sinsert(info)
			if res {
				this.Data["json"] = types.Successre{Status: 200, Message: "存入成功", Code: 1}
			} else {
				this.Data["json"] = types.Successre{Status: 400, Message: "服务器内部错误", Code: -1}
			}

		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "店铺已存在，请直接授权", Code: 3}
		}

	}
	this.ServeJSON()
}

/*
--店铺列表
*/
func (this *StroreController) Shoplist() {

	var where *types.WshopList
	//当前用户 -- 个人信息
	userinfo := comm.GetTokeninfo(this.Ctx)
	company_id := userinfo.UserId

	//按照条件 查询
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &where); err != nil {
		fmt.Println(err)
	}
	where.CompanyId = company_id // 公司id

	//验证器
	v := validation.Validation{}
	v.Required(where.Page, "page").Message("请传入页码")
	v.Required(where.PageSize, "pagesize").Message("页数不可以为空")

	//错误信息
	if v.HasErrors() {

		//错误提醒
		var message string
		for _, err := range v.Errors {
			message = message + "," + err.Message
		}
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}

	} else {

		//数据库数据 返回
		info, count := models.NewShoplist().Sshoplist(where)
		companyinfos := models.NewUser().IdGetIntInfo(company_id)

		if count == 0 {
			this.Data["json"] = types.Successre{Status: 400, Message: "暂无数据", Code: -1}
		} else {
			this.Data["json"] = map[string]interface{}{"Status": 200, "Message": "返回成功", "Data": info, "Count": count, "Code": 1, "CompanyName": companyinfos[0].CompanyName}
		}
	}

	this.ServeJSON()
}

/*
--店铺信息完善
*/
func (this *StroreController) Changeshop() {

	var infos map[string]interface{}
	//接受数据
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &infos); err != nil {
		fmt.Println(err)
	}

	//必须传入 店铺id
	shop_id := infos["id"]

	//验证器
	v := validation.Validation{}
	v.Required(shop_id, "shop_id").Message("未知修改店铺")

	if v.HasErrors() {

		var message string
		for _, err := range v.Errors {
			message = message + "," + err.Message
		}
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}

	} else {

		//判断是否存在 店铺信息
		result := models.NewShoplist().Iidinfo(shop_id.(float64))

		if result != nil {
			//修改--方法
			res := models.NewShoplist().Ushoplist(infos)

			if res {
				this.Data["json"] = types.Successre{Status: 200, Message: "完善成功", Code: 1}
			} else {
				this.Data["json"] = types.Successre{Status: 400, Message: "完善失败", Code: -1}
			}

		} else {

			this.Data["json"] = types.Successre{Status: 400, Message: "店铺信息不存在", Code: -1}
		}
	}

	this.ServeJSON()
}
