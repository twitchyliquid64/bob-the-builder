package builder

import (
	//"bobthebuilder/logging"
	"os"
	"os/exec"
	"time"
)

type GitClonePhase struct {
	BasicPhase
	GitSrcPath string
	CanFail    bool

	run      *Run
	builder  *Builder
	defIndex int
}

func (p *GitClonePhase) init(index int) {
	p.Type = "GIT-CLONE"
	p.StatusString = PHASE_STATUS_READY
	p.Index = index
}

func (p *GitClonePhase) Run(r *Run, builder *Builder, defIndex int) int {
	p.run = r
	p.builder = builder
	p.defIndex = defIndex

	//make sure build dir exists
	if exists, _ := exists(BuildDir); !exists {
		os.MkdirAll(BuildDir, 0777)
	}

	cmd := exec.Command("git", "clone", p.GitSrcPath, ".")
	cmd.Dir = BuildDir

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
		if ok {
			p.WriteOutput("Process Error: "+e.String(), r, builder, defIndex)
		} else {
			p.WriteOutput("Command setup error. Are you sure the command exists on this system?", r, builder, defIndex)
		}
		if p.CanFail {
			return 0
		}
		return -1
	} else {
		p.ErrorCode = 0
		p.StatusString = "Completed successfully"
		return 0
	}
}

func (p *GitClonePhase) Write(in []byte) (n int, err error) {
	//logging.Info("command-phase", string(in))
	p.WriteOutput(string(in), p.run, p.builder, p.defIndex)
	return len(in), nil
}
