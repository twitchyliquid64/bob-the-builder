package builder

import (
	"bytes"
	"html"
	"strings"
	"text/template"
	"time"
)

type TemplateInformation struct {
	Day    int
	Month  int
	Year   int
	Minute int
	Hour   int

	Phase   interface{}
	Builder *Builder
	Run     *Run
}

func hasTag(tInfo TemplateInformation, wantedTag string) bool {
	for _, tag := range tInfo.Run.Tags {
		if tag == wantedTag {
			return true
		}
	}
	return false
}

func allTags(tInfo TemplateInformation) string {
	out := ""
	for _, tag := range tInfo.Run.Tags {
		out += " " + tag
	}
	return out
}

func getParameter(tInfo TemplateInformation, wantedVar string) string {
	v, _ := tInfo.Run.buildVariables[wantedVar]
	return v
}

func htmlEscape(in string) string {
	return html.EscapeString(in)
}

func getOutput(tInfo TemplateInformation, id string) string {
	for _, p := range tInfo.Run.Phases {
		if p.GetID() == id {
			return strings.Join(p.GetOutputs(), "\n")
		}
	}
	return ""
}

func getBaseTemplateInfoStruct() TemplateInformation {
	return TemplateInformation{
		Day:    time.Now().Day(),
		Month:  int(time.Now().Month()),
		Year:   time.Now().Year(),
		Minute: time.Now().Minute(),
		Hour:   time.Now().Hour(),
	}
}

func ExecTemplate(templ string, phase interface{}, r *Run, builder *Builder) (string, error) {
	tinfo := getBaseTemplateInfoStruct()
	tinfo.Phase = phase
	tinfo.Builder = builder
	tinfo.Run = r

	funcMap := template.FuncMap{
		"hasTag": func(tname string) bool {
			return hasTag(tinfo, tname)
		},
		"getParameter": func(pname string) string {
			return getParameter(tinfo, pname)
		},
		"allTags": func() string {
			return allTags(tinfo)
		},
		"getOutput": func(id string) string {
			return getOutput(tinfo, id)
		},
		"htmlEscape":   htmlEscape,
		"strContains":  strings.Contains,
		"strHasPrefix": strings.HasPrefix,
		"strSplit":     strings.Split,
		"strHasSuffix": strings.HasSuffix,
		"strReplace":   strings.Replace,
	}

	resultBuf := new(bytes.Buffer)
	t, err := template.New("t").Funcs(funcMap).Parse(templ)
	if err != nil {
		return "", err
	}

	err = t.Execute(resultBuf, tinfo)
	return resultBuf.String(), err
}
