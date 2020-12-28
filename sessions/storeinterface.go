package sessions

//session存储需要实现的接口列表
type StoreInterface interface {
	GetSession(sessionID string) (memCell SessionInterface)
	SetSession()(sessionID string)
	DelSession(sessionID string)
}
