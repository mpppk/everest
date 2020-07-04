package main

import (
	"fmt"

	"github.com/mpppk/everest/lib"
)

const cmdPkgPath = "github.com/mpppk/everest/cmd"

func main() {
	conf := &lib.AppConfig{}

	buildOption := &lib.BuildOption{
		Option: lib.Option{
			Dir: ".",
		},
		OutputPath: "./everest",
		BuildPath:  ".",
		LdFlags:    []string{fmt.Sprintf("-X %s.appMode=true", cmdPkgPath)},
	}
	buildLog, err := lib.GoBuild(buildOption)
	if err != nil {
		panic(err)
	}
	fmt.Println(buildLog)
	if _, err := lib.BuildMacOsApp(conf, "everest", "."); err != nil {
		panic(err)
	}
}
