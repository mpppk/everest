package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
	"github.com/shurcooL/vfsgen"
)

func main() {
	rootPkgName := "self"

	os.Remove(rootPkgName)

	if err := os.MkdirAll(rootPkgName, 0777); err != nil {
		panic(err)
	}

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	srcDirs := []string{"cmd", "command", "internal", "server"}
	if err := copyDirs(srcDirs, tempDir); err != nil {
		panic(err)
	}

	srcFiles := []string{"main.go", "go.mod", "go.sum"}
	for _, file := range srcFiles {
		baseName := filepath.Base(file)
		if err := copyFile(file, filepath.Join(tempDir, baseName)); err != nil {
			panic(err)
		}
	}

	testFs := http.Dir(tempDir)
	opt := vfsgen.Options{
		PackageName:  rootPkgName,
		Filename:     filepath.Join(rootPkgName, "assets_vfsdata.go"),
		VariableName: "Assets",
	}

	if err := vfsgen.Generate(testFs, opt); err != nil {
		panic(err)
	}
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
