package main

import "github.com/mpppk/everest/cmd"

//go:generate statik -src .
//go:generate statik -src defaultembedded/build -p embedded

func main() {
	cmd.Execute()
}
