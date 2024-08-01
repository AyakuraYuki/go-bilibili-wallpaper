package downloader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/cavaliergopher/grab/v3"
	"github.com/samber/lo"

	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/downloader/internal/colors"
	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/downloader/internal/file"
	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/downloader/internal/filenamify"
	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/downloader/internal/misc"
)

func (d *Downloader) Download() {
	// 请求最新的壁纸列表
	d.requestDocList()
	// 加载任务
	res := d.loadTasks()
	if len(res) == 0 {
		log.Println(colors.Yellow("没有需要下载的壁纸"))
		return
	}
	defer func() { _ = os.Remove(d.dataFilePath) }()

	tasks := make([]*Task, 0)
	for _, item := range res {
		ok, err := file.Exist(item.FullPath)
		if err != nil {
			d.verboseLog(colors.Red("err: %v", err))
			continue
		}
		if ok {
			d.verboseLog(colors.White("跳过已存在的壁纸: %q", item.FullPath))
			continue
		}
		tasks = append(tasks, item)
	}

	tasksCount := len(tasks)
	d.verboseLog(colors.Green("有 %d 张壁纸需要下载", tasksCount))

	if d.blocking {
		// 一次下载一张

		counter := 0
		for _, task := range tasks {
			task := task
			if _, err := grab.Get(task.FullPath, task.Url); err == nil {
				counter++
				fmt.Printf("downloading [%v / %v] file: %q \r", counter, tasksCount, task.FullPath)
			}
		}

	} else {
		// 并行下载

		size := tasksCount / coroutines
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
					if _, err := grab.Get(task.FullPath, task.Url); err == nil {
						fmt.Printf("downloading [%v / %v] file: %q \r", counter.Add(1), tasksCount, task.FullPath)
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
}

func (d *Downloader) loadTasks() (tasks []*Task) {
	tasks = make([]*Task, 0)

	docs, err := d.readDataFile()
	if err != nil {
		return tasks
	}

	for _, doc := range docs {
		if len(doc.Pictures) == 0 {
			continue
		}
		description := fmt.Sprintf("%s-%s", doc.Title, doc.Description)
		filename, _ := filenamify.FilenamifyV2(description)
		for index, picture := range doc.Pictures {
			ext := filepath.Ext(picture.ImgSrc)
			taskFilename := fmt.Sprintf("[%d] [%d] %s [%dx%d]%s", doc.DocId, index, filename, picture.ImgWidth, picture.ImgHeight, ext)
			task := &Task{
				Filename: taskFilename,
				FullPath: filepath.Join(d.distDir, taskFilename),
				Url:      picture.ImgSrc,
			}
			d.verboseLog(colors.White("已加载任务 %q (%q)", task.Filename, task.Url))
			tasks = append(tasks, task)
		}
	}
	return tasks
}
