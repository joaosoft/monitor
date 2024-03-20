package monitor

import (
	"database/sql/driver"

	"github.com/joaosoft/errors"
)

func (s *Status) Scan(src interface{}) error {
	if s == nil {
		return nil
	}

	switch value := src.(type) {
	case string:
		*s = Status(value)
		return nil
	case nil:
		s = nil
		return nil
	}

	return errors.New(errors.LevelError, 0, "pq: cannot convert %T to %T", src, *s)
}

func (s Status) Value() (driver.Value, error) {
	if s == "" {
		return nil, nil
	}

	return string(s), nil
}
