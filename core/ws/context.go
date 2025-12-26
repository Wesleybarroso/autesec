package ws

import "secure-core/core/session"

type WSContext struct {
	UserID    string
	Session   *session.Session
}
