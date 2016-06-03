package builder

import (
  ring "github.com/zfjagann/golang-ring"
  "github.com/stianeikeland/go-rpio"
  "bobthebuilder/logging"
  "bobthebuilder/config"
  "bobthebuilder/util"
  "encoding/json"
  "io/ioutil"
  "errors"
  "path"
  "time"
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

  //subscribers to events
  subscribers map[chan BuilderEvent]bool
}



// (Re)load all build definitions.
func (b *Builder)Init()error{
  b.Lock.Lock()
  defer b.Lock.Unlock()

  b.Definitions = []*BuildDefinition{} //clear our definitions
  b.CurrentRun = nil

  b.publishEvent(EVT_DEF_REFRESH, time.Now().Unix(), -1)

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
  b.publishEvent(EVT_RUN_QUEUED, run, index)
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

      run := event.(*Run)
      logging.Info("builder-worker", "Got Run to execute from queue: ", run.Definition.Name)
      go b.ledBeaconFlashLoop(run)
      go b.ledFlasherLoop(run)
      index, _ := b.findDefinitionIndex(run.Definition.Name)
      run.SetupForRun()
      b.publishEvent(EVT_RUN_STARTED, run, index)
      b.CurrentRun = run
      b.Lock.Unlock()

      run.Run(b, index)

      b.Lock.Lock()
      b.CompletedBacklog.Enqueue(run)
      b.publishEvent(EVT_RUN_FINISHED, run, index)
      b.Lock.Unlock()

    } else {
      b.Lock.Unlock()
    }
  }
}


func (b* Builder)ledBeaconFlashLoop(run *Run){
  if config.All().RaspberryPi.Enable && config.All().RaspberryPi.BuildLedPin > 0 {
    for run.IsRunning() {
      time.Sleep(time.Millisecond * 400)
      led := rpio.Pin(config.All().RaspberryPi.BuildLedPin)
      for i := 0; i < 3; i++ {
        led.High()
        time.Sleep(time.Millisecond * 7)
        led.Low()
        time.Sleep(time.Millisecond * 50)
      }
      time.Sleep(time.Millisecond * 350)
    }
  }
}

func (b* Builder)ledFlasherLoop(run *Run){
  if config.All().RaspberryPi.Enable && len(config.All().RaspberryPi.CycleFlashers) > 0 {
    for run.IsRunning() {
      for _, pin := range config.All().RaspberryPi.CycleFlashers {
        p := rpio.Pin(pin)
        p.High()
        time.Sleep(time.Millisecond * 500)
        p.Low()
      }
    }
  }
}

//Returns a list of recent runs.
func (b* Builder)GetHistory()[]interface{}{
  b.Lock.Lock()
  defer b.Lock.Unlock()
  return b.CompletedBacklog.Values()
}

//returns the index of the currently running definition, and the currently running run. Returns -1, nil if nothing is running.
func (b *Builder)GetStatus()(currIndex int, currentRun *Run){
  b.Lock.Lock()
  defer b.Lock.Unlock()

  if !b.IsRunning(){
    return -1, nil
  }

  currentRun = b.CurrentRun
  currIndex, _ = b.findDefinitionIndex(currentRun.Definition.Name)
  return
}


//code for events subscription/publishing system is in builder_events.go

//Creates a new builder object. Run in package init(), so keep this method simple.
func New()*Builder{
  out :=  &Builder{
    CurrentRun: nil,
    EventsToProcess: &ring.Ring{},
    TriggerWorkerChan: make(chan bool, MAX_EVENT_QUEUE_SIZE),
    CompletedBacklog: &ring.Ring{},
    subscribers: map[chan BuilderEvent]bool{},
  }
  out.EventsToProcess.SetCapacity(MAX_EVENT_QUEUE_SIZE)
  out.CompletedBacklog.SetCapacity(MAX_HISTORY_BACKLOG_SIZE)
  go out.builderRunLoop()
  return out
}
