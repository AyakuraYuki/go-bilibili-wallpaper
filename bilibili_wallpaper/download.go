package bilibili_wallpaper

import (
	"fmt"
	"github.com/AyakuraYuki/bilibili-wallpaper/plugins/colors"
	"github.com/AyakuraYuki/bilibili-wallpaper/plugins/file"
	"github.com/AyakuraYuki/bilibili-wallpaper/plugins/filenamify"
	"github.com/AyakuraYuki/bilibili-wallpaper/plugins/misc"
	nhttp "github.com/AyakuraYuki/bilibili-wallpaper/plugins/net/http"
	"github.com/AyakuraYuki/bilibili-wallpaper/plugins/part"
	"log"
	"net/http"
	"os"
	"path"
)

func getUrlAndFilename() (res []*UrlAndFilename) {
	res = make([]*UrlAndFilename, 0)
	data, err := readExistData()
	if err != nil {
		log.Println(colors.Red("读取壁纸列表暂存数据失败 %v", err))
		return res
	}
	for _, item := range data.Data.Items {
		if len(item.Pictures) == 0 {
			continue
		}
		description := fmt.Sprintf("%s%s", item.Title, item.Description)
		safetyFilename, _ := filenamify.Filenamify(description, filenamify.Options{})
		for index, picture := range item.Pictures {
			fileExt := path.Ext(picture.ImgSrc)
			v := &UrlAndFilename{
				Url:      picture.ImgSrc,
				Filename: fmt.Sprintf("[%d] [%d] %s [%dx%d]%s", item.DocId, index, safetyFilename, picture.ImgWidth, picture.ImgHeight, fileExt),
			}
			res = append(res, v)
		}
	}
	return res
}

func Download() {
	res := getUrlAndFilename()
	if len(res) == 0 {
		return
	}

	tasks := make([]*UrlAndFilename, 0)
	for _, item := range res {
		filename := item.Filename
		fullPath := path.Join(DistDir, filename)
		ok, err := file.IsPathExist(fullPath)
		if err != nil {
			log.Println(colors.Red("err: %v", err))
			continue
		}
		if ok {
			log.Println(colors.White("已存在文件 %s，跳过", fullPath))
			continue
		}
		tasks = append(tasks, item)
	}

	taskAmount := len(tasks)
	counterChan := make(chan bool)
	defer close(counterChan)
	counter := uint64(0)
	go func() {
		for range counterChan {
			counter += 1
			if Verbose {
				fmt.Printf("downloading [%v / %v] \r", counter, taskAmount)
			}
		}
	}()

	if MultiThread {

		partitionSize := taskAmount / Coroutines
		funcs := make([]misc.WorkFunc, 0)

		for indexRange := range part.Partition(len(tasks), partitionSize) {
			bulkTasks := tasks[indexRange.Low:indexRange.High]
			funcs = append(funcs, func() error {
				var err0 error
				for _, task := range bulkTasks {
					filename := task.Filename
					fullPath := path.Join(DistDir, filename)
					if err := downloadInternal(task.Url, fullPath); err == nil {
						counterChan <- true
					}
				}
				return err0
			})
		}
		if err := misc.MultiRun(funcs...); err != nil {
			log.Println(colors.Red("多线下载遇到异常: %v", err))
			return
		}

	} else {

		for _, task := range tasks {
			filename := task.Filename
			fullPath := path.Join(DistDir, filename)
			if err := downloadInternal(task.Url, fullPath); err == nil {
				counterChan <- true
			}
		}

	}

	return
}

func downloadInternal(downloadUrl, fullPath string) error {
	headers := http.Header{}
	headers.Add("cookie", Cookie)
	bs, _, _, err := nhttp.GetRaw(nil, downloadUrl, headers, nil, 30*1000, 2)
	if err != nil {
		return err
	}
	return os.WriteFile(fullPath, bs, os.ModePerm)
}
