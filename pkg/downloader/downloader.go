package downloader

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"

	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/downloader/internal/colors"
	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/downloader/internal/encoding/json"
	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/downloader/internal/file"
)

const (
	dataFilename     = "data_list.json"
	wallpaperDirname = "images"
	coroutines       = 10
)

type Option func(*Downloader)

func Cookie(cookie string) Option { return func(d *Downloader) { d.cookie = cookie } }
func Dist(absPath string) Option  { return func(d *Downloader) { d.distDir = absPath } }
func Verbose() Option             { return func(d *Downloader) { d.verbose = true } }
func DisableAsync() Option        { return func(d *Downloader) { d.blocking = true } }

type Downloader struct {
	cookie       string
	workDir      string
	distDir      string
	dataFilePath string
	verbose      bool
	blocking     bool
	client       *resty.Client
}

func New(opts ...Option) (*Downloader, error) {
	downloader := &Downloader{}
	for _, o := range opts {
		o(downloader)
	}

	if downloader.cookie == "" {
		return nil, errors.New("缺少 cookie 会拿不到完整的图片列表，本程序停止工作")
	}

	var err error
	downloader.workDir, err = os.Getwd()
	if err != nil {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("尝试获取本程序的工作目录失败: %v", err)
		}
		downloader.workDir = filepath.Join(home, "bilibili-wallpaper")
	}

	if downloader.distDir == "" {
		downloader.distDir = filepath.Join(downloader.workDir, wallpaperDirname)
	}

	downloader.dataFilePath = filepath.Join(downloader.workDir, dataFilename)

	downloader.client = resty.New()
	downloader.client.SetRetryCount(3)
	downloader.client.SetHeader("Cookie", downloader.cookie)

	_ = os.MkdirAll(downloader.distDir, os.ModePerm)
	log.Println(colors.Green("工作路径: %q", downloader.workDir))
	log.Println(colors.Green("保存路径: %q", downloader.distDir))

	return downloader, nil
}

func (d *Downloader) readDataFile() (docs []*Doc, err error) {
	docs = make([]*Doc, 0)
	content := file.ReadString(d.dataFilePath)
	err = json.JSON.UnmarshalFromString(content, &docs)
	if err != nil {
		log.Println(colors.Red("读取接口数据临时文件失败: %v", err))
	}
	return docs, err
}

func (d *Downloader) writeDataFile(docs []*Doc) error {
	bs, err := json.JSON.MarshalIndent(&docs, "", "    ")
	if err != nil {
		log.Println(colors.Red("保存接口数据到临时文件失败: %v", err))
		return err
	}
	file.SaveString(d.dataFilePath, string(bs))
	return nil
}

func (d *Downloader) verboseLog(a ...any) {
	if d.verbose {
		log.Println(a...)
	}
}
