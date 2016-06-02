package builder

import (
  "time"
)



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
  Run(*Run,*Builder,int)int
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
  pOut := struct{
    Phase *BasicPhase `json:"phase"`
    Content string `json:"content"`
  }{
    Phase: p,
    Content: info,
  }
  builder.publishEvent(EVT_PHASE_DATA_UPDATE, pOut, defIndex)
}
