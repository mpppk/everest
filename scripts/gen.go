package main

import (
	"path/filepath"

	"github.com/mpppk/everest/lib"
)

func main() {
	if err := lib.GenerateSelfPackage(".", "."); err != nil {
		panic(err)
	}

	if err := generateDefaultEmbeddedPackage(); err != nil {
		panic(err)
	}
}

func generateDefaultEmbeddedPackage() error {
	return lib.GenerateEmbeddedPackage(filepath.Join("defaultembedded", "build"), ".")
}
