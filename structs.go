package monitor

import (
	"github.com/joaosoft/types"
	"github.com/joaosoft/web"

	"time"
)

type Status string

type ErrorResponse struct {
	Code    web.Status `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Cause   string     `json:"cause,omitempty"`
}

type GetProcessRequest struct {
	IdProcess string `json:"id" validate:"notzero"`
}

type CreateProcessRequest struct {
	Body struct {
		IdProcess   string         `json:"id_process" validate:"notzero"`
		Type        string         `json:"type" validate:"notzero"`
		Name        string         `json:"name" validate:"notzero"`
		Description string         `json:"description"`
		DateFrom    *types.Date    `json:"date_from" validate:"special={date}"`
		DateTo      *types.Date    `json:"date_to" validate:"special={date}"`
		TimeFrom    *types.Time    `json:"time_from" validate:"special={time}"`
		TimeTo      *types.Time    `json:"time_to" validate:"special={time}"`
		DaysOff     *types.ListDay `json:"days_off" validate:"options=monday;tuesday;wednesday;thursday;friday;saturday;sunday"`
		Monitor     string         `json:"monitor"`
		Status      *Status        `json:"status" validate:"options=stopped;running"`
	}
}

type UpdateProcessRequest struct {
	IdProcess string `json:"id_process" validate:"notzero"`

	Body struct {
		Type        string         `json:"type" validate:"notzero"`
		Name        string         `json:"name" validate:"notzero"`
		Description string         `json:"description"`
		DateFrom    *types.Date    `json:"date_from" validate:"special={date}"`
		DateTo      *types.Date    `json:"date_to" validate:"special={date}"`
		TimeFrom    *types.Time    `json:"time_from" validate:"special={time}"`
		TimeTo      *types.Time    `json:"time_to" validate:"special={time}"`
		DaysOff     *types.ListDay `json:"days_off" validate:"options=monday;tuesday;wednesday;thursday;friday;saturday;sunday"`
		Monitor     string         `json:"monitor"`
		Status      *Status        `json:"status" validate:"options=stopped;running"`
	}
}

type UpdateProcessStatusRequest struct {
	IdProcess string `json:"id_process" validate:"notzero"`
	Status    Status `json:"status" validate:"options=stopped;running"`
}

type DeleteProcessRequest struct {
	IdProcess string `json:"id_process" validate:"notzero"`
}

type Process struct {
	IdProcess   string         `json:"id_process"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	DateFrom    *types.Date    `json:"date_from"`
	DateTo      *types.Date    `json:"date_to"`
	TimeFrom    *types.Time    `json:"time_from"`
	TimeTo      *types.Time    `json:"time_to"`
	DaysOff     *types.ListDay `json:"days_off"`
	Monitor     string         `json:"monitor"`
	Status      *Status        `json:"status"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

type ListProcess []*Process
