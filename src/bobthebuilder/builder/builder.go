package builder

import (
  "bobthebuilder/logging"
  "bobthebuilder/util"
  "encoding/json"
  "io/ioutil"
  "path"
  "sync"
)

const DEFINITIONS_FOLDER_NAME = "definitions"
const DEFINITIONS_FILE_SUFFIX = ".json"
const BUILD_TEMP_FOLDER_NAME = "build"

type Builder struct {
  Lock sync.Mutex
  Definitions []*BuildDefinition
}



// (Re)load all build definitions.
func (b *Builder)Init()error{
  b.Lock.Lock()
  defer b.Lock.Unlock()

  b.Definitions = []*BuildDefinition{} //clear our definitions

  defFiles, err := util.GetFilenameListInFolder(DEFINITIONS_FOLDER_NAME, DEFINITIONS_FILE_SUFFIX)
  if err != nil {
    return err
  }

  for _, fpath := range defFiles { //parse each configuration file
      bData, err := ioutil.ReadFile(fpath)
      if err != nil{
        logging.Error("builder-init", "Error reading build definition for " + path.Base(fpath) + ". Skipping.")
        logging.Error("builder-init", err)
        continue
      }

      var def BuildDefinition
      err = json.Unmarshal(bData, &def)
      if err != nil{
        logging.Error("builder-init", "Error parsing build definition for " + path.Base(fpath) + ". Skipping.")
        logging.Error("builder-init", err)
        continue
      }

      valid := def.Validate()
      if valid{
        b.Definitions = append(b.Definitions, &def)
        logging.Info("builder-init", def.Name + " - definition ready.")
      }else {
        logging.Warning("builder-init", "Skipping " + path.Base(fpath) + " (" + def.Name + ").")
      }
  }

  return nil
}



//Creates a new builder object. Run in init(), so keep this method simple.
func New()*Builder{
  return &Builder{}
}
