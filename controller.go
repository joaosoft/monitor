package monitor

import (
	"github.com/joaosoft/errors"
	"github.com/joaosoft/validator"
	"github.com/joaosoft/webserver"
)

type Controller struct {
	interactor *Interactor
}

func NewDbMigration(interactor *Interactor) *Controller {
	return &Controller{
		interactor: interactor,
	}
}

func (controller *Controller) GetProcessHandler(ctx *webserver.Context) error {
	request := GetProcessRequest{
		IdProcess: ctx.Request.GetParam("id"),
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		return ctx.Response.JSON(webserver.StatusBadRequest, errs)
	}

	if process, err := controller.interactor.GetProcess(request.IdProcess); err != nil {
		return ctx.Response.JSON(webserver.StatusInternalServerError, ErrorResponse{Code: webserver.StatusInternalServerError, Message: err.Error()})
	} else if process == nil {
		return ctx.Response.NoContent(webserver.StatusNotFound)
	} else {
		return ctx.Response.JSON(webserver.StatusOK, process)
	}
}

func (controller *Controller) GetProcessesHandler(ctx *webserver.Context) error {
	if processes, err := controller.interactor.GetProcesses(ctx.Request.UrlParms); err != nil {
		return ctx.Response.JSON(webserver.StatusInternalServerError, ErrorResponse{Code: webserver.StatusInternalServerError, Message: err.Error()})
	} else if processes == nil {
		return ctx.Response.NoContent(webserver.StatusNotFound)
	} else {
		return ctx.Response.JSON(webserver.StatusOK, processes)
	}
}

func (controller *Controller) CreateProcessHandler(ctx *webserver.Context) error {
	request := CreateProcessRequest{}
	if err := ctx.Request.Bind(&request.Body); err != nil {
		err = log.WithFields(map[string]interface{}{"error": err}).
			Error("error getting body").ToError()
		return ctx.Response.JSON(webserver.StatusBadRequest, ErrorResponse{Code: webserver.StatusBadRequest, Message: err.Error()})
	}

	if errs := validator.Validate(request.Body); len(errs) > 0 {
		err := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body request").ToError()
		return ctx.Response.JSON(webserver.StatusBadRequest, ErrorResponse{Code: webserver.StatusBadRequest, Message: err.Error()})
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
		return ctx.Response.JSON(webserver.StatusBadRequest, ErrorResponse{Code: webserver.StatusBadRequest, Message: err.Error()})
	} else {
		return ctx.Response.NoContent(webserver.StatusCreated)
	}
}

func (controller *Controller) UpdateProcessHandler(ctx *webserver.Context) error {
	request := UpdateProcessRequest{
		IdProcess: ctx.Request.GetParam("id"),
	}
	if err := ctx.Request.Bind(&request.Body); err != nil {
		err := errors.New("0", err)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error getting body").ToError()
		return ctx.Response.JSON(webserver.StatusBadRequest, ErrorResponse{Code: webserver.StatusBadRequest, Message: err.Error()})
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body request").ToError()
		return ctx.Response.JSON(webserver.StatusBadRequest, ErrorResponse{Code: webserver.StatusBadRequest, Message: err.Error()})
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
		return ctx.Response.JSON(webserver.StatusInternalServerError, ErrorResponse{Code: webserver.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.Response.NoContent(webserver.StatusOK)
	}
}

func (controller *Controller) UpdateProcessStatusHandler(ctx *webserver.Context) error {
	request := UpdateProcessStatusRequest{
		IdProcess: ctx.Request.GetParam("id"),
		Status:    Status(ctx.Request.GetParam("status")),
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating query request").ToError()
		return ctx.Response.JSON(webserver.StatusBadRequest, ErrorResponse{Code: webserver.StatusBadRequest, Message: err.Error()})
	}

	if errs := controller.interactor.UpdateProcessStatus(request.IdProcess, request.Status); errs != nil {
		return ctx.Response.JSON(webserver.StatusInternalServerError, errs)
	} else {
		return ctx.Response.NoContent(webserver.StatusOK)
	}
}

func (controller *Controller) UpdateProcessStatusCheckHandler(ctx *webserver.Context) error {
	request := UpdateProcessStatusRequest{
		IdProcess: ctx.Request.GetParam("id"),
		Status:    Status(ctx.Request.GetParam("status")),
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating query request").ToError()
		return ctx.Response.JSON(webserver.StatusBadRequest, ErrorResponse{Code: webserver.StatusBadRequest, Message: err.Error()})
	}

	if errs := controller.interactor.UpdateProcessStatus(request.IdProcess, request.Status); errs != nil {
		return ctx.Response.JSON(webserver.StatusInternalServerError, errs)
	} else {
		return ctx.Response.NoContent(webserver.StatusOK)
	}
}

func (controller *Controller) DeleteProcessHandler(ctx *webserver.Context) error {
	request := DeleteProcessRequest{
		IdProcess: ctx.Request.GetParam("id"),
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New("0", errs)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body request").ToError()
		return ctx.Response.JSON(webserver.StatusBadRequest, ErrorResponse{Code: webserver.StatusBadRequest, Message: err.Error()})
	}

	if err := controller.interactor.DeleteProcess(request.IdProcess); err != nil {
		err := errors.New("0", err)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error deleting process by id %s", request.IdProcess).ToError()
		return ctx.Response.JSON(webserver.StatusBadRequest, ErrorResponse{Code: webserver.StatusBadRequest, Message: err.Error()})
	} else {
		return ctx.Response.NoContent(webserver.StatusOK)
	}
}

func (controller *Controller) DeleteProcessesHandler(ctx *webserver.Context) error {
	if err := controller.interactor.DeleteProcesses(); err != nil {
		return ctx.Response.JSON(webserver.StatusInternalServerError, ErrorResponse{Code: webserver.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.Response.NoContent(webserver.StatusOK)
	}
}
