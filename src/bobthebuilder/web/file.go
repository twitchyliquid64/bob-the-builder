package web

import (
	"bobthebuilder/builder"
	"bobthebuilder/logging"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/hoisie/web"
)

// /api/file/definitions
func getDefinitionJSONHandler(ctx *web.Context) {
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

	did, _ := strconv.Atoi(ctx.Params["did"])
	def := builder.GetInstance().Definitions[did]
	d, _ := ioutil.ReadFile(def.AbsolutePath)
	ctx.ContentType("text/plain")
	ctx.ResponseWriter.Write(d)
}

func saveDefinitionJSONHandler(ctx *web.Context) {
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

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
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

	relPathUnsafe := strings.TrimPrefix(ctx.Params["path"], "/base/")
	safe, absPath := sanitizePath(builder.BaseDir, relPathUnsafe)
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

func downloadWorkspaceFileHandler(ctx *web.Context) {
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

	relPathUnsafe := strings.TrimPrefix(ctx.Params["path"], "/build/")
	safe, absPath := sanitizePath(builder.BuildDir, relPathUnsafe)
	if safe {
		fi, err := os.Open(absPath)
		if err != nil {
			logging.Error("web-file-api", "downloadWorkspaceFileHandler() read error: "+err.Error())
			ctx.Abort(403, string(fileError(err.Error())))
			return
		}
		defer fi.Close()

		ctx.SetHeader("Content-Disposition", "attachment; filename=\""+path.Base(absPath)+"\"", true)
		io.Copy(ctx.ResponseWriter, fi)
	} else {
		//attempted LFI attack - return error
		logging.Error("web-file-api", "downloadWorkspaceFileHandler() rejected request for: "+relPathUnsafe)
		ctx.Abort(403, "only base files are accessible")
		return
	}
}

func saveBaseFileHandler(ctx *web.Context) {
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

	relPathUnsafe := strings.TrimPrefix(ctx.Params["path"], "/base/")
	safe, absPath := sanitizePath(builder.BaseDir, relPathUnsafe)
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
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

	baseDTO, err := iterateFolderToTreeviewJSON(builder.BaseDir, path.Dir(builder.BaseDir))
	if err != nil {
		logging.Error("web-definitions-api", err)
		ctx.Abort(500, string(fileError(err.Error())))
		return
	}
	defDTO, err := iterateFolderToTreeviewJSON(builder.DefinitionsDir, path.Dir(builder.DefinitionsDir))
	if err != nil {
		logging.Error("web-definitions-api", err)
		ctx.Abort(500, string(fileError(err.Error())))
		return
	}
	buildDTO, err := iterateFolderToTreeviewJSON(builder.BuildDir, path.Dir(builder.BuildDir))
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
		ctx.ResponseWriter.Write(fileError(err.Error()))
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
		if stat.Name() == ".git" {
			return TreeviewFileDTO{}, errors.New("Ignore")
		}
		var cFiles []os.FileInfo
		cFiles, err = f.Readdir(-1)
		if err != nil {
			logging.Error("web-file-api", "iterateFolderToTreeviewJSON() error: ", err)
			return
		}
		for _, fi := range cFiles {
			fDTO, err := iterateFolderToTreeviewJSON(path.Join(absPath, fi.Name()), base)
			if err != nil {
				if err.Error() != "Ignore" {
					logging.Error("web-file-api", "iterateFolderToTreeviewJSON("+absPath+") error: ", err)
				}
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
	path = strings.ToLower(path)
	if strings.HasSuffix(path, ".json") {
		return "JSON"
	}
	if strings.HasSuffix(path, ".sh") {
		return "Unix script"
	}
	if strings.HasSuffix(path, ".py") {
		return "Python script"
	}
	if strings.HasSuffix(path, ".c") {
		return "C source"
	}
	if strings.HasSuffix(path, ".go") {
		return "Golang source"
	}
	if strings.HasSuffix(path, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
		return "image/jpg"
	}
	return "-"
}

func newFolderHandler(ctx *web.Context) {
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

	relPathUnsafe := strings.TrimPrefix(ctx.Params["path"], "/base/")
	safe, absPath := sanitizePath(builder.BaseDir, relPathUnsafe)
	if safe {
		err := os.Mkdir(absPath, 0770)
		if err != nil {
			ctx.Abort(200, string(fileError(err.Error())))
			return
		}
		ctx.Abort(200, "{\"success\": true}")
		logging.Info("web-file-api", "Created new folder: "+absPath)
	}
}

func newFileHandler(ctx *web.Context) {
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

	relPathUnsafe := strings.TrimPrefix(ctx.Params["path"], "/base/")
	safe, absPath := sanitizePath(builder.BaseDir, relPathUnsafe)
	if safe {
		f, err := os.OpenFile(absPath, os.O_EXCL|os.O_CREATE, 0770)
		if err != nil {
			ctx.Abort(200, string(fileError(err.Error())))
			return
		}
		f.Close()
		ctx.Abort(200, "{\"success\": true}")
		logging.Info("web-file-api", "Created new file: "+absPath)
	}
}

func newDefFileHandler(ctx *web.Context) {
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

	relPathUnsafe := strings.TrimPrefix(ctx.Params["path"], "/definitions/")
	if !strings.HasSuffix(relPathUnsafe, builder.DEFINITIONS_FILE_SUFFIX) {
		relPathUnsafe = relPathUnsafe + builder.DEFINITIONS_FILE_SUFFIX
	}

	safe, absPath := sanitizePath(builder.DefinitionsDir, relPathUnsafe)
	if safe {
		f, err := os.OpenFile(absPath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0770)
		if err != nil {
			ctx.Abort(200, string(fileError(err.Error())))
			return
		}
		_, err = f.Write([]byte(`
			{
			  "name": "New Definition",
			  "icon": "rocket",
			  "steps": []
			}
			`))
		if err != nil {
			ctx.Abort(200, string(fileError(err.Error())))
			return
		}
		err = f.Close()
		if err != nil {
			ctx.Abort(200, string(fileError(err.Error())))
			return
		}
		ctx.Abort(200, "{\"success\": true}")
		logging.Info("web-file-api", "Created new definition file: "+absPath)
		builder.GetInstance().EnqueueReloadEvent()
	}
}

func fileError(error string) []byte {
	out := map[string]interface{}{
		"error":   error,
		"success": false,
	}
	b, _ := json.Marshal(&out)
	return b
}

func deleteHandler(ctx *web.Context) {
	if needAuthChallenge(ctx) {
		requestAuth(ctx)
		return
	}

	relPathUnsafe := strings.TrimPrefix(ctx.Params["path"], "/base/")
	safe, absPath := sanitizePath(builder.BaseDir, relPathUnsafe)
	if safe {
		err := os.RemoveAll(absPath)
		if err != nil {
			ctx.Abort(200, string(fileError(err.Error())))
			return
		}
		ctx.Abort(200, "{\"success\": true}")
		logging.Info("web-file-api", "Delete: "+absPath)
	}
}
