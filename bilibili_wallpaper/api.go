package bilibili_wallpaper

import (
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cast"

	"github.com/AyakuraYuki/bilibili-wallpaper/plugins/colors"
	"github.com/AyakuraYuki/bilibili-wallpaper/plugins/file"
	cjson "github.com/AyakuraYuki/bilibili-wallpaper/plugins/json"
	nhttp "github.com/AyakuraYuki/bilibili-wallpaper/plugins/net/http"
)

func getApi(pageNum int) string {
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

func GetInfo() {
	pageNum := 0
	for {
		api := getApi(pageNum)
		verbosePrintln(colors.White("请求壁纸列表: %s", api))

		headers := http.Header{}
		headers.Add("cookie", Cookie)

		bs, rspHeader, code, err := nhttp.GetRaw(nil, api, headers, nil, 5000, 3)
		if err != nil {
			log.Println(colors.Red("请求壁纸列表失败: %v", err))
			continue
		}
		verbosePrintln(colors.White("rsp => (%d) headers: %v", code, rspHeader))

		rsp := &Rsp{}
		if err = cjson.JSON.Unmarshal(bs, &rsp); err != nil {
			log.Println(colors.Red("解析壁纸列表响应结果失败: %v", err))
			continue
		}

		// EOF
		if rsp.Data == nil || len(rsp.Data.Items) == 0 {
			persistedRsp, err0 := readExistData()
			if err0 != nil {
				log.Println(colors.Red("读取壁纸列表暂存数据失败 %v", err0))
				return
			}
			writeDataJsonFile(persistedRsp)
			return
		}

		if pageNum == 0 {
			writeDataJsonFile(rsp)
		} else {
			persistedRsp, _ := readExistData()
			persistedRsp.Data.Items = append(persistedRsp.Data.Items, rsp.Data.Items...)
			persistedRsp.Code = rsp.Code
			persistedRsp.Message = rsp.Message
			persistedRsp.Ttl = rsp.Ttl
			writeDataJsonFile(persistedRsp)
		}

		pageNum++
	}
}

func readExistData() (*Rsp, error) {
	rsp := &Rsp{}
	content := file.ReadFile(DataJsonFilePath)
	err := cjson.JSON.UnmarshalFromString(content, &rsp)
	return rsp, err
}

func writeDataJsonFile(rsp *Rsp) {
	bs, _ := cjson.JSON.MarshalIndent(&rsp, "", "    ")
	content := string(bs)
	file.WriteLines(DataJsonFilePath, []string{content})
}
