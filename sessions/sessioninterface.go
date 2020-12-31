package sessions

type SessionInterface interface {
	Save()
	Load()(err error)
	Get(key string) (value interface{})
	Set(key string, value interface{})
	Del(key string)
	GetKey()(key string)
}
