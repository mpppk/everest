package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mpppk/everest/lib"

	"github.com/otiai10/copy"
	"github.com/shurcooL/vfsgen"
)

func main() {
	if err := generateSelfPackage(); err != nil {
		panic(err)
	}

	if err := generateDefaultEmbeddedPackage(); err != nil {
		panic(err)
	}
}

func generateSelfPackage() error {
	srcDirs := []string{"cmd", "command", "internal", "server"}
	srcFiles := []string{"main.go", "go.mod", "go.sum"}
	return generateAssetPackage("self", srcDirs, srcFiles, "self.go", "Self")

}

func generateDefaultEmbeddedPackage() error {
	return lib.GenerateEmbeddedPackage(filepath.Join("defaultembedded", "build"), ".")
}

func generateAssetPackage(targetPkgName string, srcDirs, srcFiles []string, assetFileName, assetVariableName string) error {
	os.Remove(targetPkgName)

	if err := os.MkdirAll(targetPkgName, 0777); err != nil {
		return err
	}

	tempDir, err := ioutil.TempDir("", "")
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
		Filename:     filepath.Join(targetPkgName, assetFileName),
		VariableName: assetVariableName,
	}

	if err := vfsgen.Generate(testFs, opt); err != nil {
		return err
	}
	return nil
}

func copyDirs(dirs []string, dst string) error {
	for _, dir := range dirs {
		baseDir := filepath.Base(dir)
		if err := copy.Copy(dir, filepath.Join(dst, baseDir)); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}
