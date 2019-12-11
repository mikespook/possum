package possum

import (
	"net/http"

	"github.com/mikespook/possum/session"
)

const sessionCookieName = "possum"

// SessionFacotryFunc is the default factory function of session.
var SessionFacotryFunc = session.NewFactory(session.CookieStorage(sessionCookieName, nil))

// Session extracts data from the request and returns session instance.
// It uses SessionFacotryFunc to initialise session if no session has been set yet.
func Session(w http.ResponseWriter, req *http.Request) (sn *session.Session, err error) {
	sn, err = SessionFacotryFunc(w, req)
	if err != nil {
		return nil, err
	}
	return sn, nil
}
