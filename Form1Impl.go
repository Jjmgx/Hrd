// 由res2go自动生成。
// 在这里写你的事件。

package main

import (
	"time"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

//::private::
type TForm1Fields struct {
}

var kHW = int32(107) //每小块的宽高
var xPy = int32(66)  //角色图片在背景的偏移x
var yPy = int32(66)  //角色图片在背景的偏移y
var selX = int32(-1) //被选择的角色位置x
var selY = int32(-1) //被选择的角色位置y
var seled = false    //是否有选择的角色
var zhenFa TZHF      //当前刷新的角色数据

func (f *TForm1) OnFormCreate(vcl.IObject) {
	//默认为第一个阵法
	zhenFa = zhenFaList[0].ZhFa
	//窗口双缓冲刷新，这样不抖动
	f.SetDoubleBuffered(true)
	//窗口标题
	f.SetCaption("华容道")
	//初始化阵法列表
	font := f.ListBox1.Font()
	font.SetSize(18)
	f.ListBox1.SetFont(font)
	for _, z := range zhenFaList {
		f.ListBox1.Items().Add(z.Name)
	}
}
func (f *TForm1) OnListBox1DblClick(sender vcl.IObject) {
	//阵法列表双击时选择新的阵法
	n := f.ListBox1.ItemIndex()
	if n == -1 {
		return
	}
	zhenFa = zhenFaList[int(n)].ZhFa
	seled = false
	selX = int32(-1)
	selY = int32(-1)
	f.OnPaintBox1Paint(f.PaintBox1)
}

func (f *TForm1) OnPaintBox1Paint(sender vcl.IObject) {
	//画背景及角色
	p := f.PaintBox1.Canvas()
	p.Draw(0, 0, pngBg)
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if zhenFa[i][j] > 0 { //0空格及-1补位,不用画
				x := xPy + int32(j)*kHW
				y := yPy + int32(i)*kHW
				p.Draw(x, y, pngHash[zhenFa[i][j]])
				if seled && i == int(selX) && j == int(selY) { //如果是被选择的，画框
					switch zhenFa[i][j] {
					case 1:
						p.FrameRect(types.TRect{
							x, y, x + kHW*2, y + kHW*2,
						})
					case 2, 7, 8, 9, 10:
						p.FrameRect(types.TRect{
							x, y, x + kHW*2, y + kHW,
						})
					case 3, 4, 5, 6:
						p.FrameRect(types.TRect{
							x, y, x + kHW, y + kHW*2,
						})
					case 11:
						p.FrameRect(types.TRect{
							x, y, x + kHW, y + kHW,
						})
					}

				}

			}
		}
	}
}

func (f *TForm1) OnPaintBox1MouseDown(sender vcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	if button == types.MbLeft {
		//取出点击位置最左边索引及角色名
		indexX, indexY, lx := getDjInfo(x, y, zhenFa)
		if lx == -1 { //点击了非正常位置，放弃
			return
		}
		if lx > 0 { //点击了正常位置时，标边框
			selX = indexX
			selY = indexY
			seled = true
			f.PaintBox1.Repaint()
			return
		}
		if lx == 0 && seled { //点击了空位,且已经选择了角色时，开始移动
			if tz, ok := yd(selX, selY, indexX, indexY, zhenFa); ok {
				zhenFa = tz
				f.PaintBox1.Repaint()
				if zhenFa[3][1] == 1 { //成功出逃
					vcl.ShowMessage("曹瞒兵败走华容,正与关公狭路逢,\r\n只为当初恩义重,放开金锁走蛟龙。")
				}
			}
		}

	}
}

//OnButton1Click 电脑算法
func (f *TForm1) OnButton1Click(sender vcl.IObject) {
	hashjg := make(map[int]TZHF) //存储走过的步骤，用来防止重复走及回溯解法
	sz := []TZHF{}               //存储每一次要尝试走法的数据列表
	sz = append(sz, zhenFa)      //把当前数据加入到尝试列表
	jg := TZHF{}                 //记录最终成功时的数据
	for {                        //循环至无解或者得到解法
		newsz := []TZHF{}       //缓存下一次要尝试的列表
		for _, sj := range sz { //遍历本次要尝试的列表
			ydhsz := getYdsz(sj)        //得到可以移动到的结果
			for _, ydh := range ydhsz { //遍历可移动的结果
				key := getKey(ydh)  //转换成key
				if ydh[3][1] == 1 { //判断是否已是成功解，如果是，直接跳出
					hashjg[key] = sj
					jg = ydh
					goto END
				}
				if _, ok := hashjg[key]; !ok { //判断这个结果以前没有出现过时，加入到下次要尝试的列表
					hashjg[key] = sj
					newsz = append(newsz, ydh)
				}
			}
		}
		if len(newsz) == 0 { //如果下次要尝试的列表为空，证明没有解
			vcl.ShowMessage("无语2,电脑也没解出来！")
			return
		}
		//将下次要尝试的列表保存，为下次遍历作准备
		sz = newsz
	}
END: //到此处时,证明求到了解
	jgsz := []TZHF{jg}     //先把最后一次的结构保存进结果数组，为后面演示做准备
	key := getKey(jg)      //取出最后一次数据的key,往回追溯父节点
	ykey := getKey(zhenFa) //保存求解前的key,为是否追溯到根判断用
	for {                  //循环直到追溯到根
		if jf, ok := hashjg[key]; ok {
			jgsz = append(jgsz, jf)
			key = getKey(jf)
			if key == ykey {
				break
			}
		}
	}
	go func() { //开始演示
		//倒着从最后一个数据开始
		for i := len(jgsz) - 1; i >= 0; i-- {
			vcl.ThreadSync(func() { //赋值后刷新画板再延时即可
				zhenFa = jgsz[i]
				f.PaintBox1.Repaint()
				time.Sleep(500 * time.Millisecond)
			})
		}
	}()
}

func getKey(sj TZHF) int {
	//key是一个int类型，一是节约空间，二是可以更好的判断重复的不用再尝试
	jg := 1
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			sj := int(sj[i][j])
			if sj == 4 || sj == 5 || sj == 6 { //为竖着的块时，统一成3
				sj = 3
			} else if sj == 7 || sj == 8 || sj == 9 || sj == 10 { //为横块时，统一成2
				sj = 2
			} else if sj == 11 { //为小兵时,转换成4
				sj = 4
			} else if sj == 0 { //为空格时，转换成5
				sj = 5
			}
			if sj > -1 { //-1点位块时，不做处理
				jg = jg*10 + sj //合并数据得到一个int型的key
			}
		}
	}
	return jg
}

func getYdsz(sj TZHF) []TZHF {
	//保存可以走动的数据
	ls := []TZHF{}
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if sj[i][j] == 0 { //当块为空时，判断块的上下左右能否移动到此块
				sja := [3]int32{int32(i - 1), int32(i), int32(i + 1)}
				sjb := [3]int32{int32(j - 1), int32(j), int32(j + 1)}
				for _, a := range sja {
					for _, b := range sjb {
						if a >= 0 && b >= 0 && a < 5 && b < 4 && //位置合法
							(a+b-int32(i+j) == 1 || int32(i+j)-a-b == 1) { //不是斜角
							x, y := a, b
							if sj[a][b] == -1 { //如果为点位块，找到主块
								x, y = getZhukuai(a, b, sj)
							}
							//尝试移动，得到结果时保存
							if jg, ok := yd(x, y, int32(i), int32(j), sj); ok {
								ls = append(ls, jg)
								if jg[3][1] == 1 { //成功时直接返回，节约性能
									return ls
								}
							}
						}
					}
				}
			}
		}
	}
	return ls //没成功时返回所有可能性，为下次尝试做准备
}

func yd(yx, yy, nx, ny int32, TempZF TZHF) (TZHF, bool) {
	if yx == nx && yy == ny { //如果点击的是自己，直接返回
		return TempZF, false
	}
	if TempZF[nx][ny] != 0 { //如果目的位不是空，直接返回
		return TempZF, false
	}
	switch TempZF[yx][yy] {
	case 1: //曹操
		if (yy-ny == 1 && yx == nx && TempZF[nx+1][ny] == 0) || (yy-ny == 1 && yx == nx-1 && TempZF[nx-1][ny] == 0) { //左移
			TempZF[yx][yy-1] = 1
			TempZF[yx+1][yy-1] = -1
			TempZF[yx][yy] = -1
			TempZF[yx+1][yy] = -1
			TempZF[yx][yy+1] = 0
			TempZF[yx+1][yy+1] = 0
			selX, selY = yx, yy-1
			return TempZF, true
		}
		if (ny-yy == 2 && yx == nx && TempZF[nx+1][ny] == 0) || (ny-yy == 2 && yx == nx-1 && TempZF[nx-1][ny] == 0) { //右移
			TempZF[yx][yy] = 0
			TempZF[yx+1][yy] = 0
			TempZF[yx][yy+1] = 1
			TempZF[yx+1][yy+1] = -1
			TempZF[yx][yy+2] = -1
			TempZF[yx+1][yy+2] = -1
			selX, selY = yx, yy+1
			return TempZF, true
		}
		if (yx-nx == 1 && yy == ny && TempZF[nx][ny+1] == 0) || (yx-nx == 1 && yy == ny-1 && TempZF[nx][ny-1] == 0) { //上移
			TempZF[yx-1][yy] = 1
			TempZF[yx-1][yy+1] = -1
			TempZF[yx][yy] = -1
			TempZF[yx][yy+1] = -1
			TempZF[yx+1][yy] = 0
			TempZF[yx+1][yy+1] = 0
			selX, selY = yx-1, yy
			return TempZF, true
		}
		if (nx-yx == 2 && yy == ny && TempZF[nx][ny+1] == 0) || (nx-yx == 2 && yy == ny-1 && TempZF[nx][ny-1] == 0) { //下移
			TempZF[yx][yy] = 0
			TempZF[yx][yy+1] = 0
			TempZF[yx+1][yy] = 1
			TempZF[yx+1][yy+1] = -1
			TempZF[yx+2][yy] = -1
			TempZF[yx+2][yy+1] = -1
			selX, selY = yx+1, yy
			return TempZF, true
		}
	case 2, 7, 8, 9, 10: //关羽等横块
		dqjs := TempZF[yx][yy]
		if yx == nx && yy-ny == 1 { //左移
			TempZF[yx][yy-1] = dqjs
			TempZF[yx][yy] = -1
			TempZF[yx][yy+1] = 0
			selX, selY = yx, yy-1
			return TempZF, true
		}
		if yx == nx && ny-yy == 2 { //右移
			TempZF[yx][yy] = 0
			TempZF[yx][yy+1] = dqjs
			TempZF[yx][yy+2] = -1
			selX, selY = yx, yy+1
			return TempZF, true
		}
		if yx == nx && yy-ny == 2 && TempZF[yx][yy-1] == 0 { //左二
			TempZF[yx][yy-2] = dqjs
			TempZF[yx][yy-1] = -1
			TempZF[yx][yy] = 0
			TempZF[yx][yy+1] = 0
			selX, selY = yx, yy-2
			return TempZF, true
		}
		if yx == nx && ny-yy == 3 && TempZF[yx][yy+2] == 0 { //右二
			TempZF[yx][yy] = 0
			TempZF[yx][yy+1] = 0
			TempZF[yx][yy+2] = dqjs
			TempZF[yx][yy+3] = -1
			selX, selY = yx, yy+2
			return TempZF, true
		}
		if (yx-nx == 1 && yy == ny && TempZF[nx][ny+1] == 0) || (yx-nx == 1 && yy == ny-1 && TempZF[nx][ny-1] == 0) { //上移
			TempZF[yx-1][yy] = dqjs
			TempZF[yx-1][yy+1] = -1
			TempZF[yx][yy] = 0
			TempZF[yx][yy+1] = 0
			selX, selY = yx-1, yy
			return TempZF, true
		}
		if (nx-yx == 1 && yy == ny && TempZF[nx][ny+1] == 0) || (nx-yx == 1 && yy == ny-1 && TempZF[nx][ny-1] == 0) { //下移
			TempZF[yx][yy] = 0
			TempZF[yx][yy+1] = 0
			TempZF[yx+1][yy] = dqjs
			TempZF[yx+1][yy+1] = -1
			selX, selY = yx+1, yy
			return TempZF, true
		}
	case 3, 4, 5, 6: //张等竖块
		dqjs := TempZF[yx][yy]
		if yy == ny && yx-nx == 1 { //上移
			TempZF[yx-1][yy] = dqjs
			TempZF[yx][yy] = -1
			TempZF[yx+1][yy] = 0
			selX, selY = yx-1, yy
			return TempZF, true
		}
		if yy == ny && nx-yx == 2 { //下移
			TempZF[yx][yy] = 0
			TempZF[yx+1][yy] = dqjs
			TempZF[yx+2][yy] = -1
			selX, selY = yx+1, yy
			return TempZF, true
		}
		if yy == ny && yx-nx == 2 && TempZF[yx-1][yy] == 0 { //上二
			TempZF[yx-2][yy] = dqjs
			TempZF[yx-1][yy] = -1
			TempZF[yx][yy] = 0
			TempZF[yx+1][yy] = 0
			selX, selY = yx-2, yy
			return TempZF, true
		}
		if yy == ny && nx-yx == 3 && TempZF[yx+3][yy] == 0 { //下二
			TempZF[yx][yy] = 0
			TempZF[yx+1][yy] = 0
			TempZF[yx+2][yy] = dqjs
			TempZF[yx+3][yy] = -1
			selX, selY = yx+2, yy
			return TempZF, true
		}
		if (yy-ny == 1 && yx == nx && TempZF[nx+1][ny] == 0) || (yy-ny == 1 && yx == nx-1 && TempZF[nx-1][ny] == 0) { //左移
			TempZF[yx][yy-1] = dqjs
			TempZF[yx+1][yy-1] = -1
			TempZF[yx][yy] = 0
			TempZF[yx+1][yy] = 0
			selX, selY = yx, yy-1
			return TempZF, true
		}
		if (ny-yy == 1 && yx == nx && TempZF[nx+1][ny] == 0) || (ny-yy == 1 && yx == nx-1 && TempZF[nx-1][ny] == 0) { //右移
			TempZF[yx][yy] = 0
			TempZF[yx+1][yy] = 0
			TempZF[yx][yy+1] = dqjs
			TempZF[yx+1][yy+1] = -1
			selX, selY = yx, yy+1
			return TempZF, true
		}

	case 11: //小兵
		if (yx-nx == 1 && yy == ny) || //上移
			(nx-yx == 1 && yy == ny) || //下移
			(yy-ny == 1 && yx == nx) || //左移
			(ny-yy == 1 && yx == nx) || //右移
			(yx-nx == 2 && yy == ny && TempZF[yx-1][ny] == 0) || //上移二位
			(nx-yx == 2 && yy == ny && TempZF[yx+1][ny] == 0) || //下移二位
			(yy-ny == 2 && yx == nx && TempZF[yx][yy-1] == 0) || //左移二位
			(ny-yy == 2 && yx == nx && TempZF[yx][yy+1] == 0) || //右移二位
			(yx-nx == 1 && yy-ny == 1 && (TempZF[yx-1][yy] == 0 || TempZF[yx][yy-1] == 0)) || //左上
			(nx-yx == 1 && yy-ny == 1 && (TempZF[yx][yy-1] == 0 || TempZF[yx+1][yy] == 0)) || //左下
			(ny-yy == 1 && yx-nx == 1 && (TempZF[yx][yy+1] == 0 || TempZF[yx-1][yy] == 0)) || //右上
			(ny-yy == 1 && nx-yx == 1 && (TempZF[yx][yy+1] == 0 || TempZF[yx+1][yy] == 0)) { //右下
			TempZF[yx][yy], TempZF[nx][ny] = TempZF[nx][ny], TempZF[yx][yy]
			selX, selY = nx, ny
			return TempZF, true
		}
	}
	return TempZF, false
}

func getDjInfo(x, y int32, TempZF TZHF) (int32, int32, int8) {
	//判断点击的位置是否合法
	if x < xPy || y < yPy || x > xPy+kHW*4 || y > xPy+kHW*5 { //不在角色范围时返回
		return -1, -1, -1
	}
	j := (x - xPy) / kHW
	i := (y - yPy) / kHW
	if i > 4 || j > 3 { //不在角色范围时返回
		return -1, -1, -1
	}
	if TempZF[i][j] == -1 { //如果为-1点位块时，要找到其主块的角色
		i, j = getZhukuai(i, j, TempZF)
	}
	return i, j, TempZF[i][j]
}

func getZhukuai(x, y int32, TempZF TZHF) (int32, int32) {
	//根据点击的位置，来判断点击的主块是什么角色
	if x-1 >= 0 && TempZF[x-1][y] > 0 && TempZF[x-1][y] < 7 && TempZF[x-1][y] != 2 { //横块的点位
		return x - 1, y
	}
	if y-1 >= 0 && TempZF[x][y-1] > 0 && TempZF[x][y-1] < 11 { //竖块的点位
		return x, y - 1
	}
	if x-1 >= 0 && y-1 >= 0 && TempZF[x-1][y-1] > 0 && TempZF[x-1][y-1] < 11 { //曹操的占位
		return x - 1, y - 1
	}
	return x, y //都没有，只好返回自己
}
