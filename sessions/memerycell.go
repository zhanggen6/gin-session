package sessions

import (
	"fmt"
	"sync"
)

/*
下数据结构

 {
 sessionID1:{sessionID1:{"islogin":true,"username":"张三" }，
 sessionID2:{sessionID2:{"islogin":true,"username":"李四" }，
 sessionID3:{sessionID3:{"islogin":true,"username":"王五" }，

 }
*/

//1条session记录对应的内存结构体
type MemCell struct {
	sessionID string
	session   map[string]interface{}
	//行锁：相当于数据库中行锁
	rowlock sync.RWMutex
}

//初始化1条session记录
func NewMemCell(cellID string) (cell *MemCell) {
	return &MemCell{
		//cell中的sessionID是为了能在中间件创建sessioncell之后getkey返回给cookie！
		sessionID: cellID,
		//存储用户session数据的map
		session: make(map[string]interface{}, 20),
	}

}

func (mc *MemCell) Init(sessionID string) (cell *MemCell) {
	mc.sessionID = sessionID
	cell = mc
	return
}

//获取 1条session记录中 key对应的值
func (sc *MemCell) Get(key string) (value interface{}) {
	sc.rowlock.RLock()
	defer sc.rowlock.RUnlock()
	value = sc.session[key]
	fmt.Println(value)
	return
}

//设置 1条session记录中 key对应的值
func (mc *MemCell) Set(key string, value interface{}) {
	mc.rowlock.Lock()
	defer mc.rowlock.Unlock()
	mc.session[key] = value
	fmt.Println(mc.session)
}

//删除 1条session记录中 key对应的值
func (mc *MemCell) Del(key string) {
	mc.rowlock.Lock()
	defer mc.rowlock.Unlock()
	delete(mc.session, key)

}
func (mc *MemCell) Save() {
	fmt.Println("内存版session，不需要save")
}

func (mc *MemCell) GetKey() (key string) {
	key = mc.sessionID
	return
}

func (mc *MemCell) Load() (err error) {
	return
}
