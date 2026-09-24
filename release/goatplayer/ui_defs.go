package main

/// Focus keys

const (
	/// MainScreen
	FOCUS_VIEW1 = "1"
	FOCUS_VIEW2 = "2"
	FOCUS_PLAY  = "3"
	// FOCUS_EDIT         = "4"
	FOCUS_GENRES       = "4"
	FOCUS_ARTISTS      = "5"
	FOCUS_ALBUMS       = "6"
	FOCUS_YEARS        = "7"
	FOCUS_ALBUM_ARTIST = "8"
	FOCUS_COMPOSER     = "9"
	FOCUS_COMMENT      = "C"
	FOCUS_FILEPROP     = "P"
	FOCUS_DIRS         = "D"
	/// SettingsScreen
	FOCUS_HELP       = KEYB_F1
	FOCUS_SETTINGS   = KEYB_F2
	FOCUS_EQUALIZER  = KEYB_F3
	FOCUS_EQUALIZER2 = "E"
	FOCUS_COLORS     = KEYB_F4
	FOCUS_COLORS2    = "c"
	FOCUS_REPORT     = KEYB_F5
	// FOCUS_BIOS   = KEYB_F2
	/// FindScreen
	FOCUS_FIND = "F"
	///
	FOCUS_TAGS        = "T"
	FOCUS_PFLISTS     = "$"
	FOCUS_PLS         = "6"
	FOCUS_FOLDED      = "*"
	FOCUS_UNFOCUSABLE = ""
	MOUSE_MODE        = "MM"
)

const (
	/// MainScreen
	HELPFILE_VIEW1    = "Browsers"
	HELPFILE_VIEW2    = "Browsers"
	HELPFILE_PLAY     = "Player"
	HELPFILE_TAGS     = "Tag fields"
	HELPFILE_FILEPROP = "Properties"
	HELPFILE_DIRS     = "Directories"
	/// SettingsScreen
	HELPFILE_SETTINGS  = "Settings"
	HELPFILE_EQUALIZER = "Equalizer"
	HELPFILE_COLORS    = "Colors"
	HELPFILE_REPORT    = "Report"
	///
	HELPFILE_PARAMETERS = "Parameters"
	HELPFILE_PFLISTS    = "Lists"
	HELPFILE_FIND       = "Find"
)

/// Keyboard input

const (
	KEYB_ESC           = "\x1b"
	KEYB_CSI           = "\x1b["
	KEYB_OSC           = "\x1b]"
	KEYB_BACKSPACE     = "\x7f"
	KEYB_ENTER         = "\n"
	KEYB_TAB           = "\t"
	KEYB_SHIFT_TAB     = "\x1b[Z"
	KEYPAD_ARROW_UP    = "\x1b[A"
	KEYPAD_ARROW_DOWN  = "\x1b[B"
	KEYPAD_ARROW_LEFT  = "\x1b[D"
	KEYPAD_ARROW_RIGHT = "\x1b[C"

	KEYPAD_INS       = "\x1b[2~"
	KEYPAD_DEL       = "\x1b[3~"
	KEYPAD_HOME      = "\x1b[H"
	KEYPAD_END       = "\x1b[F"
	KEYPAD_PAGE_UP   = "\x1b[5~"
	KEYPAD_PAGE_DOWN = "\x1b[6~"

	KEYPAD_ALT_ARROW_UP    = "\x1b[1;3A"
	KEYPAD_ALT_ARROW_DOWN  = "\x1b[1;3B"
	KEYPAD_ALT_ARROW_LEFT  = "\x1b[1;3D"
	KEYPAD_ALT_ARROW_RIGHT = "\x1b[1;3C"

	KEYPAD_ALT_INS       = "\x1b[2;3~"
	KEYPAD_ALT_DEL       = "\x1b[3;3~"
	KEYPAD_ALT_HOME      = "\x1b[1;3H"
	KEYPAD_ALT_END       = "\x1b[1;3F"
	KEYPAD_ALT_PAGE_UP   = "\x1b[5;3~"
	KEYPAD_ALT_PAGE_DOWN = "\x1b[6;3~"

	KEYB_ALT_A       = "\x1ba"
	KEYB_ALT_B       = "\x1bb"
	KEYB_ALT_C       = "\x1bc"
	KEYB_ALT_D       = "\x1bd"
	KEYB_ALT_E       = "\x1be"
	KEYB_ALT_F       = "\x1bf"
	KEYB_ALT_G       = "\x1bg"
	KEYB_ALT_H       = "\x1bh"
	KEYB_ALT_I       = "\x1bi"
	KEYB_ALT_J       = "\x1bj"
	KEYB_ALT_K       = "\x1bk"
	KEYB_ALT_L       = "\x1bl"
	KEYB_ALT_M       = "\x1bm"
	KEYB_ALT_N       = "\x1bn"
	KEYB_ALT_O       = "\x1bo"
	KEYB_ALT_P       = "\x1bp"
	KEYB_ALT_Q       = "\x1bq"
	KEYB_ALT_R       = "\x1br"
	KEYB_ALT_S       = "\x1bs"
	KEYB_ALT_T       = "\x1bt"
	KEYB_ALT_U       = "\x1bu"
	KEYB_ALT_V       = "\x1bv"
	KEYB_ALT_X       = "\x1bx"
	KEYB_ALT_Y       = "\x1by"
	KEYB_ALT_Z       = "\x1bz"
	KEYB_ALT_SHIFT_Q = "\x1bQ"

	// 	"\x1bw" ( 2)    "\x1bw" ( 2)
	// "\x1be" ( 2)    "\x1be" ( 2)
	//    "\x1br" ( 2)    "\x1br" ( 2)
	//    "\x1bè" ( 3)    "\x1bè" ( 3)
	//    "\x1b+" ( 2)    "\x1b+" ( 2)

	//    "\x1bW" ( 2)    "\x1bW" ( 2)
	//    "\x1bE" ( 2)    "\x1bE" ( 2)
	//    "\x1bR" ( 2)    "\x1bR" ( 2)
	//    "\x1bé" ( 3)    "\x1bé" ( 3)
	//    "\x1b*" ( 2)    "\x1b*" ( 2)

	KEYB_F1 = "\x1bOP"
	KEYB_F2 = "\x1bOQ"
	KEYB_F3 = "\x1bOR"
	KEYB_F4 = "\x1bOS"
	KEYB_F5 = "\x1b[15~"
	KEYB_F6 = "\x1b[17~"
	KEYB_F7 = "\x1b[18~"
	KEYB_F8 = "\x1b[19~"

	KEYB_CTRL_F1 = "\x1b[1;5P"
	KEYB_CTRL_F2 = "\x1b[1;5Q"
	KEYB_CTRL_F3 = "\x1b[1;5R"
	KEYB_CTRL_F4 = "\x1b[1;5S"
	KEYB_CTRL_F5 = "\x1b[15;5~"
	KEYB_CTRL_F6 = "\x1b[17;5~"
	KEYB_CTRL_F7 = "\x1b[18;5~"
	KEYB_CTRL_F8 = "\x1b[19;5~"
)

/// Mouse input

const (
	MOUSE_B1_DOWN   = "\x1b[<0M"
	MOUSE_B1_UP     = "\x1b[<0m"
	MOUSE_B2_DOWN   = "\x1b[<1M"
	MOUSE_B2_UP     = "\x1b[<1m"
	MOUSE_B3_DOWN   = "\x1b[<2M"
	MOUSE_B3_UP     = "\x1b[<2m"
	MOUSE_WHEELUP   = "\x1b[<64M"
	MOUSE_WHEELDOWN = "\x1b[<65M"

	MOUSE_CTRL_B1_DOWN   = "\x1b[<16M"
	MOUSE_CTRL_B1_UP     = "\x1b[<16m"
	MOUSE_CTRL_B2_DOWN   = "\x1b[<17M"
	MOUSE_CTRL_B2_UP     = "\x1b[<17m"
	MOUSE_CTRL_B3_DOWN   = "\x1b[<18M"
	MOUSE_CTRL_B3_UP     = "\x1b[<18m"
	MOUSE_CTRL_WHEELUP   = "\x1b[<80M"
	MOUSE_CTRL_WHEELDOWN = "\x1b[<81M"
)

/// Input type

const (
	INPUT_MOUSE = iota
	INPUT_KEYBOARD
	INPUT_KEYBOARD_ALT
	INPUT_KEYBOARD_F
	INPUT_KEYBOARD_CTRL
)

///

const (
	//	keypad 2-chars
	MOUSE_ON           = "\x1b[?1000;1006;1015h"
	MOUSE_OFF          = "\x1b[?1000;1006;1015l"
	MOUSE_TRACKING_ON  = "\x1b[?1002h"
	MOUSE_TRACKING_OFF = "\x1b[?1002l"
)

// var KeypadMap = map[string]string{
// 	KEYPAD_ARROW_LEFT:  "◄",
// 	KEYPAD_ARROW_RIGHT: "►",
// 	KEYPAD_ARROW_UP:    "▲",
// 	KEYPAD_ARROW_DOWN:  "▼",
// 	KEYPAD_HOME:        "Home",
// 	KEYPAD_END:         "End",
// 	KEYPAD_PAGE_UP:     "PgUp",
// 	KEYPAD_PAGE_DOWN:   "PgDn",
// 	MOUSE_B3_DOWN:      "Mo-R"}

// >>>>  MOUSE ALT + press/unpress 1, 2, 3

// "\x1b[<8;68;26M" (11) "\x1b[<8;68;26M" (11)
// "\x1b[<8;68;26m" (11) "\x1b[<8;68;26m" (11)
// "\x1b[<9;68;26M" (11) "\x1b[<9;68;26M" (11)
// "\x1b[<9;68;26m" (11) "\x1b[<9;68;26m" (11)
// "\x1b[<10;68;26M" (12) "\x1b[<10;68;26M" (12)
// "\x1b[<10;68;26m" (12) "\x1b[<10;68;26m" (12)

// >>>>  MOUSE CTRL + ALT + press/unpress 1, 2, 3

// "\x1b[<24;73;21M" (12) "\x1b[<24;73;21M" (12)
// "\x1b[<24;73;21m" (12) "\x1b[<24;73;21m" (12)
// "\x1b[<25;73;21M" (12) "\x1b[<25;73;21M" (12)
// "\x1b[<25;73;21m" (12) "\x1b[<25;73;21m" (12)
// "\x1b[<26;73;21M" (12) "\x1b[<26;73;21M" (12)
// "\x1b[<26;73;21m" (12) "\x1b[<26;73;21m" (12)

