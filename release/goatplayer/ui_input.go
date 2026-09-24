// https://stackoverflow.com/questions/5966903/how-can-i-get-mousemove-and-mouseclick-in-bash
// echo -e "\e[?1000;1006;1015h"
// echo -e "\e[?1000;1006;1015l"

// cursot position
//  echo -ne "\033[6n" ; read aaa

package main

import (
	"fmt"
	"os"
	"strings"
	// "strings"
	// "os/exec"
)

func (this *uInput) DisplayInput() {
// var s string /*c*/
	var t string
	switch this.inputType {
	case INPUT_MOUSE:
// s = fmt.Sprintf("mouse, row: %d, col: %d, ", this.iRow, this.iCol) /*c*/
		t = fmt.Sprintf("Mouse: %d/%d, ", this.iRow, this.iCol)
	case INPUT_KEYBOARD:
// s = "keyboard, " /*c*/
		t = "Keyboard: "
	}
// var key_s string /*c*/
	var key_t string
	if v, ok := this.mapKeypad[this.inputKey]; ok {
// key_s = fmt.Sprintf("%s  (%q)", v, this.inputKey) /*c*/
		key_t = fmt.Sprintf("%s", v)
	} else {
// key_s = this.inputKey /*c*/
		key_t = this.inputKey
	}
	// s3 := fmt.Sprintf("%4s", key)
// s += key_s /*c*/
	t += key_t
// trace.PrintFocus("input: %s", s) // t //
	Term.AddStringWithPos(MainScreen.CommandPanelSection.marginBottom+2, MainScreen.CommandPanelSection.marginLeft, styleBlackWhite.EchoStyle(t, 25)).Flush()
	// Screen.Add(10, FileContainers.marginRight+23, Pen.StatusLine+s3+Pen.Normal)
	// Screen.Print()
}

func ReadLineFromConsole() string {
	CR_EL := "\r\x1b[0K"
	var line, s string
	var b9 []byte
	var char rune
	b9 = make([]byte, 9)
	for {
		fmt.Scanf("%c", &char)
		s = string(char)
		switch s {
		case KEYB_BACKSPACE:
			// fmt.Println(" - back")
			if len(line) > 0 {
				rr := []rune(line)
				rr = rr[:len(rr)-1]
				line = string(rr)
				// fmt.Printf("%s%s", CR_EL, line)
			}
		case KEYB_ESC:
			//! If only esc is pressed, waits for another input
			// fmt.Println(" - esc")
			os.Stdin.Read(b9)
			//fmt.Scanf("%s", &ttt) // dump the escape-sequence
		case "\n":
			fmt.Println()
			return line
		default:
			// fmt.Printf(" - s: %s\n", s)
			line += s
		}
		fmt.Printf("%s%s", CR_EL, line)
	}
}

func IsValidMouseInput() bool {
	if Input.inputKey == MOUSE_B1_UP {
		if Term.activeCommand != nil {
			PrintUnpressedCommand(Term.activeCommand)
			Term.activeCommand = nil
		}
		if Term.pPressedButton != nil {
			Term.pPressedButton.PrintUnpressed()
			Term.pPressedButton = nil
		}
		return false
	}
	if Input.inputKey == MOUSE_B2_UP {
		return false
	}
	if Input.inputKey == MOUSE_B3_UP {
		return false
	}
	return true
}

func (this *uInput) ReadNumber(suffix string) string {
	var num string
	for {
		this.GetNextChar()
		if this.s == suffix {
			// fmt.Printf("trovato num %s\n", num)
			return num
		} else {
			num += this.s
		}
	}
}

func (this *uInput) ReadNumber2(suffix1, suffix2 string) (string, string) {
	var num string
	for {
		this.GetNextChar()
		switch this.s {
		case suffix1:
			// fmt.Printf("trovato num %s\n", num)
			return num, suffix1
		case suffix2:
			return num, suffix2
		default:
			num += this.s
		}
	}
}

/* Asks the user to input a line. */
func (this *uInput) ReadLine() string {
	var command string
	var b9 []byte
	b9 = make([]byte, 9)
	fmt.Print(TPUT_SC)
	fmt.Print(TPUT_CNORM)
	for {
		this.GetNextChar()
		switch this.s {
		case KEYB_BACKSPACE:
			if len(command) > 0 {
				command = command[:len(command)-1]
			}
		case "\n":
			fmt.Print(TPUT_CIVIS)
			return command
		case KEYB_ESC:
			os.Stdin.Read(b9)
		default:
			command += this.s
		}
		// trace.Print("%q %q", this.s, command)
		le := len(command)
		em := ""
		AdjustWidth(&em, uint16(le)+2)
		Term.AddSimple(TPUT_RC + em + TPUT_RC + command)
		Term.PrintSimple()
	}
}

/* (new) Asks the user to input a line at a position. Uses the Backspace to delete and Enter to finish. */
func (this *uInput) ReadLineAtPos(iRow, iCol, lineLen uint16) string {
// trace.Begin("ReadLineAtPos") // t //

	var command string
	var b9 []byte
	b9 = make([]byte, 9)
	fmt.Print(TPUT_CNORM)
	em := styleFindPattern.EchoStyle("", lineLen)
	Term.AddStringWithPos(iRow, iCol, em).Flush()
	for {
		this.GetNextChar()
		switch this.s {
		case KEYB_BACKSPACE:
			if len(command) > 0 {
				rr := []rune(command)
				rr = rr[:len(rr)-1]
				command = string(rr)
			}
		case "\n":
			fmt.Print(TPUT_CIVIS)
			Term.AddStringWithPos(iRow, iCol, em).Flush()
// trace.ReturnAdd(command) // t //
			return command
		case KEYB_ESC:
			os.Stdin.Read(b9)
		default:
			command += this.s
		}
		// trace.Print("%q %q", this.s, command)
		s := styleFindPattern.EchoStyle(command, lineLen)
		Term.AddStringWithPos(iRow, iCol, em).AddStringWithPos(iRow, iCol, s).Flush()
	}
}

// Reads an input line and splits it in arguments.
func (this *uInput) ReadCommand(iRow uint16, iCol uint16, prompt string) string {
	fmt.Print(TputCup(iRow, iCol))
	fmt.Printf("%s", prompt)
	command := this.ReadLine()
	le := len(command) + len(prompt)
	em := ""
	AdjustWidth(&em, uint16(le))
	fmt.Print(TputCup(iRow, iCol))
	Term.AddSimple(em)
	Term.PrintSimple()
	return command
	// sttyEchoOff()
	// sttyEchoOn()
	// fmt.Println()
	// fmt.Printf("command: %s\n", command)
	// *out = strings.Split(command, " ")
}

// Reads a number from the user.
func (this *uInput) ReadInt(prompt string) int {
	var out int
	fmt.Printf("%s", prompt)
	fmt.Scanln(&out)
	return out
}

var Button vButton

func (this *uInput) GetKeyName(key string) string {
	if desc, ok := this.mapKeypad[key]; ok {
		return desc
	}
	return key
}

func (this *uInput) InitInput() {
	this.xSelect.SetOuter()
	this.b = make([]byte, 1)

	this.mapKeypad = map[string]string{
		KEYPAD_ARROW_LEFT:      "◄",
		KEYPAD_ARROW_RIGHT:     "►",
		KEYPAD_ARROW_UP:        "▲",
		KEYPAD_ARROW_DOWN:      "▼",
		KEYPAD_HOME:            "Home",
		KEYPAD_END:             "End",
		KEYPAD_PAGE_UP:         "PgUp",
		KEYPAD_PAGE_DOWN:       "PgDn",
		KEYPAD_ALT_PAGE_UP:     "A+PgUp",
		KEYPAD_ALT_PAGE_DOWN:   "A+PgDn",
		KEYPAD_ALT_HOME:        "A+Home",
		KEYPAD_ALT_END:         "A+End",
		KEYPAD_ALT_ARROW_UP:    "A+▲",
		KEYPAD_ALT_ARROW_DOWN:  "A+▼",
		KEYPAD_ALT_ARROW_LEFT:  "A+◄",
		KEYPAD_ALT_ARROW_RIGHT: "A+►",

		KEYB_ALT_A:       "A+a",
		KEYB_ALT_B:       "A+b",
		KEYB_ALT_C:       "A+c",
		KEYB_ALT_D:       "A+d",
		KEYB_ALT_E:       "A+e",
		KEYB_ALT_F:       "A+f",
		KEYB_ALT_G:       "A+g",
		KEYB_ALT_H:       "A+h",
		KEYB_ALT_I:       "A+i",
		KEYB_ALT_J:       "A+j",
		KEYB_ALT_K:       "A+k",
		KEYB_ALT_L:       "A+l",
		KEYB_ALT_M:       "A+m",
		KEYB_ALT_N:       "A+n",
		KEYB_ALT_O:       "A+o",
		KEYB_ALT_P:       "A+p",
		KEYB_ALT_Q:       "A+q",
		KEYB_ALT_R:       "A+r",
		KEYB_ALT_S:       "A+s",
		KEYB_ALT_T:       "A+t",
		KEYB_ALT_U:       "A+u",
		KEYB_ALT_V:       "A+v",
		KEYB_ALT_X:       "A+x",
		KEYB_ALT_Y:       "A+y",
		KEYB_ALT_Z:       "A+z",
		KEYB_ALT_SHIFT_Q: "A+Q",

		KEYB_ESC:       "Esc",
		KEYB_ENTER:     "Enter",
		KEYB_TAB:       "Tab",
		KEYB_SHIFT_TAB: "S+Tab",

		KEYB_F1:      "F1",
		KEYB_F2:      "F2",
		KEYB_F3:      "F3",
		KEYB_F4:      "F4",
		KEYB_F5:      "F5",
		KEYB_CTRL_F1: "C+F1",
		KEYB_CTRL_F2: "C+F2",
		KEYB_CTRL_F3: "C+F3",
		KEYB_CTRL_F4: "C+F4",
		KEYB_CTRL_F5: "C+F5",
		KEYB_CTRL_F6: "C+F6",
		KEYB_CTRL_F7: "C+F7",
		KEYB_CTRL_F8: "C+F8",

		MOUSE_B1_DOWN:      "M1",
		MOUSE_B1_UP:        "M1▲",
		MOUSE_B3_DOWN:      "M3",
		MOUSE_CTRL_B3_DOWN: "C+M3",
		MOUSE_WHEELUP:      "MW▲",
		MOUSE_WHEELDOWN:    "MW▼",
		MOUSE_B3_UP:        "M3▲"}

	this.KeypadMapLong = map[string]string{
		KEYPAD_ARROW_LEFT:    "◄",
		KEYPAD_ARROW_RIGHT:   "►",
		KEYPAD_ARROW_UP:      "▲",
		KEYPAD_ARROW_DOWN:    "▼",
		KEYPAD_HOME:          "Home",
		KEYPAD_END:           "End",
		KEYPAD_PAGE_UP:       "PgUp",
		KEYPAD_PAGE_DOWN:     "PgDn",
		KEYPAD_ALT_PAGE_UP:   "A+PgUp",
		KEYPAD_ALT_PAGE_DOWN: "A+PgDn",
		KEYB_F1:              "F1",
		KEYB_F2:              "F2",
		KEYB_F3:              "F3",
		KEYB_F4:              "F4",
		MOUSE_B1_DOWN:        "Mouse Left Down",
		MOUSE_B1_UP:          "Mouse Left Up",
		MOUSE_B3_DOWN:        "Mouse Right Down"}
}

func (this *uInput) MouseOn() {
	fmt.Print(MOUSE_ON)
}

func (this *uInput) MouseOff() {
	fmt.Print(MOUSE_OFF)
}

func (this *uInput) GetNextChar() {
// n, err := os.Stdin.Read(this.b) /*c*/
os.Stdin.Read(this.b)
	this.s = string(this.b)
// trace.PrintGreen(">>>>>> %4X   (n: %d, err:%v)", this.s, n, err) // t //
	// fmt.Printf(">>>> %q\n", this.s)
}

func (this *uInput) ReadInput() bool {
// trace.PrintNoTab("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ") // t //
// trace.Begin("ReadInput") // t //
	for {
		var s string
		var b []byte
		b = make([]byte, 20)
		le, _ := os.Stdin.Read(b)
		b = b[0:le]
		s = string(b)
		// trace.Print("string: %q (%d)", s, len(s))
		if len(s) >= 3 && s[0:3] == "\x1b[<" {
			/// Mouse input
			this.inputType = INPUT_MOUSE
			/// Scan coordinates
			i1 := strings.Index(s, ";")
			i2 := strings.LastIndex(s, ";")
			i3 := len(s) - 1
			// trace.PrintGreen("i1 %d i2 %d i3 %d", i1, i2, i3)
			fmt.Sscanf(s[i1+1:i2], "%d", &this.iCol)
			fmt.Sscanf(s[i2+1:i3], "%d", &this.iRow)
			this.iRow--
			this.iCol--
			this.inputKey = s[0:i1] + s[i3:i3+1]
			if !IsValidMouseInput() {
				continue
			}
			this.DisplayInput()
// trace.ReturnAdd("mouse") // t //
			return true
			// } else if len(s) >= 4 && s[0:4] == "\x1b[8;" {
			// 	trace.Print("sequence 8")
			// } else if len(s) >= 4 && s[0:4] == "\x1b[9;" {
			// 	trace.Print("sequence 9")
		}
		/// Keyboard input
		this.inputType = INPUT_KEYBOARD
		this.inputKey = s
		this.DisplayInput()
// trace.ReturnAdd("keyboard") // t //
		return true
	}
}

func (this *uInput) TrackMouse() bool {
	// trace.Begin("TrackMouse")
	var s string
	var b []byte
	var x uint16
	b = make([]byte, 20)
	le, _ := os.Stdin.Read(b)
	b = b[0:le]
	s = string(b)
	this.inputType = INPUT_MOUSE
	i0 := 3
	i1 := strings.Index(s, ";")
	i2 := strings.LastIndex(s, ";")
	i3 := len(s) - 1
	fmt.Sscanf(s[i0:i1], "%d", &x)
	fmt.Sscanf(s[i1+1:i2], "%d", &this.iCol)
	fmt.Sscanf(s[i2+1:i3], "%d", &this.iRow)
	this.iRow--
	this.iCol--
	this.inputKey = s[0:i1] + s[i3:i3+1]
	if x == 0 {
		return false
	}
	return true
}

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

func (this *vButton) PrintPressed() {
	bs := styleButton_Frame.EchoStyle(this.text, this.len)
	Term.AddStringWithPos(this.iRow, this.left, bs)
	Term.Flush()
}

func (this *vButton) PrintUnpressed() {
	bs := styleButton_SelFrame.EchoStyle(this.text, this.len)
	Term.AddStringWithPos(this.iRow, this.left, bs)
	Term.Flush()
}

func (this *vButton) Clicked(row uint16, col uint16) bool {
	right := this.left + this.len - 1
	if this.iRow == row {
		if this.left <= col && right >= col {
			this.PrintPressed()
			// trace.Print("Clicked %s row %d, col %d, %d %d", this.text, row, col, this.iRow, this.left)
			return true
		}
	}
	// trace.Print("Not Clicked %s row %d, col %d, %d %d", this.text, row, col, this.iRow, this.left)
	return false
}

// type t__InputCommand struct {
// 	key     string
// 	margins v__Button
// }

type uInput struct {
	b             []byte
	s             string // next char
	inputType     int
	inputKey      string
	iRow          uint16
	iCol          uint16
	X             uint16
	Y             uint16
	mapKeypad     map[string]string
	KeypadMapLong map[string]string
	xSelect       mDialog
}

