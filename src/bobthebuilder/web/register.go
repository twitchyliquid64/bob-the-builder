package web

import (
	"bobthebuilder/config"
	"bobthebuilder/logging"

	"github.com/hoisie/web"
	"golang.org/x/net/websocket"
)

// ### THIS FILE SHOULD CONTAIN ALL INITIALISATION CODE FOR BOTH TEMPLATES AND URL HANDLERS ###

func Initialise() {
	logging.Info("web", "Registering page handlers")
	registerCoreHandlers()
	registerApiHandlers()

	logging.Info("web", "Registering templates")
	registerCoreTemplates()
	web.SetDefaultDomain(config.All().Web.Domain)
}

func registerCoreHandlers() {
	web.Get("/", indexMainPage, config.All().Web.Domain)
	web.Get("/ws/events", websocket.Handler(ws_EventServer), config.All().Web.Domain)
	web.Get("/documentation/readme", documentationHandler, config.All().Web.Domain)
}

func registerApiHandlers() {
	web.Get("/api/definitions/reload", enqueueReloadHandler, config.All().Web.Domain)
	web.Get("/api/definitions", getDefinitionHandler, config.All().Web.Domain)
	web.Get("/api/history", getHistoryHandler, config.All().Web.Domain)
	web.Get("/api/status", getStatusHandler, config.All().Web.Domain)
	web.Get("/api/queue/new", enqueueBuildHandler, config.All().Web.Domain)
	web.Get("/api/lookup/buildparam", getBuildParamsLookupHandler, config.All().Web.Domain)
	web.Get("/api/definition/getIdByName", getDefIndexByIdHandler, config.All().Web.Domain)

	web.Get("/api/file/definitions", getDefinitionJSONHandler, config.All().Web.Domain)
	web.Get("/api/file/base", getBaseFileHandler, config.All().Web.Domain)
	web.Get("/api/files", getBrowserFilesData, config.All().Web.Domain)
	web.Get("/api/file/new/folder", newFolderHandler, config.All().Web.Domain)
	web.Get("/api/file/delete", deleteHandler, config.All().Web.Domain)
	web.Get("/api/file/new/file", newFileHandler, config.All().Web.Domain)
	web.Get("/api/file/new/definition", newDefFileHandler, config.All().Web.Domain)
	web.Get("/api/file/download/workspace", downloadWorkspaceFileHandler, config.All().Web.Domain)
	web.Post("/api/file/definitions/save", saveDefinitionJSONHandler, config.All().Web.Domain)
	web.Post("/api/file/base/save", saveBaseFileHandler, config.All().Web.Domain)

	web.Post("/api/queue/newWithOptions", enqueueBuildHandlerWithOptions, config.All().Web.Domain)

	web.Get("/api/cron", getCronHandler, config.All().Web.Domain)
	web.Post("/api/queue/cron", updateCronHandler, config.All().Web.Domain)
}

func registerCoreTemplates() {
	logError(registerTemplate("modals.tpl", "modals"), "Template load error: ")
	logError(registerTemplate("tailcontent.tpl", "tailcontent"), "Template load error: ")
	logError(registerTemplate("headcontent.tpl", "headcontent"), "Template load error: ")
	logError(registerTemplate("index.tpl", "index"), "Template load error: ")
	logError(registerTemplate("topnav.tpl", "topnav"), "Template load error: ")
}

func logError(e error, prefix string) {
	if e != nil {
		logging.Error("web", prefix, e.Error())
	}
}
