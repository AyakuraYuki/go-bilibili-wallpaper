package downloader

import (
	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/downloader/internal/colors"
	"github.com/spf13/cast"
	"log"
	"net/url"

	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/downloader/internal/encoding/json"
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
		if err = json.JSON.Unmarshal(bs, &apiRsp); err != nil {
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
