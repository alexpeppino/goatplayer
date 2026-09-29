package utils

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

func RuneLen(s *string) uint16 {
	return uint16(len([]rune(*s)))
}

func AdjustWidthSeparated(str, em *string, finalLen uint16) {
	var le = RuneLen(str)
	// var le = uint16(len(*str))
	if le > finalLen {
		*str = string([]rune(*str)[0:finalLen])
	} else {
		*em = fmt.Sprintf("%-*s", finalLen-le, *em)
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

func Bigstring(s string) string {
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

func Sleep(delay_seconds uint16) {
	time.Sleep(time.Second * time.Duration(delay_seconds))
}

func SleepMilli(n time.Duration) {
	time.Sleep(time.Millisecond * n)
}

func Str2uint(s string) uint16 {
	var i uint16
	fmt.Sscanf(s, "%d", &i)
	return i
}

func Str2uint32(s string) uint32 {
	var i uint32
	fmt.Sscanf(s, "%d", &i)
	return i
}

func Uint2str(i uint16) string {
	return fmt.Sprintf("%d", i)
}

func Int2str(i int) string {
	return fmt.Sprintf("%d", i)
}

func S2u(s string) uint16 {
	i, _ := strconv.Atoi(s)
	return uint16(i)
}

const (
	ALIGN_LEFT = iota
	ALIGN_RIGHT
)

func FilterStrings(s1, s2 *[]string) {
	var newS []string
	for _, s := range *s1 {
		if slices.Contains(*s2, s) {
// trace.Print("contains %s", s) // t //
			newS = append(newS, s)
		}
	}
	*s1 = nil
	*s1 = make([]string, len(newS))
	copy(*s1, newS)
}

/* Prints a slice of strings line by line. */
func PrintSlice(lines *[]string) {
	for i, line := range *lines {
		fmt.Printf("%3d  %s\n", i, line)
	}
}

func catch(err error) {
	if err != nil {
		fmt.Println("ERROR", err)
	}
}

