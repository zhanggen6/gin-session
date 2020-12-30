package sessions

import (
	"fmt"
	"sync"
)

/*
session就是以下数据结构

 {
 sessionID1:{sessionID1:{"islogin":true,"username":"张三" }，
 sessionID2:{sessionID2:{"islogin":true,"username":"李四" }，
 sessionID3:{sessionID3:{"islogin":true,"username":"王五" }，

 }

	go doc builtin.delete
	sc.session[key].var


*/
//1条session记录对应的内存结构体
type MemCell struct {
	sessionID string
	session   map[string]interface{}
	//行锁：相当于数据库中行锁
	rowlock sync.RWMutex
}

//初始化1条session记录
func NewMemCell() (cell *MemCell) {
	return &MemCell{
		session: make(map[string]interface{}, 20),
	}

}

func (sc *MemCell)Init(sessionID string) (cell *MemCell) {
	sc.sessionID = sessionID
	cell = sc
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
func (sc *MemCell) Set(key string, value interface{}) {
	sc.rowlock.Lock()
	defer sc.rowlock.Unlock()
	sc.session[key] = value
}

//删除 1条session记录中 key对应的值
func (sc *MemCell) Del(key string) {
	sc.rowlock.Lock()
	defer sc.rowlock.Unlock()
	delete(sc.session, key)

}

func (sc *MemCell)Save()(){
}
