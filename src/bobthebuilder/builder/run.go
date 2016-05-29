package builder

import (
  "sync"
  "time"
)

//A run represents a build - either planned, in progress, or finished.
type Run struct {
  sync.Mutex
  GUID string `json:"guid"`
  ExecType string `json:"type"`     //EG: build, run-action
  Version string `json:"version"`

  HasStarted bool `json:"haveStarted"`
  StartTime time.Time `json:"startTime"`
  HasFinished bool `json:"haveFinished"`
  EndTime time.Time `json:"endTime"`
  Definition *BuildDefinition `json:"definition"`

  Phases []phase `json:"phases"`    //Phases which are part of the run.
  Tags []string `json:"tags"`
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
