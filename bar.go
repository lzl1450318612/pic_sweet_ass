package main

import (
	"fmt"
)

type Bar struct {
	Percent int64  //百分比
	Cur     int64  //当前进度位置
	Total   int64  //总进度
	Rate    string //进度条
	Graph   string //显示符号
}

func (bar *Bar) NewOption(start, total int64) {
	bar.Cur = start
	bar.Total = total
	if bar.Graph == "" {
		bar.Graph = "█"
	}
	bar.Percent = bar.getPercent()
	for i := 0; i < int(bar.Percent); i += 2 {
		bar.Rate += bar.Graph //初始化进度条位置
	}
}

func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.Cur) / float32(bar.Total) * 100)
}

func (bar *Bar) NewOptionWithGraph(start, total int64, graph string) {
	bar.Graph = graph
	bar.NewOption(start, total)
}

func (bar *Bar) Play(cur int64) {
	bar.Cur = cur
	last := bar.Percent
	bar.Percent = bar.getPercent()
	if bar.Percent != last {
		bar.Rate = ""
		for i := int64(0); i < bar.Percent; i++ {
			bar.Rate += bar.Graph
		}
	}
	fmt.Printf("\r[%-100s]%3d%%  %8d/%d", bar.Rate, bar.Percent, bar.Cur, bar.Total)
}

func (bar *Bar) Finish() {
	fmt.Println("  Done!")
}
