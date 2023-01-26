package main

import "github.com/Natcel0711/gouser/app"

func main() {
	err := app.SetupApp()
	if err != nil {
		panic("error setting up app")
	}
}
