# go-bilibili-wallpaper

从B站的[壁纸喵](https://space.bilibili.com/6823116)账号下载同步壁纸的工具

## 使用方式

从 [release](https://github.com/AyakuraYuki/go-bilibili-wallpaper/releases) 页面下载对应系统的二进制文件包，解压到你想要同步壁纸的路径。这个程序会把壁纸同步到跟它同一个目录内的 `images` 文件夹里，重复执行则会下载新增的壁纸。（也就是说，如果不删除`images`文件夹，并且保留文件夹里的壁纸，那么这个程序就会进行增量下载）

```text
从B站的 壁纸喵(https://space.bilibili.com/6823116) 账号下载壁纸的工具

Usage:
  bilibili-wallpaper-downloader [flags]

Flags:
  -c, --cookie string     bilibili 用户登录浏览器 cookie，可以通过浏览器开发者工具的控制台输入 document.cookie 获得
  -o, --dist-dir string   希望保存壁纸的文件夹的绝对路径。本程序默认下载壁纸到程序工作目录的 images 文件夹
  -b, --serial            单线下载模式，如果默认的多线下载模式频繁出错，可以指定单线模式进行顺序下载
  -v, --verbose           调试模式，输出详细内容
  -h, --help              help for bilibili-wallpaper-downloader
      --version           version for bilibili-wallpaper-downloader

```
