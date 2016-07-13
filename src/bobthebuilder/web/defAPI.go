package web

import (
  "bobthebuilder/builder"
  "bobthebuilder/logging"
  "github.com/hoisie/web"
  "encoding/json"
  "strconv"
  "regexp"
  "time"
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


func calcNextVersionNumber(defName string)string{
  re := regexp.MustCompile("\\d+$")
  bVersion := re.ReplaceAllFunc([]byte(builder.GetInstance().GetDefinition(defName).LastVersion), func(match []byte)[]byte{
    iMatch, err := strconv.Atoi(string(match))
    if err != nil{
      return match
    } else {
      return []byte(strconv.Itoa(iMatch+1))
    }
  })
  candidateVersion := string(bVersion)
  if candidateVersion == "" {
    return "0.0.1"
  }
  return candidateVersion
}

func enqueueBuildHandler(ctx *web.Context){
  if ctx.Params["version"] == ""{
    ctx.Params["version"] = calcNextVersionNumber(ctx.Params["name"])
  }

  builder.GetInstance().EnqueueBuildEvent(ctx.Params["name"], []string{"web", "default"}, ctx.Params["version"])
}


type BuildOptionsDTO struct {
  Name string `json:"name"`
  Version string `json:"version"`
  Tags []string `json:"tags"`
  IsPhysDisabled bool `json:"isPhysDisabled"`
}

func enqueueBuildHandlerWithOptions(ctx *web.Context){
  decoder := json.NewDecoder(ctx.Request.Body)
  var data BuildOptionsDTO
  err := decoder.Decode(&data)
  if err != nil {
      logging.Error("web-definitions-api", "enqueueBuildHandlerWithOptions() failed to decode JSON:", err)
      ctx.Abort(500, "JSON error")
      return
  }

  if (len(data.Tags) == 1 && data.Tags[0] == ""){
    data.Tags = nil
  }

  builder.GetInstance().EnqueueBuildEventEx(data.Name, data.Tags, data.Version, data.IsPhysDisabled)
}


func enqueueReloadHandler(ctx *web.Context){
  time.Sleep(time.Millisecond * 350)
  builder.GetInstance().EnqueueReloadEvent()
}
