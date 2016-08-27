package web

import (
	"bobthebuilder/builder"
	"bobthebuilder/logging"
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

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

func sanitizePath(base, inPath string) (safe bool, absPath string) {
	absPathUnsafe := path.Clean(path.Join(base, inPath))
	if strings.HasPrefix(absPathUnsafe, base) {
		safe = true
		absPath = absPathUnsafe
	} else {
		safe = false
		absPath = ""
	}
	return
}

func getBaseFileHandler(ctx *web.Context) {
	relPathUnsafe := ctx.Params["path"]

	pwd, _ := os.Getwd()
	baseFolder := path.Join(pwd, builder.BASE_FOLDER_NAME)

	safe, absPath := sanitizePath(baseFolder, relPathUnsafe)
	if safe {
		d, err := ioutil.ReadFile(absPath)
		if err != nil {
			logging.Error("web-file-api", "getBaseFileHandler() read error: ", err)
			ctx.Abort(500, "read error")
		} else {
			ctx.ContentType("text/plain")
			ctx.ResponseWriter.Write(d)
		}
	} else {
		//attempted LFI attack - return error
		logging.Error("web-file-api", "getBaseFileHandler() rejected request for: "+relPathUnsafe)
		ctx.Abort(403, "only base files are accessible")
		return
	}
}

func saveBaseFileHandler(ctx *web.Context) {
	relPathUnsafe := ctx.Params["path"]

	pwd, _ := os.Getwd()
	baseFolder := path.Join(pwd, builder.BASE_FOLDER_NAME)

	safe, absPath := sanitizePath(baseFolder, relPathUnsafe)
	if safe {

		data := bytes.Buffer{}
		data.ReadFrom(ctx.Request.Body)

		err := ioutil.WriteFile(absPath, data.Bytes(), 0777)
		if err != nil {
			logging.Error("web-file-api", "saveBaseFileHandler() write error: ", err)
			ctx.Abort(500, "write error")
		}
	} else {
		//attempted LFI attack - return error
		logging.Error("web-file-api", "saveBaseFileHandler() rejected request for: "+relPathUnsafe)
		ctx.Abort(403, "only base files are accessible")
		return
	}
}
