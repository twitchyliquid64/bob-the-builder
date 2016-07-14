package builder

import (
  "time"
  "strings"
)

var CARRIAGE_RETURN_CONTROL_SEQUENCE = "CONTROL<CHAR-RETURN>"

//A phase represents the saved status of a phase.
//A phase is any routine (such as copy files, install an apt package etc) run in the context of a build definition
//A phase's result is stored in the context of a run, which in turn is in the context of a build definition.
type phase interface {
  GetType() string                         //CLEAN,APT-CHECK,GIT CLONE,BASE INSTALL,EXEC
  GetStatusString() string                 //One sentence summary of status.

  GetErrorCode() int                       //0 == success. Other codes are dependent on type.

  GetStart() time.Time
  GetEnd() time.Time
  GetDuration() time.Duration
  HasFinished() bool
  String()string
  SetConditional(string)
  Run(*Run,*Builder,int)int
  EvaluateShouldSkip(*Run,*Builder,int)bool
  ShouldSkip(*Run,*Builder,int)bool
}



type BasicPhase struct {
  Type string `json:"type"`
  StatusString string `json:"status"`
  ErrorCode int `json:"errorCode"`
  Start time.Time `json:"start"`
  End time.Time `json:"end"`
  Duration time.Duration `json:"duration"`
  Index int `json:"index"`
  Outputs []string `json:"-"`
  Conditional string `json:"-"`
}

func (p * BasicPhase)SetConditional(c string){
  p.Conditional = c
}
func (p * BasicPhase)GetType()string{
  return p.Type
}
func (p * BasicPhase)GetStatusString()string{
  return p.StatusString
}
func (p * BasicPhase)GetErrorCode()int{
  return p.ErrorCode
}
func (p * BasicPhase)GetStart()time.Time{
  return p.Start
}
func (p * BasicPhase)GetEnd()time.Time{
  return p.End
}
func (p * BasicPhase)GetDuration()time.Duration{
  return p.Duration
}
func (p * BasicPhase)HasFinished()bool{
  return !(p.End == time.Time{})
}
func (p * BasicPhase)String()string{
  return "(phase)" + p.Type + " -- " + p.StatusString
}
func (p * BasicPhase)Run(r* Run, builder *Builder, defIndex int)int{
  p.Start = time.Now()
  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)
  p.ErrorCode = -949
  p.StatusString = "Phase not implemented"
  return -949
}
func (p * BasicPhase)WriteOutput(info string, r* Run, builder *Builder, defIndex int){
  p.Outputs = append(p.Outputs, info)

  //to survive JSON parsing - so the frontend knows how to transform the display
  if strings.Contains(info, "\r") {
    info = strings.Replace(info, "\r", CARRIAGE_RETURN_CONTROL_SEQUENCE, -1)
  }

  pOut := struct{
    Phase *BasicPhase `json:"phase"`
    Content string `json:"content"`
  }{
    Phase: p,
    Content: info,
  }
  builder.publishEvent(EVT_PHASE_DATA_UPDATE, pOut, defIndex)
}

func (p * BasicPhase)ShouldSkip(r* Run, builder *Builder, defIndex int)bool{
  if len(p.Conditional) == 0{
    return false
  }

  o, err := ExecTemplate(p.Conditional, p, r, builder)
  if err != nil{
    p.WriteOutput( "Template Error (step conditional): " + err.Error() + "\n", r, builder, defIndex)
    return true
  }
  p.WriteOutput( "Skip conditional: " + o + "\n", r, builder, defIndex)

  if o == "false"{
    return false
  }

  if len(o) > 0 {
    return true
  }

  return false
}


func (p * BasicPhase)EvaluateShouldSkip(r* Run, builder *Builder, defIndex int)bool{
  skip := p.ShouldSkip(r, builder, defIndex)
  if skip{
    p.Start = time.Now()
    p.End = time.Now()
    p.Duration = p.End.Sub(p.Start)
    p.ErrorCode = 954321
    p.StatusString = "Phase skipped"
  }
  return skip
}
