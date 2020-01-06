package main

import "github.com/mpppk/everest/cmd"

//go:generate go run scripts/gen.go

func main() {
	cmd.Execute()
}
