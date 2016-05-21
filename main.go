package main

import (
	"bobthebuilder/config"
	"bobthebuilder/web"
)



func main() {
	config.Load("testconfig.json")
	web.Initialise()
	web.Run()
}
