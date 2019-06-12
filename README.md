# gokgj

### 简介：基于 go beego搭建erp
#### 注意：配置文件模板

需要在项目根目录创建  /conf/app.conf 文件

```conf
appname = finance
httpport = 
runmode = dev
autorender = false
copyrequestbody = true
EnableDocs = true

  mysql 配置
kg_mysql_name =
kg_mysql_password = 
kg_mysql_host = 
kg_mysql_dbname = 

   redis 配置
kg_redis_host = 
kg_redis_password = 

   短信接口 信息
smsconf_key = f2f996a99ffbb73e52e235e0d9beabe1
smsconf_tpl_id = 135197
smsconf_url = "http://v.juhe.cn/sms/send"
```
##### 2019-06-10 
````
主子帐号管理，权限管理，权限分配，
````

##### 2019-06-12
````
权限分组分配维护
````
