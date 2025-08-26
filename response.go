package possum

import (
	"encoding/json"
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"

	"github.com/mikespook/possum/config"
)

const (
	// UUIDKey is the key used to store the UUID in the context
	UUIDKey ContextKey = "uuid"
)

var (
	// Predefined error responses for common HTTP status codes.
	ErrMethodNotAllowed = errors.New("method not allowed")

	// ErrNoRows is returned when a requested resource is not found
	InternalServerErrorResponse = Response{
		Error: &Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		},
	}

	UnauthorizedResponse = Response{
		Error: &Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		},
	}

	MethodNotAllowedResponse = Response{
		Error: &Error{
			Code:    http.StatusMethodNotAllowed,
			Message: "Method Not Allowed",
		},
	}
	BadRequestResponse = Response{
		Error: &Error{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
		},
	}
	NotFoundResponse = Response{
		Error: &Error{
			Code:    http.StatusNotFound,
			Message: "Not Found",
		},
	}
	NotImplementedResponse = Response{
		Error: &Error{
			Code:    http.StatusNotImplemented,
			Message: "Not Implemented",
		},
	}
	ConflictResponse = Response{
		Error: &Error{
			Code:    http.StatusConflict,
			Message: "Conflict",
		},
	}

	ForbiddenResponse = Response{
		Error: &Error{
			Code:    http.StatusForbidden,
			Message: "Forbidden",
		},
	}
)

type Response struct {
	UUID  uuid.UUID `json:"uuid"`
	Data  any       `json:"data,omitempty"`
	Error *Error    `json:"error,omitempty"`

	code int
}

// NewResponse creates a new Response object with a UUID from the request context.
func NewResponse(r *http.Request) *Response {
	id, ok := r.Context().Value(UUIDKey).(uuid.UUID)
	if !ok {
		id = uuid.New()
	}
	return &Response{
		UUID: id,
	}
}

// CloneResponse creates a deep copy of a Response object, including error details.
func CloneResponse(r Response) Response {
	var clonedError *Error
	if r.Error != nil {
		cloned := *r.Error
		clonedError = &cloned
	}

	return Response{
		UUID:  r.UUID,
		Data:  r.Data, // note: this won't deep copy if Data is a pointer or complex type
		Error: clonedError,
	}
}

// WriteResponse writes a Response object to an HTTP response writer, handling errors appropriately.
func WriteResponse(w http.ResponseWriter, resp Response, err error) {
	if err != nil {
		newResp := CloneResponse(resp)
		newResp.Error.Message = err.Error()
		if config.IsDebug() {
			newResp.Error.Stack = debug.Stack()
		}
		newResp.Write(w)
		return
	}
	resp.Write(w)
}

// SetError configures the Response object with error details, including stack trace in debug mode.
func (resp *Response) SetError(code int, message string) {
	resp.Error = &Error{
		Code:    code,
		Message: message,
		Stack:   nil, // 确保默认为 nil
	}
	if config.IsDebug() {
		resp.Error.Stack = debug.Stack()
	}
}

// SetData configures the Response object with response data.
func (resp *Response) SetData(data any) {
	resp.Data = data
}

// WriteHeader sets the HTTP status code for the Response object.
func (resp *Response) WriteHeader(code int) {
	resp.code = code
}

// Write serializes the Response object to the HTTP response writer with proper headers.
func (resp *Response) Write(w http.ResponseWriter) {
	resp.UUID = uuid.New() // Ensure UUID is set for each response
	w.Header().Set("X-Response-ID", resp.UUID.String())
	w.Header().Set("Content-Type", "application/json")
	if resp.Error == nil {
		if resp.code == 0 {
			resp.code = http.StatusOK
		}
		w.WriteHeader(resp.code)
	} else {
		if logResp, ok := w.(*logResponseWriter); ok {
			logResp.err = errors.New(resp.Error.Message)
			w = logResp
		}
		w.WriteHeader(resp.Error.Code)
	}
	if resp.code == http.StatusNoContent {
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Stack   []byte `json:"stack,omitempty"`
}
