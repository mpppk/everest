package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rakyll/statik/fs"

	"github.com/otiai10/copy"
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

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return os.Remove(dir)
}

func WriteFs(fileSystem http.FileSystem, dst string) error {
	return fs.Walk(fileSystem, "/", func(path string, info os.FileInfo, err error) error {
		dstPath := filepath.Join(dst, path)
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0777)
		}
		file, err := fileSystem.Open(path)
		if err != nil {
			return err
		}
		// FIXME use io.Pipe
		contents, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(dstPath, contents, 0777)
	})
}
