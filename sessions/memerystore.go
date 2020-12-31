package sessions

import (
	uuid "github.com/satori/go.uuid"
	"sync"
)

//全局的session存储管理者,包含着每1条session记录
type MemStore struct {
	//这样就把所以的session记录组织在一起了！
	store map[string]SessionInterface
	//表锁：相当于数据库中的表锁
	tablelock sync.RWMutex
}

//初始化store空间
func NewMemStore() (stor *MemStore) {
	store := &MemStore{store: make(map[string]SessionInterface, 200)}
	return store
}

func (ms *MemStore) InitStore(option Option) (store StoreInterface) {
	return
}

//获取 1条session记录中 从session store中
func (sm *MemStore) GetSession(sessionID string) (memCell SessionInterface) {
	sm.tablelock.RLock()
	defer sm.tablelock.RUnlock()
	memCell = sm.store[sessionID]
	return memCell
}

//服务端创建1个sessionID(唯一的字符串)设置 1条session记录中 到session store中,
//返回给用户写到cookies
func (sm *MemStore) CreateSession() (cell SessionInterface) {
	sm.tablelock.Lock()
	defer sm.tablelock.Unlock()
	cellID := uuid.NewV4().String()
	cell = NewMemCell(cellID)
	sm.store[cellID] = cell
	return
}

//删除 1条session记录中 从session store中
func (sm *MemStore) DelSession(sessionID string) {
	sm.tablelock.Lock()
	defer sm.tablelock.Unlock()
	delete(sm.store, sessionID)
}
