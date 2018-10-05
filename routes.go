package monitor

import (
	"github.com/joaosoft/manager"
	"github.com/joaosoft/web"
)

func (controller *Controller) RegisterRoutes(w manager.IWeb) error {
	return w.AddRoutes(
		manager.NewRoute(string(web.MethodOptions), "*", controller.DoNothing, web.MiddlewareOptions()),
		manager.NewRoute(string(web.MethodGet), "/api/v1/processes/:id", controller.GetProcessHandler),
		manager.NewRoute(string(web.MethodGet), "/api/v1/processes", controller.GetProcessesHandler),
		manager.NewRoute(string(web.MethodPost), "/api/v1/processes", controller.CreateProcessHandler),
		manager.NewRoute(string(web.MethodPut), "/api/v1/processes/:id", controller.UpdateProcessHandler),
		manager.NewRoute(string(web.MethodPut), "/api/v1/processes/:id/status/:status", controller.UpdateProcessStatusHandler),
		manager.NewRoute(string(web.MethodPut), "/api/v1/processes/:id/status/:status/check", controller.UpdateProcessStatusCheckHandler),
		manager.NewRoute(string(web.MethodDelete), "/api/v1/processes/:id", controller.DeleteProcessHandler),
		manager.NewRoute(string(web.MethodDelete), "/api/v1/processes", controller.DeleteProcessesHandler),
	)
}
