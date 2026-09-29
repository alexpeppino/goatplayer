package main

import (
	"fmt"
	"goatplayer/internal/K"
	AF "goatplayer/internal/audiofiles"
	"goatplayer/internal/echo"
	"goatplayer/internal/mpv"
	"goatplayer/internal/small/utils"
	"goatplayer/internal/style"
	"os"
	"os/exec"
	"strings"
)

func (this *vTerm) CheckButtons() {
	if Term.activeCommand != nil {
		PrintUnpressedCommand(Term.activeCommand)
		Term.activeCommand = nil
	}
	if Term.pPressedButton != nil {
		Term.pPressedButton.PrintUnpressed()
		Term.pPressedButton = nil
	}
}

func (this *vTerm) MouseOn() {
	fmt.Print(K.MOUSE_ON)
}

func (this *vTerm) MouseOff() {
	fmt.Print(K.MOUSE_OFF)
}

func (this *vTerm) SetTermTitle(s string) {
	fmt.Printf("%s0;%s\007", K.KEYB_OSC, s)
}

func (this *vTerm) TputCNorm() {
	fmt.Print(TPUT_CNORM)
}

func TputCivis() {
	fmt.Print(TPUT_CIVIS)
}

func (this *vTerm) Quit() {
	if this.doSaveAtQuit {
		PFListsData.SaveDataToFile()
		Set.SaveSettings()
		mpv.Stop()
	}
	os.Remove("/tmp/echo.sock")
	Term.MouseOff()
	this.TputCNorm()
	this.SttyEchoOn()
	// Term.Clear()
	if this.isAltScreenOn {
		this.AltScreenOff()
	}
	// fmt.Print("\x1b[0;0H")
	if AF.DoSortAndSaveAtQuit {
		fmt.Println("Sorting and saving the database...")
		AF.SortAudiofiles(AF.SortAudiofilesAATT)
		AF.SaveDataToFile()
	}
	if Term.quitMessage != "" {
		fmt.Println(Term.quitMessage)
	}
	fmt.Println()
	fmt.Println("Bye. Thank you for using GoatPlayer.")
	os.Exit(0)
}

func (this *vTerm) GetLongName() string { return this.name }

func (this *vTerm) ResumeLastSession() {
	// trace.Begin_n(this, "ResumeLastSession") // t //
	var w1, w2 *vWindow

	//!!
	// MainScreen.Browser1Section.activeWindow = &BrowserWsp1.GenresCon.w
	// MainScreen.Browser2Section.activeWindow = &BrowserWsp2.GenresCon.w
	BrowserWsp1.GenresCon.w.SetActiveToSection()
	BrowserWsp2.GenresCon.w.SetActiveToSection()

	///
	name := Set.LastSectionLeftWindow
	// trace.Print("Set.LastSectionLeftWindow: '%s' [%p]", name, AllMaps.mapBasicModels[name]) // t //
	// trace.Print("%+v", AllMaps.mapBasicModels)
	if name != "" {
		w1 = &AllMaps.mapBasicModels[name].w
	} else {
		w1 = &View1.w
	}
	name = Set.LastSectionRightWindow
	// trace.Print("Set.LastSectionRightWindow: '%s' [%p]", name, AllMaps.mapBasicModels[name]) // t //
	if name != "" {
		w2 = &AllMaps.mapBasicModels[name].w
	} else {
		w2 = &BrowserWsp1.GenresCon.w
	}
	///
	MainScreen.activeWorkspace = &BrowserWsp1.maWorkspace
	MainScreen.activeSection = &MainScreen.LeftSection
	MainScreen.activeWindow = w1
	MainScreen.LeftSection.SetActiveWindow(w1)
	MainScreen.pRightSection = w2.pLayer.pSection
	MainScreen.pRightSection.SetActiveWindow(w2)
	MainScreen.GetFocusMan_o().activeKey = w1.pModel.focusKey
	///
	w1.BringWindowToFront(true)
	// w2.BringWindowToFront(false)
	///
	Play.LoadPFList(PFLISTS_PLAYER)
	Play.w.PrintWindow()
	// trace.Print("Set.LastPlayedFile: %s", Set.LastPlayedFile) // t //
	if Set.LastPlayedFile != "" {
		Play.keys.Select(Play.keys.GetIndexByFilename(Set.LastPlayedFile))
		Play.w.lines.Select(Play.keys.GetSelected().LineI)
		ui_pl.startPoint = Set.LastTime
		Play.PlaySelected()
		index := AF.MepFilenameIndex[Set.LastPlayedFile]
		audiof := AF.GetAudiof(index)
		Report.AddReportLine("Resuming last session: %s - %s", audiof.Artist, audiof.Title)
	}
	// trace.End() // t //
}

// .
// func (this *vTerm) PrintEverything() {
// 	trace.Begin_n(this, "PrintEverything")
// 	dump.TraceTerm()
// 	if this.activeScreen == nil {
// 		trace.Error("activeScreen")
// 		return
// 	}
// 	if this.activeScreen.isDialog {
// 		this.activeScreen.activeWindow.BringWindowToFront(true)
// 	} else {
// 		// Screen.Print
// 	}
// }

func (this *vTerm) InitTerm() {
	this.allButtons = append(this.allButtons, &CommandPanel.butPageMinus, &CommandPanel.butPagePlus)
}

func (this *vTerm) GetTermSize() {
	var s string
	var b []byte
	b = make([]byte, 20)
	fmt.Printf("%s18t", K.KEYB_CSI)
	le, _ := os.Stdin.Read(b)
	b = b[0:le]
	s = string(b)
	// trace.Print("string: %q (%d)", s, len(s)) // t //
	i1 := strings.Index(s, ";")
	i2 := strings.LastIndex(s, ";")
	i3 := len(s) - 1
	fmt.Sscanf(s[i1+1:i2], "%d", &Term.nRows)
	fmt.Sscanf(s[i2+1:i3], "%d", &Term.nCols)
	// trace.Print("-- size of the text area: %d rows x %d cols --", Term.nRows, Term.nCols) // t //
	// SleepMilli(200)
	// trace.Print("-- size of the screen in characters --")
	// fmt.Printf("%s19t", K.KEYB_CSI)
	// SleepMilli(100)
	// fmt.Printf("%s13;6\007", K.KEYB_OSC)
	// fmt.Printf("%s14;6\007", K.KEYB_OSC)
}

func (this *vTerm) PrintError(s string) {
	echo.AddStringWithPos(this.nRows-2, 150, style.Error.EchoStyle(s, uint16(len(s))))
	echo.Flush()
	// trace.Error(s) // t //
}

// .
func (this *vTerm) SetNormal() {
	echo.AddSimple(K.DC_RESET + style.Pen.Normal)
}

func (this *vTerm) ReprintAll() {
	// trace.Begin_n(this, "ReprintAll") // t //

	/// Save values
	iActive := Play.keys.ActiveI

	/// Prepare
	for _, in := range AllMaps.allModels_in {
		in.PrepareView()
	}
	// for _, in := range AllMaps.allScreens_in {
	// 	in.PrepareScreen()
	// }
	for _, se := range AllMaps.allSections {
		if se.tobePrepared {
			se.PrepareFramesAndButtons()
		}
	}
	/// Print
	// MainScreen.LeftSection.GetActiveWindow().BringWindowToFront(true)
	BrowserWsp1.View.SetFocusToModel()
	View1.LoadAndPrepare(nil)
	View2.LoadAndPrepare(nil)
	Play.LoadAndPrepare(nil)

	/// Restore values
	Play.keys.ActiveI = iActive

	MusicDisplay.Print()
	speedFormat()
	volumeFormat()
	// trace.End() // t //
}

func DrawRect() {
	var frame string
	frame = Term.MakeFrame_Rect(MainScreen.DisplaySection.tRect, &style.ProgressBar)
	// trace.Print("DrawRect:  %+v", MainScreen.DisplaySection.tRect) // t //
	// trace.Print("DrawRect:  %#v", *frame)
	echo.AddStringWithPos(MainScreen.DisplaySection.marginTop, MainScreen.DisplaySection.marginLeft, frame)
	echo.Flush()
}

func (this *vTerm) MakeFrame_Rect(rect tRect, style *style.GStyle) string {
	return this.MakeFrame_HW(rect.marginBottom-rect.marginTop+1, rect.marginRight-rect.marginLeft+1, style)
}

func (this *vTerm) MakeFrame_HW(height, width uint16, style *style.GStyle) string {
	// trace.Begin_n(this, "MakeRectFrame") // t //
	var sUp, sDown, sLeft, sRight, sFinal string

	sFinal = style.JustTheColors()
	for range height {
		sLeft += "█"
		sLeft += Term.MoveLeft(1)
		sLeft += Term.MoveDown(1)
	}
	for range height {
		sRight += "█"
		sRight += Term.MoveLeft(1)
		sRight += Term.MoveDown(1)
	}
	sUp = strings.Repeat(K.CHAR_UPPER_HALF_BLOCK, int(width-2))
	sDown = strings.Repeat(K.CHAR_LOWER_HALF_BLOCK, int(width-2))

	sFinal += sLeft
	sFinal += Term.MoveUp(height)
	sFinal += Term.MoveRight(width - 1)

	sFinal += sRight
	sFinal += Term.MoveUp(height)
	sFinal += Term.MoveLeft(width - 2)

	sFinal += sUp
	sFinal += Term.MoveLeft(width - 2)
	sFinal += Term.MoveDown(height - 1)

	sFinal += sDown
	// trace.End() // t //
	return sFinal
}

func (this *vTerm) MoveUp(nRows uint16) string {
	return fmt.Sprintf("\x1b[%dA", nRows)
}

func (this *vTerm) MoveDown(nRows uint16) string {
	return fmt.Sprintf("\x1b[%dB", nRows)
}

func (this *vTerm) MoveLeft(nCols uint16) string {
	return fmt.Sprintf("\x1b[%dD", nCols)
}

func (this *vTerm) MoveRight(nCols uint16) string {
	return fmt.Sprintf("\x1b[%dC", nCols)
}

/* Adds a frame to the window. Use vWindow.Print to print the frame. */
func (this *vTerm) MakeSectionFrame(re *tRect, style *style.GStyle, frameElements *[6]string) *vTerm {
	// fmt.Println("AddFrame call debug...")
	var (
		frame_upleft  = frameElements[0]
		frame_upright = frameElements[1]
		// frame_downleft   = frameElements[2]
		// frame_downright  = frameElements[3]
		frame_vertical   = frameElements[4]
		frame_horizonzal = frameElements[5]
	)
	firstRow := re.marginTop + 1
	firstCol := re.marginLeft
	lastRow := re.marginBottom
	lastCol := re.marginRight
	///	Upper corner
	echo.AddSimple(style.JustTheColors())
	echo.AddStringWithPos(firstRow, firstCol, frame_upleft)
	echo.AddStringWithPos(firstRow, lastCol, frame_upright)
	///	Vertical lines
	// echo.AddStringWithPos(lastRow - 1, firstCol, "▄")
	// echo.AddStringWithPos(lastRow - 1, lastCol, "▄")
	for iRow := firstRow + 1; iRow < lastRow; iRow++ {
		echo.AddStringWithPos(iRow, firstCol, frame_vertical)
		echo.AddStringWithPos(iRow, lastCol, frame_vertical)
	}
	///	Lower corner
	echo.AddStringWithPos(lastRow, firstCol, " ")
	echo.AddStringWithPos(lastRow, lastCol, " ")
	/// Horizontal lines
	echo.AddStringWithPos(firstRow-1, firstCol, Frames.Empty[0])
	for col_i := firstCol + 2; col_i < lastCol+2; col_i++ {
		echo.AddSimple(Frames.Empty[0])
	}
	echo.AddStringWithPos(firstRow, firstCol+1, Frames.Empty[0])
	for col_i := firstCol + 2; col_i < lastCol; col_i++ {
		echo.AddSimple(Frames.Empty[0])
	}
	echo.AddStringWithPos(lastRow, firstCol+1, frame_horizonzal)
	for col_i := firstCol + 2; col_i < lastCol; col_i++ {
		echo.AddSimple("▀")
	}
	return this
}

func (this *vTerm) MakeLayerFrame(l *vLayer, style *style.GStyle, frameElements *[6]string) *vTerm {
	se := l.pSection
	firstRow := se.marginTop + 4
	lastRow := se.marginBottom
	var frame_vertical = frameElements[4]
	// trace.AddToPrint("%s %d %d", l.name, firstRow, lastRow) // t //

	echo.AddSimple(style.JustTheColors())
	for i := 0; i < len(l.windows)-1; i++ {
		w := l.windows[i]
		var col = w.marginRight
		echo.AddStringWithPos(firstRow, col, "▀")
		for iRow := firstRow + 1; iRow < lastRow; iRow++ {
			echo.AddStringWithPos(iRow, col, frame_vertical)
		}
	}
	// trace.BeginEnd_n(se, "PrepareLayerFrame") // t //
	return this
}

func (this *vTerm) Clear() {
	fmt.Print(style.Pen.Normal + TPUT_CLEAR)
}

func (this *vTerm) SetTerminalTitle(s string) {
	fmt.Printf("%s0;%s\007", K.KEYB_OSC, s)
}

func (this *vTerm) SetCursorPos(iRow, iCol uint16) {
	fmt.Printf("\x1b[%d;%dH", iRow+1, iCol+1)
}

func (this *vTerm) SetCursorNormal() {
	fmt.Print(TPUT_CNORM)
}

func (this *vTerm) SetCursorInvisible() {
	fmt.Print(TPUT_CIVIS)
}

func (this *vTerm) SaveCursor() {
	fmt.Print(TPUT_SC)
}

func (this *vTerm) RestoreCursor() {
	fmt.Print(TPUT_RC)
}

// func (this *vTerm) EraseLineToBeginning() {
// 	fmt.Print(TPUT_EL1)
// }

/*
stty -F /dev/tty -echo
*/
func (this *vTerm) SttyEchoOff() {
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

/*
stty -F /dev/tty echo
*/
func (this *vTerm) SttyEchoOn() {
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

/*
stty -F /dev/tty cbreak min 1
*/
func (this *vTerm) SttyCBreakMin1() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
}

func (this *vTerm) AltScreenOn() {
	fmt.Print(TPUT__ALTSCREEN_ON)
	this.isAltScreenOn = true
}

func (this *vTerm) AltScreenOff() {
	fmt.Print(TPUT__ALTSCREEN_OFF)
	this.isAltScreenOn = false
}

func (this *vTerm) Dump() {
	// trace.BeginAdd_n(this, "Dump", x__) // t //
	// trace.Print("this.name             %s", this.name) // t //
	// trace.Print("this.nRows            %d", this.nRows) // t //
	// trace.Print("this.nCols            %d", this.nCols) // t //
	// trace.Print("this.activeScreen     %s", GetNilOrName(this.activeScreen)) // t //
	// trace.Print("this.oldActiveScreen  %s", GetNilOrName(this.oldActiveScreen)) // t //
	// trace.Print("allScreens            %s", GetLongNames(AllMaps.allScreens)) // t //
	// trace.EndAdd(x__) // t //
}

func (this *cStatusLine) PrepareStatusLine(iRow, width, delay uint16) {
	this.iRow = iRow
	this.width = width
	this.delay = delay
}

func (this *cStatusLine) SetContent(s string) {
	this.content = s
	if this.counter == 0 {
		this._printContent()
	}
}

func (this *cStatusLine) _printContent() {
	echo.AddStringWithPos(this.iRow, 0, style.StatusLine.EchoStyle("  "+this.content, this.width))
	echo.Flush()
}

func (this *cStatusLine) PrintMessage(s string) {
	echo.AddStringWithPos(this.iRow, 0, style.StatusLineMessage.EchoStyle("  "+s, this.width))
	echo.Flush()
	this._deleteMessage()
}

func (this *cStatusLine) _deleteMessage() {
	go func() {
		this.counter++
		utils.Sleep(this.delay)
		this.counter--
		if this.counter == 0 {
			this._printContent()
		}
	}()
}

type vTerm struct {
	name string
	// content           string
	// beforeContent     string
	// afterContent      string
	nRows             uint16
	nCols             uint16
	activeScreen      *vScreen
	oldActiveScreen   *vScreen
	activeCommand     *uCommand
	isAltScreenOn     bool
	doSaveAtQuit      bool
	focusedModel      *mBasic
	StatusLine        cStatusLine
	ProgressBar_left  uint16
	ProgressBar_right uint16
	ProgressBar_top   uint16
	allButtons        []*vButton
	pPressedButton    *vButton
	quitMessage       string
}

type cStatusLine struct {
	content string
	counter uint16 // how many messages are running
	iRow    uint16
	width   uint16
	delay   uint16
}
