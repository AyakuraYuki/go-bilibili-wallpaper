package bilibili_wallpaper

import jsoniter "github.com/json-iterator/go"

type Rsp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    *Data  `json:"data"`
}

type Data struct {
	Items []*Item `json:"items"`
}

type Item struct {
	DocId       int        `json:"doc_id"`
	PosterUid   int        `json:"poster_uid"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Pictures    []*Picture `json:"pictures"`
	Count       int        `json:"count"`
	Ctime       int        `json:"ctime"`
	View        int        `json:"view"`
	Like        int        `json:"like"`
	DynId       string     `json:"dyn_id"`
}

type Picture struct {
	ImgSrc    string              `json:"img_src"`
	ImgWidth  int                 `json:"img_width"`
	ImgHeight int                 `json:"img_height"`
	ImgSize   float64             `json:"img_size"`
	ImgTags   jsoniter.RawMessage `json:"img_tags"`
}

type UrlAndFilename struct {
	Url      string
	Filename string
}
