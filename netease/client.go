package netease

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/buger/jsonparser"
	"github.com/imroc/req/v3"
	"github.com/mdp/qrterminal/v3"
	cookiejar "github.com/orirawlings/persistent-cookiejar"
	"github.com/skip2/go-qrcode"
)

type NeteaseClient struct {
	*req.Client
	isLogged bool
}

var Client = newNeteaseClient()

func newNeteaseClient() *NeteaseClient {
	jar, err := cookiejar.New(&cookiejar.Options{
		Filename: "cookies_netease.json",
	})
	if err != nil {
		log.Fatalf("failed to create persistent cookiejar: %s\n", err.Error())
	}
	c := req.C().
		SetCookieJar(jar).
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
			if !resp.IsSuccessState() {
				return fmt.Errorf("bad response, raw dump:\n%s", resp.Dump())
			}
			code, err := jsonparser.GetInt(resp.Bytes(), "code")
			if err != nil {
				log.Panic(err)
			}
			if (code != 200) && (code < 800 || code > 803) {
				msg, _ := jsonparser.GetString(resp.Bytes(), "message")
				return fmt.Errorf("netease API error: %s", msg)
			}
			if code == 803 {
				err := jar.Save()
				if err != nil {
					log.Println(err)
				}
			}
			return nil
		})

	return &NeteaseClient{
		Client: c,
	}
}

func LoginGen(genQrFile bool) string {
	resp, err := Client.R().SetBodyString(`{"type":1}`).Post(`weapi/login/qrcode/unikey`)
	if err != nil {
		log.Println(err.Error())
	}
	unikey, err := jsonparser.GetString(resp.Bytes(), "unikey")
	if err != nil {
		log.Panic(err)
	}

	if genQrFile {
		_ = qrcode.WriteFile("https://music.163.com/login?codekey="+unikey, qrcode.Medium, 256, "qr.png")

	} else {
		qrterminal.GenerateWithConfig("https://music.163.com/login?codekey="+unikey, qrterminal.Config{
			Level:     qrterminal.L,
			Writer:    os.Stdout,
			BlackChar: qrterminal.BLACK,
			WhiteChar: qrterminal.WHITE,
			QuietZone: 1,
		})
	}
	return unikey
}

func LoginCheck(unikey string) (bool, string) {
	for i := 0; i < 20; i++ {
		time.Sleep(time.Duration(3) * time.Second)
		resp, err := Client.R().SetBodyString(`{"key":"` + unikey + `","type":1}`).Post(`weapi/login/qrcode/client/login`)
		if err != nil {
			log.Println(err.Error())
		}
		code, _ := jsonparser.GetInt(resp.Bytes(), "code")
		if code == 801 || code == 802 {
			fmt.Print("...")
		} else if code == 803 {
			return true, "登录成功"
		} else if code == 800 {
			return false, "用户拒绝登录"
		}
	}
	return false, "超时"
}

// 检测登录状态。用户名 (ID), 是否登录
func LoginStatus() (string, bool) {
	resp, err := Client.R().SetBodyString(`{}`).Post(`weapi/w/nuser/account/get`)
	if err != nil {
		log.Println(err.Error())
	}
	nickname, err := jsonparser.GetString(resp.Bytes(), "profile", "nickname")
	if err != nil {
		return "未登录", false
	}
	userId, err := jsonparser.GetInt(resp.Bytes(), "profile", "userId")
	if err != nil {
		return "未登录", false
	}
	return fmt.Sprintf("%s (%d)", nickname, userId), true
}
