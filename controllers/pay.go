package controllers

import (
	"finance/models"
	"finance/types"
	_ "fmt"
	"github.com/astaxie/beego"
)

type PayController struct {
	beego.Controller
}

//授权 -- 回调
func (this *PayController) Callback() {

	//授权成功 整合信息 
	parameter := types.Pcallback{
		Seesion:        this.Input().Get("session"),
		RefreshToken:   this.Input().Get("refresh_token"),
		TaobaoUserNick: this.Input().Get("taobao_user_nick"),
		SiteName: "TB",
	}

	//修改
	message,res := models.NewShoplist().Sauthorize(parameter)

	if res {
		//接受数据
		this.Data["json"] = types.Successre{Status: 200, Message: message, Code: 1}
	} else {
		//接受数据
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
	}
	this.ServeJSON()
}

//微信支付
func (this *PayController) Payweixin() {

	this.Data["json"] = types.SuccessLogin{Status: 200, Message: "0", Code: 0}
	this.ServeJSON()
}
