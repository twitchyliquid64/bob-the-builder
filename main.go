package main

import (
	"bobthebuilder/builder"
	"bobthebuilder/config"
	"bobthebuilder/web"
)



func main() {
	config.Load("testconfig.json")

	b := builder.GetInstance()
	e := b.Init()
	if e != nil{
		return
	}

	web.Initialise()
	web.Run()
}
