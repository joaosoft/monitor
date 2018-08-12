package monitor

import (
	"net/http"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/validator"
	"github.com/labstack/echo"
)

type Controller struct {
	interactor *Interactor
}

func NewController(interactor *Interactor) *Controller {
	return &Controller{
		interactor: interactor,
	}
}

func (controller *Controller) GetProcessHandler(ctx echo.Context) error {
	request := GetProcessRequest{
		IdProcess: ctx.Param("id"),
	}

	if errs := validator.Validate(request); !errs.IsEmpty() {
		return ctx.JSON(http.StatusBadRequest, errs)
	}

	if process, err := controller.interactor.GetProcess(request.IdProcess); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if process == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK, process)
	}
}

func (controller *Controller) GetProcessesHandler(ctx echo.Context) error {
	if processes, err := controller.interactor.GetProcesses(ctx.QueryParams()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if processes == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK, processes)
	}
}

func (controller *Controller) CreateProcessHandler(ctx echo.Context) error {
	request := CreateProcessRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		newErr := errors.New("0", err)
		log.WithFields(map[string]interface{}{"error": err, "cause": newErr.Cause()}).
			Error("error getting body").ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if errs := validator.Validate(request.Body); !errs.IsEmpty() {
		newErr := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
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
		newErr := errors.New("0", err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Errorf("error creating process %s", request.Body.IdProcess).ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	} else {
		return ctx.NoContent(http.StatusCreated)
	}
}

func (controller *Controller) UpdateProcessHandler(ctx echo.Context) error {
	request := UpdateProcessRequest{
		IdProcess: ctx.Param("id"),
	}
	if err := ctx.Bind(&request.Body); err != nil {
		newErr := errors.New("0", err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if errs := validator.Validate(request); !errs.IsEmpty() {
		newErr := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
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
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
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
		newErr := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating query request").ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
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
		newErr := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := controller.interactor.DeleteProcess(request.IdProcess); err != nil {
		newErr := errors.New("0", err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Errorf("error deleting process by id %s", request.IdProcess).ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (controller *Controller) DeleteProcessesHandler(ctx echo.Context) error {
	if err := controller.interactor.DeleteProcesses(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}
