package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

//创建Redis 连接池
func NewRdisPool(addr, pwd string) (pool *redis.Pool) {
	return &redis.Pool{
		MaxIdle:     60,
		MaxActive:   200,
		IdleTimeout: time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr)
			if err != nil {
				fmt.Println("连接redis数据库失败", err)
				_ = conn.Close()
				return nil, err
			}
			//如果需要密码！
			if len(pwd) > 1 {
				if _, err := conn.Do("AUTH", pwd); err != nil {
					_ = conn.Close()
					fmt.Println("密码错误", err)
					return conn, err
				}
			}
			return conn, err
		},
		//连接redis 连接池测试
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("ping")
			return err
		},
	}
}

type Person struct {
	Name string `json:"name" redis:"name"`
	Age  int    `json:"age" redis:"age"`
}

//1.可选参数:调用f1函数是option可以不传，传参之后函数接收到1个slice。
func f1(option ...string) {
	fmt.Println(option)
}

//2.必须传参:参数可以是任意类型，如果没有参数可以传nil
func f2(option interface{}) {
	fmt.Println(option)
}
//3.兼容以上2种情况option可选参数，参数可以是任意类型
func f3(option ...interface{})  {
	fmt.Println(option)
}


func main() {
	/*
		//初始化1个连接池
		pool := NewRdisPool("192.168.56.18:6379", "123.com")
		//从连接池中获取连接
		conn := pool.Get()
		//最后关闭连接、关闭连接池
		defer conn.Close()
		pool.Close()
		//set结构体 值
		marshalData, err := json.Marshal(Person{Name: "张根", Age: 27})
		if err != nil {
			fmt.Println(err)
		}
		_, err = conn.Do("set", "martin", marshalData)
		if err != nil {
			fmt.Println("set值失败")
		}
		//获取值:获取不到返回n
		data, err1 := conn.Do("get", "martin")
		if err1 != nil {
			fmt.Println("获取值失败")
		}
		stringData, err2 := redis.String(data, err)
		if err != nil {
			fmt.Println(err2)
		}
		var u Person
		json.Unmarshal([]byte(stringData), &u)
		fmt.Printf("%#v\n",u)
	*/
	f1()
	f2(nil)
	f3([]string{"1","2","3"})
	f3(map[string]interface{}{})
	f3()

}
