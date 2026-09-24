package main

type pf_hund uint32 // unit is a hundred

/// Paths

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
)

///

const (
	UNSET       uint16 = 65535
	NEVERSET    uint16 = 65534
	OUT_OF_PAGE uint16 = 65535
	UNDEF              = "UNDEFINED"
	UNSET_S            = "UNSET"
	NONE               = "NONE"
)

const (
	CONSOLE_PROMPT = DC_BOLD + "console>" + DC_RESET
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

//////////////////////////////////////////////////////////////////////////////h1
//	DB

// Values for unknown tags.
const (
	UNKNOWN_GENRE  = "Unknown genre"
	UNKNOWN_ARTIST = "Unknown artist"
	UNKNOWN_ALBUM  = "Unknown album"
	UNKNOWN_TITLE  = "Unknown title"
	UNKNOWN_YEAR   = "Unknown year"
)

//////////////////////////////////////////////////////////////////////////////h1
//

