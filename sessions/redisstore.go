package sessions

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

//我们把session保存在内存中，无法持久化。所以现在我们可把redis当做1个保存session的大map，用来持久保存session！

type RedisStore struct {
	//全局锁
	tablelock sync.RWMutex
	//Redis连接池
	connetpool *redis.Pool
}

//初始化redis数据库
func NewRedisStore(ip, pwd string) *RedisStore {
	return &RedisStore{
		connetpool: InitRedis(ip, pwd),
	}

}

func (rs *RedisStore) GetSession(sessionID string) (cell SessionInterface) {
	rs.tablelock.RLock()
	rs.tablelock.RUnlock()
	conn := rs.connetpool.Get()
	bdata, err := conn.Do("get", sessionID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if bdata == nil {
		fmt.Println("没有从redis获取到用户的session信息！")
		return
	}
	sdata, err := redis.String(bdata, err)
	if err != nil {
		fmt.Println("session字节转成map失败", err)
		return
	}
	//每次从redis中获取到session记录，都新创建1个在内存的变量存储
	redisCell := NewRedisCell()
	err = json.Unmarshal([]byte(sdata), &redisCell.SessionMap)
	if err != nil {
		fmt.Println("session data 序列化失败!")
	}
	return
}

func (rs *RedisStore) SetSession() (sessionID string) {
	rs.tablelock.Lock()
	defer rs.tablelock.Unlock()
	conn := rs.connetpool.Get()
	uid := uuid.NewV4().String()
	cell := NewRedisCell()
	marshalData, err := json.Marshal(cell.SessionMap)
	if err != nil {
		fmt.Println("ression data 序列化失败", err)
	}
	_, err = conn.Do("set", uid, marshalData)
	if err != nil {
		fmt.Println("set session信息到redis失败", err)
		return
	}
	return sessionID
}

func (rs *RedisStore) DelSession(sessionID string) {
	rs.tablelock.Lock()
	defer rs.tablelock.Unlock()
	conn := rs.connetpool.Get()
	_, err := conn.Do("del", sessionID)
	if err != nil {
		fmt.Println(err)
	}

}

//初始化数据库
func (rs *RedisStore) InitStore(option ...string) (store StoreInterface) {
	return
}

//创建Redis 连接池
func InitRedis(option ...string) (pool *redis.Pool) {
	return &redis.Pool{
		MaxIdle:     60,
		MaxActive:   200,
		IdleTimeout: time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", option[0])
			if err != nil {
				fmt.Println(err)
			}
			//如果需要密码！
			if len(option[1]) > 1 {
				if _, err := conn.Do("auth", option[1]); err != nil {
					_ = conn.Close()
					fmt.Println("redis密码错误！", err)
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
