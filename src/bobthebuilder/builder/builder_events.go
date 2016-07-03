package builder

import (
  "github.com/stianeikeland/go-rpio"
  "bobthebuilder/config"
  "time"
)

type BuilderEvent struct {
  Type string
  Data interface{}
  Index int
}

const EVT_DEF_REFRESH = "DEF-REFRESH"
const EVT_DEF_REFRESH_FINISHED = "DEF-REFRESH-COMPLETED"
const EVT_RUN_QUEUED = "RUN-QUEUED"
const EVT_RELOAD_QUEUED = "RELOAD-QUEUED"
const EVT_RUN_STARTED = "RUN-STARTED"
const EVT_RUN_FINISHED = "RUN-FINISHED"
const EVT_PHASE_STARTED = "PHASE-STARTED"
const EVT_PHASE_FINISHED = "PHASE-FINISHED"
const EVT_PHASE_DATA_UPDATE = "PHASE-DATA"




func (b *Builder)SubscribeToEvents(in chan BuilderEvent){
  b.Lock.Lock()
  defer b.Lock.Unlock()
  b.subscribers[in] = true
}

func (b *Builder)UnsubscribeFromEvents(in chan BuilderEvent){
  b.Lock.Lock()
  defer b.Lock.Unlock()
  delete(b.subscribers, in)
}

//assumes caller holds the lock.
func (b *Builder)publishEvent(t string, d interface{}, index int){
  if config.All().RaspberryPi.Enable && config.All().RaspberryPi.DataLedPin > 0 {
    led := rpio.Pin(config.All().RaspberryPi.DataLedPin)
    led.High()
    time.Sleep(time.Millisecond * 4)
    led.Low()
  }


  pkt := BuilderEvent{
    Type: t,
    Data: d,
    Index: index,
  }

  for ch, _ := range b.subscribers {
    select { //prevents blocking if a channel is full
      case ch <- pkt:
      default:
    }
  }
}
