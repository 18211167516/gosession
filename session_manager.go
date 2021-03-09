package gosession

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var Stores = make(map[string]StoreInterface)

type Manager struct {
	CookieName  string
	store       StoreInterface
	maxlifetime int
	idFunc      IdFunc
	lock        sync.Mutex
}

type IdFunc = func() string

/*默认生成sessionid（全局用户uid唯一标识）*/
func defaultSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	} else {
		return fmt.Sprintf("%s", md5.Sum(b))
	}
}

func NewDefaultIDManager(storeName string, cookieName string, maxlifetime int) (*Manager, error) {
	storeName = strings.ToLower(storeName)
	if store, ok := Stores[storeName]; ok {
		return &Manager{
			CookieName:  cookieName,
			maxlifetime: maxlifetime,
			store:       store,
			idFunc:      defaultSessionId,
		}, nil
	} else {
		return nil, fmt.Errorf("Session:unknown %s store please import", storeName)
	}

}

func NewManager(storeName string, cookieName string, maxlifetime int, f IdFunc) (*Manager, error) {
	storeName = strings.ToLower(storeName)
	if store, ok := Stores[storeName]; ok {
		return &Manager{
			CookieName:  cookieName,
			maxlifetime: maxlifetime,
			store:       store,
			idFunc:      f,
		}, nil
	} else {
		return nil, fmt.Errorf("Session:unknown %s store please import", storeName)
	}

}

/*生成session ID*/
func (manager *Manager) sessionId() string {
	return manager.idFunc()
}

/*生成session结构体*/
func (manager *Manager) NewSesssion(sid string) *Session {
	return &Session{
		ID:      sid,
		store:   manager.store,
		mamager: manager,
		ttl:     manager.maxlifetime,
	}
}

/*获取session*/
func (manager *Manager) GetSession(w http.ResponseWriter, r *http.Request) (s *Session) {
	/*加锁*/
	manager.lock.Lock()
	defer manager.lock.Unlock()

	cookie, err := r.Cookie(manager.CookieName)
	if err != nil || cookie.Value == "" {
		/*获取sid*/
		sid := manager.sessionId()
		s = manager.NewSesssion(sid)

		cookie := http.Cookie{
			Name:     manager.CookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   manager.maxlifetime,
		}

		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		s = manager.NewSesssion(sid)
	}

	return s
}

/*注册存储器*/
func Register(name string, store StoreInterface) {
	if store == nil {
		log.Panic("Session: Register store is nil")
	}

	name = strings.ToLower(name)

	if _, ok := Stores[name]; ok {
		log.Panic("Session: Register store is exist")
	}

	Stores[name] = store
}
