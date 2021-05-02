package session

import (
	"github.com/gorilla/sessions"
)

func NewCookieStore() *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(nil))
}
