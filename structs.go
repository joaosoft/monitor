package monitor

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Date string
type Time string

// Scan implements the sql.Scanner interface.
func (a *Date) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		*a = Date(value)
		return nil
	case time.Time:
		*a = Date(value.Format("02-01-2006"))
		return nil
	case nil:
		a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to Date", src)
}

// Value implements the driver.Valuer interface.
func (a *Date) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return string(*a), nil
}

// Scan implements the sql.Scanner interface.
func (a *Time) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		*a = Time(value)
		return nil
	case time.Time:
		*a = Time(value.Format("15:04:05"))
		return nil
	case nil:
		a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to Time", src)
}

// Value implements the driver.Valuer interface.
func (a *Time) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return string(*a), nil
}

type ErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Cause   string `json:"cause,omitempty"`
}

type GetProcessRequest struct {
	IdProcess string `json:"id" validate:"nonzero"`
}

type CreateProcessRequest struct {
	Body struct {
		IdProcess   string   `json:"id_process" validate:"nonzero"`
		Type        string   `json:"type" validate:"nonzero"`
		Name        string   `json:"name" validate:"nonzero"`
		Description string   `json:"description"`
		DateFrom    Date     `json:"date_from" validate:"special={date}"`
		DateTo      Date     `json:"date_to" validate:"special={date}"`
		TimeFrom    Time     `json:"time_from" validate:"special={time}"`
		TimeTo      Time     `json:"time_to" validate:"special={time}"`
		Days        []string `json:"days" validate:"options=monday;tuesday;wednesday;thursday;friday;saturday;sunday"`
		Status      string   `json:"status" validate:"options=stopped;running;disabled"`
	}
}

type UpdateProcessRequest struct {
	IdProcess string `json:"id_process" validate:"nonzero"`

	Body struct {
		Type        string   `json:"type" validate:"nonzero"`
		Name        string   `json:"name" validate:"nonzero"`
		Description string   `json:"description"`
		DateFrom    Date     `json:"date_from" validate:"special={date}"`
		DateTo      Date     `json:"date_to" validate:"special={date}"`
		TimeFrom    Time     `json:"time_from" validate:"special={time}"`
		TimeTo      Time     `json:"time_to" validate:"special={time}"`
		Days        []string `json:"days" validate:"options=monday;tuesday;wednesday;thursday;friday;saturday;sunday"`
		Status      string   `json:"status" validate:"options=stopped;running;disabled"`
	}
}

type DeleteProcessRequest struct {
	IdProcess string `json:"id_process" validate:"nonzero"`
}

type Process struct {
	IdProcess   string    `json:"id_process"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	DateFrom    Date      `json:"date_from"`
	DateTo      Date      `json:"date_to"`
	TimeFrom    Time      `json:"time_from"`
	TimeTo      Time      `json:"time_to"`
	Days        []string  `json:"days"`
	Status      string    `json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListProcess []*Process
