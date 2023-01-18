package main

import (
	"context"
	"flag"
	"fmt"
)

const AssModeResize = "resize"
const AssModeCompress = "compress"

type AssMode string

type AssContext struct {
	ctx  context.Context
	conf *Config
	mode AssMode
}

var confStr = flag.String("conf", "", "配置")
var compressBool = flag.Bool("c", false, "图片压缩")
var resizeBool = flag.Bool("r", false, "图片裁剪")

func main() {
	flag.Parse()

	switch *confStr {
	case "init":
		err := CreateConfFile()
		if err != nil {
			fmt.Println("CreateConfFile failed, please check auth")
			return
		}
	default:
	}

	if *compressBool {
		args := flag.Args()
		if len(args) != 1 {
			fmt.Println("please input origin path & target path")
			return
		}

		err := Compress(args[0], AssModeCompress)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Done!")
	}

	if *resizeBool {
		args := flag.Args()
		if len(args) != 1 {
			fmt.Println("please input origin path & target path")
			return
		}

		err := Compress(args[0], AssModeResize)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Done!")
	}
}
