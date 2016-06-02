package builder

import (
  //"bobthebuilder/logging"
  "os/exec"
  "path"
  "time"
  "os"
)



type CommandPhase struct{
  BasicPhase
  Command string
  Args []string
  CanFail bool

  run *Run `json:"-"`
  builder *Builder `json:"-"`
  defIndex int `json:"-"`
}
func (p * CommandPhase)init(index int){
  p.Type = "COMMAND"
  p.StatusString = PHASE_STATUS_READY
  p.Index = index
}




func (p * CommandPhase)Run(r* Run, builder *Builder, defIndex int)int{
  p.Start = time.Now()

  p.run = r
  p.builder = builder
  p.defIndex = defIndex

  pwd, _ := os.Getwd()


  //make sure build dir exists
  if exists, _ := exists(path.Join(pwd, BUILD_TEMP_FOLDER_NAME)); !exists {
    os.MkdirAll(path.Join(pwd, BUILD_TEMP_FOLDER_NAME), 700)
  }

  cmd := exec.Command(p.Command, p.Args...)
  cmd.Dir = path.Join(pwd, BUILD_TEMP_FOLDER_NAME)

  cmd.Stdout = p
  cmd.Stderr = p
  cmd.Start()
  err := cmd.Wait()

  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)

  if err != nil {
    p.ErrorCode = -1
    p.StatusString = err.Error()
    e, ok := err.(*exec.ExitError)
    if ok{
      p.WriteOutput("Process Error: " + e.String(), r, builder, defIndex)
    }else {
      p.WriteOutput("Command setup error. Are you sure the command exists on this system?", r, builder, defIndex)
    }
    if p.CanFail{
      return 0
    }
    return -1
  }else{
    p.ErrorCode = 0
    p.StatusString = "Completed successfully"
    return 0
  }
}

func (p * CommandPhase)Write(in []byte)(n int, err error){
  //logging.Info("command-phase", string(in))
  p.WriteOutput(string(in), p.run, p.builder, p.defIndex)
  return len(in), nil
}




// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}
