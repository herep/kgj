package routers

import (
	"encoding/json"
	_ "finance/comm"
	"finance/controllers"
	_"finance/models"
	"finance/types"
	_ "fmt"
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

				data, _ := json.Marshal(types.Return{Status: 400, Message: "That's not even a token", Code: -1})
				ctx.ResponseWriter.Write(data)

			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				data, _ := json.Marshal(types.Return{Status: 400, Message: "Token is either expired or not active yet", Code: -1})
				ctx.ResponseWriter.Write(data)

			} else {
				data, _ := json.Marshal(types.Return{Status: 400, Message: "Couldn't handle this token", Code: -1})
				ctx.ResponseWriter.Write(data)
			}
		} else {
			data, _ := json.Marshal(types.Return{Status: 400, Message: "Couldn't handle this token", Code: -1})
			ctx.ResponseWriter.Write(data)
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
	beego.Router("/v1/jurisdiction/ijuinfo", &controllers.JurisdictionController{}, "post:Ijuinfo")
	beego.Router("/v1/jurisdiction/ujuinfo", &controllers.JurisdictionController{}, "post:Ujurinfo")
	beego.Router("/v1/jurisdiction/djuinfo", &controllers.JurisdictionController{}, "post:Djurinfo")
	beego.Router("/v1/jurisdiction/julist", &controllers.JurisdictionController{}, "post:Julist")
	beego.Router("/v1/jurisdiction/rolelist", &controllers.JurisdictionController{}, "post:Rolelist")

	beego.InsertFilter("/v1/*", beego.BeforeRouter, FilterUser) 
}
