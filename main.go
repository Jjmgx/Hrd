package main

import (
	"io/ioutil"
	"strconv"
	"strings"

	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
)

//初始化基本数据
func init() {
	//从配置txt读取阵法
	bb, _ := ioutil.ReadFile("config.txt")
	strsz := strings.Split(string(bb), "\r\n")
	for _, str := range strsz {
		if z, err := strToZhenFa(str); err == nil {
			zhenFaList = append(zhenFaList, z)
		}

	}
	//加载背景图片
	pngBg.LoadFromFile("png/0.png")
	//加载11个角色图片
	for i := 1; i < 12; i++ {
		pngHash[int8(i)] = vcl.NewPngImage()
		pngHash[int8(i)].LoadFromFile("png/" + strconv.Itoa(i) + ".png")
	}
}

func main() {
	//退出时释放资源
	defer func() {
		for i := 1; i < 12; i++ {
			pngHash[int8(i)].Free()
		}
	}()
	//加载窗口
	vcl.Application.Initialize()
	vcl.Application.CreateForm(&Form1)
	vcl.Application.Run()
}
