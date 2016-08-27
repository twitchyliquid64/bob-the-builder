package web

import (
	"bobthebuilder/builder"
	"bobthebuilder/logging"
	"bytes"
	"encoding/json"
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

//TreeviewFileDTO is in a format that angular-treeview understands
type TreeviewFileDTO struct {
	Label     string            `json:"label"`
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Children  []TreeviewFileDTO `json:"children"`
	Size      int64             `json:"size"`
	Collapsed bool              `json:"collapsed"`
	FileMode  os.FileMode       `json:"bits"`
	Media     string            `json:"media"`
}

func getBrowserFilesData(ctx *web.Context) {
	pwd, _ := os.Getwd()
	baseFolder := path.Join(pwd, builder.BASE_FOLDER_NAME)
	baseDTO, err := iterateFolderToTreeviewJSON(baseFolder, pwd)
	if err != nil {
		logging.Error("web-definitions-api", err)
		ctx.Abort(500, "{error: '"+err.Error()+"'}")
		return
	}
	defFolder := path.Join(pwd, builder.DEFINITIONS_FOLDER_NAME)
	defDTO, err := iterateFolderToTreeviewJSON(defFolder, pwd)
	if err != nil {
		logging.Error("web-definitions-api", err)
		ctx.Abort(500, "{error: '"+err.Error()+"'}")
		return
	}
	buildFolder := path.Join(pwd, builder.BUILD_TEMP_FOLDER_NAME)
	buildDTO, err := iterateFolderToTreeviewJSON(buildFolder, pwd)
	if err != nil {
		//this one is not important
	}

	b, err := json.Marshal(map[string]interface{}{
		"base":        baseDTO.Children,
		"definitions": defDTO.Children,
		"build":       buildDTO.Children,
	})
	if err != nil {
		logging.Error("web-definitions-api", err)
		ctx.ResponseWriter.Write([]byte("{error: '" + err.Error() + "'}"))
	} else {
		ctx.ResponseWriter.Write(b)
	}
}

func iterateFolderToTreeviewJSON(absPath string, base string) (out TreeviewFileDTO, err error) {
	var f *os.File
	var stat os.FileInfo

	f, err = os.Open(absPath)
	if err != nil {
		logging.Error("web-file-api", "iterateFolderToTreeviewJSON() error: ", err)
		return
	}
	defer f.Close()

	stat, err = os.Stat(absPath)
	if err != nil {
		logging.Error("web-file-api", "iterateFolderToTreeviewJSON() error: ", err)
		return
	}

	out.Label = stat.Name()
	out.ID = strings.TrimPrefix(absPath, base)
	out.Collapsed = false
	out.Size = stat.Size()
	out.FileMode = stat.Mode()
	out.Media = getMediaType(absPath)

	if stat.IsDir() {
		out.Type = "folder"
		var cFiles []os.FileInfo
		cFiles, err = f.Readdir(-1)
		if err != nil {
			logging.Error("web-file-api", "iterateFolderToTreeviewJSON() error: ", err)
			return
		}
		for _, fi := range cFiles {
			fDTO, err := iterateFolderToTreeviewJSON(path.Join(absPath, fi.Name()), base)
			if err != nil {
				logging.Error("web-file-api", "iterateFolderToTreeviewJSON("+absPath+") error: ", err)
				continue
			}
			out.Children = append(out.Children, fDTO)
		}
	} else {
		out.Type = "file"
	}

	return
}

func getMediaType(path string) string {
	if strings.HasSuffix(path, ".json") {
		return "JSON"
	}
	if strings.HasSuffix(path, ".sh") {
		return "Unix script"
	}
	return "-"
}
