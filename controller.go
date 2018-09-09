package monitor

import (
	"net/http"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/validator"
	"github.com/joaosoft/webserver"
	"github.com/labstack/echo"
)

type Controller struct {
	s.aaa
	interactor *Interactor
}

func NewDbMigration(interactor *Interactor) *Controller {
	return &Controller{
		interactor: interactor,
	}
}

func (controller *Controller) GetProcessHandler(ctx web.Context) error {
	request := GetProcessRequest{
		IdProcess: ctx.Request.GetParam("id"),
	}

	if errs := validator.Validate(request); !errs.IsEmpty() {
		return ctx.JSON(http.StatusBadRequest, errs)
	}

	if process, err := controller.interactor.GetProcess(request.IdProcess); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if process == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK, process)
	}
}

func (controller *Controller) GetProcessesHandler(ctx echo.Context) error {
	if processes, err := controller.interactor.GetProcesses(ctx.QueryParams()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if processes == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK, processes)
	}
}

func (controller *Controller) CreateProcessHandler(ctx echo.Context) error {
	request := CreateProcessRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		err = log.WithFields(map[string]interface{}{"error": err}).
			Error("error getting body").ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if errs := validator.Validate(request.Body); !errs.IsEmpty() {
		err := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body request").ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	newProcess := Process{
		IdProcess:   request.Body.IdProcess,
		Type:        request.Body.Type,
		Name:        request.Body.Name,
		Description: request.Body.Description,
		DateFrom:    request.Body.DateFrom,
		DateTo:      request.Body.DateTo,
		TimeFrom:    request.Body.TimeFrom,
		TimeTo:      request.Body.TimeTo,
		DaysOff:     request.Body.DaysOff,
		Status:      request.Body.Status,
	}
	if err := controller.interactor.CreateProcess(&newProcess); err != nil {
		err := errors.New("0", err)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error creating process %s", request.Body.IdProcess).ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	} else {
		return ctx.NoContent(http.StatusCreated)
	}
}

func (controller *Controller) UpdateProcessHandler(ctx echo.Context) error {
	request := UpdateProcessRequest{
		IdProcess: ctx.Param("id"),
	}
	if err := ctx.Bind(&request.Body); err != nil {
		err := errors.New("0", err)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error getting body").ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if errs := validator.Validate(request); !errs.IsEmpty() {
		err := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body request").ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	updProcess := Process{
		IdProcess:   request.IdProcess,
		Type:        request.Body.Type,
		Name:        request.Body.Name,
		Description: request.Body.Description,
		DateFrom:    request.Body.DateFrom,
		DateTo:      request.Body.DateTo,
		TimeFrom:    request.Body.TimeFrom,
		TimeTo:      request.Body.TimeTo,
		DaysOff:     request.Body.DaysOff,
		Status:      request.Body.Status,
	}
	if err := controller.interactor.UpdateProcess(&updProcess); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (controller *Controller) UpdateProcessStatusHandler(ctx echo.Context) error {
	request := UpdateProcessStatusRequest{
		IdProcess: ctx.Param("id"),
		Status:    Status(ctx.Param("status")),
	}

	if errs := validator.Validate(request); !errs.IsEmpty() {
		err := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating query request").ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if errs := controller.interactor.UpdateProcessStatus(request.IdProcess, request.Status); errs != nil {
		return ctx.JSON(http.StatusInternalServerError, errs)
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (controller *Controller) UpdateProcessStatusCheckHandler(ctx echo.Context) error {
	request := UpdateProcessStatusRequest{
		IdProcess: ctx.Param("id"),
		Status:    Status(ctx.Param("status")),
	}

	if errs := validator.Validate(request); !errs.IsEmpty() {
		err := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating query request").ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if errs := controller.interactor.UpdateProcessStatus(request.IdProcess, request.Status); errs != nil {
		return ctx.JSON(http.StatusInternalServerError, errs)
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (controller *Controller) DeleteProcessHandler(ctx echo.Context) error {
	request := DeleteProcessRequest{
		IdProcess: ctx.Param("id"),
	}

	if errs := validator.Validate(request); !errs.IsEmpty() {
		err := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body request").ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := controller.interactor.DeleteProcess(request.IdProcess); err != nil {
		err := errors.New("0", err)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error deleting process by id %s", request.IdProcess).ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (controller *Controller) DeleteProcessesHandler(ctx echo.Context) error {
	if err := controller.interactor.DeleteProcesses(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}
