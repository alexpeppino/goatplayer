package dura

import (
	"fmt"
)

// width: width of the first
func FormatSeconds(totSeconds uint32, width int) string {
	var seconds, minutes, hours uint32
	hours = totSeconds / 3600
	totSeconds -= hours * 3600
	minutes = totSeconds / 60
	totSeconds -= minutes * 60
	seconds = totSeconds
	if hours == 0 {
		return fmt.Sprintf("%*d:%02d", width, minutes, seconds)
	} else {
		return fmt.Sprintf("%*dh%02d", width, hours, minutes)
	}
}

func FormatWithCents(hundredsSeconds uint32) string {
	var h, m, s uint32
	var hundreds = hundredsSeconds % 100
	h, m, s = GetHMS((hundredsSeconds - hundreds) / 100)
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%02d", h, m, s, hundreds)
	} else {
		return fmt.Sprintf("%02d:%02d.%02d", m, s, hundreds)
	}
}

func GetHMS(s uint32) (uint32, uint32, uint32) {
	var h, m uint32
	h = s / 3600
	s -= h * 3600
	m = s / 60
	s -= m * 60
	return h, m, s
}

func FormatPartialAndTotalBoth(partial, total uint32, prefix string) (string, uint32) {
	if total == 0 {
		return "", 0
	}
	h1, m1, s1 := GetHMS(partial)
	h2, m2, s2 := GetHMS(total)
	h3, m3, s3 := GetHMS(total - partial)
	perc := min(partial*100/total, 99)

	if h2 > 0 {
		return fmt.Sprintf("  %s %d:%02d:%02d / %d:%02d:%02d  (%d%%)  -%d:%02d:%02d ",
			prefix, h1, m1, s1, h2, m2, s2, perc, h3, m3, s3), perc
	} else {
		return fmt.Sprintf("  %s %d:%02d / %d:%02d  (%d%%)  -%d:%02d ",
			prefix, m1, s1, m2, s2, perc, m3, s3), perc
	}
}

