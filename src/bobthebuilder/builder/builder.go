package builder

import (
	"bobthebuilder/config"
	"bobthebuilder/logging"
	"bobthebuilder/util"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio"
	ring "github.com/zfjagann/golang-ring"
)

const DEFINITIONS_FOLDER_NAME = "definitions"
const DEFINITIONS_FILE_SUFFIX = ".json"
const BUILD_TEMP_FOLDER_NAME = "build"
const MAX_EVENT_QUEUE_SIZE = 5000
const MAX_HISTORY_BACKLOG_SIZE = 15

var DefNotFoundErr = errors.New("Definition not found")
var BuildRunningErr = errors.New("Build already running")

type Builder struct {
	Lock        sync.Mutex //should be used to lock EventsToProcess,Definitions (iteration & modify)
	Definitions []*BuildDefinition

	CurrentRun        *Run //represents any currently running build at any instant (but check if it is != nil and IsRunning())
	EventsToProcess   *ring.Ring
	TriggerWorkerChan chan bool
	CompletedBacklog  *ring.Ring

	//subscribers to events
	subscribers map[chan BuilderEvent]bool
}

func (b *Builder) loadDefinitionFile(fpath string) error {
	bData, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	var def BuildDefinition
	err = json.Unmarshal(bData, &def)
	if err != nil {
		return err
	}
	def.AbsolutePath = fpath

	valid := def.Validate()
	if valid {
		b.Definitions = append(b.Definitions, &def)
		b.CurrentRun = nil
		logging.Info("builder-init", def.Name+" - definition ready.")
	} else {
		return errors.New("Definition failed validation")
	}
	return nil
}

// Init (Re)load all build definitions.
func (b *Builder) Init() error {
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
		err := b.loadDefinitionFile(fpath)
		if err != nil {
			logging.Error("builder-init", "Error reading build definition for "+path.Base(fpath)+". Skipping.")
			logging.Error("builder-init", err)
			continue
		}
	}

	time.Sleep(time.Millisecond * 20)
	b.publishEvent(EVT_DEF_REFRESH_FINISHED, time.Now().Unix(), -1)
	return nil
}

// IsRunning Returns true if a Run (build) is currently executing.
func (b *Builder) IsRunning() bool {
	if b.CurrentRun != nil && b.CurrentRun.IsRunning() {
		return true
	}
	b.CurrentRun = nil
	return false
}

// EnqueueDefinitionUpdateEvent updates the existing definition at defID with the JSON blob contained in jsonData.
func (b *Builder) EnqueueDefinitionUpdateEvent(defID int, jsonData []byte) {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	b.EventsToProcess.Enqueue(DefUpdateEvent{Def: defID, JsonData: jsonData})
	b.publishEvent(EVT_RELOAD_QUEUED, time.Now().Unix(), -1)
	b.TriggerWorkerChan <- true
}

// EnqueueBuildEvent enqueues a build based on the build definition with the given name.
// returns DefNotFoundErr if the build definition does not exist.
func (b *Builder) EnqueueBuildEvent(buildDefinitionName string, tags []string, version string) (*Run, error) {
	return b.EnqueueBuildEventEx(buildDefinitionName, tags, version, false, nil)
}

func (b *Builder) EnqueueBuildEventEx(buildDefinitionName string, tags []string, version string, physDisabled bool, parameterOverrides map[string]string) (*Run, error) {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	index, err := b.findDefinitionIndex(buildDefinitionName)
	if err != nil {
		return nil, err
	}
	run := b.Definitions[index].genRun(tags, version, physDisabled)
	run.SetDefaultVariables(parameterOverrides)
	b.EventsToProcess.Enqueue(run)
	b.publishEvent(EVT_RUN_QUEUED, run, index)
	b.TriggerWorkerChan <- true

	return run, nil
}

//Forces a reload of all definitions. Any further builds in the queue are discarded.
func (b *Builder) EnqueueReloadEvent() {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	b.EventsToProcess.Enqueue("RELOAD")
	b.publishEvent(EVT_RELOAD_QUEUED, time.Now().Unix(), -1)
	b.TriggerWorkerChan <- true
}

func (b *Builder) GetDefinition(definitionName string) *BuildDefinition {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	i, _ := b.findDefinitionIndex(definitionName)
	if i >= 0 {
		return b.Definitions[i]
	}
	return nil
}

//returns the array index of the build definition with the given name.
//returns DefNotFoundErr if it does not exist.
//caller should hold b.Lock.
func (b *Builder) findDefinitionIndex(definitionName string) (int, error) {
	for i, def := range b.Definitions {
		if def.Name == definitionName {
			return i, nil
		}
	}
	return -1, DefNotFoundErr
}

//returns an interface{} that can be serialized to JSON for the dashboard.
func (b *Builder) GetDefinitionsSerialisable() interface{} {
	return b.Definitions
}

func (b *Builder) processOutOfBandEvent(event string) {
	switch event {
	case "RELOAD":
		logging.Info("builder-worker", "Now refreshing definitions")
		time.Sleep(time.Millisecond * 150)
		for b.EventsToProcess.Dequeue() != nil {
		} //delete all of the existing items in the queue
		b.Lock.Unlock()
		b.Init() //reinit
	}
}

func (b *Builder) processDefinitionUpdateEvent(event DefUpdateEvent) {
	logging.Info("builder-worker", "Now updating definition")
	time.Sleep(time.Millisecond * 150)
	for b.EventsToProcess.Dequeue() != nil {
	} //delete all of the existing items in the queue
	b.updateDefinition(event.Def, event.JsonData)
	b.Lock.Unlock()
	b.Init() //reinit
}

func (b *Builder) executeRun(run *Run) {
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
	run.Definition.LastVersion = run.Version
	run.Definition.LastRunTime = int64(run.EndTime.Sub(run.StartTime).Seconds() * 1000)
	run.Definition.Flush()
	b.publishEvent(EVT_RUN_FINISHED, run, index)
}

//Do not call. This is run once on init.
func (b *Builder) builderRunLoop() {
	for <-b.TriggerWorkerChan {
		b.Lock.Lock()
		event := b.EventsToProcess.Dequeue()
		if event != nil {

			_, isOutOfBandEvent := event.(string) //TODO: refactor out of this method
			if isOutOfBandEvent {
				b.processOutOfBandEvent(event.(string))
				continue
			}

			_, isDefUpdate := event.(DefUpdateEvent)
			if isDefUpdate {
				b.processDefinitionUpdateEvent(event.(DefUpdateEvent))
				continue
			}

			//only remaining option - event is a run
			run := event.(*Run)
			b.executeRun(run)
		}
		b.Lock.Unlock()
	}
}

func (b *Builder) updateDefinition(defID int, jsonData []byte) {
	def := b.Definitions[defID]
	ioutil.WriteFile(def.AbsolutePath, jsonData, 777)
}

func (b *Builder) ledBeaconFlashLoop(run *Run) {
	if (!run.PhysDisabled) && config.All().RaspberryPi.Enable && config.All().RaspberryPi.BuildLedPin > 0 {
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

func (b *Builder) ledFlasherLoop(run *Run) {
	if (!run.PhysDisabled) && config.All().RaspberryPi.Enable && len(config.All().RaspberryPi.CycleFlashers) > 0 {
		for run.IsRunning() {
			for _, pin := range config.All().RaspberryPi.CycleFlashers {
				p := rpio.Pin(pin)
				p.High()
				time.Sleep(time.Millisecond * 300)
				p.Low()
			}
		}
	}
}

// GetHistory returns a list of recent runs.
func (b *Builder) GetHistory() []interface{} {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	return b.CompletedBacklog.Values()
}

// GetStatus returns the index of the currently running definition, and the currently running run. Returns -1, nil if nothing is running.
func (b *Builder) GetStatus() (currIndex int, currentRun *Run) {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	if !b.IsRunning() {
		return -1, nil
	}

	currentRun = b.CurrentRun
	currIndex, _ = b.findDefinitionIndex(currentRun.Definition.Name)
	return
}

//code for events subscription/publishing system is in builder_events.go

// New creates a new builder object. Run in package init(), so keep this method simple.
func New() *Builder {
	out := &Builder{
		CurrentRun:        nil,
		EventsToProcess:   &ring.Ring{},
		TriggerWorkerChan: make(chan bool, MAX_EVENT_QUEUE_SIZE),
		CompletedBacklog:  &ring.Ring{},
		subscribers:       map[chan BuilderEvent]bool{},
	}
	out.EventsToProcess.SetCapacity(MAX_EVENT_QUEUE_SIZE)
	out.CompletedBacklog.SetCapacity(MAX_HISTORY_BACKLOG_SIZE)
	go out.builderRunLoop()
	return out
}
