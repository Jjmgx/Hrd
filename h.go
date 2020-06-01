package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ying32/govcl/vcl"
)

//TZHF 阵法的二维数组
type TZHF [5][4]int8

//ZhFaPz 阵法数据结构
type ZhFaPz struct {
	Name string //阵法名称
	ZhFa TZHF   //阵法二维数组
}

var pngBg = vcl.NewPngImage()
var pngHash = make(map[int8]*vcl.TPngImage)
var zhenFaList []ZhFaPz

//从文本行解析出阵法名及结构
func strToZhenFa(str string) (ZhFaPz, error) {
	zp := ZhFaPz{}
	sz := strings.Split(str, "=")
	if len(sz) != 2 {
		return zp, errors.New("no data!")
	}
	zf, err := strToSz(sz[1])
	if err != nil {
		return zp, err
	}
	zp.Name = sz[0]
	zp.ZhFa = zf
	return zp, nil
}

//从文本行解析阵法
func strToSz(str string) (TZHF, error) {
	sz := strings.Split(str, ",")
	zf := TZHF{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
	if len(sz) != 20 {
		return zf, errors.New("no data！")
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			sj, _ := strconv.Atoi(strings.Trim(sz[i*4+j], " "))
			zf[i][j] = int8(sj)
		}
	}
	return zf, nil
}
