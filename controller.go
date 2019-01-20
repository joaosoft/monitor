package monitor

import (
	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
	"github.com/joaosoft/validator"
	"github.com/joaosoft/web"
)

type Controller struct {
	interactor *Interactor
	logger logger.ILogger
}

func (monitor *Monitor) NewController(interactor *Interactor) *Controller {
	return &Controller{
		interactor: interactor,
		logger: monitor.logger,
	}
}

func (controller *Controller) DoNothing(ctx *web.Context) error {
	return nil
}

func (controller *Controller) GetProcessHandler(ctx *web.Context) error {
	request := GetProcessRequest{
		IdProcess: ctx.Request.GetUrlParam("id"),
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		return ctx.Response.JSON(web.StatusBadRequest, errs)
	}

	if process, err := controller.interactor.GetProcess(request.IdProcess); err != nil {
		return ctx.Response.JSON(web.StatusInternalServerError, ErrorResponse{Code: web.StatusInternalServerError, Message: err.Error()})
	} else if process == nil {
		return ctx.Response.NoContent(web.StatusNotFound)
	} else {
		return ctx.Response.JSON(web.StatusOK, process)
	}
}

func (controller *Controller) GetProcessesHandler(ctx *web.Context) error {
	if processes, err := controller.interactor.GetProcesses(ctx.Request.Params); err != nil {
		return ctx.Response.JSON(web.StatusInternalServerError, ErrorResponse{Code: web.StatusInternalServerError, Message: err.Error()})
	} else if processes == nil {
		return ctx.Response.NoContent(web.StatusNotFound)
	} else {
		return ctx.Response.JSON(web.StatusOK, processes)
	}
}

func (controller *Controller) CreateProcessHandler(ctx *web.Context) error {
	request := CreateProcessRequest{}
	if err := ctx.Request.Bind(&request.Body); err != nil {
		err = controller.logger.WithFields(map[string]interface{}{"error": err}).
			Error("error getting body").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	}

	if errs := validator.Validate(request.Body); len(errs) > 0 {
		err := errors.New(errors.ErrorLevel, 0, errs)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body request").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
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
		err := errors.New(errors.ErrorLevel, 0, err)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error creating process %s", request.Body.IdProcess).ToError()
		return ctx.Response.JSON(web.StatusBadRequest, ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	} else {
		return ctx.Response.NoContent(web.StatusCreated)
	}
}

func (controller *Controller) UpdateProcessHandler(ctx *web.Context) error {
	request := UpdateProcessRequest{
		IdProcess: ctx.Request.GetUrlParam("id"),
	}
	if err := ctx.Request.Bind(&request.Body); err != nil {
		err := errors.New(errors.ErrorLevel, 0, err)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error getting body").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New(errors.ErrorLevel, 0, errs)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body request").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
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
		return ctx.Response.JSON(web.StatusInternalServerError, ErrorResponse{Code: web.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.Response.NoContent(web.StatusOK)
	}
}

func (controller *Controller) UpdateProcessStatusHandler(ctx *web.Context) error {
	request := UpdateProcessStatusRequest{
		IdProcess: ctx.Request.GetUrlParam("id"),
		Status:    Status(ctx.Request.GetUrlParam("status")),
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New(errors.ErrorLevel, 0, errs)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating query request").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	}

	if errs := controller.interactor.UpdateProcessStatus(request.IdProcess, request.Status); errs != nil {
		return ctx.Response.JSON(web.StatusInternalServerError, errs)
	} else {
		return ctx.Response.NoContent(web.StatusOK)
	}
}

func (controller *Controller) UpdateProcessStatusCheckHandler(ctx *web.Context) error {
	request := UpdateProcessStatusRequest{
		IdProcess: ctx.Request.GetUrlParam("id"),
		Status:    Status(ctx.Request.GetUrlParam("status")),
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New(errors.ErrorLevel, 0, errs)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating query request").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	}

	if errs := controller.interactor.UpdateProcessStatusCheck(request.IdProcess, request.Status); errs != nil {
		return ctx.Response.JSON(web.StatusInternalServerError, errs)
	} else {
		return ctx.Response.NoContent(web.StatusOK)
	}
}

func (controller *Controller) DeleteProcessHandler(ctx *web.Context) error {
	request := DeleteProcessRequest{
		IdProcess: ctx.Request.GetUrlParam("id"),
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New(errors.ErrorLevel, 0, errs)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body request").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	}

	if err := controller.interactor.DeleteProcess(request.IdProcess); err != nil {
		err := errors.New(errors.ErrorLevel, 0, err)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error deleting process by id %s", request.IdProcess).ToError()
		return ctx.Response.JSON(web.StatusBadRequest, ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	} else {
		return ctx.Response.NoContent(web.StatusOK)
	}
}

func (controller *Controller) DeleteProcessesHandler(ctx *web.Context) error {
	if err := controller.interactor.DeleteProcesses(); err != nil {
		return ctx.Response.JSON(web.StatusInternalServerError, ErrorResponse{Code: web.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.Response.NoContent(web.StatusOK)
	}
}
