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

func ExecTemplate(templ string, phase interface{})(string,error){
  resultBuf := new(bytes.Buffer)
  t, err := template.New("t").Parse(templ)
  if err != nil{
    return "", err
  }

  tinfo := getBaseTemplateInfoStruct()
  tinfo.Phase = phase
  err = t.Execute(resultBuf, tinfo)
  return resultBuf.String(), err
}
