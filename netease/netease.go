package netease

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ludoux/cloudlrc/cfg"
	"github.com/ludoux/cloudlrc/utils"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cast"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var responseChannel = make(chan string, 15)

func responseController() {
	for rc := range responseChannel {
		fmt.Println(rc)
	}
}

func convertEncoding(s string, newEncoding string) []byte {
	if strings.ToLower(newEncoding) == "utf8" || strings.ToLower(newEncoding) == "utf-8" {
		return []byte(s)
	}
	if strings.ToLower(newEncoding) == "gbk" || strings.ToLower(newEncoding) == "gb2312" {
		reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewEncoder())
		b, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Panicln(newEncoding, "字符编码转换失败:", err.Error())
		}
		return b
	}
	log.Fatalln("字符编码转换失败:", newEncoding)
	return []byte("")
}

/**
 * @description: Album 和 Playlist 都会调用此方法来获取歌词
 * @return {*}
 */
func (musics NeteaseSingleMusics_t) fetchLrcsAsync() {
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(2, func(music_i interface{}) {
		i := cast.ToInt(music_i)
		if musics[i].needDownload {
			time.Sleep(time.Millisecond * time.Duration(330))
			err := musics[i].fetchLrc()
			if err != nil {
				responseChannel <- fmt.Sprintf("第 %d 首<%s>(%d): 发生错误: %s。", i+1, musics[i].title, musics[i].id, err.Error())
			} else {
				if musics[i].genlyric {
					responseChannel <- fmt.Sprintf("第 %d 首<%s>(%d): 下载成功。", i+1, musics[i].title, musics[i].id)
				} else {
					responseChannel <- fmt.Sprintf("第 %d 首<%s>(%d): 无歌词。", i+1, musics[i].title, musics[i].id)
				}
			}

		}
		wg.Done()
	})
	defer p.Release()
	go responseController()
	//startTime := time.Now()
	for i := range musics {
		wg.Add(1)
		_ = p.Invoke(i)
	}
	wg.Wait()
	//elapsedTime := time.Since(startTime) / time.Millisecond
	//fmt.Printf("下载歌词总耗时: %dms\n", elapsedTime)
}

func DownloadSingleMusicLrc(id int64) {
	cfg := cfg.GetCfgFile()
	DownloadSingleMusicLrcWCfg(id, cfg)
}

func DownloadSingleMusicLrcWCfg(id int64, config *cfg.CfgFile_s) {
	nsm := newNeteaseSingleMusic(id)
	if nsm.genlyric {

		aligntext_trackno := "%02d"

		filename := config.Filename + "." + strings.ToLower(config.Format)
		filename = strings.ReplaceAll(filename, "<TITLE>", nsm.title)
		filename = strings.ReplaceAll(filename, "<AUTO_NO>", "1")
		filename = strings.ReplaceAll(filename, "<DISC_NO>", cast.ToString(nsm.discNo))
		filename = strings.ReplaceAll(filename, "<TRACK_NO>", fmt.Sprintf(aligntext_trackno, nsm.trackNo))
		filename = strings.ReplaceAll(filename, "<ARTIST>", nsm.getArtistsStr())
		filename = utils.ToSaveFilename(filename)
		err := os.WriteFile(filename, convertEncoding(nsm.lyric.GetLyrics(), config.Encoding), 0666)
		if err != nil {
			log.Fatalln(err.Error())
		} else {
			fmt.Printf("已生成: %s\n", filename)
		}
	} else {
		fmt.Println("无歌词以供生成")
	}

}

func DownloadPlaylistLrc(id int64) {
	cfg := cfg.GetCfgFile()
	DownloadPlaylistLrcWCfg(id, cfg)
}

func DownloadPlaylistLrcWCfg(id int64, config *cfg.CfgFile_s) {
	np := newNeteasePlaylist(id)
	path := utils.ToSaveFilename(np.title)
	os.MkdirAll(path, os.ModePerm)
	fmt.Println("音乐总数:", len(np.musics))

	np.musics.fetchLrcsAsync()

	aligncount_autono := len([]rune(cast.ToString(len(np.musics))))
	if aligncount_autono == 1 {
		aligncount_autono = 2
	}
	aligntext_autono := "%0" + cast.ToString(aligncount_autono) + "d"

	aligncount_trackno := len([]rune(cast.ToString(len(np.musics))))
	if aligncount_trackno == 1 {
		aligncount_trackno = 2
	}
	aligntext_trackno := "%02d"

	for _, v := range np.musics {
		if !v.genlyric {
			//无需下载，continue
			continue
		}

		filename := config.Filename + "." + strings.ToLower(config.Format)
		filename = strings.ReplaceAll(filename, "<TITLE>", v.title)
		filename = strings.ReplaceAll(filename, "<AUTO_NO>", fmt.Sprintf(aligntext_autono, v.listI+1))
		filename = strings.ReplaceAll(filename, "<DISC_NO>", cast.ToString(v.discNo))
		filename = strings.ReplaceAll(filename, "<TRACK_NO>", fmt.Sprintf(aligntext_trackno, v.trackNo))
		filename = strings.ReplaceAll(filename, "<ARTIST>", v.getArtistsStr())
		filename = utils.ToSaveFilename(filename)
		err := os.WriteFile(path+"/"+filename, convertEncoding(v.lyric.GetLyrics(), config.Encoding), 0666)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func DownloadAlbumLrc(id int64) {
	cfg := cfg.GetCfgFile()
	DownloadAlbumLrcWCfg(id, cfg)
}

func DownloadAlbumLrcWCfg(id int64, config *cfg.CfgFile_s) {
	np := newNeteaseAlbum(id)
	path := utils.ToSaveFilename(np.title)
	os.MkdirAll(path, os.ModePerm)
	fmt.Println("音乐总数:", len(np.musics))

	np.musics.fetchLrcsAsync()

	aligncount_autono := len([]rune(cast.ToString(len(np.musics))))
	if aligncount_autono == 1 {
		aligncount_autono = 2
	}
	aligntext_autono := "%0" + cast.ToString(aligncount_autono) + "d"

	aligncount_trackno := len([]rune(cast.ToString(len(np.musics))))
	if aligncount_trackno == 1 {
		aligncount_trackno = 2
	}
	aligntext_trackno := "%02d"
	for _, v := range np.musics {
		if !v.genlyric {
			//无需下载，continue
			continue
		}

		filename := config.Filename + "." + strings.ToLower(config.Format)
		filename = strings.ReplaceAll(filename, "<TITLE>", v.title)
		filename = strings.ReplaceAll(filename, "<AUTO_NO>", fmt.Sprintf(aligntext_autono, v.listI+1))
		filename = strings.ReplaceAll(filename, "<DISC_NO>", cast.ToString(v.discNo))
		filename = strings.ReplaceAll(filename, "<TRACK_NO>", fmt.Sprintf(aligntext_trackno, v.trackNo))
		filename = strings.ReplaceAll(filename, "<ARTIST>", v.getArtistsStr())
		filename = utils.ToSaveFilename(filename)
		err := os.WriteFile(path+"/"+filename, convertEncoding(v.lyric.GetLyrics(), config.Encoding), 0666)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
