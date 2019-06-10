package comm

import (
	"finance/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//随机 生成六位数
func GetSix() (code string) {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

// post 网络请求 ,params 是url.Values类型
func Post(apiURL string, params url.Values) (rs []byte, err error) {
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//发短信接口
func Request(phone string, code string) (info []byte) {

	//read -- conf
	urls := beego.AppConfig.String("smsconf_url")
	tpl_id := beego.AppConfig.String("smsconf_tpl_id")
	key := beego.AppConfig.String("smsconf_key")
	//请求地址
	juheURL := urls

	//初始化参数
	param := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("mobile", phone)             //接收短信的手机号码
	param.Set("tpl_id", tpl_id)            //短信模板ID，请参考个人中心短信模板设置
	param.Set("tpl_value", "#code#="+code) //变量名和变量值对。如果你的变量名或者变量值中带有#&amp;=中的任意一个特殊符号，请先分别进行urlencode编码后再传递，&lt;a href=&quot;http://www.juhe.cn/news/index/id/50&quot; target=&quot;_blank&quot;&gt;详细说明&gt;&lt;/a&gt;
	param.Set("key", key)                  //应用APPKEY(应用详细页查询)
	param.Set("dtype", "json")             //返回数据的格式,xml或json，默认json

	//发送请求
	data, err := Post(juheURL, param)

	if err != nil {
		fmt.Errorf("请求失败,错误信息:\r\n%v", err)
	} else {
		return data
	}
	return data
}

//token---获取个人信息
func GetTokeninfo(ctx *context.Context) (message []models.KgUser) {

	authString := ctx.Input.Header("Token")
	kv := strings.Split(authString, " ")
	tokenString := kv[0]
	//
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//返回 后端加密的 严
		return []byte("cpf"), nil
	})

	//获取 请求头中 信息 -- 查询个人信息
	tokeninfo, _ := token.Claims.(jwt.MapClaims)
	userid := tokeninfo["userid"].(float64)
	Newinfo := models.NewUser().IdGetInfo(userid)

	return Newinfo
}
