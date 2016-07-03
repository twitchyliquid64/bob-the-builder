package builder
import (
  "bobthebuilder/logging"
  "time"
  "os"
)

type SetEnvPhase struct{
  BasicPhase

  Key string
  Value string
}

func (p * SetEnvPhase)init(index int){
  p.Type = "SET_ENV"
  p.StatusString = PHASE_STATUS_READY
  p.Index = index
}

func (p * SetEnvPhase)String()string{
  return "(phase)" + p.Type + " -- " + p.StatusString + " (" + p.Key + " = " + p.Value + ")"
}

func (p* SetEnvPhase)phaseError(eCode int, statusString string)int{
  p.ErrorCode = eCode
  logging.Error("phase-set-env", statusString)
  p.StatusString = statusString
  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)
  return eCode
}



func (p * SetEnvPhase)Run(r* Run, builder *Builder, defIndex int)int{
  var err error
  p.Start = time.Now()

  //run templates to sub in any variable information like dates etc
  p.Key, err = ExecTemplate(p.Key, p, r, builder)
  if err != nil{
    p.WriteOutput( "Template Error (key): " + err.Error() + "\n", r, builder, defIndex)
    return p.phaseError(-1, "Template error")
  }
  p.Value, err = ExecTemplate(p.Value, p, r, builder)
  if err != nil{
    p.WriteOutput( "Template Error (value): " + err.Error() + "\n", r, builder, defIndex)
    return p.phaseError(-2, "Template error")
  }

  p.WriteOutput( "Setting environment variable '" + p.Key + "' to '" + p.Value + "'", r, builder, defIndex)
  os.Setenv(p.Key, p.Value)

  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)
  p.ErrorCode = 0
  p.StatusString = "Upload successful"
  return 0
}
