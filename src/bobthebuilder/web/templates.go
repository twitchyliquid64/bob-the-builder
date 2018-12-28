package web

import (
	"bobthebuilder/logging"
	"github.com/hoisie/web"
	"html/template"
	"io/ioutil"
	"path"
)

var TEMPLATE_FOLDER = "templates"
var TEMPLATE_LEFT_DELIMITER = "{!{"
var TEMPLATE_RIGHT_DELIMITER = "}!}"

var templates *template.Template

type templateRecord struct {
	name string
	file string
}

var templateRecords []templateRecord

func init() {
	templates = template.New("__unused__")
}

func reloadTemplatesHandler(ctx *web.Context) {
	templateReInit()
}

// destroys memory structures for already loaded templates, re-parsing them
// such that any changes to them can now be seen.
func templateReInit() {
	logging.Info("web", "Now reloading all templates.")
	templates = template.New("__unused__")
	for _, tempFile := range templateRecords {
		logging.Info("web", "Loading template: ", tempFile.name)
		if err := newTemplateFromFile(tempFile.file, tempFile.name); err != nil {
			logging.Error("web", "Template error: ", err)
		}
	}
}

// Registers a template with a given filename and template name into the system.
// Is immediately available for use upon returning.
func registerTemplate(fname, templateName string) error {
	fname = path.Join(TEMPLATE_FOLDER, fname)
	templateRecords = append(templateRecords, templateRecord{name: templateName, file: fname})

	return newTemplateFromFile(fname, templateName)
}

// Helper function to parse templates from a file.
func newTemplateFromFile(fname, templateName string) error {
	templ := templates.New(templateName)
	templ.Delims(TEMPLATE_LEFT_DELIMITER, TEMPLATE_RIGHT_DELIMITER)

	fdata, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	_, err = templ.Parse(string(fdata))
	if err != nil {
		return err
	}

	return nil
}
