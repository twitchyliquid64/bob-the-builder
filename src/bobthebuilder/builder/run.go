package builder

import (
  "sync"
  "time"
)

//A run represents a build - either planned, in progress, or finished.
type Run struct {
  sync.Mutex
  GUID string
  ExecType string                    //EG: build, run-action

  HasStarted bool
  StartTime time.Time
  HasFinished bool
  EndTime time.Time
  Definition *BuildDefinition

  Phases []phase                    //Phases which are part of the run.
}

func (r *Run)IsRunning()bool{
  r.Lock()
  defer r.Unlock()

  if !r.HasStarted{
    return false
  }
  if r.HasStarted && !r.HasFinished{
    return true
  }
  return false
}


func (r *Run)SetupForRun(){
  r.StartTime = time.Now()
  r.HasStarted = true
  r.HasFinished = false
}

func (r *Run)Run(){
  // ...
  r.EndTime = time.Now()
  r.HasFinished = true
}
