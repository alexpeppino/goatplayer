package K

/// Files and dirs

const (
	DIR_SETTINGS           = "settings"
	DIR_DATA               = "data"
	DIR_LOG                = "log"
	FILE_MUSIC_DIR         = "settings/music_dir.txt"
	FILE_SETTINGS_INDENTED = "settings/settings_indented.txt"
	FILE_AUDIOFILES        = "data/audiofiles.txt"
	FILE_PFLISTS           = "data/pflists.txt"
	FILE_LOG               = "log/log.txt"
	FILE_ERRLOG            = "log/errlog.txt"
	FILE_SETTINGS_UIPL     = "settings/ui_pl.txt"
	FILE_SETTINGS_EQU      = "settings/equalizer.txt"
	FILE_SETTINGS_STYLES   = "settings/styles.txt"
	FILE_SETTINGS_SETS     = "settings/sets.txt"
	PERM_DIR               = 0774
	PERM_OPEN_FILE         = 0774
)

/// Input

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
	//	keypad 2-chars
	MOUSE_ON           = "\x1b[?1000;1006;1015h"
	MOUSE_OFF          = "\x1b[?1000;1006;1015l"
	MOUSE_TRACKING_ON  = "\x1b[?1002h"
	MOUSE_TRACKING_OFF = "\x1b[?1002l"
)

/// Input type

const (
	INPUT_MOUSE = iota
	INPUT_KEYBOARD
	INPUT_KEYBOARD_ALT
	INPUT_KEYBOARD_F
	INPUT_KEYBOARD_CTRL
)

/// Tput

const (
	TPUT_SC    = "\x1b7"
	TPUT_RC    = "\x1b8"
	TPUT_CIVIS = "\x1b[?25l"
	TPUT_CNORM = "\x1b[?12l\x1b[?25h"
)

/// Focus

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

/// Help files

const (
	HELPFILE_BROWSERS      = "Browsers"
	HELPFILE_COLORS        = "Colors"
	HELPFILE_DIRS          = "Directories"
	HELPFILE_DOWNLOAD      = "Download"
	HELPFILE_EQUALIZER     = "Equalizer"
	HELPFILE_FIND          = "Find"
	HELPFILE_HELP          = "Help"
	HELPFILE_LISTS         = "Lists"
	HELPFILE_PARAMETERS    = "Parameters"
	HELPFILE_PLAYER        = "Player"
	HELPFILE_PROGRESS_BARS = "Progress bars"
	HELPFILE_PROPERTIES    = "Properties"
	HELPFILE_SECTIONS      = "02 Sections"
	HELPFILE_SETTINGS      = "Settings"
	HELPFILE_TAG_FIELDS    = "Tag fields"
	HELPFILE_TERMINOLOGY   = "Terminology"
	HELPFILE_REPORT        = "Report"
)

/// Styles

const (
	ALIGN_LEFT = iota
	ALIGN_RIGHT
	CLR = "\x1b[38;5;%d;48;5;%dm" // color

	FONT_NORMAL         = ""
	FONT_BOLD           = "\x1b[1m"
	FONT_BOLD_UNDERLINE = "\x1b[1m\x1b[4m"
	FONT_ITALIC         = "\x1b[3m"
	FONT_UNDERLINE      = "\x1b[4m"
	FONT_RESET          = "\x1b[0m"
	FONT_BLINK          = "\x1b[5m"

	DC_RESET     = "\x1b[0m"
	DC_BOLD      = "\x1b[1m"
	DC_ITALIC    = "\x1b[3m"
	DC_UNDERLINE = "\x1b[4m"
)

/// Chars

const (
	CHAR_UPPER_HALF_BLOCK = "▀"
	CHAR_LOWER_HALF_BLOCK = "▄"
)

/// Settings

const (
	SET_ON         = "On"
	SET_OFF        = "Off"
	SET_YES        = "Yes"
	SET_NO         = "No"
	SET_REPEAT_OFF = "No"
	SET_REPEAT_1   = "1"
	SET_REPEAT_ALL = "All"

	SET_POS_REPEAT_OFF uint16 = 0
	SET_POS_REPEAT_1   uint16 = 1
	SET_POS_REPEAT_ALL uint16 = 2
)

///

const (
	UNKNOWN_GENRE  = "Unknown genre"
	UNKNOWN_ARTIST = "Unknown artist"
	UNKNOWN_ALBUM  = "Unknown album"
	UNKNOWN_TITLE  = "Unknown title"
	UNKNOWN_YEAR   = "Unknown year"
)

///

const (
	UNSET       uint16 = 65535
	NEVERSET    uint16 = 65534
	OUT_OF_PAGE uint16 = 65535
	UNDEF              = "UNDEFINED"
	UNSET_S            = "K.UNSET"
	NONE               = "NONE"
)

const (
	CONSOLE_PROMPT = DC_BOLD + "console>" + DC_RESET
)

