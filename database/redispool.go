package database

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	Pool *redis.Pool //创建redis连接池
)

func init() {
	//read conf
	redis_host := beego.AppConfig.String("kg_redis_host")
	redis_password := beego.AppConfig.String("kg_redis_password")

	Pool = &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redis_host)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			if _, err := c.Do("AUTH", redis_password); err != nil {
				c.Close()
				fmt.Println(err)
				return nil, err
			}
			if _, err := c.Do("SELECT", "10"); err != nil {
				c.Close()
				fmt.Println(err)
				return nil, err
			}
			return c, nil
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
