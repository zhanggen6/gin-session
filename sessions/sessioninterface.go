package sessions

type SessionInterface interface {
	Init(sessionID string) (cell *SessionCell)
	Get(key string) (value interface{})
	Set(key string, value interface{})
	Del(key string)
}