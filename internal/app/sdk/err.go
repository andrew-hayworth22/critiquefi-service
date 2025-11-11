package sdk

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func ToError(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	if err != nil {
		return nil, true
	}
	return nil, false
}

func HandleError(w http.ResponseWriter, err error) {
	e, hasErr := ToError(err)
	if !hasErr {
		return
	}

	if e == nil {
		fmt.Println(err)
		e = NewError(http.StatusInternalServerError, "internal server error")
	}

	_ = Respond(w, e, e.Code)
}
