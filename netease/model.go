package netease

import (
	"log"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/ludoux/cloudlrc/lyric"
	"github.com/spf13/cast"
)

// 单曲相关
type NeteaseSingleMusic_s struct {
	no      int
	id      int64
	title   string
	artists []string
	album   string
	retry   int
	status  int
	msg     string
	lyric   *lyric.Lyric_s
}

func NewNeteaseSingleMusic(id int64) *NeteaseSingleMusic_s {
	rt := NeteaseSingleMusic_s{id: id}
	rt.lyric = lyric.NewLyric()
	rt.fetch()
	rt.fetchLrc()
	return &rt
}
func NewNeteaseSingleMusicNofetch(id int64) *NeteaseSingleMusic_s {
	rt := NeteaseSingleMusic_s{id: id}
	rt.lyric = lyric.NewLyric()
	return &rt
}

func (it *NeteaseSingleMusic_s) fetch() {
	resp, err := Client.R().Get(`api/v3/song/detail?c=[{"id":"` + cast.ToString(it.id) + `"}]`)
	if err != nil {
		log.Println(err.Error())
	}
	//log.Println(resp)
	value, err := jsonparser.GetString(resp.Bytes(), "songs", "[0]", "name")
	if err != nil {
		log.Panic(err)
	}
	it.title = value

	valueInt, err := jsonparser.GetInt(resp.Bytes(), "songs", "[0]", "no")
	if err != nil {
		log.Panic(err)
	}
	it.no = cast.ToInt(valueInt)

	value, err = jsonparser.GetString(resp.Bytes(), "songs", "[0]", "al", "name")
	if err != nil {
		log.Panic(err)
	}
	it.album = value

	jsonparser.ArrayEach(resp.Bytes(), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		artist, err := jsonparser.GetString(value, "name")
		if err != nil {
			log.Panic(err)
		}
		it.artists = append(it.artists, artist)
	}, "songs", "[0]", "ar")
}

func (it *NeteaseSingleMusic_s) fetchLrc() {
	//获取原文歌词
	resp, err := Client.R().Get(`api/song/media?id=` + cast.ToString(it.id))
	if err != nil {
		log.Println(err.Error())
	}
	//TODO: 分析歌词状态 https://github.com/ludoux/LrcHelper/blob/223eaf8b3dc11f13ccc61371a8a222729f402aef/LrcHelper/NeteaseMusic.cs#L93

	if resp.String() == `{"code":200}` {
		it.lyric.LyricMsg = "未提供歌词"
		return
	}

	nolyric, err := jsonparser.GetBoolean(resp.Bytes(), "nolyric")
	if err == nil && nolyric {
		//比如 纯音乐
		it.lyric.LyricMsg = "无需歌词"
		return
	}

	value, err := jsonparser.GetString(resp.Bytes(), "lyric")
	if err != nil {
		it.lyric.LyricMsg = "无歌词"
		return
	} else {
		it.lyric.LyricMsg = "有歌词"
	}

	value = strings.ReplaceAll(value, "\\r", "")
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\\n", "\n")

	txtLines := strings.Split(value, "\n")
	for _, val := range txtLines {
		//原文歌词的优先级更高
		it.lyric.AppendLyricTextLine(val, 1)
	}

	//===后面翻译
	resp, err = Client.R().Get(`api/song/lyric?os=pc&id=` + cast.ToString(it.id) + `&tv=-1`)
	if err != nil {
		log.Println(err.Error())
	}
	value, err = jsonparser.GetString(resp.Bytes(), "tlyric", "lyric")
	if err != nil && value == "" {
		it.lyric.LyricMsg += ",无翻译"
		return
	}
	it.lyric.LyricMsg += ",有翻译"

	value = strings.ReplaceAll(value, "\\r", "")
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\\n", "\n")

	txtLines = strings.Split(value, "\n")
	for _, val := range txtLines {
		//译文优先级较低
		it.lyric.AppendLyricTextLine(val, 0)
	}
}

// 调换原文和翻译的优先级
func (it *NeteaseSingleMusic_s) ChangeTransOrder() {
	it.lyric.SwapPriority(0, 1)
}

// 专辑相关
type NeteaseSingleMusics_t []*NeteaseSingleMusic_s
type NeteaseAlbum_s struct {
	id     int64
	title  string
	musics NeteaseSingleMusics_t
}

func NewNeteaseAlbum(id int64) *NeteaseAlbum_s {
	rt := NeteaseAlbum_s{id: id}
	rt.fetch()
	return &rt
}

func (it *NeteaseAlbum_s) fetch() {
	resp, err := Client.R().Get(`api/album/` + cast.ToString(it.id))
	if err != nil {
		log.Println(err.Error())
	}
	value, err := jsonparser.GetString(resp.Bytes(), "album", "name")
	if err != nil {
		log.Panic(err)
	}
	it.title = value

	//遍历内部的歌曲
	jsonparser.ArrayEach(resp.Bytes(), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		musicId, err := jsonparser.GetInt(value, "id")
		if err != nil {
			log.Panic(err)
		}
		newMusic := NewNeteaseSingleMusicNofetch(musicId)
		name, err := jsonparser.GetString(value, "name")
		if err != nil {
			log.Panic(err)
		}
		newMusic.title = name

		valueInt, err := jsonparser.GetInt(value, "no")
		if err != nil {
			log.Panic(err)
		}
		newMusic.no = cast.ToInt(valueInt)

		newMusic.album = it.title

		jsonparser.ArrayEach(value, func(inValue []byte, dataType jsonparser.ValueType, offset int, err error) {
			artist, err := jsonparser.GetString(inValue, "name")
			if err != nil {
				log.Panic(err)
			}
			newMusic.artists = append(newMusic.artists, artist)
		}, "artists")

		it.musics = append(it.musics, newMusic)
	}, "album", "songs")
}
