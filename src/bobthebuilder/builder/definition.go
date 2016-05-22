package builder

import (
  "bobthebuilder/logging"
  "path"
  "os"
)

const BASE_FOLDER_NAME = "base"

type BuildDefinition struct {
  Name string `json:"name"`
  AptPackagesRequired []string `json:"apt-packages-required"`
  BaseFolder string `json:"base-folder"`
  Steps []struct {
    Command string `json:"command"`
    CanFail bool `json:"can-fail"`
  } `json:"steps"`
}

func (d *BuildDefinition)Validate()bool{
  if d.BaseFolder != "" {
    pwd, _ := os.Getwd() //cant have error - would have failed in file/util.go
    if _, err := os.Stat(path.Join(pwd, BASE_FOLDER_NAME, d.BaseFolder)); os.IsNotExist(err) {
      logging.Error("definition-validate", d.Name + " base folder does not exist.")
      return false// base folder does not exist
    }
  }
  return true
}
