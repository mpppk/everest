package lib

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/shurcooL/vfsgen"
)

func GenerateEmbeddedPackage(targetDirPath, outputPath string) error {
	targetPkgName := "embedded"
	os.Remove(targetPkgName)

	if err := os.MkdirAll(targetPkgName, 0777); err != nil {
		return err
	}
	testFs := http.Dir(filepath.Join(targetDirPath))
	opt := vfsgen.Options{
		PackageName:  targetPkgName,
		Filename:     filepath.Join(outputPath, targetPkgName, targetPkgName+".go"),
		VariableName: "Embedded",
	}

	if err := vfsgen.Generate(testFs, opt); err != nil {
		return err
	}
	return nil
}
