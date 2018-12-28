package web

import (
	"bobthebuilder/builder"
	"bobthebuilder/config"
	"bobthebuilder/logging"
	"github.com/hoisie/web"
)

func indexMainPage(ctx *web.Context) {
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

	mdl := modelBasic{Config: config.All(), Builder: builder.GetInstance()}

	if gAuth != nil {
		info, err := gAuth.AuthInfo(ctx)
		if err == nil && info != nil {
			mdl.Auth = info
		}
	}

	t := templates.Lookup("index")
	if t == nil {
		logging.Error("web", "No template found.")
		return
	}

	err := t.Execute(ctx.ResponseWriter, mdl)
	if err != nil {
		logging.Error("views-index", err)
	}
}

func loginMainPage(ctx *web.Context) {
	t := templates.Lookup("login")
	if t == nil {
		logging.Error("web", "No template found.")
		return
	}

	err := t.Execute(ctx.ResponseWriter, modelBasic{Config: config.All(), Builder: builder.GetInstance()})
	if err != nil {
		logging.Error("views-login", err)
	}
}
