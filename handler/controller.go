package handler

import "github.com/jschaefer-io/IDaaS/server"

type baseController struct {
	components *server.Components
	settings   *server.Settings
}

func newBaseController(components *server.Components, settings *server.Settings) baseController {
	return baseController{
		components: components,
		settings:   settings,
	}
}
