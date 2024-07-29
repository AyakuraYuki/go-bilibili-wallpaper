package downloader

import (
	jsoniter "github.com/json-iterator/go"
)

type ApiResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Ttl     int          `json:"ttl"`
	Data    *DocListData `json:"data"`
}

type DocListData struct {
	Items []*Doc `json:"items"`
}

type Doc struct {
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

type Task struct {
	Filename string
	FullPath string
	Url      string
}
