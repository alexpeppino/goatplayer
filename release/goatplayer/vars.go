package main

import "os"

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
	Display      tDisplaySection
// RemConsole   uRemote /*c*/
	FocusToModel *uFocusKeys
	Input        uInput
// trace        cTrace /*c*/
	AF           aAudiofiles
	MusicScanner cMusicScanner
	Mpv          cMpv
	Mrg          t__Margins
	Pen          t__Pencils
	PFListsData  aPFLists
	// PFLists      mPFListsCon
	Set cSettings
	// trace2       *cTrace
	// MusicPlayer  mMusicPlayer
	//
	activeSection    *vSection
	oldActiveSection *vSection
	logFile          *os.File
	ui_pl            tUiToPlayer
	pl_ft            tPlayerToFiletime
	dialogMode       bool
)

type fContainerFilter func(audiof *aAudiofile) bool

type tTableKey struct {
	v              string
	iLine          uint16
	dbIndex        uint16
	isActiveFilter bool
}

type tMpvJson struct {
	Data       float32 `json:"data"`
	Request_id int     `json:"request_id"`
	Error      string  `json:"error"`
}

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

type f_SortSlice func(i int, j int) bool
type f_FilelistFilter func(audiof *aAudiofile, i int) bool
type f_FilelistFilterUi func(audiof *aAudiofile, i uint16) bool
type f_FieldFilter func(audiof *aAudiofile) string
type f_AudiofFilter func(audiof *aAudiofile) bool

type tRect struct {
	marginTop    uint16
	marginBottom uint16
	marginLeft   uint16
	marginRight  uint16
}

