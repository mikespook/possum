package possum

import (
	"context"
	"net/http"

	"github.com/mikespook/possum/session"
)

// SessionName is the type of name for context
type SessionName string

const sessionCookieName = "possum"
const sessionName SessionName = "session"

// SessionFacotryFunc is the default factory function of session.
var SessionFacotryFunc session.FactoryFunc = session.NewFactory(session.CookieStorage(sessionCookieName, nil))

// Session extracts data from the request and returns session instance.
// It uses SessionFacotryFunc to initialise session if no session has been set yet.
func Session(w http.ResponseWriter, req *http.Request) (session *session.Session, err error) {
	var ok bool
	if session, ok = req.Context().Value(sessionName).(*Session); ok {
		return session, nil
	}

	session, err = SessionFacotryFunc(w, req)
	if err != nil {
		return nil, err
	}
	context.WithValue(req.Context(), sessionName, session)
	return session, nil
}
