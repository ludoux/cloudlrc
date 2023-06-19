package netease

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ludoux/cloudlrc/utils"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cast"
)

var responseChannel = make(chan string, 15)

type Config_s struct {
	TransFirst    bool
	OriDelayMs    int64
	TransDelayMs  int64
	FileNameStyle string
}

const FILENAME_STYLE_1 = "<AUTONO>. <TITLE>.lrc"
const FILENAME_STYLE_2 = "<AUTONO>. <TITLE> - <ARTIST>.lrc"

func (nsm *NeteaseSingleMusic_s) applyConfig(config *Config_s) {
	if !config.TransFirst && config.TransDelayMs != 0 {
		nsm.lyric.DelayLyricLine(0, config.TransDelayMs)
	} else if !config.TransFirst && config.OriDelayMs != 0 {
		nsm.lyric.DelayLyricLine(1, config.OriDelayMs)
	} else if config.TransFirst {
		nsm.ChangeTransOrder()
	}
}

func responseController() {
	for rc := range responseChannel {
		fmt.Println(rc)
	}
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
	cfg := Config_s{FileNameStyle: FILENAME_STYLE_1}
	DownloadSingleMusicLrcWCfg(id, cfg)
}

func DownloadSingleMusicLrcWCfg(id int64, config Config_s) {
	nsm := newNeteaseSingleMusic(id)
	if nsm.genlyric {
		nsm.applyConfig(&config)

		filename := config.FileNameStyle
		filename = strings.ReplaceAll(filename, "<TITLE>", nsm.title)
		filename = strings.ReplaceAll(filename, "<AUTONO>", "1")
		filename = strings.ReplaceAll(filename, "<ARTIST>", nsm.getArtistsStr())
		filename = utils.ToSaveFilename(filename)
		err := os.WriteFile(filename, []byte(nsm.lyric.GetLyrics()), 0666)
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
	cfg := Config_s{FileNameStyle: FILENAME_STYLE_1}
	DownloadPlaylistLrcWCfg(id, cfg)
}

func DownloadPlaylistLrcWCfg(id int64, config Config_s) {
	np := newNeteasePlaylist(id)
	path := utils.ToSaveFilename(np.title)
	os.MkdirAll(path, os.ModePerm)
	fmt.Println("音乐总数:", len(np.musics))

	np.musics.fetchLrcsAsync()

	aligncount := len([]rune(cast.ToString(len(np.musics))))
	aligntext := "%0" + cast.ToString(aligncount) + "d"
	for _, v := range np.musics {
		if !v.genlyric {
			//无需下载，continue
			continue
		}
		v.applyConfig(&config)

		filename := config.FileNameStyle
		filename = strings.ReplaceAll(filename, "<TITLE>", v.title)
		filename = strings.ReplaceAll(filename, "<AUTONO>", fmt.Sprintf(aligntext, v.listI+1))
		filename = strings.ReplaceAll(filename, "<ARTIST>", v.getArtistsStr())
		filename = utils.ToSaveFilename(filename)
		err := os.WriteFile(path+"/"+filename, []byte(v.lyric.GetLyrics()), 0666)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func DownloadAlbumLrc(id int64) {
	cfg := Config_s{FileNameStyle: FILENAME_STYLE_1}
	DownloadAlbumLrcWCfg(id, cfg)
}

func DownloadAlbumLrcWCfg(id int64, config Config_s) {
	np := newNeteaseAlbum(id)
	path := utils.ToSaveFilename(np.title)
	os.MkdirAll(path, os.ModePerm)
	fmt.Println("音乐总数:", len(np.musics))

	np.musics.fetchLrcsAsync()

	aligncount := len([]rune(cast.ToString(len(np.musics))))
	aligntext := "%0" + cast.ToString(aligncount) + "d"
	for _, v := range np.musics {
		if !v.genlyric {
			//无需下载，continue
			continue
		}
		v.applyConfig(&config)

		filename := config.FileNameStyle
		filename = strings.ReplaceAll(filename, "<TITLE>", v.title)
		filename = strings.ReplaceAll(filename, "<AUTONO>", fmt.Sprintf(aligntext, v.listI+1))
		filename = strings.ReplaceAll(filename, "<ARTIST>", v.getArtistsStr())
		filename = utils.ToSaveFilename(filename)
		err := os.WriteFile(path+"/"+filename, []byte(v.lyric.GetLyrics()), 0666)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
