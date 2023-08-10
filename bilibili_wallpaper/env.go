package bilibili_wallpaper

var (
	Cookie           string
	Workdir          string
	DistDir          string
	DataJsonFilePath string
	Verbose          = false
	Serial           = false
)

const (
	JsonFile     = "data_list.json"
	WallpaperDir = "images"
	Coroutines   = 10
)
