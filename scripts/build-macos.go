package main

import (
	"github.com/mpppk/everest/lib"
)

func main() {
	conf := &lib.MacOsAppConfig{
		AppName:        "everest.app",
		ExecutablePath: "./everest",
		IconPath:       "./defaultembedded/src/everest.icns",
		Identifier:     "com.github.mpppk.everest",
	}

	if _, err := lib.BuildMacOsApp(conf, "."); err != nil {
		panic(err)
	}
}
