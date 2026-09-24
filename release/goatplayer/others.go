package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// func GetTerminalSize() {
// 	var s string
// 	var b []byte
// 	b = make([]byte, 20)
// 	fmt.Printf("%s18t", KEYB_CSI)
// 	le, _ := os.Stdin.Read(b)
// 	b = b[0:le]
// 	s = string(b)
// 	trace.Print("string: %q (%d)", s, len(s))
// 	i1 := strings.Index(s, ";")
// 	i2 := strings.LastIndex(s, ";")
// 	i3 := len(s) - 1
// 	fmt.Sscanf(s[i1+1:i2], "%d", &nRows)
// 	fmt.Sscanf(s[i2+1:i3], "%d", &nCols)
// 	trace.Print("-- size of the text area: %d rows x %d cols --", nRows, nCols)
// 	// SleepMilli(200)
// 	// trace.Print("-- size of the screen in characters --")
// 	// fmt.Printf("%s19t", KEYB_CSI)
// 	// SleepMilli(100)
// 	// fmt.Printf("%s13;6\007", KEYB_OSC)
// 	// fmt.Printf("%s14;6\007", KEYB_OSC)
// }

func SetTerminalTitle(s string) {
	fmt.Printf("%s0;%s\007", KEYB_OSC, s)
}

// func sttyEchoOff() {
// 	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
// }

// func sttyEchoOn() {
// 	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
// }

// func sttyCBreakMin1() {
// 	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
// }

func EchoTputCup(row uint16, col uint16, s string) {
	fmt.Printf("%s%s", TputCup(row, col), s)
}

func TputCup(row uint16, col uint16) string {
	return fmt.Sprintf("\x1b[%d;%dH", row, col)
}

func TputCNorm() {
	fmt.Print(TPUT_CNORM)
}

func TputCivis() {
	fmt.Print(TPUT_CIVIS)
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
		case KEYB_BACKSPACE: // tab
			if len(input) > 0 {
				input = input[0 : len(input)-1]
			}
		default:
			input += s
		}
// trace.Print("%s", input) // t //
		Term.AddStringWithPos(iRow, col_i, input+" ")
		Term.Flush()
	}

	// var s string
	// exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	// tput_cup_echo(54, 100, "Enter: ")

	// reader := bufio.NewReader(os.Stdin)
	// s, _ = reader.ReadString('\n')
	// s, _ = strings.CutSuffix(s, "\n")
	// trace.Print("input: %q", s)
	// tput_cup_echo(54, 100, "                         ")
	// tput_cup_echo(54, 100, "")
	// exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// return s
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

// Get the terminal size using command "stty size".
// Set values of parameters cols and lines.
func GetSize(cols_ *uint16, lines_ *uint16) {

	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	var cols, lines uint16
	fmt.Sscan(string(out), &lines, &cols)
	fmt.Printf("cols= %v lines= %v\n", cols, lines)
	*cols_ = cols
	*lines_ = lines

	Mrg.ui_Cols = cols
	Mrg.ui_Rows = lines
}

