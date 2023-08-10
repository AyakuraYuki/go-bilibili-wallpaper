package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/AyakuraYuki/bilibili-wallpaper/bilibili_wallpaper"
	"log"
	"os"
	"path"
)

func init() {
	var err error
	bilibili_wallpaper.Workdir, err = os.Getwd()
	if err != nil {
		bilibili_wallpaper.Workdir, _ = os.UserHomeDir()
		bilibili_wallpaper.Workdir = path.Join(bilibili_wallpaper.Workdir, "bili-wallpaper")
		log.Printf("获取不到工作路径，切换到 %v", bilibili_wallpaper.Workdir)
	}
	bilibili_wallpaper.DistDir = path.Join(bilibili_wallpaper.Workdir, bilibili_wallpaper.WallpaperDir)
	bilibili_wallpaper.DataJsonFilePath = path.Join(bilibili_wallpaper.Workdir, bilibili_wallpaper.JsonFile)

	// parse flags
	flag.StringVar(&bilibili_wallpaper.Cookie, "c", "", "bilibili 用户登录浏览器 cookie，可以通过浏览器开发者工具的控制台输入 document.cookie 获得")
	flag.StringVar(&bilibili_wallpaper.Cookie, "cookie", "", "bilibili 用户登录浏览器 cookie，可以通过浏览器开发者工具的控制台输入 document.cookie 获得")

	flag.BoolVar(&bilibili_wallpaper.MultiThread, "multi-thread", false, "多线下载模式")
	flag.BoolVar(&bilibili_wallpaper.Verbose, "verbose", false, "调试模式，输出详细内容")

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(w, "\n")
	}
}

func main() {
	flag.Parse()
	if bilibili_wallpaper.Cookie == "" {
		panic(errors.New("缺少 cookie 会拿不到完整的图片列表，停止爬取动作"))
	}

	_ = os.MkdirAll(bilibili_wallpaper.DistDir, os.ModePerm)
	log.Printf("工作路径 %s", bilibili_wallpaper.Workdir)
	log.Printf("下载路径 %s", bilibili_wallpaper.DistDir)

	bilibili_wallpaper.GetInfo()
	bilibili_wallpaper.Download()

	log.Println("完成")
}
