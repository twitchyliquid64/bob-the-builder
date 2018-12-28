package builder

import (
	"bobthebuilder/config"
	"bobthebuilder/logging"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type BuilderEvent struct {
	Type  string
	Data  interface{}
	Index int
}

type DefUpdateEvent struct {
	Def      int
	JsonData []byte
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
const EVT_SERVER_STATS = "SERVER-STATS"

func (b *Builder) SubscribeToEvents(in chan BuilderEvent) {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	b.subscribers[in] = true
}

func (b *Builder) UnsubscribeFromEvents(in chan BuilderEvent) {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	delete(b.subscribers, in)
	close(in)
}

//assumes caller holds the lock.
func (b *Builder) publishEvent(t string, d interface{}, index int) {
	if config.All().RaspberryPi.Enable && config.All().RaspberryPi.DataLedPin > 0 {
		led := rpio.Pin(config.All().RaspberryPi.DataLedPin)
		led.High()
		time.Sleep(time.Millisecond * 4)
		led.Low()
	}

	pkt := BuilderEvent{
		Type:  t,
		Data:  d,
		Index: index,
	}

	for ch, _ := range b.subscribers {
		select { //prevents blocking if a channel is full
		case ch <- pkt:
		default:
		}
	}

	if config.All().Events.Enable && config.All().Events.EventTopic != "" {
		if t == EVT_PHASE_DATA_UPDATE && !config.All().Events.PublishDataEvents {
			return
		}
		select {
		case b.pubsubMsgs <- pkt:
		default:
		}
	}
}

func handlePubsubQueue(c chan BuilderEvent) {
	logging.Info("builder-pubsub", "Starting pubsub transmission routine.")
	defer logging.Info("builder-pubsub", "Transmission routine stopping.")

	uri, err := url.Parse("https://pubsub.googleapis.com/v1/projects/" + config.All().Events.Project + "/topics/" + config.All().Events.EventTopic)
	if err != nil {
		logging.Error("builder-pubsub", "Error constructing pubsub url: ", err)
		return
	}

	isClosed := false
	for !isClosed {
		time.Sleep(time.Second)
		// each run, iterate through the whole channel and gather all the objects
		// if its closed, set a flag so we exit at the end of the iteration
		var eventsToSend []BuilderEvent
		for {
			select {
			case evt, ok := <-c:
				if !ok {
					isClosed = true
					goto allEvents
				}
				eventsToSend = append(eventsToSend, evt)
			default:
				goto allEvents
			}
		}
	allEvents:
		if len(eventsToSend) == 0 {
			continue
		}

		// build up the list of structures necessary for the pubsub API
		var msgs []map[string]interface{}
		for _, evt := range eventsToSend {
			eventData, err := json.Marshal(evt)
			if err != nil {
				logging.Error("builder-pubsub", "Error constructing event data: ", err)
				continue
			}
			msgs = append(msgs, map[string]interface{}{
				"data": eventData,
			})
		}

		// marshal and send
		data, err := json.Marshal(map[string]interface{}{
			"messages": msgs,
		})
		if err != nil {
			logging.Error("builder-pubsub", "Marshal error: ", err)
			return
		}

		resp, err := config.Pubsub().Post(uri.String()+":publish", "application/json", bytes.NewBuffer(data))
		if err != nil {
			logging.Error("builder-pubsub", "Pubsub POST error: ", err)
			resp.Body.Close()
			return
		}
		if resp.StatusCode != 200 {
			e, _ := ioutil.ReadAll(resp.Body)
			logging.Error("builder-pubsub", "Pubsub publish error: ", err, string(e))
			resp.Body.Close()
			return
		}
		resp.Body.Close()
	}
}
