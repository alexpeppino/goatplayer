/*  */
package main

import (
	"bufio"
	"fmt"
	"goatplayer/internal/K"
	AF "goatplayer/internal/audiofiles"
	"goatplayer/internal/echo"
	"goatplayer/internal/help"
	"goatplayer/internal/input"
	"goatplayer/internal/mpv"
	"goatplayer/internal/musicscanner"
	"goatplayer/internal/small/utils"
	"goatplayer/internal/style"
	"os"
	"os/signal"
	"syscall"
)

type fContainerFilter func(audiof *AF.TAudiofile) bool
type f_SortSlice func(i int, j int) bool
type f_FilelistFilter func(audiof *AF.TAudiofile, i int) bool
type f_FilelistFilterUi func(audiof *AF.TAudiofile, i uint16) bool
type f_FieldFilter func(audiof *AF.TAudiofile) string
type f_AudiofFilter func(audiof *AF.TAudiofile) bool

func main() {
	_ = help.Sections
	var info os.FileInfo
	var err error
	defer E()
	Term.SttyEchoOff()

	//// * * * * * * * * * * * * * * * *
	/// Init before db

	Term.SetTermTitle("GoatPlayer")
	Term.SttyCBreakMin1()

	info, err = os.Stat(K.FILE_MUSIC_DIR)

	/// Check if InitFile exists
	if err != nil {

		/// InitFile does not exist yet. Start initialization
		fmt.Printf("Initializing...\n")

		/// Ask for music library path
		fmt.Println("Type the Music Library directory's name (without the final slash):")
		musicDir := input.ReadLineFromConsole()
		info, err = os.Stat(musicDir)

		/// Check if musicDir from user is valid
		if err == nil {
			if info.IsDir() {

				/// MusicDir is a valid dir
				/// I N I T I A L I Z E

				fmt.Println("This is a valid directory.")
				AF.MusicDir = musicDir

				/// Making subdirs
				fmt.Printf("Making subdir '%s'...\n", K.DIR_SETTINGS)
				os.Mkdir(K.DIR_SETTINGS, K.PERM_DIR)
				fmt.Printf("Making subdir '%s'...\n", K.DIR_DATA)
				os.Mkdir(K.DIR_DATA, K.PERM_DIR)
				/* // b //
				fmt.Printf("Making subdir '%s'...\n", K.DIR_LOG)
				os.Mkdir(K.DIR_LOG, K.PERM_DIR)
				*/ // e //

				/// Save musicDir into init file
				fd, _ := os.OpenFile(K.FILE_MUSIC_DIR, os.O_CREATE|os.O_WRONLY, K.PERM_OPEN_FILE)
				fd.WriteString(musicDir)
				fd.Close()
				// fmt.Printf("File '%s' created.\n", K.FILE_MUSIC_DIR)
				fmt.Println("Setting initial values for the equalizer...")
				Equalizer.ResetValues()
			}
		} else {
			/// MusicDir is not valid
			fmt.Println("This is not a valid directory. Exiting...")
			Term.Quit()
		}
	} else {
		/// InitFile exists
		fd, _ := os.OpenFile(K.FILE_MUSIC_DIR, os.O_CREATE|os.O_RDONLY, K.PERM_OPEN_FILE)
		sc := bufio.NewScanner(fd)
		sc.Scan()
		AF.MusicDir = sc.Text()
		fmt.Printf("Music Library directory is '%s'\n", AF.MusicDir)
		fd.Close()
	}

	AllMaps.Init()

	/// Trap ctrl-c and SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// os.Remove("/tmp/echo.sock")
		// trace.Print("QUIT (CTRL-C / SIGTERM)") // t //
		Term.Quit()
	}()

	/* // b //
	logFile, err = os.OpenFile(K.FILE_LOG, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, K.PERM_OPEN_FILE)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logFileErr, err := os.OpenFile(K.FILE_ERRLOG, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, K.PERM_OPEN_FILE)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	defer logFileErr.Close()
	log.SetOutput(&errorBuffer)
	// log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	*/ // e //

	// trace.InitTrace(logFile, logFileErr) // t //
	// trace2 = &trace

	var print1 = func(s string) {
		// trace.PrintNoTab("* * * * * * * * " + strings.ToUpper(s)) // t //
	}
	var print2 = func(s string) {
		// trace.PrintNoTab("* * * * " + s) // t //
	}

	/// * * * * * * * * * * * * * * * *
	print1("init before db")

	PFListsData.InitPFList("PFLists.data")

	print2("init all")
	for _, in := range AllMaps.allScreens_in {
		in.InitScreen()
	}
	View1 = &BrowserWsp1.View
	View2 = &BrowserWsp2.View
	Play = &PlayerWsp.Player
	// Edit = &EditorWsp.Edit
	InitAll()
	Play.FileProp = &PlayerWsp.AudioPropCon
	AllMaps.allBrowsers = append(AllMaps.allBrowsers, &BrowserWsp1, &BrowserWsp2)
	activeBrowser = &BrowserWsp1
	Set.consoleModeOn = false
	/// trace init
	//	data init
	ui_pl.playerIsWaiting = true
	ui_pl.stopPlaying = false
	ui_pl.isVerbose = false
	ui_pl.playlistHasChanged = false
	ui_pl.playSelectedFile = false
	ui_pl.RepeatList = true
	ui_pl.Speed_i = 112
	ui_pl.Volume_i = 100
	pl_ft.filetimeIsWaiting = true
	pl_ft.stop = false
	ui_pl.signal = make(chan uint16)
	pl_ft.signal = make(chan uint16)
	/// Objects init

	Term.GetTermSize()
	/// Check if term size is large enough
	var minCols = 238
	var minRows = 55
	if Term.nRows < uint16(minRows) || Term.nCols < uint16(minCols) {
		fmt.Printf("Terminal size is too small (%d x %d). At least %d x %d is required.\n", Term.nRows, Term.nCols, minRows, minCols)
		Term.Quit()
	}
	fmt.Printf("Terminal size is %d x %d.", Term.nRows, Term.nCols)

	Set.InitSettings()
	Set.LoadSettings() // after init
	// Set.TraceSettings()
	// Set.RepeatMode.value = "ccc"
	// Set.TraceSettings()
	// FocusToModel.InitFocus()
	// RemConsole.InitRemote() /*c*/
	input.InitInput(&Term)
	Term.MouseOn()
	style.Pen.InitPencils()
	style.InitStyles()
	SetOuters()
	musicscanner.InitMusicScanner()
	// SetWindows(50)
	MusicDisplay.marginTop = 53
	MusicDisplay.marginLeft = 1
	_f()

	/// * * * * * * * * * * * * * * * *
	print1("db: load music library data")

	// Term.SttyEchoOn()
	Term.SttyCBreakMin1()
	Term.SttyEchoOff()

	// os.Mkdir("settings", 0774)
	// fd, _ := os.OpenFile("settings/file1", os.O_CREATE|os.O_RDWR, 0774)
	// fd.Close()

	if err := AF.CheckDb(); err != nil {
		Term.quitMessage = err.Error()
		Term.Quit()
	}
	Term.SttyEchoOff()
	/// { db-is-ready } ///

	/// * * * * * * * * * * * * * * * *
	/// Init after db

	print1("init after db")
	MainScreen.TUIMode()
	Term.AltScreenOn()
	// Term.Clear()
	BrowserWsp1.InitBrowserWsp("BrowserWsp1", "Browser1")
	BrowserWsp2.InitBrowserWsp("BrowserWsp1", "Browser2")
	PlayerWsp.InitPlayerWsp("PlayerWsp", "Player")
	// EditorWsp.InitEditorWsp("Editor", "Editor")

	/// * * * * * * * * * * * * * * * *
	print1("prepare")

	// DB.SaveAudiofDB()
	PFListsData.LoadDataFromFile()
	Term.StatusLine.PrepareStatusLine(Term.nRows-1, Term.nCols-28, 5)
	// trace.Print("%d %s", DB.filenamesMapI[DB.filenames[0]], DB.filenames[0])
	// trace.Print("%d %s", DB.filenamesMapI[DB.filenames[1]], DB.filenames[1])
	// os.Exit(0)

	for _, ie := range AllMaps.allModels_in {
		ie.PrepareModel()
	}

	View1.pRightSection = &MainScreen.Browser1Section
	View2.pRightSection = &MainScreen.Browser2Section
	Play.pRightSection = &MainScreen.PlayerSection
	// Edit.pRightSection = &MainScreen.EditorSection

	BrowserWsp1.PrepareContainers(&AF.AllDbIndexes)
	BrowserWsp2.PrepareContainers(&AF.AllDbIndexes)

	print2("prepare screens")
	for _, in := range AllMaps.allScreens_in {
		in.PrepareScreen()
	}
	Term.activeScreen = &MainScreen.vScreen
	input.PrepareInput(MainScreen.CommandPanelSection.marginBottom+2, MainScreen.CommandPanelSection.marginLeft)

	CommandPanel.PrepareCommandPanel()

	// PrepareScreens()

	// Albums.Prepare("Albums", &FileContainers, K.FOCUS_ALBUMS)
	// Artists.Prepare("Artists", &FileContainers, K.FOCUS_ARTISTS)
	// Folded.Prepare("Folded", K.FOCUS_FOLDED)
	// PfLists.Prepare("PfLists", K.FOCUS_PFLISTS)
	// Pls.Prepare("Pls", K.FOCUS_PLS)
	// Tags.Prepare("Genres", K.FOCUS_TAGS)
	// Browser1.Artists.keys.IActive = 0
	// Browser2.Artists.keys.IActive = 0
	// activePlt = &FileContainers

	SetCommands()
	// Equalizer.w.fGetWidth = func() uint16 { return 20 }

	for _, p := range AllMaps.allBasicModels {
		p.CheckDoubleCommandKeys()
	}

	for _, p := range AllMaps.allBasicModels {
		p.PrepareCommandsEcho()
	}

	/// * * * * * * * * * * * * * * * *
	print1("dump everything")

	// dump.TraceAllScreens() // t //
	// dump.DumpAllModels_in()
	// dump.DumpAllScreens_in()

	/// * * * * * * * * * * * * * * * *
	print1("load data into workspaces")

	BrowserWsp1.GenresCon.FilterAndPrint__new()
	BrowserWsp2.GenresCon.FilterAndPrint__new()
	BrowserWsp1.AudioPropCon.LoadFile(View1.keys.GetSelected().V)
	BrowserWsp2.AudioPropCon.LoadFile(View2.keys.GetSelected().V)
	PFListsData.LoadDataFromFile()
	// PFLists.LoadListnames()
	// EditorWsp.PFLists.LoadListnames()

	//! PlayerWsp.PFListsCon.LoadListnames()

	// BrowserWsp1.ArtistsCon.keys.Select(BrowserWsp1.ArtistsCon.keys.GetIndexByFilename("Pink Floyd"))
	// BrowserWsp2.ArtistsCon.keys.Select(BrowserWsp2.ArtistsCon.keys.GetIndexByFilename("Pink Floyd"))
	// PLView1.LoadArtist()
	// PLView2.LoadArtist()

	// Folded.Load()
	// Tags.LoadGenres()
	//

	/// * * * * * * * * * * * * * * * *
	print1("prepare views")

	// ActivePL.PrepareLines()
	PrepareViews()
	/// { windows-are-ready } ///
	/// CLEAR SCREEN

	/// * * * * * * * * * * * * * * * *
	print1("prepare term")

	// Term.Clear()

	/// Go
	go Player_go(ui_pl.signal)
	go PlayerTime_go(pl_ft.signal)
	// go RemoteInterface_go() /*c*/
	// Monitor.InitMonitor()
	// chMonitor = make(chan uint16)
	// go Monitor_go(chMonitor)
	// go func() {
	// 	for {
	// 		Sleep(2)
	// 		Monitor.UpdateFile()
	// 	}
	// }()
	///	Print
	Report.AddReportLine("Music Library: %s", AF.MusicDir)
	speedFormat()
	volumeFormat()
	Term.name = "Term"
	Term.ResumeLastSession()
	///
	MainScreen.LeftSection.PrintFocusButtons(true)
	// FileContainers.PrintCaptions()
	// butSpeedM.PrintUnpressed()
	// butSpeedP.PrintUnpressed()
	// butPlayNext.PrintUnpressed()
	Report.AddReportLine("Audiofiles list created with %d entries.", len(AF.Sl))
	Report.AddReportLine("Size of the text area: %d x %d (minimum is %d x %d)", Term.nRows, Term.nCols, minRows, minCols)
	if AF.NewN > 0 {
		Report.AddReportLine("%d new files found in the Music Library", AF.NewN)
	}
	if AF.MissingN > 0 {
		Report.AddReportLine("%d files missing in the Music Library", AF.MissingN)
	}
	// get_input(56 ,56)
	/// Focus
	// Focus.activeKey = K.FOCUS_VIEW2
	///
	echo.AddStringWithPos(Term.nRows, Term.nCols-21, utils.Bigstring("GOATPLAYER"))
	echo.Flush()

	Report.AddReportLine("Found %d directories in the Music Library.", len(musicscanner.ValidDirs))
	SetWorkspaceAndDialogFocusKeys()

	/// * * * * * * * * * * * * * * * *
	print1("start mainscreen")

	Term.doSaveAtQuit = true
	MainScreen.KeyboardMode()
	// for {
	// 	MouseMode()
	// 	Focus.activeKey = Screen.actWindow.focus
	// }
}

func SetOuters() {
	// PLView1.SetOuter()
	// PLView2.SetOuter()
	BrowserWsp1.AudioPropCon.SetOuter()
	// BrowserWsp1.Dirs.SetOuter()
	BrowserWsp2.AudioPropCon.SetOuter()
	// BrowserWsp2.Dirs.SetOuter()
	Find.SetOuter()
	MainScreen.SetOuter()
	// HelpDialog.SetOuter()
	SettingsDialog.SetOuter()
	FindDialog.SetOuter()
	// ActivePL.SetOuter()
	// Edit.SetOuter()
	PlayerWsp.AudioPropCon.SetOuter()
	Report.SetOuter()
	// PFLists.SetOuter()
	Colors.SetOuter()
	Colors.m2.SetOuter()
	//
	SetM.SetOuter()
	// Help.SetOuter()
	// Bios.SetOuter()
	Equalizer.SetOuter()
	Equalizer.m2.SetOuter()
}

func SetCommands() {
	// trace.Begin("SetCommands") // t //
	for _, ie := range AllMaps.allModels_in {
		ie.SetCommands()
	}
	// trace.End() // t //
}

func InitAll() {
	Term.InitTerm()
	/// BrowserWsp1
	help.Init()
	BrowserWsp1.View.InitModel("View1", K.FOCUS_VIEW1, K.HELPFILE_BROWSERS)
	BrowserWsp1.GenresCon.InitModel("Genres", K.FOCUS_GENRES, K.HELPFILE_TAG_FIELDS)
	BrowserWsp1.ArtistsCon.InitModel("Artists", K.FOCUS_ARTISTS, K.HELPFILE_TAG_FIELDS)
	BrowserWsp1.AlbumsCon.InitModel("Albums", K.FOCUS_ALBUMS, K.HELPFILE_TAG_FIELDS)
	BrowserWsp1.YearsCon.InitModel("Years", K.FOCUS_YEARS, K.HELPFILE_TAG_FIELDS)
	BrowserWsp1.AlbumArtistsCon.InitModel("AlbumArtists", K.FOCUS_ALBUM_ARTIST, K.HELPFILE_TAG_FIELDS)
	BrowserWsp1.ComposersCon.InitModel("Composers", K.FOCUS_COMPOSER, K.HELPFILE_TAG_FIELDS)
	// BrowserWsp1.CommentsCon.InitModel("Comments", K.FOCUS_COMMENT, K.HELPFILE_TAGS)
	BrowserWsp1.Dirs.InitModel("Dirs", K.FOCUS_DIRS, K.HELPFILE_TAG_FIELDS)
	BrowserWsp1.AudioPropCon.InitModel("Properties", K.FOCUS_FILEPROP, K.HELPFILE_COLORS)
	/// BrowserWsp2
	BrowserWsp2.View.InitModel("View2", K.FOCUS_VIEW2, K.HELPFILE_BROWSERS)
	BrowserWsp2.GenresCon.InitModel("Genres", K.FOCUS_GENRES, K.HELPFILE_TAG_FIELDS)
	BrowserWsp2.ArtistsCon.InitModel("Artists", K.FOCUS_ARTISTS, K.HELPFILE_TAG_FIELDS)
	BrowserWsp2.AlbumsCon.InitModel("Albums", K.FOCUS_ALBUMS, K.HELPFILE_TAG_FIELDS)
	BrowserWsp2.YearsCon.InitModel("Years", K.FOCUS_YEARS, K.HELPFILE_TAG_FIELDS)
	BrowserWsp2.AlbumArtistsCon.InitModel("AlbumArtists", K.FOCUS_ALBUM_ARTIST, K.HELPFILE_TAG_FIELDS)
	BrowserWsp2.ComposersCon.InitModel("Composers", K.FOCUS_COMPOSER, K.HELPFILE_TAG_FIELDS)
	// BrowserWsp2.CommentsCon.InitModel("Comments", K.FOCUS_COMMENT, K.HELPFILE_TAGS)
	BrowserWsp2.Dirs.InitModel("Dirs", K.FOCUS_DIRS, K.HELPFILE_TAG_FIELDS)
	BrowserWsp2.AudioPropCon.InitModel("Properties", K.FOCUS_FILEPROP, K.HELPFILE_COLORS)
	/// PlayerWsp
	PlayerWsp.Player.InitModel("Player", K.FOCUS_PLAY, K.HELPFILE_TAG_FIELDS)
	PlayerWsp.AudioPropCon.InitModel("Properties", K.FOCUS_FILEPROP, K.HELPFILE_COLORS)
	//! PlayerWsp.Param.InitModel("Parameters", "X", K.HELPFILE_PARAMETERS)
	// PlayerWsp.PFListsCon.InitModel("Lists", "L", K.HELPFILE_PFLISTS)

	/// EditorWsp
	// EditorWsp.Edit.InitModel("Edit", K.FOCUS_EDIT)
	// EditorWsp.Param.InitModel("Param", "B")
	// EditorWsp.PFLists.InitModel("PFLists", "L")
	/// SettingsDialog
	SetM.InitModel("Settings", K.FOCUS_SETTINGS, K.HELPFILE_SETTINGS)
	Equalizer.InitModel("Equalizer", K.FOCUS_EQUALIZER, K.HELPFILE_EQUALIZER)
	Equalizer.m2.InitModel("Equalizer2", K.FOCUS_EQUALIZER2, K.HELPFILE_EQUALIZER)
	Colors.InitModel("Colors", K.FOCUS_COLORS, K.HELPFILE_COLORS)
	Colors.m2.InitModel("Colors2", K.FOCUS_COLORS2, K.HELPFILE_COLORS)
	Colors.m2.w.isCanvas = true
	/// HelpDialog
	Help.InitModel("Help", K.FOCUS_HELP, "help")
	// Bios.InitModel("Bios", K.FOCUS_BIOS, "bios")
	Report.InitModel("Report", K.FOCUS_REPORT, K.HELPFILE_COLORS)
	/// FindDialog
	Find.InitModel("Find", K.FOCUS_FIND, K.HELPFILE_COLORS)
	///
	// PFLists.InitModel("PFLists", K.FOCUS_PFLISTS, K.HELPFILE_PFLISTS)
	Help.w2.pModel = &Help.mBasic
	Help.w2.name = Help.name + ".w2"
	// Bios.w2.pModel = &Bios.mBasic
	// Bios.w2.name = Bios.name + ".w2"
}

func PrepareViews() {
	for _, ie := range AllMaps.allModels_in {
		// trace.Print("prepare view for: %s", ie.GetName()) // t //
		ie.PrepareView()
	}
}

func PrepareScreens() {
}

func E() {
	// trace.Print("-- deferred E --") // t //
	Term.MouseOff()
	Term.SttyEchoOn()
	mpv.Stop()
}
