package utils

import "strings"

func ToSaveFilename(in string) string {
	//https://stackoverflow.com/questions/1976007/what-characters-are-forbidden-in-windows-and-linux-directory-names
	rp := strings.NewReplacer(
		"/", " ",
		"\\", " ",
		"<", " ",
		">", " ",
		":", " ",
		"\"", " ",
		"|", " ",
		"?", " ",
		"*", " ",
	)
	rt := rp.Replace(in)
	return rt
}
