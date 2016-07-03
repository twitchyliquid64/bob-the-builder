package builder

import (
  "text/template"
  "time"
  "bytes"
)

type TemplateInformation struct {
  Day int
  Month int
  Year int
  Minute int
  Hour int


  Phase interface{}
  Builder *Builder
  Run *Run
}

func getBaseTemplateInfoStruct()TemplateInformation{
  return TemplateInformation{
    Day: time.Now().Day(),
    Month: int(time.Now().Month()),
    Year: time.Now().Year(),
    Minute: time.Now().Minute(),
    Hour: time.Now().Hour(),
  }
}

func ExecTemplate(templ string, phase interface{}, r* Run, builder *Builder)(string,error){
  resultBuf := new(bytes.Buffer)
  t, err := template.New("t").Parse(templ)
  if err != nil{
    return "", err
  }

  tinfo := getBaseTemplateInfoStruct()
  tinfo.Phase = phase
  tinfo.Builder = builder
  tinfo.Run = r
  err = t.Execute(resultBuf, tinfo)
  return resultBuf.String(), err
}
