package web

import (
	"bobthebuilder/config"
	"bobthebuilder/logging"
	"bobthebuilder/web/auth"
	"io/ioutil"

	"github.com/hoisie/web"
	"github.com/russross/blackfriday"
	//"errors"
)

var gAuth auth.Auther

//Run() initialises the web server based on the configuration package.
func Run() {

	if config.All().Web.RequireAuth {
		if config.All().Web.PamAuth {
			gAuth = auth.MultiAuth(auth.CookieAuth(config.All()), &auth.PAMAuther{})
		} else {
			gAuth = auth.CookieAuth(config.All())
		}
	}

	if config.All().TLS.PrivateKey == "" {
		logging.Info("web", "Initialised HTTP server on ", config.All().Web.Listener)
		web.Run(config.All().Web.Listener)
	} else {
		logging.Info("web", "Initialised HTTPS server on ", config.All().Web.Listener)
		web.RunTLS(config.All().Web.Listener, config.TLS())
	}

}

func documentationHandler(ctx *web.Context) {
	d, _ := ioutil.ReadFile("README.md")
	output := blackfriday.MarkdownCommon(d)
	ctx.ResponseWriter.Write([]byte("<h1 style='margin-top: 14px;'>Documentation</h1>"))
	ctx.ResponseWriter.Write(output)
}
