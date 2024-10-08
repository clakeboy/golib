package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/clakeboy/golib/utils"
)

// 创建项目模板
func main() {
	InitCommand()
	projPath := fmt.Sprintf("%s/%s", strings.TrimSuffix(CmdOut, "/"), CmdName)
	if utils.Exist(projPath) {
		fmt.Println("已经存在项目名称的目录")
		return
	}
	savePath := fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(CmdOut, "/"), CmdName, "cache")
	fmt.Println("开始获取 golang 模板...")
	fetchGolangFiles(savePath)
	if CmdFront {
		fmt.Printf("开始获取 react-%s 模板...\n", CmdFrontType)
		fetchFrontFiles(savePath)
	}
	fmt.Println("清理下载文件...")
	os.RemoveAll(savePath)
	fmt.Printf("完成项目 [%s] 初始化\n", CmdName)
	fmt.Println("进入项目目录执行 go mod tidy 完成 golang 项目依赖安装")
	if CmdFront {
		fmt.Println("进入项目目录 frontend 执行 npm i 完成前端项目依赖安装")
	}
}

var (
	CmdFront     bool   //是否初始化前端
	CmdFrontType string //前端项目模板类型 vite 默认,webpack
	CmdName      string //项目名称
	CmdOut       string //输出目录 默认当前执行目录
	CmdProxy     bool   //是否使用代理下载
	CmdGh        string //代理下载地址
)

func InitCommand() {
	flag.BoolVar(&CmdFront, "front", true, "是否初始化前端项目")
	flag.StringVar(&CmdFrontType, "front-type", "vite", "前端项目模板类型: vite (默认), webpack")
	flag.StringVar(&CmdName, "name", "CKCDemo", "项目名称")
	flag.StringVar(&CmdOut, "out", "./", "输出目录 默认当前执行目录")
	flag.BoolVar(&CmdProxy, "proxy", false, "是否用代理下载")
	flag.StringVar(&CmdGh, "gh", "https://gh.vramcc.com/", "代理下载地址")
	flag.Parse()
}

func fetchGolangFiles(savePath string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("下载 golang 项目模板出错:", err)
		}
	}()

	if !utils.Exist(savePath) {
		err := os.MkdirAll(savePath, 0755)
		if err != nil {
			panic(fmt.Errorf("创建文件夹失败 %v", err))
		}
	}
	filePath := fmt.Sprintf("%s/%s", savePath, "golang.zip")
	urlStr := "https://github.com/clakeboy/cc_template/archive/refs/heads/main.zip"
	err := downloadFile(filePath, urlStr, nil)
	if err != nil {
		panic(fmt.Errorf("打开远程文件错误 %v", err))
	}

	zp, err := zip.OpenReader(filePath)
	if err != nil {
		panic(fmt.Errorf("读取zip文件出错 %v", err))
	}
	projectPath := fmt.Sprintf("%s/%s", strings.TrimSuffix(CmdOut, "/"), CmdName)
	for _, item := range zp.File {
		zfPath := fmt.Sprintf("%s%s", projectPath, strings.ReplaceAll(item.Name, "cc_template-master", ""))
		fmt.Printf("正在解压文件 %s -> %s:\n", item.Name, zfPath)

		if item.FileInfo().IsDir() {
			if !utils.Exist(zfPath) {
				err = os.MkdirAll(zfPath, 0755)
				if err != nil {
					panic(fmt.Errorf("创建zip文件夹错误: %s,err: %v", zfPath, err))
				}
			}
			continue
			// return PkgMsg(item.Name)
		}
		rc, err := item.Open()
		if err != nil {
			panic(fmt.Errorf("打开压缩文件失败: %s,err: %v", item.Name, err))
		}
		content, err := io.ReadAll(rc)
		if err != nil {
			panic(fmt.Errorf("读取 zip 文件出错: %s,err: %v", item.Name, err))
		}
		newCon := bytes.ReplaceAll(content, []byte("cc_template"), []byte(CmdName))
		err = os.WriteFile(zfPath, newCon, 0755)
		// _, err = io.Copy(fs, rc)
		if err != nil {
			panic(fmt.Errorf("zip 文件写入到本地文件出错: %s,err: %v", item.Name, err))
		}
		rc.Close()
		// return PkgMsg(item.Name)
	}

	zp.Close()
}

func fetchFrontFiles(savePath string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("下载 react 项目模板出错:", err)
		}
	}()

	if !utils.Exist(savePath) {
		err := os.MkdirAll(savePath, 0755)
		if err != nil {
			panic(fmt.Errorf("创建文件夹失败 %v", err))
		}
	}

	filePath := fmt.Sprintf("%s/%s", savePath, "react.zip")

	urlStr := "https://github.com/clakeboy/cc_react_template/archive/refs/heads/main.zip"
	if CmdFrontType == "vite" {
		urlStr = "https://github.com/clakeboy/cc_react_template/archive/refs/heads/vite.zip"
	}
	err := downloadFile(filePath, urlStr, nil)
	if err != nil {
		panic(fmt.Errorf("打开远程文件错误 %v", err))
	}

	zp, err := zip.OpenReader(filePath)
	if err != nil {
		panic(fmt.Errorf("读取zip文件出错 %v", err))
	}
	projectPath := fmt.Sprintf("%s/%s/frontend", strings.TrimSuffix(CmdOut, "/"), CmdName)
	for _, item := range zp.File {
		zfPath := fmt.Sprintf("%s%s", projectPath, strings.ReplaceAll(item.Name, "cc_react_template-vite", ""))
		fmt.Printf("正在解压文件 %s -> %s:\n", item.Name, zfPath)

		if item.FileInfo().IsDir() {
			if !utils.Exist(zfPath) {
				err = os.MkdirAll(zfPath, 0755)
				if err != nil {
					panic(fmt.Errorf("创建zip文件夹错误: %s,err: %v", zfPath, err))
				}
			}
			continue
		}
		rc, err := item.Open()
		if err != nil {
			panic(fmt.Errorf("打开压缩文件失败: %s,err: %v", item.Name, err))
		}
		content, err := io.ReadAll(rc)
		if err != nil {
			panic(fmt.Errorf("读取 zip 文件出错: %s,err: %v", item.Name, err))
		}
		newCon := bytes.ReplaceAll(content, []byte("CCTP"), []byte(CmdName))
		err = os.WriteFile(zfPath, newCon, 0755)
		// _, err = io.Copy(fs, rc)
		if err != nil {
			panic(fmt.Errorf("zip 文件写入到本地文件出错: %s,err: %v", item.Name, err))
		}
		rc.Close()
	}
	zp.Close()
}

type httpProgress struct {
	r        io.Reader
	callbank func(n int)
}

func (hp *httpProgress) Read(p []byte) (n int, err error) {
	n, err = hp.r.Read(p)
	hp.callbank(n)
	return
}

func downloadFile(filepath string, url string, prog func(n int, total int)) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Get the data

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "CKC/1.0")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	length := resp.Header.Get("Content-Length")
	total, err := strconv.Atoi(length)
	// if err != nil {
	// 	return err
	// }
	_, err = io.Copy(out, &httpProgress{r: resp.Body, callbank: func(n int) {
		if prog != nil {
			prog(n, total)
		}
	}})
	if err != nil {
		return err
	}
	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	// Writer the body to file
	return nil
}
