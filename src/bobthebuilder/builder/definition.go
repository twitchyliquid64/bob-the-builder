package builder

import (
  "bobthebuilder/logging"
  "bobthebuilder/util"
  "path"
  "os"
)

const BASE_FOLDER_NAME = "base"

type BuildDefinition struct {
  Name string `json:"name"`
  AptPackagesRequired []string `json:"apt-packages-required"`
  BaseFolder string `json:"base-folder"`
  GitSrc string `json:"git-src"`
  Steps []struct {
    Command string `json:"command"`
    CanFail bool `json:"can-fail"`
  } `json:"steps"`

  CurrentRun *Run
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


func (d *BuildDefinition)genRun()*Run{
  out := &Run{
    Definition: d,
    GUID: util.RandAlphaKey(32),
    ExecType: "BUILD",
    Version: "?",
    Tags: []string{
      "auto",
    },
  }
  pwd, _ := os.Getwd() //cant have error - would have failed in file/util.go

  //clean up the build folder
  delPhase := &CleanPhase{
    DeletePath: path.Join(pwd, BUILD_TEMP_FOLDER_NAME),
  }
  delPhase.init()
  out.Phases = append(out.Phases, delPhase)

  if d.GitSrc != "" {//if we are sourcing files from git, that needs to happen first for reasons.
    p := &GitClonePhase{
      GitSrcPath: d.GitSrc,
    }
    p.init()
    out.Phases = append(out.Phases, p)
  }

  if d.BaseFolder != "" {//next, copy in any static files specified.
    p := &BaseInstallPhase{
      BaseAbsPath: path.Join(pwd, BASE_FOLDER_NAME, d.BaseFolder),
    }
    p.init()
    out.Phases = append(out.Phases, p)
  }
  logging.Info("definition-run-generate", out) //TODO: Remove
  return out
}
