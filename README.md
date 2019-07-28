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
smsconf_key = 
smsconf_tpl_id = 
smsconf_url = 
```
##### 2019-06-10 
````
主子帐号管理，权限管理，权限分配，
````

##### 2019-06-12
````
权限分组分配维护
````

##### 2019-06-13
````
权限分组 限制只可以主帐号，token中获取数据 区分 子主帐号登录信息
````

##### 2019-06-16
````
主子帐号 - 相同电话 用户名 邮箱 禁止新建
````
##### 2019-06-17
````
通过路由 用户对应权限  是否合法访问1
````