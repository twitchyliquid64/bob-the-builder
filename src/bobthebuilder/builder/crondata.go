package builder

import (
	"bobthebuilder/logging"
	"encoding/json"
	"github.com/robfig/cron"
	"io/ioutil"

	"regexp"
	"strconv"
)

const CRONTAB_FILE_NAME = "cronentries.json"

type CronRecord struct {
	Spec             string
	TargetDefinition string
	Tags             []string
	Params           map[string]string
	DisablePhys      bool
}

func updateCron(data []CronRecord) {
	d, err := json.Marshal(&data)
	if err != nil {
		logging.Error("cron", "Could not marshal crontab: ", err)
		return
	}

	err = ioutil.WriteFile(CRONTAB_FILE_NAME, d, 0755)
	if err != nil {
		logging.Error("cron", "Could write crontab: ", err)
	}
}

func readCron() []CronRecord {
	bData, err := ioutil.ReadFile(CRONTAB_FILE_NAME)
	if err != nil {
		logging.Warning("cron", "Could not open crontab file.")
		return nil
	}

	var records []CronRecord
	err = json.Unmarshal(bData, &records)
	if err != nil {
		logging.Error("cron", "Could not read crontab: ", err)
		return nil
	}
	return records
}

func initCron(cronManager *cron.Cron, builder *Builder) {
	for i, rec := range readCron() {
		logging.Info("cron", "Inserting cron entry (", i, "): [", rec.Spec, "] ", rec.TargetDefinition)
		err := cronManager.AddFunc(rec.Spec, buildCronFunc(rec, builder))
		if err != nil {
			logging.Error("cron", "Invalid cronspec: ", err)
		}
	}
}

func buildCronFunc(rec CronRecord, builder *Builder) func() {
	return func() {
		logging.Info("cron", "Starting job for ", rec.TargetDefinition)
		builder.EnqueueBuildEventEx(rec.TargetDefinition, rec.Tags, calcNextVersionNumber(rec.TargetDefinition, builder), rec.DisablePhys, rec.Params)
	}
}

func calcNextVersionNumber(defName string, builder *Builder) string {
	re := regexp.MustCompile("\\d+$")
	bVersion := re.ReplaceAllFunc([]byte(builder.GetDefinition(defName).LastVersion), func(match []byte) []byte {
		iMatch, err := strconv.Atoi(string(match))
		if err != nil {
			return match
		} else {
			return []byte(strconv.Itoa(iMatch + 1))
		}
	})
	candidateVersion := string(bVersion)
	if candidateVersion == "" {
		return "0.0.1"
	}
	return candidateVersion
}
