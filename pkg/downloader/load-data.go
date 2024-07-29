package downloader

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/spf13/cast"

	"github.com/AyakuraYuki/bilibili-wallpaper/internal/colors"
	"github.com/AyakuraYuki/bilibili-wallpaper/internal/file"
	"github.com/AyakuraYuki/bilibili-wallpaper/internal/filenamify"
	cjson "github.com/AyakuraYuki/bilibili-wallpaper/internal/json"
)

func (d *Downloader) requestDocList() {
	var (
		pageNum         = 0
		docs            []*Doc
		maxRetryPerPage = 0
		maxRetry        = 0
	)
	for {
		if maxRetryPerPage > 5 {
			pageNum++
			maxRetryPerPage = 0
			d.verboseLog(colors.White("第 %d 个页面已重试超过10次均失败了，跳过这一页"))
			continue
		}
		if maxRetry > 1000 {
			break // 避免无限循环
		}

		api := assembleApiUrl(pageNum)
		d.verboseLog(colors.White("请求壁纸列表: %q", api))

		rsp, err := d.client.R().Get(api)
		if err != nil {
			log.Println(colors.Red("请求壁纸列表失败: %v", err))
			maxRetryPerPage++
			maxRetry++
			continue
		}
		d.verboseLog(colors.White("(%d) headers: %v", rsp.StatusCode(), rsp.Header()))

		bs := rsp.Body()
		d.verboseLog(colors.White("response data: %s", string(bs)))
		apiRsp := &ApiResponse{}
		if err = cjson.JSON.Unmarshal(bs, &apiRsp); err != nil {
			log.Println(colors.Red("解析壁纸列表响应结果失败: %v", err))
			maxRetryPerPage++
			maxRetry++
			continue
		}

		// EOF
		if apiRsp.Data == nil || len(apiRsp.Data.Items) == 0 {
			break
		}

		docs = append(docs, apiRsp.Data.Items...)
		pageNum++
		maxRetryPerPage = 0
	}
	_ = d.writeDataFile(docs)
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

func (d *Downloader) readDataFile() (docs []*Doc, err error) {
	docs = make([]*Doc, 0)
	content := file.ReadFile(d.dataFilePath)
	err = cjson.JSON.UnmarshalFromString(content, &docs)
	if err != nil {
		log.Println(colors.Red("读取接口数据临时文件失败: %v", err))
	}
	return docs, err
}

func (d *Downloader) writeDataFile(docs []*Doc) error {
	bs, err := cjson.JSON.MarshalIndent(&docs, "", "    ")
	if err != nil {
		log.Println(colors.Red("保存接口数据到临时文件失败: %v", err))
		return err
	}
	file.WriteLines(d.dataFilePath, []string{string(bs)})
	return nil
}

func assembleApiUrl(pageNum int) string {
	u := url.URL{
		Scheme: "https",
		Host:   "api.vc.bilibili.com",
		Path:   "/link_draw/v1/doc/doc_list",
	}
	q := u.Query()
	q.Set("uid", "6823116")
	q.Set("page_num", cast.ToString(pageNum))
	q.Set("page_size", "20")
	q.Set("biz", "all")
	u.RawQuery = q.Encode()
	return u.String()
}
