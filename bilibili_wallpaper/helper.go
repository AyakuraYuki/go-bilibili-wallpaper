package bilibili_wallpaper

import "log"

func verbosePrintln(a ...any) {
	if Verbose {
		log.Println(a...)
	}
}
