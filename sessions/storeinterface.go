package sessions

//session存储需要实现的接口列表
type StoreInterface interface {
	InitStore(option Option)(store StoreInterface)
	GetSession(sessionID string) (cell SessionInterface)
	CreateSession()(cell SessionInterface)
	DelSession(sessionID string)
}
