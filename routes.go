package monitor

import (
	"net/http"

	"github.com/joaosoft/manager"
)

func (controller *Controller) RegisterRoutes(web manager.IWeb) error {
	return web.AddRoutes(
		manager.NewRoute(http.MethodGet, "/api/v1/processes/:id", controller.GetProcessHandler),
		manager.NewRoute(http.MethodGet, "/api/v1/processes", controller.GetProcessesHandler),
		manager.NewRoute(http.MethodPost, "/api/v1/processes", controller.CreateProcessHandler),
		manager.NewRoute(http.MethodPut, "/api/v1/processes/:id", controller.UpdateProcessHandler),
		manager.NewRoute(http.MethodDelete, "/api/v1/processes/:id", controller.DeleteProcessHandler),
		manager.NewRoute(http.MethodDelete, "/api/v1/processes", controller.DeleteProcessesHandler),
	)
}
