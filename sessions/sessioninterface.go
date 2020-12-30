package sessions

type SessionInterface interface {
	Save()
	Get(key string) (value interface{})
	Set(key string, value interface{})
	Del(key string)
}
