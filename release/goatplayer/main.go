/*  */
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var info os.FileInfo
	var err error
	defer E()
	Term.SttyEchoOff()

	//// * * * * * * * * * * * * * * * *
	/// Init before db

	SetTerminalTitle("GoatPlayer")
	Term.SttyCBreakMin1()

	info, err = os.Stat(FILE_MUSIC_DIR)

	/// Check if InitFile exists
	if err != nil {

		/// InitFile does not exist yet. Start initialization
		fmt.Printf("File '%s' does not exist yet. Initializing...\n", FILE_MUSIC_DIR)

		/// Ask for music library path
		fmt.Println("Type the Music Library directory's name (without the final slash):")
		musicDir := ReadLineFromConsole()
		info, err = os.Stat(musicDir)

		/// Check if musicDir from user is valid
		if err == nil {
			if info.IsDir() {

				/// MusicDir is a valid dir
				/// I N I T I A L I Z E

				fmt.Println("This is a valid directory.")
				AF.musicDir = musicDir

				/// Making subdirs
				// fmt.Printf("Making subdir '%s'...\n", DIR_SETTINGS)
				// os.Mkdir(DIR_SETTINGS, 0774)
				// fmt.Printf("Making subdir '%s'...\n", DIR_DATA)
				// os.Mkdir(DIR_DATA, 0774)
				// fmt.Printf("Making subdir '%s'...\n", DIR_LOG)
				// os.Mkdir(DIR_LOG, 0774)

				/// Save musicDir into init file
				fd, _ := os.OpenFile(FILE_MUSIC_DIR, os.O_CREATE|os.O_WRONLY, 0664)
				fd.WriteString(musicDir)
				fd.Close()
				fmt.Printf("File '%s' created.\n", FILE_MUSIC_DIR)
				fmt.Println("Setting initial values for the equalizer...")
				Equalizer.ResetValues()
			}
		} else {
			/// MusicDir is not valid
			fmt.Println("This is not a valid directory. Exiting...")
			Quit()
		}
	} else {
		/// InitFile exists
		fmt.Printf("File '%s' already exists.\n", FILE_MUSIC_DIR)
		fd, _ := os.OpenFile(FILE_MUSIC_DIR, os.O_CREATE|os.O_RDONLY, 0664)
		sc := bufio.NewScanner(fd)
		sc.Scan()
		AF.musicDir = sc.Text()
		fmt.Printf("Music Library directory is '%s'", AF.musicDir)
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
		Quit()
	}()

/* // b //
	logFile, err = os.OpenFile(FILE_LOG, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logFileErr, err := os.OpenFile(FILE_ERRLOG, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
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
// trace.PrintNoTab(traceStyleFocus.Echo("* * * * * * * * " + strings.ToUpper(s))) // t //
	}
	var print2 = func(s string) {
// trace.PrintNoTab(traceStyleFocus.Echo("* * * * " + s)) // t //
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
	Set.InitSettings()
	Set.LoadSettings() // after init
	// Set.TraceSettings()
	// Set.RepeatMode.value = "ccc"
	// Set.TraceSettings()
	// FocusToModel.InitFocus()
// RemConsole.InitRemote() /*c*/
	Input.InitInput()
	Input.MouseOn()
	Pen.InitPencils()
	InitStyles()
	SetOuters()
	MusicScanner.InitMusicScanner()
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

	AF.CheckDb()
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

	AF.PrepareContainers(&BrowserWsp1, &AF.allDbIndexes)
	AF.PrepareContainers(&BrowserWsp2, &AF.allDbIndexes)

	print2("prepare screens")
	for _, in := range AllMaps.allScreens_in {
		in.PrepareScreen()
	}
	Term.activeScreen = &MainScreen.vScreen

	CommandPanel.PrepareCommandPanel()

	// PrepareScreens()

	// Albums.Prepare("Albums", &FileContainers, FOCUS_ALBUMS)
	// Artists.Prepare("Artists", &FileContainers, FOCUS_ARTISTS)
	// Folded.Prepare("Folded", FOCUS_FOLDED)
	// PfLists.Prepare("PfLists", FOCUS_PFLISTS)
	// Pls.Prepare("Pls", FOCUS_PLS)
	// Tags.Prepare("Genres", FOCUS_TAGS)
	// Browser1.Artists.keys.iActive = 0
	// Browser2.Artists.keys.iActive = 0
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
	BrowserWsp1.AudioPropCon.LoadFile(View1.keys.GetSelected().v)
	BrowserWsp2.AudioPropCon.LoadFile(View2.keys.GetSelected().v)
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
	Monitor.InitMonitor()
	chMonitor = make(chan uint16)
	go Monitor_go(chMonitor)
	// go func() {
	// 	for {
	// 		Sleep(2)
	// 		Monitor.UpdateFile()
	// 	}
	// }()
	///	Print
	Report.AddReportLine("Music Library: %s", AF.musicDir)
	speedFormat()
	volumeFormat()
	Term.name = "Term"
	Term.RestoreLastSession()
	///
	MainScreen.LeftSection.PrintFocusButtons(true)
	// FileContainers.PrintCaptions()
	// butSpeedM.PrintUnpressed()
	// butSpeedP.PrintUnpressed()
	// butPlayNext.PrintUnpressed()
	Report.AddReportLine("Audiofiles list created with %d entries.", len(AF.sl))
	Report.AddReportLine("Size of the text area: %d x %d", Term.nRows, Term.nCols)
	if nNew > 0 {
		Report.AddReportLine("%d new files found in the Music Library", nNew)
	}
	if nMissing > 0 {
		Report.AddReportLine("%d files missing in the Music Library", nMissing)
	}
	// get_input(56 ,56)
	/// Focus
	// Focus.activeKey = FOCUS_VIEW2
	///
	Term.AddStringWithPos(Term.nRows, Term.nCols-21, utils.Bigstring("GOATPLAYER")).Flush()

	Report.AddReportLine("Found %d directories in the Music Library.", len(MusicScanner.validDirs))
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
	BrowserWsp1.View.InitModel("View1", FOCUS_VIEW1, HELPFILE_VIEW1)
	BrowserWsp1.GenresCon.InitModel("Genres", FOCUS_GENRES, HELPFILE_TAGS)
	BrowserWsp1.ArtistsCon.InitModel("Artists", FOCUS_ARTISTS, HELPFILE_TAGS)
	BrowserWsp1.AlbumsCon.InitModel("Albums", FOCUS_ALBUMS, HELPFILE_TAGS)
	BrowserWsp1.YearsCon.InitModel("Years", FOCUS_YEARS, HELPFILE_TAGS)
	BrowserWsp1.AlbumArtistsCon.InitModel("AlbumArtists", FOCUS_ALBUM_ARTIST, HELPFILE_TAGS)
	BrowserWsp1.ComposersCon.InitModel("Composers", FOCUS_COMPOSER, HELPFILE_TAGS)
	// BrowserWsp1.CommentsCon.InitModel("Comments", FOCUS_COMMENT, HELPFILE_TAGS)
	BrowserWsp1.Dirs.InitModel("Dirs", FOCUS_DIRS, HELPFILE_DIRS)
	BrowserWsp1.AudioPropCon.InitModel("Properties", FOCUS_FILEPROP, HELPFILE_FILEPROP)
	/// BrowserWsp2
	BrowserWsp2.View.InitModel("View2", FOCUS_VIEW2, HELPFILE_VIEW1)
	BrowserWsp2.GenresCon.InitModel("Genres", FOCUS_GENRES, HELPFILE_TAGS)
	BrowserWsp2.ArtistsCon.InitModel("Artists", FOCUS_ARTISTS, HELPFILE_TAGS)
	BrowserWsp2.AlbumsCon.InitModel("Albums", FOCUS_ALBUMS, HELPFILE_TAGS)
	BrowserWsp2.YearsCon.InitModel("Years", FOCUS_YEARS, HELPFILE_TAGS)
	BrowserWsp2.AlbumArtistsCon.InitModel("AlbumArtists", FOCUS_ALBUM_ARTIST, HELPFILE_TAGS)
	BrowserWsp2.ComposersCon.InitModel("Composers", FOCUS_COMPOSER, HELPFILE_TAGS)
	// BrowserWsp2.CommentsCon.InitModel("Comments", FOCUS_COMMENT, HELPFILE_TAGS)
	BrowserWsp2.Dirs.InitModel("Dirs", FOCUS_DIRS, HELPFILE_DIRS)
	BrowserWsp2.AudioPropCon.InitModel("Properties", FOCUS_FILEPROP, HELPFILE_FILEPROP)
	/// PlayerWsp
	PlayerWsp.Player.InitModel("Player", FOCUS_PLAY, HELPFILE_PLAY)
	PlayerWsp.AudioPropCon.InitModel("Properties", FOCUS_FILEPROP, HELPFILE_FILEPROP)
	//! PlayerWsp.Param.InitModel("Parameters", "X", HELPFILE_PARAMETERS)
	// PlayerWsp.PFListsCon.InitModel("Lists", "L", HELPFILE_PFLISTS)

	/// EditorWsp
	// EditorWsp.Edit.InitModel("Edit", FOCUS_EDIT)
	// EditorWsp.Param.InitModel("Param", "B")
	// EditorWsp.PFLists.InitModel("PFLists", "L")
	/// SettingsDialog
	SetM.InitModel("Settings", FOCUS_SETTINGS, HELPFILE_SETTINGS)
	Equalizer.InitModel("Equalizer", FOCUS_EQUALIZER, HELPFILE_EQUALIZER)
	Equalizer.m2.InitModel("Equalizer2", FOCUS_EQUALIZER2, HELPFILE_EQUALIZER)
	Colors.InitModel("Colors", FOCUS_COLORS, HELPFILE_COLORS)
	Colors.m2.InitModel("Colors2", FOCUS_COLORS2, HELPFILE_COLORS)
	Colors.m2.w.isCanvas = true
	/// HelpDialog
	Help.InitModel("Help", FOCUS_HELP, "help")
	// Bios.InitModel("Bios", FOCUS_BIOS, "bios")
	Report.InitModel("Report", FOCUS_REPORT, HELPFILE_REPORT)
	/// FindDialog
	Find.InitModel("Find", FOCUS_FIND, HELPFILE_FIND)
	///
	// PFLists.InitModel("PFLists", FOCUS_PFLISTS, HELPFILE_PFLISTS)
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
	Input.MouseOff()
	Term.SttyEchoOn()
	Mpv.Stop()
}

