package web

import (
	"bobthebuilder/config"
	"bobthebuilder/logging"
	"io/ioutil"

	"github.com/hoisie/web"
	"github.com/russross/blackfriday"
	//"errors"
)

func Run() {
	logging.Info("web", "Initialised server on ", config.All().Web.Listener)
	//web.RunTLS(config.All().Web.Listener, config.TLS())
	web.Run(config.All().Web.Listener)
}

func documentationHandler(ctx *web.Context) {
	d, _ := ioutil.ReadFile("README.md")
	output := blackfriday.MarkdownCommon(d)
	ctx.ResponseWriter.Write([]byte("<h1 style='margin-top: 14px;'>Documentation</h1>"))
	ctx.ResponseWriter.Write(output)
}
