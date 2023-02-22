package netease

import (
	"log"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/ludoux/cloudlrc/lyric"
	"github.com/spf13/cast"
)

type NeteaseSingleMusic_s struct {
	no      int
	id      int64
	titile  string
	artists []string
	album   string
	retry   int
	status  int
	msg     string
	lyric   *lyric.Lyric_s
}

func (it *NeteaseSingleMusic_s) fetch() {
	resp, err := Client.R().Get(`v3/song/detail?c=[{"id":"` + cast.ToString(it.id) + `"}]`)
	if err != nil {
		log.Println(err.Error())
	}
	//log.Println(resp)
	value, err := jsonparser.GetString(resp.Bytes(), "songs", "[0]", "name")
	if err != nil {
		log.Panic(err)
	}
	it.titile = value

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
	resp, err := Client.R().Get(`song/media?id=` + cast.ToString(it.id))
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
		it.lyric.AppendLyricTextLine(val, 1)
	}

	//===后面翻译
	resp, err = Client.R().Get(`song/lyric?os=pc&id=` + cast.ToString(it.id) + `&tv=-1`)
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
		it.lyric.AppendLyricTextLine(val, 0)
	}
}
func NewNeteaseSingleMusic(id int64) *NeteaseSingleMusic_s {
	rt := NeteaseSingleMusic_s{id: id}
	rt.lyric = lyric.NewLyric()
	rt.fetch()
	rt.fetchLrc()
	return &rt
} /*
func NewNeteaseSingleMusicNofetch(id int64) *NeteaseSingleMusic_s {

}*/
