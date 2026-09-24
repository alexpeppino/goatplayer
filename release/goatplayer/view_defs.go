package main

//////////////////////////////////////////////////////////////////////////////h1
/*

v  =>  Prefix for the "view" class of types

The Screen manages a slice of the foregrounded Windows.
The foreground/background status for all Windows is managed by tutor objects
of type Section, which contain Layers.
The Screen is partitioned into first level Sections.
Each section has multiple layers and each layer contains one or more Windows.
Or, instead of a Window, it may contain a sub-section, so that in that section there
are layers and Windows.

Print method: Only the Screen has power to actually print to the terminal.
Windows can only add-to-print: they add their content to the Screen content that will
be printed evetually. The Screen knows when it's time to re-print.

class	type
=====	====
v		Screen, has Sections
		Section, has Layers
		Layer, has Windows and Subsections
		Window
		  Is the basic type for:
		  - AFContainer type (audio-file container)
		      (Dirs, PFLists, Artists, Albums, Genres, ... objects)
		  - AFList type (audio-file list)
		  	  (ActivePL, OpenPL, ... objects)
		  - Sidepane object
		  - Player object
		...

*/

// type v__Container struct {
// 	v__Window
// }

//////////////////////////////////////////////////////////////////////////////h1
//	Platoon

// Platoon type. A window that holds many window-tabs (tables).
// type w__MultiLayer struct {
// 	v__Window
// 	lSl []t__Layer

// 	selectedTable *t__Table
// 	activeTable   *t__Table

// 	name         string
// 	marginTop    uint16
// 	marginBottom uint16
// 	marginLeft   uint16
// 	marginRight  uint16
// 	caption      string

// 	frameLight string
// 	frameBold  string
// }

type t__Pencils struct {
	Normal          string
	NormalInv       string
	LightNormal     string
	Unselected      string
	Selected        string
	Active          string
	SelectedActive  string
	StatusLine      string
	Frame           string
	EmptyFrame      string
	SingSong        string
	ShowMargins     string
	Tab             string
	SelectedTab     string
	Grey            string
	ButtonUnpressed string
	ButtonPressed   string
	Error           string
}

type t__Column struct {
	caption string
	width   uint16
	style   *cStyle
}

type t__ColumnLine struct {
	text  string
	width uint16
	style *cStyle
}

//////////////////////////////////////////////////////////////////////////////h1

// Frames
//
// 0	up-left
// 1	up-right
// 2	down-left
// 3	down-right
// 4	vertical
// 5	horizontal
type frames_ struct {
	Bold          [6]string
	Light         [6]string
	Double        [6]string
	Empty         [6]string
	Light_rounded [6]string
}

var Frames = frames_{
	Bold:          [6]string{"┏", "┓", "┗", "┛", "┃", "━"},
	Light:         [6]string{"┌", "┐", "└", "┘", "│", "─"},
	Double:        [6]string{"╔", "╗", "╚", "╝", "║", "═"},
	Empty:         [6]string{" ", " ", " ", " ", " ", " "},
	Light_rounded: [6]string{"╭", "╮", "╰", "╯", "│", "─"}}

//////////////////////////////////////////////////////////////////////////////h1

// Margins
type t__Margins struct {
	ui_Cols uint16
	ui_Rows uint16

	ui_Top    uint16
	ui_Bottom uint16
	ui_Left   uint16
	ui_Right  uint16
	ui_Width  uint16
	ui_Height uint16

	sect_1_Top                uint16
	sect_1_Bottom             uint16
	sect_1_Height             uint16
	sect_1_HeightMin          uint16
	sect_1A_Left              uint16
	sect_1A_Right             uint16
	sect_1A_Width             uint16
	sect_1A_WidthMin          uint16
	sect_1B_Left              uint16
	sect_1B_Right             uint16
	sect_1B_WidthFix          uint16
	sect_1A_SelectionBarLeft  uint16
	sect_1A_SelectionBarWidth uint16
	sect_1A_Height            uint16
	sect_1A_Bottom            uint16

	sect_2_Top       uint16
	sect_2_Bottom    uint16
	sect_2_Height    uint16
	sect_2_HeightMin uint16
	sect_2A_Left     uint16
	sect_2A_Right    uint16
	sect_2A_Width    uint16
	sect_2B_Left     uint16
	sect_2B_Right    uint16
	sect_2B_WidthFix uint16

	sect_3_Top       uint16
	sect_3_Bottom    uint16
	sect_3_HeightFix uint16
	sect_3A_Left     uint16
	sect_3A_Right    uint16
	sect_3A_Width    uint16
	sect_3B_Left     uint16
	sect_3B_Right    uint16
	sect_3B_WidthFix uint16

	sect_4_Top       uint16
	sect_4_HeightFix uint16
	sect_4A_Left     uint16
	sect_4A_Right    uint16
	sect_4A_Width    uint16
	sect_4B_Left     uint16
	sect_4B_Right    uint16
	sect_4B_WidthFix uint16
}

const (
	tputScCivis = "\x1b7\x1b[?25l"
	tputRcCnorm = "\x1b8\x1b[?12l\x1b[?25h"
	TPUT_SC     = "\x1b7"
	TPUT_RC     = "\x1b8"
	TPUT_CIVIS  = "\x1b[?25l"
	TPUT_CNORM  = "\x1b[?12l\x1b[?25h"
	TPUT_EL     = "\x1b[K" // erase line
	// TPUT_EL1      = "\x1b[1K" // erase line to beginning
	TPUT_CUP_INIT = "\x1b[57;1H"
	TPUT_CLEAR    = "\x1b[H\x1b[2J\x1b[3J"
	CLR           = "\x1b[38;5;%d;48;5;%dm" // color

	TPUT__ALTSCREEN_ON  = "\x1b[?1049h"
	TPUT__ALTSCREEN_OFF = "\x1b[?1049l"
)

//
//	\E7\E[?25l
//	tput sc; tput civis

const (
	ALIGN_LEFT = iota
	ALIGN_RIGHT
)

const (
	FONT_NORMAL         = ""
	FONT_BOLD           = "\x1b[1m"
	FONT_BOLD_UNDERLINE = "\x1b[1m\x1b[4m"
	FONT_ITALIC         = "\x1b[3m"
	FONT_UNDERLINE      = "\x1b[4m"
	FONT_RESET          = "\x1b[0m"
	FONT_BLINK          = "\x1b[5m"
)

const (
	DC_RESET     = "\x1b[0m"
	DC_BOLD      = "\x1b[1m"
	DC_ITALIC    = "\x1b[3m"
	DC_UNDERLINE = "\x1b[4m"
)

//////	Colors

