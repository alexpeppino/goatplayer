package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	// "os/exec"
	// "reflect"
)

type nUtils struct{}

var utils nUtils // namespace

type nDuration struct{}

var dura nDuration // namespace

// width: width of the first
func (nDuration) FormatSeconds(totSeconds uint32, width int) string {
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

func (nDuration) FormatWithCents(hundredsSeconds uint32) string {
	var h, m, s uint32
	var hundreds = hundredsSeconds % 100
	h, m, s = dura.GetHMS((hundredsSeconds - hundreds) / 100)
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%02d", h, m, s, hundreds)
	} else {
		return fmt.Sprintf("%02d:%02d.%02d", m, s, hundreds)
	}
}

func (nDuration) GetHMS(s uint32) (uint32, uint32, uint32) {
	var h, m uint32
	h = s / 3600
	s -= h * 3600
	m = s / 60
	s -= m * 60
	return h, m, s
}

// func (nDuration) FormatPartialAndTotal(partial, total uint32) (string, uint32) {
// 	if total == 0 {
// 		return "", 0
// 	}
// 	var s string
// 	if Set.ShowRemainingTime.Value {
// 		partial = total - partial
// 		s = "-"
// 	}
// 	h1, m1, s1 := dura.GetHMS(partial)
// 	h2, m2, s2 := dura.GetHMS(total)
// 	perc := partial * 100 / total
// 	if h2 > 0 {
// 		return fmt.Sprintf(" %s%02d:%02d:%02d / %02d:%02d:%02d (%d%%) ", s, h1, m1, s1, h2, m2, s2, perc), perc
// 	} else {
// 		return fmt.Sprintf(" %s%02d:%02d / %02d:%02d (%d%%) ", s, m1, s1, m2, s2, perc), perc
// 	}
// }

func (nDuration) FormatPartialAndTotalBoth(partial, total uint32, prefix string) (string, uint32) {
	if total == 0 {
		return "", 0
	}
	h1, m1, s1 := dura.GetHMS(partial)
	h2, m2, s2 := dura.GetHMS(total)
	h3, m3, s3 := dura.GetHMS(total - partial)
	perc := min(partial*100/total, 99)

	if h2 > 0 {
		return fmt.Sprintf("  %s %d:%02d:%02d / %d:%02d:%02d  (%d%%)  -%d:%02d:%02d ",
			prefix, h1, m1, s1, h2, m2, s2, perc, h3, m3, s3), perc
	} else {
		return fmt.Sprintf("  %s %d:%02d / %d:%02d  (%d%%)  -%d:%02d ",
			prefix, m1, s1, m2, s2, perc, m3, s3), perc
	}
}

func (nUtils) Bigstring(s string) string {
	var bigStart rune = 33
	var bigEnd rune = 127
	var bigDiff rune = 65248
	s = strings.ReplaceAll(s, " ", "  ")
	var space string = " "
	var b = space[0]
	var rr []rune
	rr = make([]rune, len(s))
	for i := range len(s) {
		c := s[i]
		if c == b {
			rr[i] = rune(b)
		} else if rune(c) < bigStart || rune(c) > bigEnd {
			rr[i] = rune(bigDiff) + bigStart + 2
		} else {
			rr[i] = rune(c) + bigDiff
		}
	}
	return string(rr)
}

// func (nUtils) Bigstring(s string) string {
// 	// s = strings.ToUpper(s)
// 	var space string = " "
// 	var b = space[0]
// 	var rr = make([]rune, len(s)*2)
// 	for i := range len(s) {
// 		c := s[i]
// 		if c == b {
// 			rr[i*2] = rune(b)
// 		} else {
// 			rr[i*2] = rune(c) + 65281
// 		}
// 		rr[i*2+1] = rune(b)
// 	}
// 	return string(rr)
// }

// .
func RuneLen(s *string) uint16 {
	return uint16(len([]rune(*s)))
}

func (nUtils) Sleep(delay uint16) {
	time.Sleep(time.Second * time.Duration(delay))
}

func (nUtils) SleepMilli(n time.Duration) {
	time.Sleep(time.Millisecond * n)
}

// Prints a slice of strings line by line.
func PrintSlice(lines *[]string) {
	for i, line := range *lines {
		fmt.Printf("%3d  %s\n", i, line)
	}
}

var err error

// // Give a struct var reference and the field index. Get the field's name.
// func refl_GetFieldName(x any, i int) string {

// 	return reflect.ValueOf(x).Elem().Type().Field(i).Name
// }

// // Give reference to var. Get var type's name.
// func refl_GetTypeName(x any) string {

// 	return reflect.ValueOf(x).Type().String()
// }

// func refl_GetType(x any) reflect.Type {

// 	return reflect.ValueOf(x).Type()
// }

func catch(err error) {
	if err != nil {
		fmt.Println("ERROR", err)
	}
}

func AdjustWidth(str *string, finalLen uint16) {
	var le = uint16(len([]rune(*str)))
	// var le = uint16(len(*str))
	if le > finalLen {
		*str = string([]rune(*str)[0:finalLen])
	} else {
		*str = fmt.Sprintf("%-*s", finalLen, *str)
	}
}

// Aligns 'str' to the left and 'right' to the right.
func AdjustWidthLeftAndRight(str *string, right string, finalLen uint16) {
	var leStr = uint16(len([]rune(*str)))
	var leRight = uint16(len([]rune(right)))
	// var le = uint16(len(*str))
	if finalLen <= leStr+leRight {
		*str = string([]rune(*str)[0:finalLen-leRight-1]) + string([]rune(right))
	} else {
		*str = fmt.Sprintf("%-*s%s ", finalLen-leRight-1, *str, right)
	}
}

// Aligns left of right.
func AdjustWidthAligned(str *string, finalLen uint16, align uint16) {
	var le = uint16(len([]rune(*str)))
	if le > finalLen {
		*str = string([]rune(*str)[0:finalLen])
	} else {
		if align == ALIGN_LEFT {
			*str = fmt.Sprintf("%-*s", finalLen, *str)
		} else {
			*str = fmt.Sprintf("%*s", finalLen, *str)
		}
	}
}

func (nUtils) s2u(s string) uint16 {
	i, _ := strconv.Atoi(s)
	return uint16(i)
}

func speedFormat() {
	// top := MainScreen.DisplaySection.marginTop
	// left := MainScreen.DisplaySection.marginLeft
	ui_pl.Speed = fmt.Sprintf("%.2f", float32(ui_pl.Speed_i)/100)
	Display.AddLineToPrint(1, 3, "Speed:   "+ui_pl.Speed)
	Term.Flush()
}

func speedPlus(x uint16) {
	ui_pl.Speed_i += x
	speedFormat()
}

func speed1() {
	ui_pl.Speed_i = 100
	speedFormat()
}

func speedMinus(x uint16) {
	ui_pl.Speed_i -= x
	speedFormat()
}

func volumeFormat() {
	ui_pl.Volume = fmt.Sprintf("%d", ui_pl.Volume_i)
	Display.AddLineToPrint(1, 4, "")
	// DisplaySection.AddLineToPrint(1, 4, "Volume:  "+ui_pl.Volume)
	Term.Flush()
}

func volume100() {
	ui_pl.Volume_i = 100
	volumeFormat()
}

func volumePlus(x uint16) {
	ui_pl.Volume_i += x
	volumeFormat()
}

func volumeMinus(x uint16) {
	ui_pl.Volume_i -= x
	volumeFormat()
}

func (nUtils) str2uint(s string) uint16 {
	var i uint16
	fmt.Sscanf(s, "%d", &i)
	return i
}

func (nUtils) str2uint32(s string) uint32 {
	var i uint32
	fmt.Sscanf(s, "%d", &i)
	return i
}

func (nUtils) uint2str(i uint16) string {
	return fmt.Sprintf("%d", i)
}

func (nUtils) int2str(i int) string {
	return fmt.Sprintf("%d", i)
}

// { Quit }
func Quit() {
	if Term.doSaveAtQuit {
		PFListsData.SaveDataToFile()
		Set.SaveSettings()
		Mpv.Stop()
	}
	os.Remove("/tmp/echo.sock")
	Input.MouseOff()
	TputCNorm()
	Term.SttyEchoOn()
	// Term.Clear()
	if Term.isAltScreenOn {
		Term.AltScreenOff()
	}
	// fmt.Print("\x1b[0;0H")
	if AF.doSortAndSaveAtQuit {
		fmt.Println("Sorting and saving the database...")
		AF.SortAudiofiles(SortAudiofilesAATT)
		AF.SaveDataToFile()
	}
	fmt.Println()
	fmt.Println("Bye. Thank you for using GoatPlayer.")
	os.Exit(0)
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func GetSortedUnique(sl *[]string, results *[]string) int {
	var aMap map[string]int
	aMap = make(map[string]int)
	for _, st := range *sl {
		aMap[st]++
	}
	for st := range aMap {
		*results = append(*results, st)
	}
	sort.Strings(*results)
	return len(*results)
	// trace.BeginEnd("DB.GetFieldSortedUnique, %d filenames, found %d audiof", len(*filenames), len(*results))
}

// func print_key(s2 string) {

// 	//	Print key
// 	var key string
// 	if v, ok := KeypadMap[s2]; ok {
// 		key = v
// 	} else {
// 		key = s2
// 	}
// 	s3 := fmt.Sprintf("%4s", key)
// 	Screen.Add(10, FileContainers.marginRight+23, Pen.StatusLine+s3+Pen.Normal)
// 	Screen.Print()
// }

