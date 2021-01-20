package gosession

import "sync"

type Session struct {
	cookieName  string
	lock        sync.Mutex
	store       StoreInterface
	maxlifetime int
}

type StoreInterface interface {
	//初始化session 返回新的session变量
	SessionInt(sid string) (Session, error)
	//返回sid对应的Session变量
	SessionGet(sid string) (Session, error)
	SessionDestory(sid string) error
}

type SessionInterface interface {
	Set(key string, value interface{}) error
	Get(key string) interface{}
	Delete(key string) error
	SessionId() string
}
