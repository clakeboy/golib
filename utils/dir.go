package utils

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"strings"
)

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func GetDirectoryCount(dir_path string) (int, error){
	length := 0

	err := filepath.Walk(dir_path,func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		length++

		return nil
	})

	if err != nil {
		return 0,err
	}

	return length, nil
}

func ReadDirRecursion(dir_path string) ([]string,error) {
	dir_path = fixDirPath(dir_path);
	dir_path,_ = filepath.Abs(dir_path)
	var files []string

	err := filepath.Walk(dir_path,func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		files = append(files,path)

		return nil
	})

	if err != nil {
		return nil,err
	}

	return files, nil
}

func ReadDir(path string) (files []string, dirs []string, err error){
	path = fixDirPath(path)
	dir,err := ioutil.ReadDir(path)
	if err != nil {
		files = nil
		dirs = nil
		return
	}

	for _,info := range dir {
		if info.IsDir() {
			dirs = append(dirs,path+"/"+info.Name()+"/")
		} else {
			files = append(files,path+"/"+info.Name())
		}
	}

	return
}

func fixDirPath(path string) string {
	return strings.TrimRight(path,"/")
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}