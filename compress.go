package main

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io/fs"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nfnt/resize"
)

const maxGoroutineNum = 10

func Compress(originPath string) error {
	originPath, err := filepath.Abs(originPath)
	if err != nil {
		return err
	}

	config, err := GetConfig()
	if err != nil {
		return err
	}
	assContext := AssContext{
		ctx:  context.TODO(),
		conf: config,
	}
	fileInfos, err := ioutil.ReadDir(originPath)
	if err != nil {
		return err
	}

	imgInfos := make([]fs.FileInfo, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() || !isImg(fileInfo) {
			continue
		}
		imgInfos = append(imgInfos, fileInfo)
	}

	// 先按每个goroutine处理5个图片计算需要多少goroutine
	goroutineNum := len(imgInfos)/5 + 1

	// 如果超过最大goroutine个数，使用最大个数
	if goroutineNum > maxGoroutineNum {
		goroutineNum = maxGoroutineNum
	}

	// 计算每个goroutine需要处理多少图片
	count := len(imgInfos) / goroutineNum
	var tmpImgNames []string
	wg := sync.WaitGroup{}
	successStrChan := make(chan string)
	errStrChan := make(chan string, len(imgInfos))

	var bar Bar
	bar.NewOption(0, int64(len(imgInfos)))
	progress := make([]struct{}, 0, len(imgInfos))

	go func() {
		bar.Play(0)
	ForLoop:
		for {
			select {
			case s := <-successStrChan:
				if s == "done" {
					bar.Finish()
					break ForLoop
				} else {
					progress = append(progress, struct{}{})
					bar.Play(int64(len(progress)))
				}
			}
		}
	}()

	if len(imgInfos)%goroutineNum != 0 {
		wg.Add(goroutineNum + 1)
	} else {
		wg.Add(goroutineNum)
	}

	for i := 0; i < len(imgInfos); i++ {
		tmpImgNames = append(tmpImgNames, originPath+"/"+(imgInfos[i]).Name())

		if (i+1)%count == 0 {
			go handleImg(assContext, &wg, tmpImgNames, originPath, successStrChan, errStrChan)
			tmpImgNames = []string{}
		}
	}

	if len(tmpImgNames) != 0 {
		go handleImg(assContext, &wg, tmpImgNames, originPath, successStrChan, errStrChan)
	}
	wg.Wait()
	successStrChan <- "done"
	close(successStrChan)
	close(errStrChan)
	for errImgName := range errStrChan {
		fmt.Printf("handle err, err file:%s\n", errImgName)
	}
	return nil
}

func handleImg(ctx AssContext, wg *sync.WaitGroup, imgNames []string, originPath string, successChan, errChan chan string) {
	defer wg.Done()
	var file os.File
	defer func(file *os.File) {
		if file != nil {
			err := file.Close()
			if err != nil {
				return
			}
		}
	}(&file)
	for _, imgName := range imgNames {
		file, err := os.Open(imgName)
		if err != nil {
			errChan <- imgName
			continue
		}
		img, _, err := image.Decode(file)
		if err != nil {
			continue
		}

		width := ctx.conf.Width
		height := ctx.conf.Height

		if ctx.conf.Scale != 0 {
			sqrt := math.Sqrt(float64(ctx.conf.Scale))
			width = uint(float64(img.Bounds().Dx()) * sqrt)
			height = uint(float64(img.Bounds().Dy()) * sqrt)
		}

		// 压缩图片
		targetImg := resize.Thumbnail(width, height, img, resize.Lanczos3)

		// 生成导出后的图片
		err = saveImg(originPath+"/output", path.Base(file.Name()), targetImg)
		if err != nil {
			continue
		}
		successChan <- file.Name()
	}
}

func saveImg(dir, fileName string, img image.Image) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	file, err := os.Create(dir + "/" + fileName) // ignore_security_alert
	defer func() {
		err := file.Close()
		if err != nil {
			return
		}
	}()
	if err != nil {
		return err
	}

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return err
	}
	return nil
}

func isImg(file fs.FileInfo) bool {
	return strings.HasSuffix(file.Name(), "jpg") || strings.HasSuffix(file.Name(), "jpeg") || strings.HasSuffix(file.Name(), "JPG")
}
