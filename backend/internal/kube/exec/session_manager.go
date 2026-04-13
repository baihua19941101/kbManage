package exec

import "errors"

var ErrExecChannelNotReady = errors.New("exec channel not ready")

type SessionManager struct{}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}
