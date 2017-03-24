package web

import (
	"bobthebuilder/builder"
	"bobthebuilder/logging"
	"bobthebuilder/util"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hoisie/web"
)


func getDefinitionHandler(ctx *web.Context) {
	if needAuthChallenge(ctx){
		requestAuth(ctx)
		return
	}
	out := builder.GetInstance().GetDefinitionsSerialisable()
	b, err := json.Marshal(out)
	if err != nil {
		logging.Error("web-definitions-api", err)
		ctx.ResponseWriter.Write([]byte("{error: '" + err.Error() + "'}"))
	} else {
		//logging.Info("web-definitions-api", string(b))
		ctx.ResponseWriter.Write(b)
	}
}

func getCronHandler(ctx *web.Context) {
	if needAuthChallenge(ctx){
		requestAuth(ctx)
		return
	}
	out := builder.GetInstance().CronEntries()
	b, err := json.Marshal(out)
	if err != nil {
		logging.Error("web-cron-api", err)
		ctx.ResponseWriter.Write([]byte("{error: '" + err.Error() + "'}"))
	} else {
		//logging.Info("web-definitions-api", string(b))
		ctx.ResponseWriter.Write(b)
	}
}

func updateCronHandler(ctx *web.Context) {
	if needAuthChallenge(ctx){
		requestAuth(ctx)
		return
	}
	decoder := json.NewDecoder(ctx.Request.Body)
	var data []builder.CronRecord
	err := decoder.Decode(&data)
	logging.Info("debug", data)
	if err != nil {
		logging.Error("web-cron-api", "updateCronHandler() failed to decode JSON:", err)
		ctx.Abort(500, "JSON error")
		return
	}

	builder.GetInstance().UpdateCron(data)
}


func getHistoryHandler(ctx *web.Context) {
	if needAuthChallenge(ctx){
		requestAuth(ctx)
		return
	}
	out := builder.GetInstance().GetHistory()
	b, err := json.Marshal(out)
	if err != nil {
		logging.Error("web-definitions-api", err)
		ctx.ResponseWriter.Write([]byte("{error: '" + err.Error() + "'}"))
	} else {
		//logging.Info("web-definitions-api", string(b))
		ctx.ResponseWriter.Write(b)
	}
}

func getStatusHandler(ctx *web.Context) {
	index, run := builder.GetInstance().GetStatus()
	out := map[string]interface{}{"index": index, "run": run}

	b, err := json.Marshal(out)
	if err != nil {
		logging.Error("web-definitions-api", err)
		ctx.ResponseWriter.Write([]byte("{error: '" + err.Error() + "'}"))
	} else {
		//logging.Info("web-definitions-api", string(b))
		ctx.ResponseWriter.Write(b)
	}
}

func getBuildParamsLookupHandler(ctx *web.Context) {
	if needAuthChallenge(ctx){
		requestAuth(ctx)
		return
	}
	paramIndex, _ := strconv.Atoi(ctx.Params["param"])

	branches, err := util.GetGitLookupData(builder.GetInstance().GetDefinition(ctx.Params["name"]).Params[paramIndex].Options)
	if err != nil {
		ctx.Abort(500, "{\"success\": true, \"error\": \"Internal Server Error\"}")
		return
	}

	out2 := map[string]interface{}{"success": true, "results": branches}
	b, err := json.Marshal(out2)
	if err != nil {
		logging.Error("web-definitions-api", err)
		ctx.ResponseWriter.Write([]byte("{success: false, error: '" + err.Error() + "'}"))
	} else {
		ctx.ResponseWriter.Write(b)
	}
}

func calcNextVersionNumber(defName string) string {
	re := regexp.MustCompile("\\d+$")
	bVersion := re.ReplaceAllFunc([]byte(builder.GetInstance().GetDefinition(defName).LastVersion), func(match []byte) []byte {
		iMatch, err := strconv.Atoi(string(match))
		if err != nil {
			return match
		} else {
			return []byte(strconv.Itoa(iMatch + 1))
		}
	})
	candidateVersion := string(bVersion)
	if candidateVersion == "" {
		return "0.0.1"
	}
	return candidateVersion
}

func enqueueBuildHandler(ctx *web.Context) {
	if ctx.Params["version"] == "" {
		ctx.Params["version"] = calcNextVersionNumber(ctx.Params["name"])
	}

	tags := []string{"web", "default"}
	if gAuth != nil{
		info, authErr := gAuth.AuthInfo(ctx)
		if authErr == nil {
			tags = append(tags, "auth")
			tags = append(tags, "startedby:" + info.User.Name())
		}
	}

	builder.GetInstance().EnqueueBuildEvent(ctx.Params["name"], tags, ctx.Params["version"])
}

type BuildOptionsDTO struct {
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	Tags           []string          `json:"tags"`
	Params         map[string]string `json:"params"`
	IsPhysDisabled bool              `json:"isPhysDisabled"`
}

func enqueueBuildHandlerWithOptions(ctx *web.Context) {
	decoder := json.NewDecoder(ctx.Request.Body)
	var data BuildOptionsDTO
	err := decoder.Decode(&data)
	if err != nil {
		logging.Error("web-definitions-api", "enqueueBuildHandlerWithOptions() failed to decode JSON:", err)
		ctx.Abort(500, "JSON error")
		return
	}

	if len(data.Tags) == 1 && data.Tags[0] == "" {
		data.Tags = nil
	}

	for i := range data.Tags {
		if data.Tags[i] == "auth" || strings.HasPrefix(data.Tags[i], "startedby:") {
			data.Tags[i] = ""
		}
	}
	if gAuth != nil{
		info, authErr := gAuth.AuthInfo(ctx)
		if authErr == nil {
			data.Tags = append(data.Tags, "auth")
			data.Tags = append(data.Tags, "startedby:" + info.User.Name())
		}
	}

	if data.Version == "" {
		data.Version = calcNextVersionNumber(data.Name)
	}

	builder.GetInstance().EnqueueBuildEventEx(data.Name, data.Tags, data.Version, data.IsPhysDisabled, data.Params)
}

func enqueueReloadHandler(ctx *web.Context) {
	time.Sleep(time.Millisecond * 120)
	builder.GetInstance().EnqueueReloadEvent()
}

func getDefIndexByIdHandler(ctx *web.Context) {
	if needAuthChallenge(ctx){
		requestAuth(ctx)
		return
	}
	def := builder.GetInstance().GetDefinitionByFilename(ctx.Params["fname"])
	ctx.ResponseWriter.Write([]byte(strconv.Itoa(def)))
}
