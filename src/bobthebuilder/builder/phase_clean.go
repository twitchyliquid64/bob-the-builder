package builder

import (
  "bobthebuilder/logging"
  "time"
  "os"
)


type CleanPhase struct{
  BasicPhase
  DeletePath string
}
func (p * CleanPhase)init(index int){
  p.Type = "CLEAN"
  p.StatusString = PHASE_STATUS_READY
  p.Index = index
}
func (p * CleanPhase)String()string{
  return "(phase)" + p.Type + " -- " + p.StatusString + " (" + p.DeletePath + ")"
}



func (p * CleanPhase)Run(r* Run, builder *Builder, defIndex int)int{
  p.Start = time.Now()

  err := os.RemoveAll(p.DeletePath)
  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)

  if err == nil{
    p.ErrorCode = 0
    p.StatusString = "Clean successful"
    return 0
  } else {
    p.ErrorCode = -1
    logging.Error("phase-clean", err)
    p.StatusString = err.Error()
    return -1
  }
}
