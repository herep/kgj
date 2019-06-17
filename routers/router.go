package routers

import (
	"bytes"
	"encoding/json"
	"finance/comm"
	_ "finance/comm"
	"finance/controllers"
	"finance/models"
	_ "finance/models"
	"finance/types"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

var FilterUser = func(ctx *context.Context) {
	//请求头信息获取
	authString := ctx.Input.Header("Token")

	kv := strings.Split(authString, " ")
	tokenString := kv[0]

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//返回 后端加密的 严
		return []byte("cpf"), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {

			if ve.Errors&jwt.ValidationErrorMalformed != 0 {

				data, _ := json.Marshal(types.Return{Status: 400, Message: "That's not even a token", Code: -2})
				ctx.ResponseWriter.Write(data)

			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				data, _ := json.Marshal(types.Return{Status: 400, Message: "Token is either expired or not active yet", Code: -2})
				ctx.ResponseWriter.Write(data)

			} else {
				data, _ := json.Marshal(types.Return{Status: 400, Message: "Couldn't handle this token", Code: -2})
				ctx.ResponseWriter.Write(data)
			}
		} else {
			data, _ := json.Marshal(types.Return{Status: 400, Message: "Couldn't handle this token", Code: -2})
			ctx.ResponseWriter.Write(data)
		}
	} else {

		url := ctx.Input.URL() //登录--权限 获取请求路由
		info := comm.GetTokeninfo(ctx)

		//主张号 权限过滤
		if info.AccountId != 0 {
			da := models.Newuafiliation().SelectRoleinfo(info)

			//请求路由 /v1/模块名/方法名  转化 模块名-方法名
			var buff bytes.Buffer
			urlinfo := strings.Split(url, "/")
			buff.WriteString(urlinfo[2])
			buff.WriteString("-")
			buff.WriteString(urlinfo[3])
			urlrole := buff.String()

			// 所属角色 所有权限 -- 是否为合法访问
			var ints int
			for _, v := range da {
				if strings.Index(v.RolePsCas, ",") != -1 {
					s := strings.Split(v.RolePsCas, ",")
					for _, v2 := range s {
						if v2 == urlrole {

							ints = ints + 1
						}
					}
				} else {
					if v.RolePsCas == urlrole {
						ints = ints + 1
					}
				}
			}

			if ints == 0 {
				data, _ := json.Marshal(types.Successre{Status: 400, Message: "无权限，禁止访问！", Code: -2})
				ctx.ResponseWriter.Write(data)
			}
		}
	}
}

func init() {
	beego.Router("/register/sendcode", &controllers.RegisterController{}, "post:SendCode")
	beego.Router("/register/logion", &controllers.RegisterController{}, "post:Register")
	beego.Router("/register/login", &controllers.RegisterController{}, "post:Login")

	beego.Router("/pay/callback", &controllers.PayController{}, "get:Callback")

	beego.Router("/v1/pay/payweixin", &controllers.PayController{}, "get:Payweixin")

	beego.Router("/v1/strore/ishopname", &controllers.StroreController{}, "post:Ishopname")
	beego.Router("/v1/strore/shoplist", &controllers.StroreController{}, "post:Shoplist")
	beego.Router("/v1/strore/changeshop", &controllers.StroreController{}, "post:Changeshop")

	beego.Router("/v1/account/iaccount", &controllers.AccountController{}, "post:Iaccount")
	beego.Router("/v1/account/uaccount", &controllers.AccountController{}, "post:Uaccount")
	beego.Router("/v1/account/daccount", &controllers.AccountController{}, "post:Daccount")
	beego.Router("/v1/account/accountlist", &controllers.AccountController{}, "post:Accountlist")

	beego.Router("/v1/jurisdiction/ijuinfo", &controllers.JurisdictionController{}, "post:Ijuinfo")
	beego.Router("/v1/jurisdiction/ujuinfo", &controllers.JurisdictionController{}, "post:Ujurinfo")
	beego.Router("/v1/jurisdiction/djuinfo", &controllers.JurisdictionController{}, "post:Djurinfo")
	beego.Router("/v1/jurisdiction/julist", &controllers.JurisdictionController{}, "post:Julist")
	beego.Router("/v1/jurisdiction/rolelist", &controllers.JurisdictionController{}, "post:Rolelist")
	beego.Router("/v1/jurisdiction/roleinsert", &controllers.JurisdictionController{}, "post:RoleInsert")
	beego.Router("/v1/jurisdiction/roleupdate", &controllers.JurisdictionController{}, "post:RoleUpdate")
	beego.Router("/v1/jurisdiction/roledelect", &controllers.JurisdictionController{}, "post:RoleDelect")

	beego.InsertFilter("/v1/*", beego.BeforeRouter, FilterUser)
}
