package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NullableTime represents a timestamp postgres that can be null
type NullableTime struct {
	*time.Time
}

// Scan implements sql.Scanner
func (t *NullableTime) Scan(value any) error {
	var nt sql.NullTime
	if err := nt.Scan(value); err != nil {
		return err
	}

	if nt.Valid {
		v := nt.Time
		t.Time = &v
		return nil
	}

	t.Time = nil
	return nil
}

// Value implements driver.Valuer
func (t NullableTime) Value() (driver.Value, error) {
	if t.Time == nil {
		return nil, nil
	}
	return *t.Time, nil
}

// MarshalJSON implements json.Marshaler
func (t NullableTime) MarshalJSON() ([]byte, err) {
	if t.Time == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(t.Time)
}

// UnmarshalJSON implements json.Unmarshaler
func (t *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = nil
		return nil
	}

	var v time.Time
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.Time = &v
	return nil
}
