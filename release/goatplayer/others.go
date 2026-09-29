package main

import (
	"fmt"
	"goatplayer/internal/K"
	"goatplayer/internal/echo"
	"goatplayer/internal/style"
	"os"
	"sort"
)

func speedFormat() {
	// top := MainScreen.DisplaySection.marginTop
	// left := MainScreen.DisplaySection.marginLeft
	ui_pl.Speed = fmt.Sprintf("%.2f", float32(ui_pl.Speed_i)/100)
	Display.AddLineToPrint(1, 3, "Speed:   "+ui_pl.Speed)
	echo.Flush()
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
	echo.Flush()
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

type iGetName interface {
	GetName() string
}

type iTrace interface {
	GetLongName() (name string)
}

// Give a slice of iTrace, get list of names.
func GetLongNames[T iTrace](names []T) string {
	var s string
	for _, name := range names {
		s += name.GetLongName() + ", "
	}
	return s
}

func GetOuterNilOrName(outer iGetName) string {
	if outer == nil {
		return "<nil>"
	} else {
		return outer.GetName()
	}
}

func GetNilOrName(a any) string {
	s := "<nil>"
	switch aa := a.(type) {
	case *vMainScreen:
		if aa == nil {
			return s
		} else {
			return aa.name
		}
	case *maWorkspace:
		if aa == nil {
			return s
		} else {
			return aa.name
		}
	case *mBrowserWsp:
		if aa == nil {
			return s
		} else {
			return aa.name
		}
	case *vSection:
		if aa == nil {
			return s
		} else {
			return aa.name
		}
	case *vLayer:
		if aa == nil {
			return s
		} else {
			return aa.name
		}
	case *vScreen:
		if aa == nil {
			return s
		} else {
			return aa.name
		}
	case *vWindow:
		if aa == nil {
			return s
		} else {
			return aa.name
		}
	}
	return "<invalid type>"
}

func get_input(iRow uint16, col_i uint16) string {
	var b []byte = make([]byte, 1)
	var s string
	var input string
	for {
		os.Stdin.Read(b)
		s = string(b)
		// trace.Print("%q", s) // t //
		switch s {
		case "\n":
			return input
		case K.KEYB_BACKSPACE: // tab
			if len(input) > 0 {
				input = input[0 : len(input)-1]
			}
		default:
			input += s
		}
		// trace.Print("%s", input) // t //
		echo.AddStringWithPos(iRow, col_i, input+" ")
		echo.Flush()
	}
}

// Set costant margins
func (mrg *t__Margins) setConstantValues() {
	//	constant values
	mrg.sect_1_Top = 0
	mrg.sect_1_HeightMin = 25
	mrg.sect_1A_Left = 0
	mrg.sect_1A_WidthMin = 120
	mrg.sect_1A_SelectionBarLeft = 1
	mrg.sect_1B_WidthFix = 32
	mrg.sect_2_HeightMin = 22
	mrg.sect_2A_Left = 0
	mrg.sect_2B_WidthFix = 32
	mrg.sect_3_HeightFix = 6
	mrg.sect_3A_Left = 0
	mrg.sect_3B_WidthFix = 50
	mrg.sect_4_HeightFix = 1
	mrg.sect_4A_Left = 0
	mrg.sect_4B_WidthFix = 50
}

// Calculates margins according to
// mrg.ui_Cols and mrg.ui_Rows
func (mrg *t__Margins) setValues() {
	//	UI
	mrg.ui_Top = 0
	mrg.ui_Bottom = mrg.ui_Rows - 1
	mrg.ui_Left = 0
	mrg.ui_Right = mrg.ui_Cols - 1
	mrg.ui_Height = mrg.ui_Rows
	mrg.ui_Width = mrg.ui_Cols
	//	widths
	mrg.sect_1A_Width = mrg.ui_Width - mrg.sect_1B_WidthFix - 1
	mrg.sect_2A_Width = mrg.ui_Width - mrg.sect_2B_WidthFix - 1
	mrg.sect_3A_Width = mrg.ui_Width - mrg.sect_3B_WidthFix - 1
	mrg.sect_4A_Width = mrg.ui_Width - mrg.sect_4B_WidthFix - 1
	mrg.sect_1A_SelectionBarWidth = mrg.sect_1A_Width - 2
	//	heights
	var remain uint16 = mrg.ui_Height - mrg.sect_3_HeightFix - mrg.sect_4_HeightFix
	mrg.sect_2_Height = mrg.sect_2_HeightMin
	mrg.sect_1_Height = remain - mrg.sect_2_Height
	mrg.sect_1A_Height = mrg.sect_1_Height
	//	tops & bottoms
	mrg.sect_1_Bottom = mrg.sect_1_Top + mrg.sect_1_Height - 1
	mrg.sect_2_Top = mrg.sect_1_Bottom + 1
	mrg.sect_2_Bottom = mrg.sect_2_Top + mrg.sect_2_Height - 1
	mrg.sect_3_Top = mrg.sect_2_Bottom + 1
	mrg.sect_3_Bottom = mrg.sect_3_Top + mrg.sect_3_HeightFix - 1
	mrg.sect_4_Top = mrg.sect_3_Bottom + 1
	//	lefts & rights
	mrg.sect_1A_Right = mrg.sect_1A_Left + mrg.sect_1A_Width - 1
	mrg.sect_1B_Left = mrg.sect_1A_Right + 2
	mrg.sect_1B_Right = mrg.ui_Right
	mrg.sect_2A_Right = mrg.sect_2A_Left + mrg.sect_2A_Width - 1
	mrg.sect_2B_Left = mrg.sect_2A_Right + 2
	mrg.sect_2B_Right = mrg.ui_Right
	mrg.sect_3A_Right = mrg.sect_3A_Left + mrg.sect_3A_Width - 1
	mrg.sect_3B_Left = mrg.sect_3A_Right + 2
	mrg.sect_3B_Right = mrg.ui_Right
	mrg.sect_4A_Right = mrg.sect_4A_Left + mrg.sect_4A_Width - 1
	mrg.sect_4B_Left = mrg.sect_4A_Right + 2
	mrg.sect_4B_Right = mrg.ui_Right
	mrg.sect_1A_Height = mrg.sect_1_Height
	mrg.sect_1A_Bottom = mrg.sect_1_Bottom
}

var Button vButton

func (this *vButton) PrintPressed() {
	bs := style.Button_Frame.EchoStyle(this.text, this.len)
	echo.AddStringWithPos(this.iRow, this.left, bs)
	echo.Flush()
}

func (this *vButton) PrintUnpressed() {
	bs := style.Button_SelFrame.EchoStyle(this.text, this.len)
	echo.AddStringWithPos(this.iRow, this.left, bs)
	echo.Flush()
}

func (this *vButton) Clicked(row uint16, col uint16) bool {
	right := this.left + this.len - 1
	if this.iRow == row {
		if this.left <= col && right >= col {
			this.PrintPressed()
			// trace.Print("Clicked %s row %d, col %d, %d %d", this.Text, row, col, this.iRow, this.left)
			return true
		}
	}
	// trace.Print("Not Clicked %s row %d, col %d, %d %d", this.Text, row, col, this.iRow, this.left)
	return false
}

// https://stackoverflow.com/questions/5966903/how-can-i-get-mousemove-and-mouseclick-in-bash
// echo -e "\e[?1000;1006;1015h"
// echo -e "\e[?1000;1006;1015l"

// cursot position
//  echo -ne "\033[6n" ; read aaa

//	ESC	[	A					arrow up
//			B					arrow down
//			C					arrow right
//			D					arrow left
//			H					home
//			F					end
//			5	~				page up
//			6	~				page down
//
//	ESC	[	<	0;				mouse button 1
//				...	{X};		button down, X coord
//				... {Y}M		button down, Y coord
//				ESC [ <	0;
//				...	{X};		button up, X coord
//				... {Y}M		button up, Y coord
//
//	ESC	[	<	1;				mouse button 2
//				...
//	ESC	[	<	2;				mouse button 3
//				...
//
//	ESC [ 	<	64;				mouse wheel up
//				...	{X};		X coord
//				... {Y}M		Y coord
//
//	ESC [ 	<	65				mouse wheel down
//
//	ESC	O	P					F1
//	ESC	O	Q					F2
//
//	ESC	a						Alt-a
//	ESC	A						Alt-A
//
//	\x7f						backspace
//
//	\x01						Ctrl-A
//	\x02						Ctrl-B

// func (this *t__Button) Len() uint16 {
// 	return this.right - this.left + 1
// }

//////	Focus

//////

// type t__InputCommand struct {
// 	key     string
// 	margins v__Button
// }
