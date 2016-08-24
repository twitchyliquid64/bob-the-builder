package web

import (
	"bobthebuilder/config"
	"bobthebuilder/logging"
	"io/ioutil"

	"github.com/hoisie/web"
	"github.com/russross/blackfriday"
	//"errors"
)

//Run() initialises the web server based on the configuration package.
func Run() {

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
