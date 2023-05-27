package netease

import (
	"fmt"
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

const FILENAME_STYLE_1 = "<AUTONO>. <TITLE> - <ARTIST>.lrc"

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

func (musics NeteaseSingleMusics_t) fetchLrcsAsync() {
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(2, func(music_i interface{}) {
		if musics[cast.ToInt(music_i)].needDownload {
			time.Sleep(time.Millisecond * time.Duration(330))
			responseChannel <- fmt.Sprintf("开始下载第 %02d 首歌词", cast.ToInt(music_i))
			musics[cast.ToInt(music_i)].fetchLrc()
		}
		wg.Done()
	})
	defer p.Release()
	go responseController()
	startTime := time.Now()
	for i := range musics {
		wg.Add(1)
		_ = p.Invoke(i)
	}
	wg.Wait()
	elapsedTime := time.Since(startTime) / time.Millisecond
	fmt.Printf("下载歌词总耗时: %dms\n", elapsedTime)
}

func DownloadSingleMusicLrc(id int64) {
	cfg := Config_s{FileNameStyle: FILENAME_STYLE_1}
	DownloadSingleMusicLrcWCfg(id, cfg)
}

func DownloadSingleMusicLrcWCfg(id int64, config Config_s) {
	nsm := newNeteaseSingleMusic(id)
	nsm.applyConfig(&config)

	filename := config.FileNameStyle
	filename = strings.ReplaceAll(filename, "<TITLE>", nsm.title)
	filename = strings.ReplaceAll(filename, "<ARTIST>", nsm.getArtistsStr())
	err := os.WriteFile(filename, []byte(nsm.lyric.GetLyrics()), 0666)
	if err != nil {
		fmt.Println(err.Error())
	}
	/**
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		_, _ = os.Create(fileName)
	}
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	_, _ = f.WriteString(history)
	**/
}

func DownloadPlaylistLrc(id int64) {
	cfg := Config_s{FileNameStyle: FILENAME_STYLE_1}
	DownloadPlaylistLrcWCfg(id, cfg)
}

func DownloadPlaylistLrcWCfg(id int64, config Config_s) {
	np := newNeteasePlaylist(id)
	path := utils.ToSaveFilename(np.title)
	os.MkdirAll(path, os.ModePerm)
	np.musics.fetchLrcsAsync()
	for _, v := range np.musics {
		if !v.genlyric {
			//无需下载，continue
			continue
		}
		v.applyConfig(&config)

		filename := config.FileNameStyle
		filename = strings.ReplaceAll(filename, "<TITLE>", v.title)
		filename = strings.ReplaceAll(filename, "<AUTONO>", cast.ToString(v.listI+1))
		filename = strings.ReplaceAll(filename, "<ARTIST>", v.getArtistsStr())
		filename = utils.ToSaveFilename(filename)
		err := os.WriteFile(path+"/"+filename, []byte(v.lyric.GetLyrics()), 0666)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
