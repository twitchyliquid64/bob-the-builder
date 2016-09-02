package builder

import (
  "bobthebuilder/logging"
  "sync"
  "time"
  "io/ioutil"
  "os"
  "path"
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
    phase.SetStartTime(time.Now())
    builder.publishEvent(EVT_PHASE_STARTED, phase, defIndex)

    shouldSkip := phase.EvaluateShouldSkip(r, builder, defIndex)
    if shouldSkip{
      builder.publishEvent(EVT_PHASE_FINISHED, phase, defIndex)
    }else{
      status := phase.Run(r, builder, defIndex)
      builder.publishEvent(EVT_PHASE_FINISHED, phase, defIndex)
      if status < STATUS_SUCCESS {
        r.Status = STATUS_FAILURE
        break
      }
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
  //logging.Info("run-variables-setdefaults", overrides)

  for _, parameter := range r.Definition.Params {
    if parameter.Varname == ""{
      continue
    }

    if overrides != nil{
      if _, ok := overrides[parameter.Varname]; ok{
        if parameter.Type == "file" {
          r.buildVariables[parameter.Varname] = parameter.Filename
          pwd, _ := os.Getwd()

          //make sure build dir exists
          if exists, _ := exists(path.Join(pwd, BUILD_TEMP_FOLDER_NAME)); !exists {
            os.MkdirAll(path.Join(pwd, BUILD_TEMP_FOLDER_NAME), 0777)
          }

          err := ioutil.WriteFile(path.Join(pwd, BUILD_TEMP_FOLDER_NAME, parameter.Filename), []byte(overrides[parameter.Varname]), 0777)
          if err != nil {
            logging.Error("run-variables-setdefaults", "Could not save file parameter " + parameter.Varname + ": " + err.Error())
          }
        } else {
          r.buildVariables[parameter.Varname] = overrides[parameter.Varname]
        }
        continue
      }
    }

    //if we are at this stage no build parameter was specified. Populate with defaults...
    var ok bool
    switch parameter.Type {
    case "check":
      r.buildVariables[parameter.Varname] = interfaceToStringyBoolean(parameter.Default)
    case "file":
      r.buildVariables[parameter.Varname] = ""
    case "select":
      fallthrough
    case "branchselect":
      fallthrough
    case "text":
      if parameter.Default != nil{
        r.buildVariables[parameter.Varname], ok = parameter.Default.(string)
        if !ok{
          logging.Warning("run-variables-setdefaults", "Could not set default for " + parameter.Varname + ", unexpected default type.")
        }
      }
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
