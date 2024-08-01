package commands

import (
	"github.com/spf13/cobra"

	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/downloader"
)

const (
	binName = "bilibili-wallpaper-downloader"
)

var (
	cookie  string
	dist    string
	serial  bool
	verbose bool
)

var Cmd = &cobra.Command{
	Use:   binName,
	Short: "从B站的 壁纸喵(https://space.bilibili.com/6823116) 账号下载壁纸的工具",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := []downloader.Option{
			downloader.Cookie(cookie),
		}
		if dist != "" {
			opts = append(opts, downloader.Dist(dist))
		}
		if serial {
			opts = append(opts, downloader.DisableAsync())
		}
		if verbose {
			opts = append(opts, downloader.Verbose())
		}
		dl, err := downloader.New(opts...)
		if err != nil {
			return err
		}
		dl.Download()
		return nil
	},
}

func init() {
	Cmd.PersistentFlags().StringVarP(&cookie, "cookie", "c", "", "bilibili 用户登录浏览器 cookie，可以通过浏览器开发者工具的控制台输入 document.cookie 获得")
	Cmd.PersistentFlags().StringVarP(&dist, "dist-dir", "o", "", "希望保存壁纸的文件夹的绝对路径。本程序默认下载壁纸到程序工作目录的 images 文件夹")
	Cmd.PersistentFlags().BoolVarP(&serial, "serial", "b", false, "单线下载模式，如果默认的多线下载模式频繁出错，可以指定单线模式进行顺序下载")
	Cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "调试模式，输出详细内容")
}
