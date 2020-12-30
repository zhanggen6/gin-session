package sessions

//session存储需要实现的接口列表
type StoreInterface interface {
	InitStore()(store StoreInterface)
	GetSession(sessionID string) (Cell SessionInterface)
	SetSession()(sessionID string)
	DelSession(sessionID string)
}
