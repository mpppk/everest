package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rakyll/statik/fs"

	"github.com/shurcooL/vfsgen"
)

func GenerateSelfPackage(srcPath, outputPath string) error {
	srcDirs := []string{"cmd", "lib", "internal"}
	var srcDirPaths []string
	for _, srcDir := range srcDirs {
		srcDirPaths = append(srcDirPaths, filepath.Join(srcPath, srcDir))
	}

	srcFiles := []string{"main.go", "go.mod", "go.sum"}
	var srcFilePaths []string
	for _, srcFile := range srcFiles {
		srcFilePaths = append(srcFilePaths, filepath.Join(srcPath, srcFile))
	}
	err := generateAssetPackage("self", outputPath, srcDirPaths, srcFilePaths, "self.go", "Self")
	if err != nil {
		return err
	}
	return nil
}

func generateAssetPackage(targetPkgName, outputPath string, srcDirs, srcFiles []string, assetFileName, assetVariableName string) error {
	targetPkgPath := filepath.Join(outputPath, targetPkgName)
	os.Remove(targetPkgPath)

	if err := os.MkdirAll(targetPkgPath, 0777); err != nil {
		return err
	}

	tempDir, err := ioutil.TempDir("", "everest-build")
	if err != nil {
		return err
	}

	if err := copyDirs(srcDirs, tempDir); err != nil {
		return err
	}

	for _, file := range srcFiles {
		baseName := filepath.Base(file)
		if err := copyFile(file, filepath.Join(tempDir, baseName)); err != nil {
			return err
		}
	}

	testFs := http.Dir(tempDir)
	opt := vfsgen.Options{
		PackageName:  targetPkgName,
		Filename:     filepath.Join(targetPkgPath, assetFileName),
		VariableName: assetVariableName,
	}

	if err := vfsgen.Generate(testFs, opt); err != nil {
		return err
	}

	if err := RemoveContents(tempDir); err != nil {
		return err
	}
	return nil
}

func GenerateEmbeddedPackage(targetDirPath, outputPath string) error {
	targetPkgName := "embedded"
	targetPkgPath := filepath.Join(outputPath, targetPkgName)
	os.Remove(targetPkgPath)

	if err := os.MkdirAll(targetPkgPath, 0777); err != nil {
		return err
	}
	testFs := http.Dir(filepath.Join(targetDirPath))
	opt := vfsgen.Options{
		PackageName:  targetPkgName,
		Filename:     filepath.Join(targetPkgPath, targetPkgName+".go"),
		VariableName: "Embedded",
	}
	if err := vfsgen.Generate(testFs, opt); err != nil {
		return err
	}
	return nil
}

func PrintFs(fileSystem http.FileSystem) error {
	return fs.Walk(fileSystem, "/", func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})
}
