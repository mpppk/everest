package lib

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
	"github.com/rakyll/statik/fs"
)

func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
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
