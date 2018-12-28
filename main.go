package main

import (
	"flag"

	"bobthebuilder/builder"
	"bobthebuilder/config"
	"bobthebuilder/web"
)

var (
	definitionsDir = flag.String("definitions_dir", "definitions", "Path to the directory containing definition files")
	baseDir       = flag.String("base_dir", "base", "Path to the directory containing base files")
	buildDir       = flag.String("build_dir", "build", "Path to the directory where builds will take place")
)

func main() {
	flag.Parse()

	if flag.Arg(0) != "" {
		e := config.Load(flag.Arg(0))
		if e != nil {
			return
		}
	} else {
		e := config.Load("config.json")
		if e != nil {
			return
		}
	}

	b := builder.GetInstance()
	builder.DefinitionsDir = *definitionsDir
	builder.BuildDir = *buildDir
	builder.BaseDir = *baseDir

	e := b.Init()
	if e != nil {
		return
	}

	web.Initialise()
	web.Run()
}
