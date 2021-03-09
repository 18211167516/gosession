package gosession

import (
	"net/http"
)

type Session struct {
	ID      string
	store   StoreInterface
	mamager *Manager
	ttl     int
}

func (s *Session) DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   s.mamager.CookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)
}

func (s *Session) Set(key string, value interface{}) error {
	return s.store.Set(s.ID, key, value, s.ttl)
}

func (s *Session) Get(key string) (interface{}, error) {
	return s.store.Get(s.ID, key)
}

func (s *Session) Delete(key string) error {
	return s.store.Remove(s.ID, key)
}
