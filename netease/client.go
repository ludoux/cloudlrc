package netease

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/imroc/req/v3"
	"github.com/mdp/qrterminal/v3"
	cookiejar "github.com/orirawlings/persistent-cookiejar"
	"github.com/skip2/go-qrcode"
)

const (
	deviceIdFile = "device_id.json"
)

type deviceConfig_s struct {
	appver      string
	versioncode string
	resolution  string
	os          string
	channel     string
	brand       string
	mobilename  string
	ua          string
}

type NeteaseClient struct {
	*req.Client
	isLogged bool
}

var (
	deviceConfigs = []deviceConfig_s{
		{
			// 红米k50
			appver:      "9.4.70",
			versioncode: "9004070",
			resolution:  "3096x1440",
			os:          "android",
			channel:     "xiaomi",
			brand:       "Redmi",
			mobilename:  "22041211AC",
			ua:          "NeteaseMusic/9.4.70.260305110614(9004070);Dalvik/2.1.0 (Linux; U; Android 14; 22041211AC Build/UP1A.231005.007)",
		},
	}
)

var (
	deviceId = getOrGenerateDeviceId()
	Client   = newNeteaseClient()
)

func getOrGenerateDeviceId() string {
	content, err := os.ReadFile(deviceIdFile)
	if err == nil {
		var storedDeviceId string
		if json.Unmarshal(content, &storedDeviceId) == nil && storedDeviceId != "" {
			return storedDeviceId
		}
	}

	newDeviceId := generateRandomDeviceId()
	data, _ := json.Marshal(newDeviceId)
	os.WriteFile(deviceIdFile, data, 0644)
	return newDeviceId
}

func generateRandomDeviceId() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 52)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}

func getRandomDeviceConfig() deviceConfig_s {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(deviceConfigs))))
	return deviceConfigs[n.Int64()]
}

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
		SetBaseURL("https://interface3.music.163.com").
		SetCommonHeader("Referer", "https://interface3.music.163.com").
		SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 Edg/96.0.1054.62").
		SetTimeout(time.Duration(5) * time.Second).

		// EnableDump at the request level in request middleware which dump content into
		// memory (not print to stdout), we can record dump content only when unexpected
		// exception occurs, it is helpful to troubleshoot problems in production.
		OnBeforeRequest(func(c *req.Client, r *req.Request) error {
			if r.RetryAttempt > 0 { // Ignore on retry.
				return nil
			}
			if len(r.RawURL) >= 5 && r.RawURL[0:5] == "weapi" {
				params, encSecKey, encErr := Encrypt(string(r.Body))
				if encErr != nil {
					log.Println(encErr)
				}
				r.SetFormData(map[string]string{
					"params":    params,
					"encSecKey": encSecKey,
				})
			} else if len(r.RawURL) >= 4 && r.RawURL[0:4] == "eapi" {
				magicStr := "36cd479b6b5"
				deviceConfig := getRandomDeviceConfig()
				path := "/" + r.RawURL[1:]
				data := string(r.Body)
				text := fmt.Sprintf("nobody%suse%smd5forencrypt", path, data)
				md5Sum := md5.Sum([]byte(text))
				md5str := fmt.Sprintf("%x", md5Sum)
				eapiData := fmt.Sprintf("%s-%s-%s-%s-%s", path, magicStr, data, magicStr, md5str)
				params := eapiEncrypt(eapiData)
				r.SetFormData(map[string]string{
					"params": params,
				})
				r.SetCookies(
					&http.Cookie{
						Name:  "buildver",
						Value: strconv.FormatInt(time.Now().Unix(), 10)[0:10],
					},
					&http.Cookie{
						Name:  "deviceId",
						Value: deviceId,
					},
					&http.Cookie{
						Name:  "appver",
						Value: deviceConfig.appver,
					},
					&http.Cookie{
						Name:  "versioncode",
						Value: deviceConfig.versioncode,
					},
					&http.Cookie{
						Name:  "resolution",
						Value: deviceConfig.resolution,
					},
					&http.Cookie{
						Name:  "os",
						Value: deviceConfig.os,
					},
					&http.Cookie{
						Name:  "channel",
						Value: deviceConfig.channel,
					},
					&http.Cookie{
						Name:  "brand",
						Value: deviceConfig.brand,
					},
					&http.Cookie{
						Name:  "mobilename",
						Value: deviceConfig.mobilename,
					},
				)
				r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
				r.SetHeader("User-Agent", deviceConfig.ua)

			}
			return nil
		}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			// Corner case: neither an error response nor a success response,
			// dump content to help troubleshoot.
			if !resp.IsSuccessState() {
				return fmt.Errorf("网络请求失败！错误信息: %s", resp.Err.Error())
			}
			// // 有些接口返回不加密
			// if len(resp.Request.RawURL) >= 4 && resp.Request.RawURL[0:4] == "eapi" {
			// 	fmt.Println(resp.String())
			// 	decrypted, err := eapiDecrypt(resp.String())
			// 	if err != nil {
			// 		return fmt.Errorf("eapi 解密失败: %s", err.Error())
			// 	}
			// 	resp.SetBody([]byte(decrypted))
			// }
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
	resp, err := Client.R().SetBodyString(`{"type":3}`).Post(`eapi/login/qrcode/unikey`)
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
	for i := 0; i < 10; i++ {

		resp, err := Client.R().SetBodyString(`{"key":"` + unikey + `","type":3}`).Post(`eapi/login/qrcode/client/login`)
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
		time.Sleep(time.Duration(3) * time.Second)
	}
	return false, "超时"
}

// 检测登录状态。用户名 (ID), 是否登录
func LoginStatus() (string, bool) {
	resp, err := Client.R().SetBodyString(`{}`).Post(`eapi/w/nuser/account/get`)
	if err != nil {
		log.Fatalln(err.Error())
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
