package input // Single-instance package

import (
	"fmt"
	"goatplayer/internal/K"
	"goatplayer/internal/echo"
	"goatplayer/internal/small/utils"
	"goatplayer/internal/style"
	"os"
	"strings"
)

var (
	b             []byte
	s             string // next char
	InputType     int
	InputKey      string
	RowI          uint16
	ColI          uint16
	MapKeypad     map[string]string
	KeypadMapLong map[string]string
	// xSelect       mDialog
	marginTop  uint16
	marginLeft uint16
	term       inTerm
)

type inTerm interface {
	CheckButtons()
}

func _isValidMouseInput() bool {
	if InputKey == K.MOUSE_B1_UP {
		term.CheckButtons()
		// if Term.activeCommand != nil {
		// 	PrintUnpressedCommand(Term.activeCommand)
		// 	Term.activeCommand = nil
		// }
		// if Term.pPressedButton != nil {
		// 	Term.pPressedButton.PrintUnpressed()
		// 	Term.pPressedButton = nil
		// }
		return false
	}
	if InputKey == K.MOUSE_B2_UP {
		return false
	}
	if InputKey == K.MOUSE_B3_UP {
		return false
	}
	return true
}

func TrackMouse() bool {
	// trace.Begin("TrackMouse")
	var s string
	var b []byte
	var x uint16
	b = make([]byte, 20)
	le, _ := os.Stdin.Read(b)
	b = b[0:le]
	s = string(b)
	InputType = K.INPUT_MOUSE
	i0 := 3
	i1 := strings.Index(s, ";")
	i2 := strings.LastIndex(s, ";")
	i3 := len(s) - 1
	fmt.Sscanf(s[i0:i1], "%d", &x)
	fmt.Sscanf(s[i1+1:i2], "%d", &ColI)
	fmt.Sscanf(s[i2+1:i3], "%d", &RowI)
	RowI--
	ColI--
	InputKey = s[0:i1] + s[i3:i3+1]
	if x == 0 {
		return false
	}
	return true
}

func PrepareInput(top, left uint16) {
	marginTop = top
	marginLeft = left
}

func DisplayInput() {
// var s string /*c*/
	var t string
	switch InputType {
	case K.INPUT_MOUSE:
// s = fmt.Sprintf("mouse, row: %d, col: %d, ", RowI, ColI) /*c*/
		t = fmt.Sprintf("Mouse: %d/%d, ", RowI, ColI)
	case K.INPUT_KEYBOARD:
// s = "keyboard, " /*c*/
		t = "Keyboard: "
	}
// var key_s string /*c*/
	var key_t string
	if v, ok := MapKeypad[InputKey]; ok {
// key_s = fmt.Sprintf("%s  (%q)", v, InputKey) /*c*/
		key_t = fmt.Sprintf("%s", v)
	} else {
// key_s = InputKey /*c*/
		key_t = InputKey
	}
	// s3 := fmt.Sprintf("%4s", key)
// s += key_s /*c*/
	t += key_t
// trace.PrintFocus("input: %s", s) // t //
	echo.AddStringWithPos(marginTop, marginLeft, style.BlackWhite.EchoStyle(t, 25))
	echo.Flush()
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
		case K.KEYB_BACKSPACE:
			// fmt.Println(" - back")
			if len(line) > 0 {
				rr := []rune(line)
				rr = rr[:len(rr)-1]
				line = string(rr)
				// fmt.Printf("%s%s", CR_EL, line)
			}
		case K.KEYB_ESC:
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

func _readNumber(suffix string) string {
	var num string
	for {
		_getNextChar()
		if s == suffix {
			// fmt.Printf("trovato num %s\n", num)
			return num
		} else {
			num += s
		}
	}
}

func _readNumber2(suffix1, suffix2 string) (string, string) {
	var num string
	for {
		_getNextChar()
		switch s {
		case suffix1:
			// fmt.Printf("trovato num %s\n", num)
			return num, suffix1
		case suffix2:
			return num, suffix2
		default:
			num += s
		}
	}
}

/* Asks the user to input a line. */
func ReadLine() string {
	var command string
	var b9 []byte
	b9 = make([]byte, 9)
	fmt.Print(K.TPUT_SC)
	fmt.Print(K.TPUT_CNORM)
	for {
		_getNextChar()
		switch s {
		case K.KEYB_BACKSPACE:
			if len(command) > 0 {
				command = command[:len(command)-1]
			}
		case "\n":
			fmt.Print(K.TPUT_CIVIS)
			return command
		case K.KEYB_ESC:
			os.Stdin.Read(b9)
		default:
			command += s
		}
		// trace.Print("%q %q", self.s, command)
		le := len(command)
		em := ""
		utils.AdjustWidth(&em, uint16(le)+2)
		echo.AddSimple(K.TPUT_RC + em + K.TPUT_RC + command)
		echo.PrintSimple()
	}
}

/* (new) Asks the user to input a line at a position. Uses the Backspace to delete and Enter to finish. */
func ReadLineAtPos(iRow, iCol, lineLen uint16) string {
// trace.Begin("ReadLineAtPos") // t //
	var command string
	var b9 []byte
	b9 = make([]byte, 9)
	fmt.Print(K.TPUT_CNORM)
	em := style.FindPattern.EchoStyle("", lineLen)
	echo.AddStringWithPos(iRow, iCol, em)
	echo.Flush()
	for {
		_getNextChar()
		switch s {
		case K.KEYB_BACKSPACE:
			if len(command) > 0 {
				rr := []rune(command)
				rr = rr[:len(rr)-1]
				command = string(rr)
			}
		case "\n":
			fmt.Print(K.TPUT_CIVIS)
			echo.AddStringWithPos(iRow, iCol, em)
			echo.Flush()
// trace.ReturnAdd(command) // t //
			return command
		case K.KEYB_ESC:
			os.Stdin.Read(b9)
		default:
			command += s
		}
		// trace.Print("%q %q", self.s, command)
		s := style.FindPattern.EchoStyle(command, lineLen)
		echo.AddStringWithPos(iRow, iCol, em)
		echo.AddStringWithPos(iRow, iCol, s)
		echo.Flush()
	}
}

// Reads an input line and splits it in arguments.
func ReadCommand(iRow uint16, iCol uint16, prompt string) string {
	fmt.Print(echo.TputCup(iRow, iCol))
	fmt.Printf("%s", prompt)
	command := ReadLine()
	le := len(command) + len(prompt)
	em := ""
	utils.AdjustWidth(&em, uint16(le))
	fmt.Print(echo.TputCup(iRow, iCol))
	echo.AddSimple(em)
	echo.PrintSimple()
	return command
	// sttyEchoOff()
	// sttyEchoOn()
	// fmt.Println()
	// fmt.Printf("command: %s\n", command)
	// *out = strings.Split(command, " ")
}

// Reads a number from the user.
func ReadInt(prompt string) int {
	var out int
	fmt.Printf("%s", prompt)
	fmt.Scanln(&out)
	return out
}

func GetKeyName(key string) string {
	if desc, ok := MapKeypad[key]; ok {
		return desc
	}
	return key
}

func InitInput(t inTerm) {
	// self.xSelect.SetOuter()
	term = t
	b = make([]byte, 1)

	MapKeypad = map[string]string{
		K.KEYPAD_ARROW_LEFT:      "◄",
		K.KEYPAD_ARROW_RIGHT:     "►",
		K.KEYPAD_ARROW_UP:        "▲",
		K.KEYPAD_ARROW_DOWN:      "▼",
		K.KEYPAD_HOME:            "Home",
		K.KEYPAD_END:             "End",
		K.KEYPAD_PAGE_UP:         "PgUp",
		K.KEYPAD_PAGE_DOWN:       "PgDn",
		K.KEYPAD_ALT_PAGE_UP:     "A+PgUp",
		K.KEYPAD_ALT_PAGE_DOWN:   "A+PgDn",
		K.KEYPAD_ALT_HOME:        "A+Home",
		K.KEYPAD_ALT_END:         "A+End",
		K.KEYPAD_ALT_ARROW_UP:    "A+▲",
		K.KEYPAD_ALT_ARROW_DOWN:  "A+▼",
		K.KEYPAD_ALT_ARROW_LEFT:  "A+◄",
		K.KEYPAD_ALT_ARROW_RIGHT: "A+►",

		K.KEYB_ALT_A:       "A+a",
		K.KEYB_ALT_B:       "A+b",
		K.KEYB_ALT_C:       "A+c",
		K.KEYB_ALT_D:       "A+d",
		K.KEYB_ALT_E:       "A+e",
		K.KEYB_ALT_F:       "A+f",
		K.KEYB_ALT_G:       "A+g",
		K.KEYB_ALT_H:       "A+h",
		K.KEYB_ALT_I:       "A+i",
		K.KEYB_ALT_J:       "A+j",
		K.KEYB_ALT_K:       "A+k",
		K.KEYB_ALT_L:       "A+l",
		K.KEYB_ALT_M:       "A+m",
		K.KEYB_ALT_N:       "A+n",
		K.KEYB_ALT_O:       "A+o",
		K.KEYB_ALT_P:       "A+p",
		K.KEYB_ALT_Q:       "A+q",
		K.KEYB_ALT_R:       "A+r",
		K.KEYB_ALT_S:       "A+s",
		K.KEYB_ALT_T:       "A+t",
		K.KEYB_ALT_U:       "A+u",
		K.KEYB_ALT_V:       "A+v",
		K.KEYB_ALT_X:       "A+x",
		K.KEYB_ALT_Y:       "A+y",
		K.KEYB_ALT_Z:       "A+z",
		K.KEYB_ALT_SHIFT_Q: "A+Q",

		K.KEYB_ESC:       "Esc",
		K.KEYB_ENTER:     "Enter",
		K.KEYB_TAB:       "Tab",
		K.KEYB_SHIFT_TAB: "S+Tab",

		K.KEYB_F1:      "F1",
		K.KEYB_F2:      "F2",
		K.KEYB_F3:      "F3",
		K.KEYB_F4:      "F4",
		K.KEYB_F5:      "F5",
		K.KEYB_CTRL_F1: "C+F1",
		K.KEYB_CTRL_F2: "C+F2",
		K.KEYB_CTRL_F3: "C+F3",
		K.KEYB_CTRL_F4: "C+F4",
		K.KEYB_CTRL_F5: "C+F5",
		K.KEYB_CTRL_F6: "C+F6",
		K.KEYB_CTRL_F7: "C+F7",
		K.KEYB_CTRL_F8: "C+F8",

		K.MOUSE_B1_DOWN:      "M1",
		K.MOUSE_B1_UP:        "M1▲",
		K.MOUSE_B3_DOWN:      "M3",
		K.MOUSE_CTRL_B3_DOWN: "C+M3",
		K.MOUSE_WHEELUP:      "MW▲",
		K.MOUSE_WHEELDOWN:    "MW▼",
		K.MOUSE_B3_UP:        "M3▲"}

	KeypadMapLong = map[string]string{
		K.KEYPAD_ARROW_LEFT:    "◄",
		K.KEYPAD_ARROW_RIGHT:   "►",
		K.KEYPAD_ARROW_UP:      "▲",
		K.KEYPAD_ARROW_DOWN:    "▼",
		K.KEYPAD_HOME:          "Home",
		K.KEYPAD_END:           "End",
		K.KEYPAD_PAGE_UP:       "PgUp",
		K.KEYPAD_PAGE_DOWN:     "PgDn",
		K.KEYPAD_ALT_PAGE_UP:   "A+PgUp",
		K.KEYPAD_ALT_PAGE_DOWN: "A+PgDn",
		K.KEYB_F1:              "F1",
		K.KEYB_F2:              "F2",
		K.KEYB_F3:              "F3",
		K.KEYB_F4:              "F4",
		K.MOUSE_B1_DOWN:        "Mouse Left Down",
		K.MOUSE_B1_UP:          "Mouse Left Up",
		K.MOUSE_B3_DOWN:        "Mouse Right Down"}
}

func _getNextChar() {
// n, err := os.Stdin.Read(b) /*c*/
os.Stdin.Read(b)
	s = string(b)
// trace.PrintGreen(">>>>>> %4X   (n: %d, err:%v)", s, n, err) // t //
	// fmt.Printf(">>>> %q\n", self.s)
}

func ReadInput() bool {
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
			InputType = K.INPUT_MOUSE
			/// Scan coordinates
			i1 := strings.Index(s, ";")
			i2 := strings.LastIndex(s, ";")
			i3 := len(s) - 1
			// trace.PrintGreen("i1 %d i2 %d i3 %d", i1, i2, i3)
			fmt.Sscanf(s[i1+1:i2], "%d", &ColI)
			fmt.Sscanf(s[i2+1:i3], "%d", &RowI)
			RowI--
			ColI--
			InputKey = s[0:i1] + s[i3:i3+1]
			if !_isValidMouseInput() {
				continue
			}
			DisplayInput()
// trace.ReturnAdd("mouse") // t //
			return true
			// } else if len(s) >= 4 && s[0:4] == "\x1b[8;" {
			// 	trace.Print("sequence 8")
			// } else if len(s) >= 4 && s[0:4] == "\x1b[9;" {
			// 	trace.Print("sequence 9")
		}
		/// Keyboard input
		InputType = K.INPUT_KEYBOARD
		InputKey = s
		DisplayInput()
// trace.ReturnAdd("keyboard") // t //
		return true
	}
}

