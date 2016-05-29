package builder

import (
  "time"
)

const PHASE_STATUS_READY = "READY"

//A phase represents the saved status of a phase.
//A phase is any routine (such as copy files, install an apt package etc) run in the context of a build definition
//A phase's result is stored in the context of a run, which in turn is in the context of a build definition.
type phase interface {
  GetGUID() string                         //Globally unique random string
  GetType() string                         //CLEAN,APT-CHECK,GIT CLONE,BASE INSTALL,EXEC
  GetStatusString() string                 //One sentence summary of status.

  GetErrorCode() int                       //0 == success. Other codes are dependent on type.

  GetStart() time.Time
  GetEnd() time.Time
  GetDuration() time.Duration
  HasFinished() bool
  String()string
}



type BasicPhase struct {
  GUID string `json:"guid"`
  Type string `json:"type"`
  StatusString string `json:"status"`
  ErrorCode int `json:"errorCode"`
  Start time.Time `json:"start"`
  End time.Time `json:"end"`
  Duration time.Duration `json:"duration"`
}

func (p * BasicPhase)GetGUID()string{
  return p.GUID
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
  return p.End == time.Time{}
}
func (p * BasicPhase)String()string{
  return "(phase)" + p.Type + " -- " + p.StatusString
}


type BaseInstallPhase struct{
  BasicPhase
  BaseAbsPath string
}
func (p * BaseInstallPhase)init(){
  p.Type = "BASE-INSTALL"
  p.StatusString = PHASE_STATUS_READY
}
func (p * BaseInstallPhase)String()string{
  return "(phase)" + p.Type + " -- " + p.StatusString + " (" + p.BaseAbsPath + ")"
}

type GitClonePhase struct{
  BasicPhase
  GitSrcPath string
}
func (p * GitClonePhase)init(){
  p.Type = "GIT-CLONE"
  p.StatusString = PHASE_STATUS_READY
}


type CleanPhase struct{
  BasicPhase
  DeletePath string
}
func (p * CleanPhase)init(){
  p.Type = "CLEAN"
  p.StatusString = PHASE_STATUS_READY
}
func (p * CleanPhase)String()string{
  return "(phase)" + p.Type + " -- " + p.StatusString + " (" + p.DeletePath + ")"
}
