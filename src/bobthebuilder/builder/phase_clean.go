package builder

import (
	"bobthebuilder/logging"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type CleanPhase struct {
	BasicPhase
	DeletePath string
}

func (p *CleanPhase) init(index int) {
	p.Type = "CLEAN"
	p.StatusString = PHASE_STATUS_READY
	p.Index = index
}
func (p *CleanPhase) String() string {
	return "(phase)" + p.Type + " -- " + p.StatusString + " (" + p.DeletePath + ")"
}

func (p *CleanPhase) Run(r *Run, builder *Builder, defIndex int) int {
	err := os.RemoveAll(p.DeletePath)
	p.End = time.Now()
	p.Duration = p.End.Sub(p.Start)

	if err == nil {
		if len(r.fileData) > 0 {
			os.MkdirAll(BuildDir, 0777)
			for fname, _ := range r.fileData {
				logging.Info("CleanPhase-run", "Writing runvar-file to: "+path.Join(BuildDir, fname))
				p.WriteOutput("Writing runvar-file to: "+path.Join(BuildDir, fname), r, builder, defIndex)
				err = ioutil.WriteFile(path.Join(BuildDir, fname), r.fileData[fname], 0777)
				if err != nil {
					logging.Error("CleanPhase-run", "Could not save file "+fname+": "+err.Error())
					p.WriteOutput("Error: "+err.Error(), r, builder, defIndex)
				} else {
					delete(r.fileData, fname)
				}
			}
		}
	}

	if err == nil {
		p.ErrorCode = 0
		p.StatusString = "Clean successful"
		return 0
	}

	p.ErrorCode = -1
	logging.Error("phase-clean", err)
	p.StatusString = err.Error()
	return -1
}
