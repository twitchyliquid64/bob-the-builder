package web

import (
  "bobthebuilder/builder"
  "bobthebuilder/logging"
  "github.com/hoisie/web"
  "encoding/json"
)


func getDefinitionHandler(ctx *web.Context) {
  out := builder.GetInstance().GetDefinitionsSerialisable()
  b, err := json.Marshal(out)
  if err != nil{
    logging.Error("web-definitions-api", err)
    ctx.ResponseWriter.Write([]byte("{error: '" + err.Error() + "'}"))
  } else {
    logging.Info("web-definitions-api", string(b))
    ctx.ResponseWriter.Write(b)
  }
}
