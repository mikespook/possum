package possum

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mikespook/possum/view"
)

const statusKey contextKey = "status"
const datakey contextKey = "data"

// Redirect performs a redirecting to the url.
// It only works with the code belongs to one of StatusMovedPermanently,
// StatusFound, StatusSeeOther, and StatusTemporaryRedirect.
func Redirect(req *http.Request, code int, url string) {
	ctx := context.WithValue(req.Context(), statusKey, code)
	ctx = context.WithValue(req.Context(), dataKey, url)
	req.WithContext(ctx)
}

func handleRedirect(w http.ResponseWriter, req *http.Request) bool {
	ctx := req.Context()
	status := getStatus(ctx)
	data := getData(ctx)
	if status != http.StatusMovedPermanently &&
		status != http.StatusFound &&
		status != http.StatusSeeOther &&
		status != http.StatusTemporaryRedirect {
		return false
	}
	http.Redirect(w, req, data, status)
	return true
}

func handleRender(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	view := getView(ctx)
	data := getData(ctx)

	body, header, err := view.Render(data)
	if err != nil {
		panic(Error{http.StatusInternalServerError, err.Error()})
	}
	if header != nil {
		for key, values := range header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}
	w.WriteHeader(status)
	if _, err = w.Write(body); err != nil {
		panic(Error{http.StatusInternalServerError, err.Error()})
	}
}

func setContextValue(key contextKey, value interface{}) {
	req.WithContext(context.WithValue(req.Context(), key, value))
}

func getStatus(ctx context.Context) int {
	status, ok := ctx.Value(statusKey).(int)
	if !ok {
		panic(Error{http.StatusInternalServerError, fmt.Sprintf("Type casting error, `int` expected, `%T` got.", ctx.Value(statusKey))})
	}
	return status
}

func getData(ctx context.Context) string {
	data, ok := ctx.Value(dataKey).(string)
	if !ok {
		panic(Error{http.StatusInternalServerError, fmt.Sprintf("Type casting error, `string` expected, `%T` got.", ctx.Value(dataKey))})
	}
	return data
}

func setData(req *http.Request, data interface{}) {
	setContextValue(req, dataKey, data)
}

func getError(ctx context.Context) error {
	err, ok := ctx.Value(errorKey).(error)
	if !ok {
		panic(Error{http.StatusInternalServerError, fmt.Sprintf("Type casting error, `error` expected, `%T` got.", ctx.Value(errorKey))})
	}
	return err
}

func getFunc(ctx context.Context) http.HandlerFunc {
	f, ok := ctx.Value(funcKey).(http.HandlerFunc)
	if !ok {
		panic(Error{http.StatusInternalServerError, fmt.Sprintf("Type casting error, `http.HandlerFunc` expected, `%T` got.", ctx.Value(funcKey))})
	}
	return f
}

func getView(ctx context.Context) view.View {
	view, ok := ctx.Value(viewKey).(view.View)
	if !ok {
		panic(Error{http.StatusInternalServerError, fmt.Sprintf("Type casting error, `view.View` expected, `%T` got.", ctx.Value(viewKey))})
	}
	return view
}
