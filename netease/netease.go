package netease

import (
	"fmt"
	"os"
	"strings"
)

type Config_s struct {
	TransFirst    bool
	OriDelayMs    int64
	TransDelayMs  int64
	FileNameStyle string
}

const FILENAME_STYLE_1 = "<TITLE> - <ARTIST>.lrc"

func (nsm *NeteaseSingleMusic_s) applyConfig(config *Config_s) {
	if config.TransFirst == false && config.TransDelayMs != 0 {
		nsm.lyric.DelayLyricLine(0, config.TransDelayMs)
	} else if config.TransFirst == false && config.OriDelayMs != 0 {
		nsm.lyric.DelayLyricLine(1, config.OriDelayMs)
	} else if config.TransFirst {
		nsm.ChangeTransOrder()
	}
}

func (musics NeteaseSingleMusics_t) fetchLrcsAsync() {

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
	os.MkdirAll(np.title, os.ModePerm)

	for _, v := range np.musics {
		v.applyConfig(&config)

		filename := fmt.Sprintf("%s/%s", np.title, config.FileNameStyle)
		filename = strings.ReplaceAll(filename, "<TITLE>", v.title)
		filename = strings.ReplaceAll(filename, "<ARTIST>", v.getArtistsStr())
		err := os.WriteFile(filename, []byte(v.lyric.GetLyrics()), 0666)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
