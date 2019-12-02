package possum

import (
	"net/http"

	"github.com/mikespook/possum/session"
)

type contextKey string

const sessionCookieName = "possum"
const sessionKey contextKey = "session"

// SessionFacotryFunc is the default factory function of session.
var SessionFacotryFunc = session.NewFactory(session.CookieStorage(sessionCookieName, nil))

// Session extracts data from the request and returns session instance.
// It uses SessionFacotryFunc to initialise session if no session has been set yet.
func Session(w http.ResponseWriter, req *http.Request) (session *session.Session, err error) {
	var ok bool
	if session, ok = req.Context().Value(sessionKey).(*Session); ok {
		return session, nil
	}

	session, err = SessionFacotryFunc(w, req)
	if err != nil {
		return nil, err
	}
	setContextValue(req, sessionKey, session)
	return session, nil
}

func handleSession(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if session, ok := ctx.Value(sessionKey).(*Session); ok {
		if err := session.Flush(); err != nil {
			panic(Error{http.StatusInternalServerError, err.Error()})
		}
	}
}
