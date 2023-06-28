package cfg

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type CfgFile_s struct {
	Format   string
	Filename string
	Encoding string
}

type CfgLrc_s struct {
	Style              int
	Delayms            int
	Split              string
	SkipEmpty          bool
	TimelineForceHour  bool
	TimelineForceFixMs bool
}

func GetCfgFile() *CfgFile_s {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	rt := CfgFile_s{}
	sec, err := cfg.GetSection("file")
	if err != nil {
		fmt.Printf("Fail to read section lrc: %v", err)
		os.Exit(1)
	}
	rt.Format = sec.Key("format").MustString("lrc")
	rt.Filename = sec.Key("filename").MustString("<AUTO_NO>. <TITLE>")
	rt.Encoding = sec.Key("encoding").MustString("utf-8")
	return &rt
}

func GetCfgLrc() *CfgLrc_s {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	rt := CfgLrc_s{}
	sec, err := cfg.GetSection("lrc")
	if err != nil {
		fmt.Printf("Fail to read section lrc: %v", err)
		os.Exit(1)
	}
	rt.Style = sec.Key("style").RangeInt(1, 0, 3)
	rt.Delayms = sec.Key("style").MustInt()
	rt.Split = sec.Key("split").MustString(" ")
	rt.SkipEmpty = sec.Key("skip_empty").MustBool(true)
	rt.TimelineForceHour = sec.Key("timeline_force_hour").MustBool(false)
	rt.TimelineForceFixMs = sec.Key("timeline_force_fix_ms").MustBool(false)
	return &rt
}
