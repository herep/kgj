package controllers

import (
	"bytes"
	"encoding/json"
	"finance/comm"
	"finance/models"
	"finance/types"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type OrderController struct {
	beego.Controller
}

//订单显示接口
func (this *OrderController) GetOrderinfo() {

	//接受条件
	var GetOrderParments types.GetOrderParamenter
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &GetOrderParments); err != nil {
		fmt.Println(err.Error())
	}

	//验证器
	v := validation.Validation{}
	v.Required(GetOrderParments.Page, "page").Message("请输入查看页数")
	v.Required(GetOrderParments.PageSize, "page_size").Message("请输入查看页数")

	if v.HasErrors() {
		var buff bytes.Buffer
		for _, err := range v.Errors {
			buff.WriteString(err.Message)
			buff.WriteString(" ")
		}
		this.Data["json"] = types.Successre{Status: 400, Message: buff.String(), Code: -1}
	} else {
		//默认赛选条
		userinfo := comm.GetTokeninfo(this.Ctx)
		result_create, oreder_info := models.NewShipments().GetShipmentsInfo(GetOrderParments, userinfo)
		if result_create {
			this.Data["json"] = types.SuccessreInfo{Status: 200, Message: "success", Data: oreder_info, Code: 1}
		}
	}

	this.ServeJSON()
}
