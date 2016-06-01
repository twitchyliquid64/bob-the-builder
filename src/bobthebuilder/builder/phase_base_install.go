package builder


import (
  "bobthebuilder/logging"
  "time"
  "path"
  "io"
  "os"
)


type BaseInstallPhase struct{
  BasicPhase
  BaseAbsPath string
}
func (p * BaseInstallPhase)init(index int){
  p.Type = "BASE-INSTALL"
  p.StatusString = PHASE_STATUS_READY
  p.Index = index
}
func (p * BaseInstallPhase)String()string{
  return "(phase)" + p.Type + " -- " + p.StatusString + " (" + p.BaseAbsPath + ")"
}


func (p * BaseInstallPhase)Run(r* Run, builder *Builder, defIndex int)int{
  p.Start = time.Now()

  pwd, _ := os.Getwd()
  err := copy_folder(p.BaseAbsPath, path.Join(pwd, BUILD_TEMP_FOLDER_NAME), true)

  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)

  if err == nil{
    p.ErrorCode = 0
    p.StatusString = "Install successful"
    p.WriteOutput("Base install successful.", r, builder, defIndex)
    return 0
  } else {
    p.ErrorCode = -1
    logging.Error("phase-clean", err)
    p.StatusString = "Error"
    p.WriteOutput("Error: " + p.StatusString, r, builder, defIndex)
    return -1
  }
}








func copy_folder(source string, dest string, start bool) (err error) {

	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if !start && err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := path.Join(source,obj.Name())

		destinationfilepointer := path.Join(dest, obj.Name())

		if obj.IsDir() {
			err = copy_folder(sourcefilepointer, destinationfilepointer, false)
			if err != nil {
				logging.Error("phase-clean", err)
			}
		} else {
			err = copy_file(sourcefilepointer, destinationfilepointer)
			if err != nil {
				logging.Error("phase-clean", err)
			}
		}

	}
	return
}






func copy_file(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}
