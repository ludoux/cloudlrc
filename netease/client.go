package netease

import (
	"fmt"
	"log"

	"github.com/buger/jsonparser"
	"github.com/imroc/req/v3"
)

type NeteaseClient struct {
	*req.Client
	isLogged bool
}

var Client = NewNeteaseClient()

func NewNeteaseClient() *NeteaseClient {
	c := req.C().
		//SetCommonHeader("Accept-Language", "zh-CN, zh-TW, en-US").
		//SetCommonHeader("Accept", "application/json").
		//SetCommonHeader("Content-Type", "application/json").
		SetBaseURL("https://music.163.com").
		SetCommonHeader("Referer", "https://music.163.com").
		SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 Edg/96.0.1054.62").

		// EnableDump at the request level in request middleware which dump content into
		// memory (not print to stdout), we can record dump content only when unexpected
		// exception occurs, it is helpful to troubleshoot problems in production.
		OnBeforeRequest(func(c *req.Client, r *req.Request) error {
			if r.RetryAttempt > 0 { // Ignore on retry.
				return nil
			}
			if r.RawURL[0:5] == "weapi" {
				//需要加密
				fmt.Println(string(r.Body))
				params, encSecKey, encErr := Encrypt(string(r.Body))
				if encErr != nil {
					log.Println(encErr)
				}
				r.SetFormData(map[string]string{
					"params":    params,
					"encSecKey": encSecKey,
				})
			}
			//r.EnableDump()
			return nil
		}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {

			// Corner case: neither an error response nor a success response,
			// dump content to help troubleshoot.
			if !resp.IsSuccess() {
				return fmt.Errorf("bad response, raw dump:\n%s", resp.Dump())
			}
			code, err := jsonparser.GetInt(resp.Bytes(), "code")
			if err != nil {
				log.Panic(err)
			}
			if code != 200 {
				msg, _ := jsonparser.GetString(resp.Bytes(), "message")
				return fmt.Errorf("Netease API Error: %s", msg)
			}
			return nil
		})

	return &NeteaseClient{
		Client: c,
	}
}

func Demo2() {
	resp, err := Client.R().Get(`api/v3/song/detail?c=[{"id":"426881480"},{"id":"426881487"}]`)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(resp)
}

func Demo3(id int64) {
	nsm := NewNeteaseSingleMusic(id)
	//nsm.lyric.DelayLyricLine(0, 500)
	nsm.ChangeTransOrder()
	fmt.Print(nsm.lyric.GetLyrics())
}
func Demo5() {
	resp, err := Client.R().SetBodyString(`{"type":1}`).Post(`weapi/login/qrcode/unikey`)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(resp)
}
