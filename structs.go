package monitor

import (
	"time"
)

type Date string
type Time string
type Status string
type Day string
type ListDay []Day

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
		DateFrom    *Date    `json:"date_from" validate:"special={date}"`
		DateTo      *Date    `json:"date_to" validate:"special={date}"`
		TimeFrom    *Time    `json:"time_from" validate:"special={time}"`
		TimeTo      *Time    `json:"time_to" validate:"special={time}"`
		DaysOff     *ListDay `json:"days_off" validate:"options=monday;tuesday;wednesday;thursday;friday;saturday;sunday"`
		Monitor     string   `json:"monitor"`
		Status      Status   `json:"status" validate:"options=stopped;running"`
	}
}

type UpdateProcessRequest struct {
	IdProcess string `json:"id_process" validate:"nonzero"`

	Body struct {
		Type        string   `json:"type" validate:"nonzero"`
		Name        string   `json:"name" validate:"nonzero"`
		Description string   `json:"description"`
		DateFrom    *Date    `json:"date_from" validate:"special={date}"`
		DateTo      *Date    `json:"date_to" validate:"special={date}"`
		TimeFrom    *Time    `json:"time_from" validate:"special={time}"`
		TimeTo      *Time    `json:"time_to" validate:"special={time}"`
		DaysOff     *ListDay `json:"days_off" validate:"options=monday;tuesday;wednesday;thursday;friday;saturday;sunday"`
		Monitor     string   `json:"monitor"`
		Status      Status   `json:"status" validate:"options=stopped;running"`
	}
}

type UpdateProcessStatusRequest struct {
	IdProcess string `json:"id_process" validate:"nonzero"`
	Status    Status `json:"status" validate:"options=stopped;running"`
}

type DeleteProcessRequest struct {
	IdProcess string `json:"id_process" validate:"nonzero"`
}

type Process struct {
	IdProcess   string    `json:"id_process"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	DateFrom    *Date     `json:"date_from"`
	DateTo      *Date     `json:"date_to"`
	TimeFrom    *Time     `json:"time_from"`
	TimeTo      *Time     `json:"time_to"`
	DaysOff     *ListDay  `json:"days_off"`
	Monitor     string    `json:"monitor"`
	Status      Status    `json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListProcess []*Process

func (l ListDay) Contains(day Day) bool {
	for _, value := range l {
		if day == Day(value) {
			return true
		}
	}

	return false
}
