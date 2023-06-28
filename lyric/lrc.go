package lyric

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/ludoux/cloudlrc/cfg"
	"github.com/spf13/cast"
)

type LyricLine_s struct {
	time_hr  int
	time_min int
	time_sec int
	time_ms  int
	text     string
	is_tag   bool
	priority int //越大，同时轴越前
	forceH   bool
	force2ms bool
}

type LyricLines_t []*LyricLine_s

type LyricConfig_s struct {
	CombineSameTimeline bool
	Split               string
	SkipEmpty           bool
	TimelineForceH      bool
	TimelineForce2ms    bool
}

type Lyric_s struct {
	//LyricStatus int
	LyricMsg   string
	lyricLines LyricLines_t
	lyricCfg   cfg.CfgLrc_s
}

func newLyricLine() *LyricLine_s {
	rt := LyricLine_s{}
	rt.is_tag = false
	rt.priority = 0
	rt.forceH = false
	rt.force2ms = false
	return &rt
}

func NewLyric(cfg *cfg.CfgLrc_s) *Lyric_s {
	rt := Lyric_s{}
	rt.lyricCfg = *cfg
	return &rt
}

func (it *LyricLine_s) getPriority() int {
	return it.priority
}

func (it *LyricLine_s) setPriority(value int) error {
	it.priority = value
	return nil
}

// [ar:artistname]
func (it *LyricLine_s) setTag(value string) error {
	reTag := regexp.MustCompile(`\[[a-zA-Z]+?:.+?\]`)
	if reTag.MatchString(value) {
		//is tag,ex [ar:artistname]
		it.is_tag = true
		it.text = value
		return nil
	}
	return errors.New("no tag: " + value)
}

func (it *LyricLine_s) setText(value string) error {
	it.is_tag = false
	it.text = value
	return nil
}

// [1:2:3] [2:3] [3] [00:01:03.990]
func (it *LyricLine_s) setTimeline(value string) error {
	reType1 := regexp.MustCompile(`\[(\d+):(\d+):(\d+)(\.(\d+))?\]`) // [1:2:3(.990)]
	reType2 := regexp.MustCompile(`\[(\d+):(\d+)(\.(\d+))?\]`)       // [2:3(.990)]
	reType3 := regexp.MustCompile(`\[(\d+)(\.(\d+))?\]`)             // [3(.990)]
	submatch := reType1.FindSubmatch([]byte(value))
	if len(submatch) > 0 {
		it.time_hr = cast.ToInt(string(submatch[1]))
		it.time_min = cast.ToInt(string(submatch[2]))
		it.time_sec = cast.ToInt(string(submatch[3]))
		if len(submatch[5]) > 0 {
			it.time_ms = cast.ToInt(string(submatch[5]))
		}
		if it.force2ms && it.time_ms > 99 {
			it.time_ms /= 10
		}
		return nil
	}

	submatch = reType2.FindSubmatch([]byte(value))
	if len(submatch) > 0 {
		it.time_hr = 0
		it.time_min = cast.ToInt(string(submatch[1]))
		it.time_sec = cast.ToInt(string(submatch[2]))
		if len(submatch[4]) > 0 {
			it.time_ms = cast.ToInt(string(submatch[4]))
		}
		if it.force2ms && it.time_ms > 99 {
			it.time_ms /= 10
		}
		return nil
	}

	submatch = reType3.FindSubmatch([]byte(value))
	if len(submatch) > 0 {
		it.time_hr = 0
		it.time_min = 0
		it.time_sec = cast.ToInt(string(submatch[1]))
		if len(submatch[3]) > 0 {
			it.time_ms = cast.ToInt(string(submatch[3]))
		}
		if it.force2ms && it.time_ms > 99 {
			it.time_ms /= 10
		}
		return nil
	}

	return errors.New("setTimeline failed. no match regex: " + value)
}

func (it *LyricLine_s) setTimeAndText(value string, priority int) error {
	/*Support value:
	[1:2:3]hello [2:3]hello [0:1:3]hello [00:01:03]hello [3]hello [00:01:03.990]hello
	[ar:artistname]
	不支持同行多个时间轴
	*/
	it.setPriority(priority)

	err := it.setTag(value)
	if err == nil {
		//is [ar:artistname] type
		return nil
	}
	reTimeline := regexp.MustCompile(`\[[0-9:.]+?\]`)
	allTimeline := reTimeline.FindAll([]byte(value), -1)
	//check timeline count, only support 1 timeline per line
	if len(allTimeline) != 1 {
		return errors.New("timeline count err: " + cast.ToString(len(allTimeline)))
	}

	err = it.setTimeline(string(allTimeline[0]))
	if err != nil {
		return err
	}
	lrcText := strings.Replace(value, string(allTimeline[0]), "", 1)
	it.setText(lrcText)

	return nil
}
func (it *LyricLine_s) GetTimeStringWithBrackets() string {
	return fmt.Sprintf("[%s]", it.GetTimeStringNoBrackets())

}
func (it *LyricLine_s) GetTimeStringNoBrackets() string {
	var build strings.Builder
	if it.forceH || it.time_hr > 0 {
		build.WriteString(fmt.Sprintf("%d:", it.time_hr))
	}
	build.WriteString(fmt.Sprintf("%02d:%02d", it.time_min, it.time_sec))
	build.WriteString(fmt.Sprintf(".%.2d", it.time_ms))
	return build.String()
}

func (it *LyricLine_s) getTimelineInMs() int64 {
	total := cast.ToInt64(it.time_ms) + cast.ToInt64(it.time_sec)*1000 + cast.ToInt64(it.time_min)*60000 + cast.ToInt64(it.time_hr)*3600000
	return total
}

func divmod64(x int64, y int64) (int, int) {
	return cast.ToInt(x / y), cast.ToInt(x % y)
}

func divmod(x int, y int) (int, int) {
	return x / y, x % y
}

// 给定总ms,重新分配timeline属性
func (it *LyricLine_s) setTimelineFromMs(ms int64) {
	it.time_sec, it.time_ms = divmod64(ms, 1000)
	it.time_min, it.time_sec = divmod(it.time_sec, 60)
	it.time_hr, it.time_min = divmod(it.time_min, 60)
}

// 输入 ms 单位的偏移量（如1000指延后1s
func (it *LyricLine_s) adjustTimelineMs(ms int64) {
	it.setTimelineFromMs(it.getTimelineInMs() + ms)
}

func (it LyricLines_t) Len() int {
	return len(it)
}

func (it LyricLines_t) Less(i, j int) bool {
	//歌词在前面的值大，false
	//若a<b,返回true,指a在后出现
	a := it[i]
	b := it[j]
	aInMs := a.getTimelineInMs()
	bInMs := b.getTimelineInMs()

	if a.is_tag && !b.is_tag {
		return false
	} else if !a.is_tag && b.is_tag {
		return true
	} else if a.is_tag && b.is_tag {
		if a.priority > b.priority {
			return false
		} else if a.priority == b.priority {
			return false //相等
		} else {
			return true
		}
	}
	if aInMs == bInMs && a.priority == b.priority {
		return false //相等
	} else if aInMs < bInMs || (aInMs == bInMs && a.priority > b.priority) {
		return false
	} else {
		return true
	}
}

func (it LyricLines_t) Swap(i, j int) {
	it[i], it[j] = it[j], it[i]
}

func (it *LyricLine_s) GetLyricLine() string {
	if it.is_tag {
		return it.text
	} else {
		return fmt.Sprintf("%s%s", it.GetTimeStringWithBrackets(), it.text)
	}
}

func (it *Lyric_s) AppendLyricTextLine(value string, priority int) {
	ll := newLyricLine()
	ll.force2ms = it.lyricCfg.TimelineForceFixMs
	ll.forceH = it.lyricCfg.TimelineForceHour
	ll.setTimeAndText(value, priority)

	it.lyricLines = append(it.lyricLines, ll)
}

func (it *Lyric_s) GetLyrics() string {
	//sort.Sort(it.lyricLines)
	sort.Sort(sort.Reverse(it.lyricLines))
	var build strings.Builder
	preTimeline := ""
	storedText := ""
	for _, val := range it.lyricLines {
		if it.lyricCfg.SkipEmpty && (val.text == "" || strings.ReplaceAll(val.text, " ", "") == "") {
			continue
		}
		if it.lyricCfg.Style == 3 {
			if preTimeline != val.GetTimeStringWithBrackets() {
				//新的开始
				if preTimeline != "" {
					//前者有数据，先输出
					build.WriteString(preTimeline)
					build.WriteString(storedText)
					build.WriteString("\n")
				}
				//存储新的开始
				preTimeline = val.GetTimeStringWithBrackets()
				storedText = val.text
			} else {
				//当前为同timeline下的其它行
				storedText = storedText + it.lyricCfg.Split + val.text
			}
		} else {
			build.WriteString(val.GetLyricLine())
			build.WriteString("\n")
		}
	}
	//导出最后一组
	if it.lyricCfg.Style == 3 {
		build.WriteString(preTimeline)
		build.WriteString(storedText)
		build.WriteString("\n")
	}
	return build.String()
}

/**
 * @description: 为指定priority的歌词行调整
 * @param {int} priority
 * @param {int64} ms
 * @return {*}
 */
func (it *Lyric_s) DelayLyricLine(priority int, ms int64) {
	for _, val := range it.lyricLines {
		if val.getPriority() == priority {
			val.adjustTimelineMs(ms)
		}
	}
}

func (it *Lyric_s) SwapPriority(a, b int) {
	for _, val := range it.lyricLines {
		if val.getPriority() == a {
			val.setPriority(b)
		} else if val.getPriority() == b {
			val.setPriority(a)
		}
	}
}
