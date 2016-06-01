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
    Type string `json:"type"`
    Command string `json:"command"`
    CanFail bool `json:"can-fail"`
    Args []string `json:"args"`
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
    Status: STATUS_NOT_YET_RUN,
    Tags: []string{
      "auto",
    },
  }
  pwd, _ := os.Getwd() //cant have error - would have failed in file/util.go

  //generate phase to clean up the build folder
  delPhase := &CleanPhase{
    DeletePath: path.Join(pwd, BUILD_TEMP_FOLDER_NAME),
  }
  delPhase.init(len(out.Phases))
  out.Phases = append(out.Phases, delPhase)

  if d.GitSrc != "" {//if we are sourcing files from git, that needs to happen first for reasons.
    p := &GitClonePhase{ //generate phase to clone in a git repo
      GitSrcPath: d.GitSrc,
    }
    p.init(len(out.Phases))
    out.Phases = append(out.Phases, p)
  }

  if d.BaseFolder != "" {//next, copy in any static files specified.
    p := &BaseInstallPhase{ //generate phase to copy in files
      BaseAbsPath: path.Join(pwd, BASE_FOLDER_NAME, d.BaseFolder),
    }
    p.init(len(out.Phases))
    out.Phases = append(out.Phases, p)
  }


  //generate a corresponding phase for each step
  for _, step := range d.Steps {
    switch step.Type{
    case "CMD":
      cmd := &CommandPhase{
        Command: step.Command,
        Args: step.Args,
        CanFail: step.CanFail,
      }
      cmd.init(len(out.Phases))
      out.Phases = append(out.Phases, cmd)
    }
  }


  return out
}
