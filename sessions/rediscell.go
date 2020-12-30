package sessions

import (
	"github.com/gomodule/redigo/redis"
	"sync"
)

/*
session就是以下数据结构

redis{
sessionID1:{sessionID1:{"islogin":true,"username":"张三" }，
sessionID2:{sessionID2:{"islogin":true,"username":"李四" }，
sessionID3:{sessionID3:{"islogin":true,"username":"王五" }，

}
 */

/

//redis版session服务
type RedisCell struct {
	//为什么我在store中加了sessionID,在这也加了呢？是以为我set的时候需要通过key去redis获取value
	SessionID  string
	SessionMap map[string]interface{}
	rowlock    sync.RWMutex //行锁
	expireTime int
	redisPool  *redis.Pool
	isSave     bool     //是否需要更新？避免用户不管有没有修改都需要，频繁save()
}

func NewRedisCell() (cell *RedisCell) {
	return &RedisCell{
		SessionMap: make(map[string]interface{}, 20),
	}
}

func (rc *RedisCell) UpdateData(SessionMap map[string]interface{}) (NewSession map[string]interface{}) {
	rc.rowlock.Lock()
	defer rc.rowlock.Unlock()
	rc.SessionMap = SessionMap
	return rc.SessionMap
}

//获取
func (rc *RedisCell) Get(key string) (value interface{}) {
	rc.rowlock.RLock()
	defer rc.rowlock.RUnlock()
	value = rc.SessionMap[key]
	return
}

//设置
func (rc *RedisCell) Set(key string, value interface{}) {
	rc.rowlock.Lock()
	defer rc.rowlock.Unlock()
	rc.SessionMap[key] = value
}

//删除
func (rc *RedisCell) Del(key string) {

}

//向redis中持久化数据
func (rc *RedisCell) Save() {
	return
}
