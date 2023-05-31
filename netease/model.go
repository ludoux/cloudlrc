package netease

import (
	"fmt"
	"log"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/ludoux/cloudlrc/lyric"
	"github.com/spf13/cast"
)

// 单曲相关
type NeteaseSingleMusic_s struct {
	no           int //CD中的顺序
	id           int64
	title        string
	artists      []string
	album        string
	retry        int
	listI        int //列表中的顺序，从0开始
	needDownload bool
	status       int
	msg          string
	lyric        *lyric.Lyric_s
	genlyric     bool //是否需要生成歌词文件
}

/**
 * @description: 信息和歌词均获得
 * @param {int64} id
 * @return {*}
 */
func newNeteaseSingleMusic(id int64) *NeteaseSingleMusic_s {
	rt := NeteaseSingleMusic_s{id: id, needDownload: true}
	rt.lyric = lyric.NewLyric()
	rt.fetch()
	rt.fetchLrc()
	return &rt
}
func NewNeteaseSingleMusicNofetch(id int64) *NeteaseSingleMusic_s {
	rt := NeteaseSingleMusic_s{id: id, needDownload: true}
	rt.lyric = lyric.NewLyric()
	return &rt
}

func (it *NeteaseSingleMusic_s) singleMusicAnalyze(resp string) {
	bytes := []byte(resp)
	value, err := jsonparser.GetString(bytes, "name")
	if err != nil {
		log.Panic(err)
	}
	it.title = value

	valueInt, err := jsonparser.GetInt(bytes, "no")
	if err != nil {
		log.Panic(err)
	}
	it.no = cast.ToInt(valueInt)

	value, err = jsonparser.GetString(bytes, "al", "name")
	if err != nil {
		log.Panic(err)
	}
	it.album = value

	jsonparser.ArrayEach(bytes, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		artist, err := jsonparser.GetString(value, "name")
		if err != nil {
			log.Panic(err)
		}
		it.artists = append(it.artists, artist)
	}, "ar")
}

func (it *NeteaseSingleMusic_s) getArtistsStr() string {
	return strings.Join(it.artists, ",")
}

func (it *NeteaseSingleMusic_s) fetch() {
	resp, err := Client.R().Get(`api/v3/song/detail?c=[{"id":"` + cast.ToString(it.id) + `"}]`)
	if err != nil {
		log.Println(err.Error())
	}
	//log.Println(resp)
	obj, _, _, _ := jsonparser.Get(resp.Bytes(), "songs", "[0]")
	it.singleMusicAnalyze(string(obj))

}

func (it *NeteaseSingleMusic_s) fetchLrc() error {
	it.genlyric = false
	resp, err := Client.R().Get(`api/song/lyric?os=pc&id=` + cast.ToString(it.id) + `&lv=-1&tv=-1`)
	if err != nil {
		return fmt.Errorf("获取服务端歌词错误: %s", err.Error())
	}
	if strings.Contains(resp.String(), "纯音乐，请欣赏") {
		it.lyric.LyricMsg = "无需歌词"
		return nil
	} else {
		sgc, _ := jsonparser.GetBoolean(resp.Bytes(), "sgc")
		sfy, _ := jsonparser.GetBoolean(resp.Bytes(), "sfy")
		qfy, _ := jsonparser.GetBoolean(resp.Bytes(), "qfy")
		if sgc && sfy && qfy {
			it.lyric.LyricMsg = "未提供歌词"
			return nil
		}
	}
	//原文歌词
	value, err := jsonparser.GetString(resp.Bytes(), "lrc", "lyric")
	if err != nil {
		it.lyric.LyricMsg = "无歌词"
		return nil
	} else {
		it.lyric.LyricMsg = "有歌词"
	}
	it.genlyric = true
	value = strings.ReplaceAll(value, "\\r", "")
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\\n", "\n")

	txtLines := strings.Split(value, "\n")
	for _, val := range txtLines {
		//原文歌词的优先级更高
		it.lyric.AppendLyricTextLine(val, 1)
	}

	//===后面翻译

	value, err = jsonparser.GetString(resp.Bytes(), "tlyric", "lyric")
	if err != nil && value == "" {
		it.lyric.LyricMsg += ",无翻译"
		return nil
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
	return nil
}

// 调换原文和翻译的优先级
func (it *NeteaseSingleMusic_s) ChangeTransOrder() {
	it.lyric.SwapPriority(0, 1)
}

// 列表大相关
type NeteaseSingleMusics_t []*NeteaseSingleMusic_s

// 专辑相关
type NeteaseAlbum_s struct {
	id     int64
	title  string
	musics NeteaseSingleMusics_t
}

func newNeteaseAlbum(id int64) *NeteaseAlbum_s {
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

	listI := 0
	//遍历内部的歌曲
	jsonparser.ArrayEach(resp.Bytes(), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		musicId, err := jsonparser.GetInt(value, "id")
		if err != nil {
			log.Panic(err)
		}
		newMusic := NewNeteaseSingleMusicNofetch(musicId)
		newMusic.listI = listI
		listI = listI + 1
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

// 歌单相关
type NeteasePlaylist_s struct {
	id     int64
	title  string
	musics NeteaseSingleMusics_t
}

/**
 * @description: 仅获得信息，不含内部歌曲的歌词
 * @param {int64} id
 * @return {*}
 */
func newNeteasePlaylist(id int64) *NeteasePlaylist_s {
	rt := NeteasePlaylist_s{id: id}
	rt.fetch()
	rt.fetchMusicDetail()
	return &rt
}

func (it *NeteasePlaylist_s) fetch() {
	resp, err := Client.R().Get(`api/v6/playlist/detail?id=` + cast.ToString(it.id) + `&c=[{"id":"` + cast.ToString(it.id) + `"}]`)
	if err != nil {
		log.Println(err.Error())
	}
	value, err := jsonparser.GetString(resp.Bytes(), "playlist", "name")
	if err != nil {
		log.Panic(err)
	}
	it.title = value

	listI := 0
	//遍历内部的歌曲
	jsonparser.ArrayEach(resp.Bytes(), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		musicId, err := jsonparser.GetInt(value, "id")
		if err != nil {
			log.Panic(err)
		}
		newMusic := NewNeteaseSingleMusicNofetch(musicId)
		newMusic.listI = listI
		listI = listI + 1
		//Only ID

		it.musics = append(it.musics, newMusic)
	}, "playlist", "trackIds")
}

/**
 * @description: 对内部的每个100个音乐ID合并发起请求，得到详细信息（不含歌词）
 * @return {*}
 */
func (it *NeteasePlaylist_s) fetchMusicDetail() {
	var build strings.Builder

	i := 0
	var listIList []int
	for _, val := range it.musics {
		if val.needDownload {
			listIList = append(listIList, val.listI)
			i = i + 1
			if i == 1 {
				tmp := fmt.Sprintf(`{"id":"%d"}`, val.id)
				build.WriteString(tmp)
			} else {
				tmp := fmt.Sprintf(`,{"id":"%d"}`, val.id)
				build.WriteString(tmp)
			}
			if i == 100 || i == len(it.musics) {
				resp, err := Client.R().Get(fmt.Sprintf(`api/v3/song/detail?c=[%s]`, build.String()))
				if err != nil {
					log.Println(err.Error())
				}
				j := 0

				jsonparser.ArrayEach(resp.Bytes(), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					musicId, _ := jsonparser.GetInt(value, "id")
					if it.musics[listIList[j]].id == musicId {
						it.musics[listIList[j]].singleMusicAnalyze(string(value))
						j = j + 1
					} else {
						log.Println("Error! Index not match")
					}
				}, "songs")
				break
			}
		}
	}
}
