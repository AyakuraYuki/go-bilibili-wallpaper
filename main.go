package main

import (
	"context"
	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/commands"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//func init() {
//	var err error
//	bilibili_wallpaper.Workdir, err = os.Getwd()
//	if err != nil {
//		bilibili_wallpaper.Workdir, _ = os.UserHomeDir()
//		bilibili_wallpaper.Workdir = path.Join(bilibili_wallpaper.Workdir, "bili-wallpaper")
//		log.Println(colors.Yellow("获取不到工作路径，切换到 %v", bilibili_wallpaper.Workdir))
//	}
//	bilibili_wallpaper.DistDir = path.Join(bilibili_wallpaper.Workdir, bilibili_wallpaper.WallpaperDir)
//	bilibili_wallpaper.DataJsonFilePath = path.Join(bilibili_wallpaper.Workdir, bilibili_wallpaper.JsonFile)
//
//	// parse flags
//	flag.StringVar(&bilibili_wallpaper.Cookie, "c", "", "bilibili 用户登录浏览器 cookie，可以通过浏览器开发者工具的控制台输入 document.cookie 获得")
//	flag.StringVar(&bilibili_wallpaper.Cookie, "cookie", "", "bilibili 用户登录浏览器 cookie，可以通过浏览器开发者工具的控制台输入 document.cookie 获得")
//	flag.BoolVar(&bilibili_wallpaper.Serial, "serial", false, "单线下载模式，如果默认的多线下载模式频繁出错，可以指定单线模式进行顺序下载")
//	flag.BoolVar(&bilibili_wallpaper.Verbose, "verbose", false, "调试模式，输出详细内容")
//
//	flag.Usage = func() {
//		w := flag.CommandLine.Output()
//		_, _ = fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
//		flag.PrintDefaults()
//		_, _ = fmt.Fprintf(w, "\n")
//	}
//}
//
//func main() {
//	flag.Parse()
//	if bilibili_wallpaper.Cookie == "" {
//		panic(errors.New("缺少 cookie 会拿不到完整的图片列表，停止爬取动作"))
//	}
//
//	_ = os.MkdirAll(bilibili_wallpaper.DistDir, os.ModePerm)
//	log.Println(colors.Green("工作路径 %s", bilibili_wallpaper.Workdir))
//	log.Println(colors.Green("下载路径 %s", bilibili_wallpaper.DistDir))
//
//	bilibili_wallpaper.RequestWallpaperData()
//	bilibili_wallpaper.Download()
//
//	fmt.Println("")
//	log.Println(colors.Green("完成"))
//}

func main() {
	ctx, cancel := notifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	commands.Cmd.Version = "2.0.0"
	if err := commands.Cmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}

// notifyContext 将信号绑定到上下文
func notifyContext(ctx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan os.Signal, 5)
	signal.Notify(ch, signals...)
	if ctx.Err() == nil {
		go func() {
			// 第一次取消上下文
			select {
			case <-ctx.Done():
			case <-ch:
				cancel()
			}
			// 第二次结束程序
			select {
			case <-ctx.Done():
			case <-ch:
				os.Exit(1)
			}
		}()
	}
	return ctx, cancel
}
