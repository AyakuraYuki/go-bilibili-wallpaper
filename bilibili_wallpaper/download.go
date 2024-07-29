package bilibili_wallpaper

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"

	"github.com/cavaliergopher/grab/v3"
	"github.com/samber/lo"

	"github.com/AyakuraYuki/bilibili-wallpaper/internal/colors"
	"github.com/AyakuraYuki/bilibili-wallpaper/internal/file"
	"github.com/AyakuraYuki/bilibili-wallpaper/internal/filenamify"
	"github.com/AyakuraYuki/bilibili-wallpaper/internal/misc"
)

func loadTasks() (res []*Task) {
	res = make([]*Task, 0)
	data, err := readDataFile()
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
			fileExt := filepath.Ext(picture.ImgSrc)
			v := &Task{
				Url:      picture.ImgSrc,
				Filename: fmt.Sprintf("[%d] [%d] %s [%dx%d]%s", item.DocId, index, safetyFilename, picture.ImgWidth, picture.ImgHeight, fileExt),
			}
			verbosePrintln(colors.White("已加载任务 %q (%q)", v.Filename, v.Url))
			res = append(res, v)
		}
	}
	return res
}

func Download() {
	res := loadTasks()
	if len(res) == 0 {
		return
	}
	defer func() { _ = os.Remove(DataJsonFilePath) }()

	tasks := make([]*Task, 0)
	for _, item := range res {
		filename := item.Filename
		fullPath := path.Join(DistDir, filename)
		ok, err := file.IsPathExist(fullPath)
		if err != nil {
			log.Println(colors.Red("err: %v", err))
			continue
		}
		if ok {
			verbosePrintln(colors.White("已存在文件 %s，跳过", fullPath))
			continue
		}
		tasks = append(tasks, item)
	}

	tasksCount := len(tasks)
	verbosePrintln(colors.Green("有 %d 张壁纸需要下载", tasksCount))

	if Serial {
		// 串行下载

		counter := 0
		for _, task := range tasks {
			task := task
			fullname := task.Filename
			dst := path.Join(DistDir, fullname)
			if _, err := grab.Get(dst, task.Url); err == nil {
				counter++
				fmt.Printf("downloading [%v / %v] file: %q \r", counter, tasksCount, fullname)
			}
		}

	} else {
		// 并行下载

		size := tasksCount / Coroutines
		var partitions [][]*Task
		if size > 0 {
			partitions = lo.Chunk(tasks, size)
		} else {
			partitions = append(partitions, tasks)
		}

		funcs := make([]misc.WorkFunc, 0)
		var counter atomic.Int32
		for _, partition := range partitions {
			partition := partition
			funcs = append(funcs, func() error {
				for _, task := range partition {
					task := task
					fullname := task.Filename
					dst := filepath.Join(DistDir, fullname)
					if _, err := grab.Get(dst, task.Url); err == nil {
						fmt.Printf("downloading [%v / %v] file: %q \r", counter.Add(1), tasksCount, fullname)
					}
				}
				return nil
			})
		}
		if err := misc.MultiRun(funcs...); err != nil {
			log.Println(colors.Red("多线下载遇到异常: %v", err))
			return
		}

	}

	return
}
