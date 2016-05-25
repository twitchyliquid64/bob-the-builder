package builder

import (
  ring "github.com/zfjagann/golang-ring"
  "bobthebuilder/logging"
  "bobthebuilder/util"
  "encoding/json"
  "io/ioutil"
  "errors"
  "path"
  "sync"
)

const DEFINITIONS_FOLDER_NAME = "definitions"
const DEFINITIONS_FILE_SUFFIX = ".json"
const BUILD_TEMP_FOLDER_NAME = "build"
const MAX_EVENT_QUEUE_SIZE = 5000
const MAX_HISTORY_BACKLOG_SIZE = 15

var DefNotFoundErr = errors.New("Definition not found")
var BuildRunningErr = errors.New("Build already running")

type Builder struct {
  Lock sync.Mutex //should be used to lock EventsToProcess,Definitions (iteration & modify)
  Definitions []*BuildDefinition

  CurrentRun *Run //represents any currently running build at any instant (but check if it is != nil and IsRunning())
  EventsToProcess *ring.Ring
  TriggerWorkerChan chan bool
  CompletedBacklog *ring.Ring
}



// (Re)load all build definitions.
func (b *Builder)Init()error{
  b.Lock.Lock()
  defer b.Lock.Unlock()

  b.Definitions = []*BuildDefinition{} //clear our definitions
  b.CurrentRun = nil

  defFiles, err := util.GetFilenameListInFolder(DEFINITIONS_FOLDER_NAME, DEFINITIONS_FILE_SUFFIX)
  if err != nil {
    return err
  }

  for _, fpath := range defFiles { //parse each configuration file - TODO: Split the parse-validate-insert into a separate function to allow a ParseDefinition([]byte) interface.
      bData, err := ioutil.ReadFile(fpath)
      if err != nil{
        logging.Error("builder-init", "Error reading build definition for " + path.Base(fpath) + ". Skipping.")
        logging.Error("builder-init", err)
        continue
      }

      var def BuildDefinition
      err = json.Unmarshal(bData, &def)
      if err != nil{
        logging.Error("builder-init", "Error parsing build definition for " + path.Base(fpath) + ". Skipping.")
        logging.Error("builder-init", err)
        continue
      }

      valid := def.Validate()
      if valid{
        b.Definitions = append(b.Definitions, &def)
        b.CurrentRun = nil
        logging.Info("builder-init", def.Name + " - definition ready.")
      }else {
        logging.Warning("builder-init", "Skipping " + path.Base(fpath) + " (" + def.Name + ").")
      }
  }

  return nil
}



//Returns true if a Run (build) is currently executing.
func (b *Builder)IsRunning()bool{
  if b.CurrentRun != nil && b.CurrentRun.IsRunning(){
    return true
  }
  b.CurrentRun = nil
  return false
}

//Enqueues a build based on the build definition with the given name.
//returns DefNotFoundErr if the build definition does not exist.
func (b *Builder)EnqueueBuildEvent(buildDefinitionName string)(*Run, error){
  b.Lock.Lock()
  defer b.Lock.Unlock()
  //if b.IsRunning(){
  //  return nil, BuildRunningErr
  //}

  index, err := b.findDefinitionIndex(buildDefinitionName)
  if err != nil{
    return nil, err
  }
  run := b.Definitions[index].genRun()
  b.EventsToProcess.Enqueue(run)
  b.TriggerWorkerChan <- true

  return run, nil
}

//returns the array index of the build definition with the given name.
//returns DefNotFoundErr if it does not exist.
//caller should hold b.Lock.
func (b* Builder)findDefinitionIndex(definitionName string)(int, error){
  for i, def := range b.Definitions{
    if def.Name == definitionName{
      return i, nil
    }
  }
  return -1, DefNotFoundErr
}

//returns an interface{} that can be serialized to JSON for the dashboard.
func (b* Builder)GetDefinitionsSerialisable()interface{}{
  return b.Definitions
}


//Do not call. This is run once on init.
func (b* Builder)builderRunLoop(){
  for <-b.TriggerWorkerChan{
    b.Lock.Lock()
    event := b.EventsToProcess.Dequeue()
    if event != nil {
      logging.Info("builder-worker", "Got Run to execute from queue: ", event)//TODO: Log name / type - not the whole structure
      run := event.(*Run)
      run.SetupForRun()
      b.CurrentRun = run
      b.Lock.Unlock()
      run.Run()
      b.CompletedBacklog.Enqueue(run)
    } else {
      b.Lock.Unlock()
    }
  }
}


//Creates a new builder object. Run in package init(), so keep this method simple.
func New()*Builder{
  out :=  &Builder{
    CurrentRun: nil,
    EventsToProcess: &ring.Ring{},
    TriggerWorkerChan: make(chan bool, MAX_EVENT_QUEUE_SIZE),
    CompletedBacklog: &ring.Ring{},
  }
  out.EventsToProcess.SetCapacity(MAX_EVENT_QUEUE_SIZE)
  out.CompletedBacklog.SetCapacity(MAX_HISTORY_BACKLOG_SIZE)
  go out.builderRunLoop()
  return out
}
