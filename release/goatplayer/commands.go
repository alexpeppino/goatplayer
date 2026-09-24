package main

import (
	"fmt"
	"os/exec"
	"strings"
)

/// Basic

func (this *mBasic) SetCommands() {}

func (this *mBasic) AppendCommands_OnlyMainScreen() {
	this.kCmm = append(this.kCmm,
		uCommand{"", "", nil, nil}, // ----------------
		uCommand{FOCUS_HELP, "Open Help", this.OpenHelpDialog, nil},
		uCommand{FOCUS_SETTINGS, "Open Settings", this.OpenSettingsDialog, nil},
		uCommand{FOCUS_EQUALIZER, "Open Equalizer", this.OpenEqualizerDialog, nil},
		uCommand{FOCUS_COLORS, "Open Colors", this.OpenColorsDialog, nil},
		uCommand{FOCUS_REPORT, "Open Report", this.OpenReportDialog, nil},
		uCommand{KEYB_CTRL_F1, "Reprint all", Term.ReprintAll, nil},
		uCommand{"b", "Open in Browser", this.OpenInBrowser, nil})
}

func (this *mBasic) AppendCommands_Basic() {
	this.AppendCommands_PlayerControls()
	this.kCmm = append(this.kCmm,
		uCommand{",", "Prev commands page", func() { Term.focusedModel.PrintPrevMenu() }, nil},
		uCommand{".", "Next commands page", func() { Term.focusedModel.PrintNextMenu() }, nil},
		uCommand{KEYB_ALT_Q, "Quit", Quit, nil})
}

//		}

func (this *mBasic) AppendCommands_PlayerControls() {
	this.kCmm = append(this.kCmm,
		uCommand{"", "Player Controls", nil, nil}, // ----------------
		uCommand{KEYPAD_ALT_ARROW_LEFT, "Seek -5s", Mpv.SeekM5s, nil},
		uCommand{KEYPAD_ALT_ARROW_RIGHT, "Seek +5s", Mpv.SeekP5s, nil},
		uCommand{KEYPAD_ALT_ARROW_DOWN, "Seek -30s", Mpv.SeekM30s, nil},
		uCommand{KEYPAD_ALT_ARROW_UP, "Seek +30s", Mpv.SeekP30s, nil},
		uCommand{KEYPAD_ALT_PAGE_DOWN, "Speed -0.01", func() { speedMinus(1); Mpv.SetSpeed(ui_pl.Speed) }, nil},
		uCommand{KEYPAD_ALT_PAGE_UP, "Speed +0.01", func() { speedPlus(1); Mpv.SetSpeed(ui_pl.Speed) }, nil},
		uCommand{"/", "Speed 1.00", func() { speed1(); Mpv.SetSpeed(ui_pl.Speed) }, nil},
		uCommand{KEYPAD_ALT_END, "Volume -2", func() { volumeMinus(2); Mpv.SetVolume(ui_pl.Volume) }, nil},
		uCommand{KEYPAD_ALT_HOME, "Volume +2", func() { volumePlus(2); Mpv.SetVolume(ui_pl.Volume) }, nil},
		uCommand{"*", "Volume 100", func() { volume100(); Mpv.SetVolume(ui_pl.Volume) }, nil})
	if this == &Play.mBasic {
		this.kCmm = append(this.kCmm,
			uCommand{KEYB_ALT_P, "Play selected", Play.PlaySelected, nil})
	}
	this.kCmm = append(this.kCmm,
		uCommand{KEYB_ALT_N, "Play next", Mpv.Stop, nil},
		uCommand{KEYB_ALT_S, "Stop", func() { ui_pl.stopPlaying = true; Mpv.Stop() }, nil},
		uCommand{KEYB_ALT_O, "Pause/Unpause", Mpv.CyclePause, nil},
		uCommand{KEYB_ALT_M, "Mute/Unmute", Mpv.CycleMute, nil},
		uCommand{"", "", nil, nil}) // ----------------
}

// mBasic > mBasicFileList > mPlaylist
func (this *mPlaylist) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mAbsFilelist.SetCommands()
	this.kCmm = append(this.kCmm,
		uCommand{"A", "Append from Find", this.AppendFromList, nil},
		uCommand{"I", "Insert from Find", this.InsertFromList, nil})
	// uCommand{"c", "Clear list", func() {
	// 	this.ClearPlayList()
	// 	this.w.PrintWindowAnyway()
	// }, nil},
	// uCommand{"", "Parameters", nil, nil}, // ----------------
	// uCommand{"w", "Set start", this.SetStart, nil},
	// uCommand{"e", "Set end", this.SetEnd, nil},
	// uCommand{"r", "Set speed", this.SetSpeed, nil},
	// uCommand{"t", "Unset start", this.UnsetStart, nil},
	// uCommand{"z", "Unset end", this.UnsetEnd, nil},
	// uCommand{"u", "Unset speed", this.UnsetSpeed, nil},
	// uCommand{"", "", nil, nil},
	// uCommand{"l", "Save list", this.MakeAndSaveList, nil},
	// uCommand{"E", "Edit name", this.RenameListAskUser, nil},
	// uCommand{"n", "New list", this.NewList, nil})
	///
	this.mCmm = append(this.mCmm,
		uCommand{MOUSE_B3_DOWN, "Play selected", func() {
			this.w.SelectLineByTRow(Input.iRow)
			this.PlaySelected()
		}, nil})
}

func (this *mBasic) AppendCommands_BasicDialog() {
	this.AppendCommands_Basic()
	this.AppendCommands_BasicSelection()
	if this != &Find.mBasic {
		this.kCmm = append(this.kCmm,
			uCommand{KEYB_ESC, "Exit dialog", func() {}, nil})
	}
}

func (this *mBasic) AppendCommands_BasicSelection() {
	this.kCmm = append(this.kCmm,
		uCommand{"", "Selection", nil, nil}, // ----------------
		uCommand{KEYPAD_ARROW_UP, "Select previous", this.SelectPrev, nil},
		uCommand{KEYPAD_ARROW_DOWN, "Select next", this.SelectNext, nil},
		uCommand{KEYPAD_HOME, "Select first", this.SelectFirst, nil},
		uCommand{KEYPAD_END, "Select next active", this.outer.SelectActive_o, nil},
		uCommand{KEYPAD_PAGE_UP, "Previous page", this.GotoPrevPage, nil},
		uCommand{KEYPAD_PAGE_DOWN, "Next page", this.GotoNextPage, nil},
		uCommand{KEYB_TAB, "Next Wsp Tab", func() {}, nil},
		uCommand{KEYB_SHIFT_TAB, "Prev Wsp Tab", func() {}, nil})
	/// Mouse
	this.mCmm = append(this.mCmm,
		uCommand{"", "Mouse", nil, nil},
		uCommand{MOUSE_B1_DOWN, "Select line", func() { this.w.SelectLineByTRow(Input.iRow) }, nil},
		uCommand{MOUSE_WHEELUP, "Prev. page", func() { this.GotoPrevPage() }, nil},
		uCommand{MOUSE_WHEELDOWN, "Next page", func() { this.GotoNextPage() }, nil})
}

///

// mBasic > mBasicFileList
func (this *mAbsFilelist) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mBasic.AppendCommands_Basic()
	this.mBasic.AppendCommands_BasicSelection()
	this.mBasic.AppendCommands_OnlyMainScreen()
	this.kCmm = append(this.kCmm,
		uCommand{"N", "Open with Nautilus", func() {
			go func() {
				exec.Command("nautilus", this.keys.GetSelected().v).Run()
			}()
		}, nil},
		uCommand{"P", "Show prop", func() {
			fileProp := &this.pWorkspace.AudioPropCon
			fileProp.LoadFile(this.keys.GetSelected().v)
			fileProp.PrepareView()
			fileProp.w.BringWindowToFront(true)
		}, nil})
}

// mBasic > mBasicFileList > mFilelist
func (this *mFilelist) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mAbsFilelist.SetCommands()
	this.kCmm = append(this.kCmm,
		uCommand{"p", "Play this list", this.PlayThisList, nil},
		uCommand{"j", "Play this item", this.PlayThisItem, nil},
		uCommand{"a", "Add selected to p", this.AddSelectedToPlaylist, nil})
	// uKeyCommand{"e", "Add selected to e", this.AddSelectedToEditor, nil})
	///
	this.mCmm = append(this.mCmm,
		uCommand{MOUSE_B3_DOWN, "Play this", func() {
			this.w.SelectLineByTRow(Input.iRow)
			this.PlayThisList()
		}, nil})
}

/// Container

func (this *mContainer) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mBasic.AppendCommands_Basic()
	this.mBasic.AppendCommands_BasicSelection()
	this.mBasic.AppendCommands_OnlyMainScreen()
	this.kCmm = append(this.kCmm,
		uCommand{"a", "Activate sel", func() { this.ActivateFilter(this.keys.iSelected) }, nil},
		uCommand{"f", "Find and activate", func() {
			Find.title = "Find and select one item"
			if this.OpenFindDialog() {
				this.outer.UseFoundItems_o()
			}
		}, nil})
	// uKeyCommand{"p", "Open in ActivePL", func() {
	// 	if this.outer == nil {
	// 		trace.Error("%s, outer not set", this.name)
	// 	} else {
	// 		this.outer.LoadTo(View2)
	// 	}
	// }, nil},
	// uKeyCommand{"o", "Open in OpenPL", func() { this.outer.LoadTo(View1) }, nil})

	/// Mouse
	//! if this != &PlayerWsp.PFListsCon.mContainer {
	this.mCmm = append(this.mCmm,
		uCommand{MOUSE_B1_DOWN, "Select line", func() {
			this.w.SelectLineByTRow(Input.iRow)
		}, nil},
		uCommand{MOUSE_B3_DOWN, "Activate line", func() {
			iCol := Input.iCol
			if iCol <= this.w.marginLeft || iCol >= this.w.marginRight {
			} else {
				this.ActivateFilterByTRow(Input.iRow)
			}
		}, nil})
	// }
}

func (this *mAudioProperties) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mContainer.SetCommands()
	if this != &PlayerWsp.AudioPropCon {
		this.kCmm = append(this.kCmm,
			uCommand{"l", "Edit from list",
				nil,
				// this.EditSelectedFromFind,
				nil},
			uCommand{"e", "Edit selected",
				nil,
				// this.EditSelected,
				nil})
	}
}

/* mBasic > mContainer > mPFListsCon */
func (this *mPFListsCon) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mContainer.SetCommands()
	this.kCmm = append(this.kCmm,
		uCommand{"d", "delete sel. list", this.DeleteSelectedList, nil},
		uCommand{"O", "load to player", this.LoadToPlayer, nil})
	this.mCmm = append(this.mCmm,
		uCommand{MOUSE_B3_DOWN, "load to player", func() {
			iLine := this.w.CalcLineIndexByTRow(Input.iRow)
			this.SelectKeyAndLine(iLine)
			this.w.PrintWindowAnyway()
			this.LoadToPlayer()
		}, nil})
}

///

func (this *mDirs) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mBasic.AppendCommands_Basic()
	this.mBasic.AppendCommands_BasicSelection()
	this.mBasic.AppendCommands_OnlyMainScreen()
	this.kCmm = append(this.kCmm,
		uCommand{"w", "Set as Wording dir.", this.SetSelectedAsWorkingDir, nil},
		uCommand{"N", "Open with Nautilus", func() {
			go func() { exec.Command("nautilus", this.keys.GetSelected().v).Run() }()
		}, nil})
}

/// Dialogs

func (this *mFind) SetCommands() {
	this.mBasic.AppendCommands_BasicDialog()
	this.kCmm = append(this.kCmm,
		uCommand{KEYB_ENTER, "OK", func() {}, nil},
		uCommand{KEYB_ESC, "Cancel", func() {}, nil})
}

func (this *mHelp) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mBasic.AppendCommands_BasicDialog()
	this.kCmm = append(this.kCmm,
		uCommand{KEYPAD_ARROW_LEFT, "Prev Filepage", func() {
			if this.w2.iPage > 0 {
				this.w2.iPage--
				this.w2.PrintPage()
			}
		}, nil},
		uCommand{KEYPAD_ARROW_RIGHT, "Next Filepage", func() {
			if this.w2.iPage < this.w2.nPages-1 {
				this.w2.iPage++
				this.w2.PrintPage()
			}
		}, nil},
		uCommand{"e", "Open in editor", func() {
			go func() {
				path := this.folderName + "/" + this.keys.GetSelected().v
// trace.Print("open %s", path) // t //
// by, err := exec.Command("xdg-open", path).Output() /*c*/
exec.Command("xdg-open", path).Output()
// trace.Print("%v %v", by, err) // t //
			}()
		}, nil})
}

func (this *mSettings) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mBasic.AppendCommands_BasicDialog()
	this.kCmm = append(this.kCmm,
		uCommand{KEYPAD_ARROW_LEFT, "Prev value", func() { this.ModifySelected(false) }, nil},
		uCommand{KEYPAD_ARROW_RIGHT, "Next value", func() { this.ModifySelected(true) }, nil})
}

func (this *mColors) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	/// Keyboard
	this.mBasic.AppendCommands_BasicDialog()
	this.kCmm = append(this.kCmm,
		uCommand{"r", "Reset colors", this.ResetColors, nil})
	/// Mouse
}

func (this *mColorsPicker) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	/// Keyboard
	/// Mouse
	this.mCmm = append(this.mCmm,
		uCommand{MOUSE_B1_DOWN, "Pick text color", func() {
			var iColor uint16
			iRow := Input.iRow
			iCol := Input.iCol
// trace.Print("mColors.CheckInput %d %d %d", iRow, iCol, Colors.iFgFirstRow) // t //
			if iRow >= Colors.iFgFirstRow && iRow < Colors.iFgFirstRow+16 && iCol >= Colors.iFgFirstCol && iCol < Colors.iFgFirstCol+16*5 {
				r := iRow - Colors.iFgFirstRow
				c := iCol - Colors.iFgFirstCol
				iColor = c/5 + r*16
// trace.Print("fg iColor: %d", iColor) // t //
				Colors.ModifyFgColor(iColor)
			}
		}, nil},
		uCommand{MOUSE_B3_DOWN, "Pick backgr. color", func() {
			var iColor uint16
			iRow := Input.iRow
			iCol := Input.iCol
// trace.Print("mColors.CheckInput %d %d %d", iRow, iCol, Colors.iFgFirstRow) // t //
			if iRow >= Colors.iFgFirstRow && iRow < Colors.iFgFirstRow+16 && iCol >= Colors.iFgFirstCol && iCol < Colors.iFgFirstCol+16*5 {
				r := iRow - Colors.iFgFirstRow
				c := iCol - Colors.iFgFirstCol
				iColor = c/5 + r*16
// trace.Print("bg iColor: %d", iColor) // t //
				Colors.ModifyBgColor(iColor)
			}
		}, nil})
}

func (this *mEqualizer) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mBasic.AppendCommands_BasicDialog()
	this.kCmm = append(this.kCmm,
		uCommand{"i", "Previous freq.", this.FreqM1, nil},
		uCommand{"k", "Next freq.", this.FreqP1, nil},
		uCommand{"o", "Gain -1", this.GainM1, nil},
		uCommand{"p", "Gain +1", this.GainP1, nil},
		uCommand{"a", "Activate sel.", func() {
			Set.Equalizer.Pos = this.keys.iSelected
			Set.Equalizer._fWhenModified()
			this.w.PrintWindowAnyway()
			SetM.PrepareView()
		}, nil},
		uCommand{"r", "Reset values", func() {
			this.ResetValues()
			this.PrepareView()
			this.w.PrintWindowAnyway()
			this.m2.w.PrintWindowAnyway()
		}, nil})
}

func (this *mReport) SetCommands() {
	// trace.N_BeginEnd(this, "SetCommands")
	this.mBasic.AppendCommands_BasicDialog()
}

///

func PrintUnpressedCommand(cmm *uCommand) {
	var key string
	if v, ok := Input.mapKeypad[cmm.key]; ok {
		key = v
	} else {
		key = cmm.key
	}
	s := styleCommandPanel.EchoTwoStyles(cmm.desc, key, 24, &styleSidepanelKey)
	Term.AddStringWithPos(cmm.rect.marginTop, cmm.rect.marginLeft, s).Flush()
}

func PrintPressedCommand(cmm *uCommand) {
	var key string
	if v, ok := Input.mapKeypad[cmm.key]; ok {
		key = v
	} else {
		key = cmm.key
	}
	s := styleNormal.EchoTwoStyles(cmm.desc, key, 24, &styleNormal)
	Term.AddStringWithPos(cmm.rect.marginTop, cmm.rect.marginLeft, s).Flush()
}

func (this *mBasic) CheckDoubleCommandKeys() {
	mapKeys := make(map[string]uint16)
	for i := range this.kCmm {
		cmm := &this.kCmm[i]
		mapKeys[cmm.key]++
	}
	for k, n := range mapKeys {
		if k != "" {
			if n > 1 {
// trace.Error("%s: key has doubles, '%s' %d", this.name, Input.GetKeyName(k), n) // t //
				for i := range this.kCmm {
					cmm := &this.kCmm[i]
					if cmm.key == k {
// trace.Error("  %s", cmm.desc) // t //
					}
				}
				Quit()
			}
		}
	}
}

/* Prepares this table's commands echo string. */
func (this *mBasic) PrepareCommandsEcho() {
	// ◄ ► ▲ ▼
	// trace.Begin_n(this, "PrepareCommandsEcho")
	var key string
	var iRow uint16 = MainScreen.CommandPanelSection.marginTop + 2
	var width uint16 = 24
	var left = MainScreen.CommandPanelSection.marginLeft
	var iLine uint16
	var iPage uint16
	var totCmm = uint16(len(this.kCmm) + len(this.mCmm))
	var nPages = (totCmm-1)/CommandPanel.nLines + 1
	var s string
	this.nCmmPages = nPages
	this.cmmEchoes = make([]string, nPages)
	// trace.Print("%d/%d -> %d", totCmm, CommandPanel.nLines, nPages)

	checkEndOfPage := func(s string) {
		if iLine != 0 && iLine%CommandPanel.nLines == 0 {
			/// End of page with kbo commands
			// trace.Print("%s slurp page %d", s, iPage)
			s = fmt.Sprintf("%s %d/%d", this.shortName, iPage+1, nPages)
			s = styleCommandPanelTitle.EchoStyle(s, width-7)
			Term.AddStringWithPos(MainScreen.CommandPanelSection.marginTop, left, s)
			Term.SlurpContents(&this.cmmEchoes[iPage])
			iPage++
			iRow -= CommandPanel.nLines
		}
	}

	makeTitle := func(cmm *uCommand) {
		var str1 string
		if cmm.desc == "" {
			str1 = "-------------------------"
		} else {
			str1 = "-- " + cmm.desc + " " + strings.Repeat("-", 20-len(cmm.desc))
		}
		s = styleCommandPanelItalic.EchoStyle(str1, width)
		cmm.rect = &tRect{0, 0, 0, 0}
	}

	/// Keyboard commands
	for i := range this.kCmm {
		checkEndOfPage("kbo")
		cmm := &this.kCmm[i]
		if cmm.key == "" {
			makeTitle(cmm)
		} else {
			if v, ok := Input.mapKeypad[cmm.key]; ok {
				key = v
			} else {
				key = cmm.key
			}
			s = styleCommandPanel.EchoTwoStyles(cmm.desc, key, width, &styleSidepanelKey)
			cmm.rect = &tRect{iRow, iRow, left, left + 24}
		}
		// trace.Print("add string at iRow: %d,  iLine%%CommandPanel.nLines: %d", iRow, iLine%CommandPanel.nLines)
		Term.AddStringWithPos(iRow, left, s)
		iRow++
		// trace.Print("kbo commands --  line: %3d/%3d  page: %3d  (iRow: %d)", iLine, CommandPanel.nLines, iPage, iRow)
		iLine++
	}
	/// Mouse commands
	for i := range this.mCmm {
		checkEndOfPage("mou")
		mcmm := &this.mCmm[i]
		if mcmm.key == "" {
			makeTitle(mcmm)
		} else {
			if v, ok := Input.mapKeypad[mcmm.key]; ok {
				key = v
			} else {
				key = mcmm.key
			}
			s = styleCommandPanel.EchoTwoStyles(mcmm.desc, key, width, &styleSidepanelKey)
			mcmm.rect = &tRect{iRow, iRow, left, left + 24}
		}
		Term.AddStringWithPos(iRow, left, s)
		iRow++
		// trace.Print("mou commands --  line: %3d/%3d  page: %3d", iLine, CommandPanel.nLines, iPage)
		iLine++
	}

	/// Fill with empty	lines
	l1 := CommandPanel.nLines - (iLine % CommandPanel.nLines)
	if l1 > 0 {
		for range l1 {
			s := ""
			// trace.Print("empty line --  iRow: %3d", iRow)
			Term.AddStringWithPos(iRow, left, styleCommandPanel.EchoStyle(s, width))
			iRow++
		}
		s := fmt.Sprintf("%s %d/%d", this.shortName, iPage+1, nPages)
		s = styleCommandPanelTitle.EchoStyle(s, width-7)
		Term.AddStringWithPos(MainScreen.CommandPanelSection.marginTop, left, s)
		Term.SlurpContents(&this.cmmEchoes[nPages-1])
	}
	// trace.End()
}

func (this *mBasic) PrintPrevMenu() {
	if this.iCmmPage > 0 {
		this.iCmmPage--
		this.PrintCommands()
	}
}

func (this *mBasic) PrintNextMenu() {
	if this.iCmmPage < this.nCmmPages-1 {
		this.iCmmPage++
		this.PrintCommands()
	}
}

func (this *cCommandPanel) PrepareCommandPanel() {
	this.butPageMinus = vButton{
		" ◀ ",
		0,
		MainScreen.CommandPanelSection.marginRight - 8,
		3,
		func() { Term.focusedModel.PrintPrevMenu() }}

	this.butPagePlus = vButton{
		" ▶ ",
		0,
		MainScreen.CommandPanelSection.marginRight - 4,
		3,
		func() { Term.focusedModel.PrintNextMenu() }}

	this.nLines = MainScreen.CommandPanelSection.marginBottom - MainScreen.CommandPanelSection.marginTop - 1
}

type cCommandPanel struct {
	butPageMinus vButton
	butPagePlus  vButton
	iPage        uint16
	nLines       uint16
}

var CommandPanel cCommandPanel

// type uMouseCommand struct {
// 	key  string
// 	desc string
// 	f    func()
// 	rect *tRect
// }

type uCommand struct {
	key  string
	desc string
	f    func()
	rect *tRect
}

