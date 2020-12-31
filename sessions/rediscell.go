package sessions

import (
	"encoding/json"
	"fmt"
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

//redis版session服务
type RedisCell struct {
	/*
		为什么我在store（数据库层面）中加了sessionID,在这也加sessionID呢？
		因为后期我们得支持用户直接拿着这个sessionID去redis/mysql保存数据
		后期可以支持用户clear()掉当前session
	*/
	sessionID string
	//内存中真正存储用户session记录的字典
	sessionMap map[string]interface{}
	//行锁
	rowlock sync.RWMutex
	//设置session记录在redis中的过期时间！
	expireTime int
	/*
		为什么我在store（数据库层面）中加了redisPoo,在这cell也要加redisPoo呢？
		因为我们要在用户save的时候连接数据库，把内存数据提交到数据库持久化！
	*/

	redisPool *redis.Pool
	/*
		是否需要更新？避免用户不管有没有修改都需要，频繁save()导致频繁连接数据库！
		只有修改了内存的数据才能save到数据库
	*/
	needFresh bool
}

//实例化1条redis类型的记录，也就是最真正存储用户session信息的map
func NewRedisCell(sessionid string, redisPool *redis.Pool) (cell *RedisCell) {
	return &RedisCell{
		sessionID:  sessionid,
		sessionMap: make(map[string]interface{}, 20),
		redisPool:  redisPool,
	}
}

//用户获取session中数据
func (rc *RedisCell) Get(key string) (value interface{}) {
	rc.rowlock.RLock()
	defer rc.rowlock.RUnlock()
	value = rc.sessionMap[key]
	return
}

//用户设置session中数据调用
func (rc *RedisCell) Set(key string, value interface{}) {
	rc.rowlock.Lock()
	defer rc.rowlock.Unlock()
	rc.sessionMap[key] = value
	rc.needFresh = true

}

//用户根据key删除session中数据调用
func (rc *RedisCell) Del(key string) {
	rc.rowlock.Lock()
	defer rc.rowlock.Unlock()
	delete(rc.sessionMap, key)
	rc.needFresh = true
}

func (rc *RedisCell) Load() (err error) {
	rc.rowlock.Lock()
	defer rc.rowlock.Unlock()
	conn := rc.redisPool.Get()
	defer conn.Close()
	//sessionID排上用场了吧！！
	bdata, err := conn.Do("get", rc.sessionID)
	if err != nil {
		fmt.Printf("使用%s  key从redis里获取数据失败！%v\n", rc.sessionID, err)
		return
	}
	//把从redis里面获取到的json字符串转换成字符串
	sdata, err := redis.String(bdata, err)
	if err != nil {
		fmt.Println("redis数据转换成内存数据失败")
	}
	err = json.Unmarshal([]byte(sdata),&rc.sessionMap)
	if err != nil {
		fmt.Println(err)
	}
	rc.needFresh = false
	return
}

//中间件用到的key
func (rc *RedisCell) GetKey() (key string) {
	key = rc.sessionID
	return
}

//用户申请向把当前session信息调用，提交到数据库（redis/mysql，，，）中持久化数据
func (rc *RedisCell) Save() {
	if rc.needFresh == true && len(rc.sessionMap) >= 1 {
		rc.rowlock.Lock()
		defer rc.rowlock.Unlock()
		if rc.redisPool == nil {
			fmt.Println("redis连接池错误！")
			return
		}
		conn := rc.redisPool.Get()
		defer conn.Close()
		marshalData, err := json.Marshal(rc.sessionMap)
		if err != nil {
			fmt.Println("你save session数据--->redis失败,session数据无法序列化成json数据！", err)
		}
		//sessionID排上用场了吧！
		_, err = conn.Do("set", rc.sessionID, marshalData)
		if err != nil {
			fmt.Println("set值失败")
		}
		//设置超时时间
		_, err = conn.Do("expire", rc.sessionID, MaxAge)
		if err != nil {
			fmt.Println(err)
		}
		//既然session数据以及保存到了redis了，在没有修改之前就不能让用户保存了哦！
		rc.needFresh = false

	}

}
