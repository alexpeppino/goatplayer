package main

import (
	"os"
)

type cFilters struct {
	sl    []fContainerFilter
	b     []bool
	ty    []string
	value []string
}

type cLogical struct {
	b []bool
}

var (
	Term vTerm
	///	Screens
	MainScreen     vMainScreen
	SettingsDialog vSettingsDialog
	// HelpDialog     vHelpDialog
	FindDialog vFindDialog
	/// Workspaces
	BrowserWsp1 mBrowserWsp
	BrowserWsp2 mBrowserWsp
	PlayerWsp   mPlayerWsp
	// EditorWsp   mEditorWsp
	// WsEditor mEditorWorkspace
	///	Models
	activeBrowser *mBrowserWsp
	View2         *mFilelist
	View1         *mFilelist
	Play          *mPlaylist
	// Edit          *mEditor

	// PlayerX  mPlayerX
	Find   mFind
	Report mReport
	SetM   mSettings
	Colors mColors
	Help   mHelp
	// Bios      mHelp
	Equalizer mEqualizer
	//
	Display tDisplaySection
	// RemConsole   uRemote /*c*/
	FocusToModel *uFocusKeys
	// input        input.GInput
	// AF aAudiofiles
	// musicscanner musicscanner.GMusicScanner
	Mrg         t__Margins
	PFListsData aPFLists
	// PFLists      mPFListsCon
	Set gSettings
	// MusicPlayer  mMusicPlayer
	//
	activeSection    *vSection
	oldActiveSection *vSection
	logFile          *os.File
	ui_pl            tUiToPlayer
	pl_ft            tPlayerToFiletime
	dialogMode       bool

// // trace            trace.gTrace /*c*/
// trace2       *cTrace
// Mpv          cMpv
)

type tPlayFaster struct {
}

type tUiToPlayer struct {
	signal             chan uint16
	playerIsWaiting    bool
	stopPlaying        bool
	playlistHasChanged bool
	playSelectedFile   bool
	RepeatList         bool
	Speed_i            uint16 // percentage
	Speed              string
	Volume_i           uint16
	Volume             string
	isVerbose          bool
	startPoint         uint32
}

type tPlayerToFiletime struct {
	signal            chan uint16
	stop              bool
	filetimeIsWaiting bool
	duration          uint32
}

type tRect struct {
	marginTop    uint16
	marginBottom uint16
	marginLeft   uint16
	marginRight  uint16
}
