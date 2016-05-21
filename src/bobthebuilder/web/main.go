package web

import (
  "bobthebuilder/logging"
  "bobthebuilder/config"
  "github.com/hoisie/web"
  //"errors"
)

func Run() {
  logging.Info("web", "Initialised server on ", config.All().Web.Listener)
  //web.RunTLS(config.All().Web.Listener, config.TLS())
  web.Run(config.All().Web.Listener)
}
