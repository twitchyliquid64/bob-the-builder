package web

import (
  "bobthebuilder/config"
  "github.com/hoisie/web"
)

func requestAuth(ctx *web.Context) {
  //requestBasicAuth(ctx)
  ctx.Redirect(302, "/login")
}

func requestBasicAuth(ctx *web.Context) {
	ctx.Header().Set("WWW-Authenticate", "Basic realm=\"Pushtart:"+config.All().Name+"\"")
	ctx.WriteHeader(401)
	ctx.Write([]byte("401 Unauthorized\n"))
}

func needAuthChallenge(ctx *web.Context) bool{
	if gAuth == nil{
		return false
	}

	info, err := gAuth.AuthInfo(ctx)
	if err == nil && info != nil{
		return false
	}
	return true
}

func processLogin(ctx *web.Context) {
  info, err := gAuth.DoLogin(ctx)
  if err != nil || info == nil {
    loginMainPage(ctx)
  }
  ctx.Write([]byte("{\"success\": true}"))
}
