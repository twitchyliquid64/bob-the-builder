package config

import (
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"github.com/stianeikeland/go-rpio"
	"bobthebuilder/logging"
	"golang.org/x/oauth2"
	"net/http"
	"crypto/tls"
)

var gConfig *Config = nil
var gTls *tls.Config = nil

var gEventCred *jwt.Config = nil
var gEventClient *http.Client = nil

func Load(fpath string)error{
	conf, err := readConfig(fpath)
	if err == nil{
		gConfig = conf
	} else {
		logging.Error("config", "config.Load() error:", err)
		return err
	}

	if gConfig.TLS.PrivateKey == ""{
		logging.Warning("config", "TLS keyfile paths omitted, skipping TLS setup")
	} else{
		tls, err := loadTLS(gConfig.TLS.PrivateKey, gConfig.TLS.Cert)
		if err == nil{
			gTls = tls
		} else {
			logging.Error("config", "config.Load() tls error:", err)
			return err
		}
	}

	if gConfig.RaspberryPi.Enable {
		err = rpio.Open()
		if err != nil {
			logging.Error("config", "Failed setup of RPi GPIO: ", err)
			return err
		}

		initRpiGPIO()
	}

	if gConfig.Events.Enable {
		b, err := ioutil.ReadFile(gConfig.Events.CredentialPath)
		if err != nil {
			logging.Error("config", "Failed setup of GCP Events: ", err)
			return err
		}
		eventsConfig, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/pubsub")
		if err != nil {
			logging.Error("config", "Failed setup of GCP Events: ", err)
			return err
		}
		gEventCred = eventsConfig
		gEventClient = gEventCred.Client(oauth2.NoContext)
	}

	return nil
}

func initRpiGPIO(){
	if gConfig.RaspberryPi.BuildLedPin > 0 {
		buildLed := rpio.Pin(gConfig.RaspberryPi.BuildLedPin)
		buildLed.Output()
		buildLed.Low()
	}

	if gConfig.RaspberryPi.DataLedPin > 0 {
		dataLed := rpio.Pin(gConfig.RaspberryPi.DataLedPin)
		dataLed.Output()
		dataLed.Low()
	}
	for _, pin := range gConfig.RaspberryPi.CycleFlashers {
		p := rpio.Pin(pin)
		p.Output()
		p.Low()
	}
}

func GetServerName()string{
	checkInitialisedOrPanic()
	return gConfig.Name
}

func TLS()*tls.Config{
	checkInitialisedOrPanic()
	return gTls
}

func Pubsub() *http.Client {
	checkInitialisedOrPanic()
	return gEventClient
}

func All()*Config{
	checkInitialisedOrPanic()
	return gConfig
}

func checkInitialisedOrPanic(){
	if gConfig == nil{
		panic("Config not initialised")
	}
	//if gTls == nil{
	//	panic("TLS not initialised")
	//}
}
