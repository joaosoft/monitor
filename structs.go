package monitor

type Date string
type Time string

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
		IdProcess   string   `json:"id" validate:"nonzero"`
		Type        string   `json:"type" validate:"nonzero"`
		Name        string   `json:"name" validate:"nonzero"`
		Description string   `json:"description"`
		DateFrom    Date     `json:"date_from" validate:"regex={date}"`
		DateTo      Date     `json:"date_to" validate:"regex={date}"`
		TimeFrom    Time     `json:"time_from" validate:"regex={time}"`
		TimeTo      Time     `json:"time_to" validate:"regex={time}"`
		Days        []string `json:"days" validate:"options=monday;tuesday;wednesday;thursday;friday;saturday;sunday"`
		Status      string   `json:"status" validate:"options=stopped;running;disabled"`
	}
}

type UpdateProcessRequest struct {
	IdProcess string `json:"id" validate:"nonzero"`

	Body struct {
		Type        string   `json:"type" validate:"nonzero"`
		Name        string   `json:"name" validate:"nonzero"`
		Description string   `json:"description"`
		DateFrom    Date     `json:"date_from" validate:"regex={date}"`
		DateTo      Date     `json:"date_to" validate:"regex={date}"`
		TimeFrom    Time     `json:"time_from" validate:"regex={time}"`
		TimeTo      Time     `json:"time_to" validate:"regex={time}"`
		Days        []string `json:"days" validate:"options=monday;tuesday;wednesday;thursday;friday;saturday;sunday"`
		Status      string   `json:"status" validate:"options=stopped;running;disabled"`
	}
}

type DeleteProcessRequest struct {
	IdProcess string `json:"id" validate:"nonzero"`
}

type Process struct {
	IdProcess   string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	DateFrom    Date     `json:"date_from"`
	DateTo      Date     `json:"date_to"`
	TimeFrom    Time     `json:"time_from"`
	TimeTo      Time     `json:"time_to"`
	Days        []string `json:"days"`
	Status      string   `json:"status"`
	CreatedAt   Time     `json:"created_at"`
	UpdatedAt   Time     `json:"updated_at"`
}

type ListProcess []*Process
