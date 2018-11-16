package components

import (
	"strings"
	"fmt"
	"io/ioutil"
	"ck_go_lib/utils"
	"regexp"
	"reflect"
	"strconv"
	"net/http"
	"os"
)

type Template struct {
	templateDir string //模板文件目录
	ext string //默认文件后缀
	assigns utils.M //模板变量
	cache *MemCache
}

type TemplateCon struct{
	Content []byte
	ModTime int64
}

//变量 regexp
var regVar = regexp.MustCompile(`\{%([0-9_a-zA-Z\[\]'"$\.\x7f-\xff]+)\}`)

func NewTemplate(tempDir string) *Template {
	return &Template{
		templateDir:tempDir,
		ext:"html",
		assigns:utils.M{},
		cache:NewMemCache(),
	}
}
//编辑输出模块内容
func (t *Template) Parse(tempName string) (string,error) {
	arr := strings.Split(tempName,".")
	allPath := fmt.Sprintf("%s/%s.%s",t.templateDir,strings.Join(arr,"/"),t.ext)

	res,err := t.checkAndGetFile(allPath)
	if err != nil {
		return "",err
	}
	content := t.replaceVariable(res)

	return content,nil
}
//检查模板文件是否被修改
func (t *Template) checkAndGetFile(filePath string) ([]byte,error) {
	fi,err := os.Stat(filePath)
	if err != nil {
		return nil,err
	}
	mt := fi.ModTime().Unix()
	fcon,err := t.cache.Get(filePath)

	if err == nil && fcon.(*TemplateCon).ModTime == mt {
		return fcon.(*TemplateCon).Content,nil
	}

	res,err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil,err
	}

	fcon = &TemplateCon{
		Content:res,
		ModTime:mt,
	}

	t.cache.Set(filePath,fcon,-1)
	return res,nil
}

//设置模板变量
func (t *Template) Assign(key string,val interface{}) {
	t.assigns[key] = val
}
//替换模板内变量
func (t *Template) replaceVariable(rawByte []byte) string {
	content := string(rawByte)
	//regVar := regexp.MustCompile(`\{%([0-9_a-zA-Z\[\]'"$\.\x7f-\xff]+)\}`)
	list := regVar.FindAllStringSubmatch(content,-1)
	for _,v := range list {
		arr := strings.Split(v[1],".")
		if len(arr) > 1 {
			content = strings.Replace(
				content,v[0],
				fmt.Sprintf("%v",t.getMapValue(t.assigns[arr[0]],arr[1])), -1)
		} else {
			content = strings.Replace(
				content,v[0],
				fmt.Sprintf("%v",t.assigns[arr[0]]), -1)
		}
	}
	return content
}
//从MAP类型里返回一个值
func (t *Template) getMapValue(obj interface{},key string) interface{} {
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Map:
		val := v.MapIndex(reflect.ValueOf(key))
		return utils.YN(val.IsNil(),nil,val.Interface())
	case reflect.Slice:
		idx,_ := strconv.Atoi(key)
		return v.Index(idx).Interface()
	case reflect.Ptr:
		return v.Elem().FieldByName(key).String()
	}
	return nil
}

//检查变量名是否为调用 数据或HASH
func (t *Template) checkVarMap(varName string) {
	regMap := regexp.MustCompile(`([0-9_a-zA-Z\x7f-\xff]+)\[(\d)\]`)
	if regMap.MatchString(varName) {

	}
}
//向浏览器输出内容
func (t *Template) Display(resp http.ResponseWriter,html []byte) {
	resp.Header().Set("content-type","text/html;charset=utf-8")
	resp.WriteHeader(200)
	resp.Write(html)
}