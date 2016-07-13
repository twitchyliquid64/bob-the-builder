package builder

import (
  "bobthebuilder/logging"
  "sync"
  "time"
)

const STATUS_NOT_YET_RUN = -42
const STATUS_SUCCESS = 0
const STATUS_FAILURE = -1
const PHASE_STATUS_READY = "READY"//textual default for phase.statusString before they have run

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

  Status int `json:"status"`

  Phases []phase `json:"phases"`    //Phases which are part of the run.
  Tags []string `json:"tags"`

  //todo: convert value to own type
  buildVariables map[string]string

  PhysDisabled bool
}

func (r *Run)IsRunning()bool{
  r.Lock()
  defer r.Unlock()

  if r.HasStarted && !r.HasFinished{
    return true
  }
  return false
}


//Called to initialise member fields ready for a run.
func (r *Run)SetupForRun(){
  r.StartTime = time.Now()
  r.HasStarted = true
  r.HasFinished = false
}

func (r *Run)Run(builder *Builder, defIndex int){
  //run event already sent to subscribers by Builder, so all we need to do is literally run

  for _, phase := range r.Phases{
    builder.publishEvent(EVT_PHASE_STARTED, phase, defIndex)
    status := phase.Run(r, builder, defIndex)
    builder.publishEvent(EVT_PHASE_FINISHED, phase, defIndex)
    if status < STATUS_SUCCESS {
      r.Status = STATUS_FAILURE
      break
    }
  }

  if r.Status == STATUS_NOT_YET_RUN {
    r.Status = STATUS_SUCCESS
  }

  // ...
  r.EndTime = time.Now()
  r.HasFinished = true
}


//iterate through given parameters and set them to default values.
func (r *Run)SetDefaultVariables(overrides map[string]string){
  logging.Info("run-variables-setdefaults", overrides)

  for _, parameter := range r.Definition.Params {
    if parameter.Varname == ""{
      continue
    }

    if overrides != nil{
      if _, ok := overrides[parameter.Varname]; ok{
        r.buildVariables[parameter.Varname] = overrides[parameter.Varname]
        continue
      }
    }

    var ok bool
    switch parameter.Type {
    case "check":
      r.buildVariables[parameter.Varname] = interfaceToStringyBoolean(parameter.Default)
    case "text":
      if parameter.Default != nil{
        r.buildVariables[parameter.Varname], ok = parameter.Default.(string)
        if !ok{
          logging.Warning("run-variables-setdefaults", "Could not set default for " + parameter.Varname + ", unexpected type.")
        }
      }
    case "select":
      r.buildVariables[parameter.Varname] = parameter.Default.(string)
    }
  }
}


func interfaceToStringyBoolean(in interface{})string{
  v, ok := in.(bool)
  if ok{
    if v{
      return "true"
    } else{
      return "false"
    }
  }

  s, ok := in.(string)
  if ok{
    if s == "true"{
      return "true"
    }
  }
  return "false"
}
