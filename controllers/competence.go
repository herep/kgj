package controllers

import "github.com/astaxie/beego"

type CompetenceController struct {
	beego.Controller
}

//权限管理
func (this *CompetenceController) CompetenceInfo(){

	//
	this.Ctx.Input.URL()
}
