package builder

import (
	//"bobthebuilder/logging"
	"os"
	"os/exec"
	"time"
)

type AptGetCheckInstallPhase struct {
	BasicPhase
	Packages []string

	run      *Run
	builder  *Builder
	defIndex int
}

func (p *AptGetCheckInstallPhase) init(index int) {
	p.Type = "APT-CHECK"
	p.StatusString = PHASE_STATUS_READY
	p.Index = index
}

func (p *AptGetCheckInstallPhase) Run(r *Run, builder *Builder, defIndex int) int {
	p.run = r
	p.builder = builder
	p.defIndex = defIndex

	//make sure build dir exists
	if exists, _ := exists(BuildDir); !exists {
		os.MkdirAll(BuildDir, 0777)
	}

	for _, pkg := range p.Packages {
		//check if it exists
		_, err := exec.Command("dpkg", "-s", pkg).Output()
		if err == nil { //non-zero exit value raises an error - means it needs to be installed
			p.Write([]byte("Package " + pkg + " is already installed.\n"))
			continue
		}
		p.Write([]byte("\n" + pkg + ":\n"))

		//run install if it doesnt
		cmd := exec.Command("apt-get", "-y", "install", pkg)
		cmd.Dir = BuildDir

		cmd.Stdout = p
		cmd.Stderr = p
		cmd.Start()
		err = cmd.Wait()

		if err != nil {
			p.ErrorCode = -1
			p.StatusString = err.Error()
			e, ok := err.(*exec.ExitError)
			if ok {
				p.WriteOutput("Process Error: "+e.String(), r, builder, defIndex)
			} else {
				p.WriteOutput("Internal error. Does the aptitude package manager exist on your system?", r, builder, defIndex)
			}
			p.End = time.Now()
			p.Duration = p.End.Sub(p.Start)
			p.ErrorCode = -1
			p.StatusString = err.Error()
			return -1
		}
	}

	p.End = time.Now()
	p.Duration = p.End.Sub(p.Start)
	p.ErrorCode = 0
	p.StatusString = "Completed successfully"
	return 0

}

func (p *AptGetCheckInstallPhase) Write(in []byte) (n int, err error) {
	//logging.Info("command-phase", string(in))
	p.WriteOutput(string(in), p.run, p.builder, p.defIndex)
	return len(in), nil
}
