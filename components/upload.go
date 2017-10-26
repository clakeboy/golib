package components

import (
	"strings"
	"encoding/base64"
	"io/ioutil"
	"ck_go_lib/utils"
	"errors"
	"os"
	"time"
	"path/filepath"
)

type Upload struct {

}

func NewUpload() *Upload{
	return &Upload{}
}

var ImageMime = map[string]string{
	"image/jpeg":"jpg",
	"image/jpg":"jpg",
	"image/png":"png"}

const saveDir = "./assets/ig/"
const httpSaveDir = "/assets/ig/"
const tmpDir = "./cache/ig/"
const httpTmpDir = "/assets/tmp/"

func (*Upload) ImageBase64(base64_str string) string {
	img_type := base64_str[5:strings.Index(base64_str,";")]

	if ImageMime[img_type] == "" {
		panic(errors.New("can not support this mime"))
	}

	img_str := base64_str[strings.Index(base64_str,",")+1:]

	img,base64_err := base64.StdEncoding.DecodeString(img_str)
	if base64_err != nil {
		panic(base64_err)
	}

	file_date := utils.FormatTime("YYMMDD",time.Now().Unix())

	file_name := file_date+"_"+utils.RandStr(10,nil) + "." + ImageMime[img_type]

	err := os.MkdirAll(tmpDir,0755)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(tmpDir+file_name,img,0755)
	if err != nil {
		panic(err)
	}

	return httpTmpDir+file_name
}

func (*Upload) MoveTmpToSave(tmp_file string) string {
	file_name := filepath.Base(tmp_file)
	list := strings.Split(file_name,"_")

	err := os.MkdirAll(saveDir+list[0],0755)
	if err != nil {
		panic(err)
	}

	old_file := tmpDir+file_name
	new_file := saveDir+list[0]+"/"+list[1]

	err = os.Rename(old_file,new_file)
	if err != nil {
		panic(err)
	}

	return httpSaveDir+list[0]+"/"+list[1]
}
