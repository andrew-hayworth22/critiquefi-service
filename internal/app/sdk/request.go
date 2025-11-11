package sdk

import (
	"errors"
	"io"
	"net/http"
	"regexp"
)

type Decoder interface {
	Decode(data []byte) error
}

type validator interface {
	Validate() error
}

// Decode reads from the HTTP request, decodes the data into a structure, and optionally validates the data
func Decode(r *http.Request, v Decoder) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return NewError(http.StatusBadRequest, "malformed request")
	}

	if err := v.Decode(data); err != nil {
		return NewError(http.StatusBadRequest, "malformed request")
	}

	if v, ok := v.(validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// ValidateText validates a text field given a regex pattern and length constraints
func ValidateText(text, pattern string, min, max int) error {
	if len(text) < min {
		return errors.New("too short")
	}

	if len(text) > max {
		return errors.New("too long")
	}

	if len(pattern) == 0 {
		return nil
	}

	if _, err := regexp.MatchString(pattern, text); err != nil {
		return errors.New("invalid characters")
	}
	return nil
}
