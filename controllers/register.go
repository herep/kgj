package controllers

import (
	"crypto/md5"
	"encoding/json"
	"finance/comm"
	"finance/models"
	"finance/types"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type RegisterController struct {
	beego.Controller
}

//发短信接口
func (this *RegisterController) SendCode() {

	//接受数据
	res := make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &res)
	fmt.Println(err)
	phone := res["phone"].(string)

	//验证器
	v := validation.Validation{}
	v.Required(phone, "phone").Message("手机号必须输入！")
	v.Mobile(phone, "phone").Message("手机号格式不正确！")

	if v.HasErrors() {
		//验证信息
		var message string
		for _, err := range v.Errors {
			message = message + "," + err.Message
		}
		this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}

	} else {
		//频繁发送信息
		data := models.NewBindcode().CheckBinde(phone)

		//此时时间
		newtime := time.Now().Unix()
		if len(data) != 0 && data[0].Expirytime > newtime {
			this.Data["json"] = types.Successre{Status: 400, Message: "发送太频繁", Code: -1}
		} else {
			Sixcode := comm.GetSix() //6 位验证码
			result := comm.Request(phone, Sixcode)

			var netReturn map[string]interface{}
			err := json.Unmarshal(result, &netReturn)
			if err != nil {
				fmt.Println(err)
			}

			if netReturn["error_code"] == 0 {
				//入库信息
				infos := map[string]interface{}{
					"code":       0,
					"phone":      phone,
					"VerifyCode": Sixcode,
					"CreateTime": newtime,
					"ExpiryTime": newtime + 120,
				}
				//入库操作
				mes := models.NewBindcode().InsertBcode(infos)
				//入库--成功 或 否
				if mes["status"] == 200 {
					this.Data["json"] = types.Successre{Status: 200, Message: "发送成功", Code: 0}
				} else {
					this.Data["json"] = types.Return{Status: mes["status"].(int), Message: mes["message"].(string), Code: -1}
				}
			} else {

				//此时时间
				ntime := time.Now().Unix()
				//入库信息集合
				infos := map[string]interface{}{
					"code":       1,
					"phone":      phone,
					"VerifyCode": Sixcode,
					"CreateTime": ntime,
					"ExpiryTime": ntime + 120,
					"ResCode":    netReturn["error_code"],
					"ResMessage": netReturn["reason"],
				}

				//入库
				mes := models.NewBindcode().InsertBcode(infos)

				//入库 是否成功
				if mes["status"] == 200 {
					this.Data["json"] = types.Return{Status: 200, Message: netReturn["reason"].(string), Code: 0, Err: netReturn}
				} else {
					this.Data["json"] = types.Return{Status: mes["status"].(int), Message: mes["message"].(string), Code: -1, Err: netReturn}
				}
			}
		}
	}
	this.ServeJSON()
}

//注册接口
func (this *RegisterController) Register() {

	//接受数
	if len(this.Ctx.Input.RequestBody) == 0 {
		this.Data["json"] = types.Successre{Status: 400, Message: "信息不可以为空", Code: -1}
	} else {

		//整合参数
		res := make(map[string]interface{})
		err := json.Unmarshal(this.Ctx.Input.RequestBody, &res)
		fmt.Println(err)

		company_name := res["company_name"].(string)
		admin_num := res["admin_num"].(string)
		password := res["password"].(string)
		admin_name := res["admin_name"].(string)
		phone_num := res["phone_num"].(string)
		mailbox := res["mailbox"].(string)
		consent := res["consent"].(string)
		code := res["code"].(string)

		//验证器
		v := validation.Validation{}
		v.Required(company_name, "company_name").Message("公司名不可以为空！")
		v.Required(admin_name, "admin_name").Message("管理员姓名不可以为空！")
		v.Required(admin_num, "admin_num").Message("管理员帐号不可以为空！")
		v.Required(phone_num, "phone_num").Message("手机号码不可以为空！")
		v.Required(mailbox, "mailbox").Message("邮箱不可以为空！")
		v.Required(password, "password").Message("密码不可以为空！")
		v.Required(consent, "consent").Message("是否同意协议")
		v.Email(mailbox, "mailbox").Message("邮箱格式不对！")
		v.Phone(phone_num, "phone_num").Message("手机帐号不正确！")
		v.Required(code, "code").Message("请输入收入的验证码！")

		//开始验证
		if v.HasErrors() {
			//错误信息 提示
			var message string
			for _, err := range v.Errors {
				message = message + "," + err.Message
			}
			this.Data["json"] = types.Successre{Status: 400, Message: message, Code: -1}
		} else {
			//手机号 -- 不可以重复注册
			companys := models.NewUser().Checkuser(admin_num)
			if len(companys) == 0 {
				//邮箱 -- 不可以重复注册
				mailboxinfo := models.NewUser().Checkmailboxuser(mailbox)
				if len(mailboxinfo) == 0 {
					//用户名 -- 不可以重复注册
					nameinfo := models.NewUser().Checknameuser(mailbox)
					if len(nameinfo) == 0 {
						//验证通过
						bindcode := models.NewBindcode().CheckBinde(phone_num)

						if len(bindcode) == 0 {
							this.Data["json"] = types.Successre{Status: 400, Message: "手机号错误", Code: -1}
						} else {
							if bindcode[0].Verifycode == code {

								ntime := time.Now().Unix() // 当前时间
								data := []byte(password)
								pas := md5.Sum(data)
								md5str := fmt.Sprintf("%x", pas)

								//验证码 过期
								if bindcode[0].Expirytime > ntime {
									//整合数据
									info := map[string]interface{}{
										"company_name": company_name,
										"admin_num":    admin_num,
										"admin_name":   admin_name,
										"password":     md5str,
										"phone_num":    phone_num,
										"mailbox":      mailbox,
										"create_time":  ntime,
										"consent":      consent,
									}

									res, id := models.NewUser().Insertv(info)

									if res {
										//维护 关系 表
										uainfo := models.Uafiliation{UserId: id, AccountId: 0, UserName: info["admin_name"].(string), UserMailbox: info["mailbox"].(string), UserPhone: info["admin_num"].(string), Status: 1}
										ua := models.Newuafiliation().Iuainfo(uainfo)
										if ua {
											this.Data["json"] = types.Successre{Status: 200, Message: "注册成功", Code: 0}
										} else {
											this.Data["json"] = types.Successre{Status: 400, Message: "关系表维护失败，请联系管理员", Code: -1}
										}
									} else {
										this.Data["json"] = types.Successre{Status: 400, Message: "服务器内部错误", Code: -1}
									}
								} else {
									this.Data["json"] = types.Successre{Status: 400, Message: "验证码过期", Code: -1}
								}

							} else {
								this.Data["json"] = types.Successre{Status: 400, Message: "验证码错误", Code: -1}
							}
						}
					} else {
						this.Data["json"] = types.Successre{Status: 400, Message: "用户名已经绑定公司", Code: -1}
					}
				} else {
					this.Data["json"] = types.Successre{Status: 400, Message: "邮箱已经绑定公司", Code: -1}
				}
			} else {
				this.Data["json"] = types.Successre{Status: 400, Message: "手机号已经绑定公司", Code: -1}
			}
		}
	}
	this.ServeJSON()
}

//登录
func (this *RegisterController) Login() {

	//接受数据
	var ob models.KgUser
	var err error
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &ob)

	if err != nil {
		fmt.Println(err)
	}

	//验证器
	v := validation.Validation{}
	v.Required(ob.AdminNum, "username").Message("用户名不可以为空！")
	v.Required(ob.Password, "password").Message("密码不可以为空！")

	if v.HasErrors() {
		var message string
		for _, err := range v.Errors {
			message = message + " " + err.Message
		}
		this.Data["json"] = types.Return{Status: 400, Message: message, Code: -1}

	} else {

		Ophone := v.Phone(ob.AdminNum, "username").Message("unphone")
		Oemail := v.Email(ob.AdminNum, "username").Message("unemail")

		var field string
		if Ophone.Ok {
			field = "user_phone"
		} else if Oemail.Ok {
			field = "user_mailbox"
		} else {
			field = "user_name"
		}

		//主 子帐号 关系表 -- 判断主子帐号
		uainfo, err := models.Newuafiliation().Suainfo(ob.AdminNum, field)
		pas := md5.Sum([]byte(ob.Password))
		md5str := fmt.Sprintf("%x", pas)

		if err {
			if uainfo[0].Status != 0 {

				//子帐号登录 status --- 1
				if field == "user_phone" {
					field = "account_phone"
				} else if Oemail.Ok {
					field = "account_mailbox"
				} else {
					field = "account_name"
				}

				//登录---查询
				Acinfo, err := models.Newaccount().Checkacinfo(ob.AdminNum, field)

				if !err {
					this.Data["json"] = types.Successre{Status: 400, Message: "用户信息不存在", Code: -1}
				} else {
					if Acinfo.AccountPas != md5str {
						this.Data["json"] = types.Successre{Status: 400, Message: "密码错误", Code: -1}
					} else {

						//创建令牌 -- token
						claims := make(jwt.MapClaims)
						claims["userid"] = Acinfo.Id
						claims["belong"] = "Subaccount"

						claims["exp"] = time.Now().Add(time.Hour * 48).Unix()
						token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

						//使用自定义string 加严
						tokenString, err := token.SignedString([]byte("cpf"))
						if err != nil {
							this.Data["json"] = types.Successre{Status: 400, Message: "token获取失败", Code: -1}
						}

						this.Data["json"] = types.SuccessLogin{Status: 200, Message: "登录成功", Code: 1, AdminName: Acinfo.AccountName, Token: tokenString}
					}
				}

			} else {
				// status -- 0 主帐号

				//用户选择登录方式
				if field == "user_phone" {
					field = "admin_num"
				} else if Oemail.Ok {
					field = "mailbox"
				} else {
					field = "admin_name"
				}

				Uinfo := models.NewUser().Checkfuser(ob.AdminNum, field)

				if len(Uinfo) == 0 {
					this.Data["json"] = types.Successre{Status: 400, Message: "用户信息不存在", Code: -1}
				} else {
					if Uinfo[0].Password != md5str {
						this.Data["json"] = types.Successre{Status: 400, Message: "密码错误", Code: -1}
					} else {

						//创建令牌 -- token
						claims := make(jwt.MapClaims)
						claims["userid"] = Uinfo[0].Id
						claims["belong"] = "BOSS"

						claims["exp"] = time.Now().Add(time.Hour * 48).Unix()
						token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

						//使用自定义string 加严
						tokenString, err := token.SignedString([]byte("cpf"))
						if err != nil {
							this.Data["json"] = types.Successre{Status: 400, Message: "token获取失败", Code: -1}
						}

						this.Data["json"] = types.SuccessLogin{Status: 200, Message: "登录成功", Code: 1, AdminName: Uinfo[0].AdminName, Token: tokenString}
					}
				}
			}

		} else {
			this.Data["json"] = types.Successre{Status: 400, Message: "主子帐号关系需要维护，联系管理员", Code: -1}
		}
	}

	this.ServeJSON()
}
