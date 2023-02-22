package lyric

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cast"
)

type LyricLine_s struct {
	time_hr  int
	time_min int
	time_sec int
	time_ms  int
	text     string
	is_tag   bool
	priority int
}

type LyricLines_t []*LyricLine_s

type Lyric_s struct {
	//LyricStatus int
	LyricMsg   string
	lyricLines LyricLines_t
}

func NewLyricLine() *LyricLine_s {
	rt := LyricLine_s{}
	rt.is_tag = false
	rt.priority = 0
	return &rt
}

func NewLyric() *Lyric_s {
	rt := Lyric_s{}
	return &rt
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
	return fmt.Sprintf("[%d:%2d:%2d.%.3d]", it.time_hr, it.time_min, it.time_sec, it.time_ms)

}
func (it *LyricLine_s) GetTimeStringNoBrackets() string {
	return fmt.Sprintf("%d:%2d:%2d.%.3d", it.time_hr, it.time_min, it.time_sec, it.time_ms)
}

func (it *LyricLine_s) getTimeLineInMs() int64 {
	total := cast.ToInt64(it.time_ms) + cast.ToInt64(it.time_sec)*1000 + cast.ToInt64(it.time_min)*60000 + cast.ToInt64(it.time_hr)*3600000
	return total
}

func (it LyricLines_t) Len() int {
	return len(it)
}

func (it LyricLines_t) Less(i, j int) bool {
	//歌词在前面的值大，false
	//若a<b,返回true,指a在后出现
	a := it[i]
	b := it[j]
	aInMs := a.getTimeLineInMs()
	bInMs := b.getTimeLineInMs()

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
		return fmt.Sprintf("[%d:%02d:%02d.%.3d]%s", it.time_hr, it.time_min, it.time_sec, it.time_ms, it.text)
	}
}

func (it *Lyric_s) AppendLyricTextLine(value string, priority int) {
	ll := NewLyricLine()
	ll.setTimeAndText(value, priority)

	it.lyricLines = append(it.lyricLines, ll)
}

func (it *Lyric_s) GetLyrics() string {
	//sort.Sort(it.lyricLines)
	sort.Sort(sort.Reverse(it.lyricLines))
	var build strings.Builder

	for _, val := range it.lyricLines {
		build.WriteString(val.GetLyricLine())
		build.WriteString("\n")
	}
	return build.String()
}
