package sessions

import (
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
	//存储整个session
	store map[string]SessionInterface
}

//初始化数据库
func (rs *RedisStore) InitStore(option Option) (store StoreInterface) {
	rs.connetpool=InitRedis(option.ip, option.pwd)
	return rs
}

//在大仓库中，创建1条用户session到内存
func (rs *RedisStore) CreateSession() (cell SessionInterface) {
	rs.tablelock.Lock()
	rs.tablelock.Unlock()
	sessionID := uuid.NewV4().String()
	//创建小cell（存用户session的）
	cell = NewRedisCell(sessionID,rs.connetpool)
	rs.store[sessionID] = cell
	return

}


//redis数据库
func NewRedisStore() *RedisStore {
	return &RedisStore{
		//创建大仓库
		store: make(map[string]SessionInterface, 200),
	}

}

//每次从redis中，获取1条用户session到保存到内存
func (rs *RedisStore) GetSession(sessionID string) (cell SessionInterface) {
	rs.tablelock.RLock()
	rs.tablelock.RUnlock()
	cell = NewRedisCell(sessionID, rs.connetpool)
	//每次都去redis里面加载！！
	err := cell.Load()
	if err != nil {
		return
	}
	rs.store[sessionID] = cell
	return

}


//从大仓库中，删除1条用户session到
func (rs *RedisStore) DelSession(sessionID string) {
	rs.tablelock.Lock()
	defer rs.tablelock.Unlock()
	delete(rs.store, sessionID)

}

//创建Redis 连接池
func InitRedis(option ...string) (pool *redis.Pool) {
	pool=&redis.Pool{
		MaxIdle:     60,
		MaxActive:   200,
		IdleTimeout: time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", option[0])
			if err != nil {
				fmt.Println("连接redis数据库失败", err)
				_ = conn.Close()
				return nil, err
			}
			//如果需要密码！
			if len(option[1]) > 1 {
				if _, err := conn.Do("AUTH", option[1]); err != nil {
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
	return pool
}
