package web

import (
	"bobthebuilder/builder"
	"bobthebuilder/logging"
	"io/ioutil"
	"strconv"

	"github.com/hoisie/web"
)

// /api/file/definitions
func getDefinitionJSONHandler(ctx *web.Context) {
	did, _ := strconv.Atoi(ctx.Params["did"])
	def := builder.GetInstance().Definitions[did]
	d, _ := ioutil.ReadFile(def.AbsolutePath)
	ctx.ContentType("text/plain")
	ctx.ResponseWriter.Write(d)
}

func saveDefinitionJSONHandler(ctx *web.Context) {
	jsonData, err := ioutil.ReadAll(ctx.Request.Body) //no need to close Body
	if err != nil {
		logging.Error("web-file-api", "saveDefinitionJSONHandler() failed read:", err)
		ctx.Abort(500, "read error")
		return
	}

	did, _ := strconv.Atoi(ctx.Params["did"])
	builder.GetInstance().EnqueueDefinitionUpdateEvent(did, jsonData)
}
