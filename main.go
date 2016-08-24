package main

import (
	"bobthebuilder/builder"
	"bobthebuilder/config"
	"bobthebuilder/web"
	"os"
)

func main() {

	if len(os.Args) > 1 {
		e := config.Load(os.Args[1])
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
	e := b.Init()
	if e != nil {
		return
	}

	web.Initialise()
	web.Run()
}
