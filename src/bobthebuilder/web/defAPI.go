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
    //logging.Info("web-definitions-api", string(b))
    ctx.ResponseWriter.Write(b)
  }
}

func getHistoryHandler(ctx *web.Context){
  out := builder.GetInstance().GetHistory()
  b, err := json.Marshal(out)
  if err != nil{
    logging.Error("web-definitions-api", err)
    ctx.ResponseWriter.Write([]byte("{error: '" + err.Error() + "'}"))
  } else {
    //logging.Info("web-definitions-api", string(b))
    ctx.ResponseWriter.Write(b)
  }
}

func getStatusHandler(ctx *web.Context){
  index, run := builder.GetInstance().GetStatus()
  out := map[string]interface{}{"index": index, "run": run}

  b, err := json.Marshal(out)
  if err != nil{
    logging.Error("web-definitions-api", err)
    ctx.ResponseWriter.Write([]byte("{error: '" + err.Error() + "'}"))
  } else {
    //logging.Info("web-definitions-api", string(b))
    ctx.ResponseWriter.Write(b)
  }
}

func enqueueBuildHandler(ctx *web.Context){
  builder.GetInstance().EnqueueBuildEvent(ctx.Params["name"])
}
