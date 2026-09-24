package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func (this *mSettings) ComposeLine(i uint16) string {
	var ifSet iSet
	ifSet = Set.allSettings[i]
	desc := ifSet.GetDesc()
	value := ifSet.GetValue()
	return fmt.Sprintf("  %-80s%s", desc, value)
}

func (this *mSettings) SelectKey_o() {
	this.w.statusLine = "  All values:  " + Set.allSettings[this.keys.iSelected].GetAllValues()
}

func (this *mSettings) PrepareView() {
// trace.Begin_n(this, "PrepareView") // t //
	w := &this.w
	// w.rowOffset = 4
	le := len(Set.allSettings)
	this.keys.Make(le)
	w.lines.Make(le)
	var ifSet iSet
	var st string
	// trace.N_Begin(this, "PrepareLines")
	for i := range le {
		ifSet = Set.allSettings[i]
		desc := ifSet.GetDesc()
		this.keys.sl[i].v = desc
		st = this.ComposeLine(uint16(i))
		w.SetLineForDialog(uint16(i), st, uint16(i), 0)
	}
	/// Empty line
	w.SetEmptyLineForDialog()
	w.columnsCaption = stylePageColCaption.EchoStyle("  Item                                                                            Value", this.w.nCols)
	this.keys.iActive = UNSET
	this.SelectKeyAndLine(0)
// trace.EndAdd("%d lines", w.lines.Len()) //@1 // t //
}

func (this *mSettings) ModifySelected(b bool) {
// trace.Begin_n(this, "ModifySelected") // t //
	i := this.keys.iSelected
	set := Set.allSettings[i]
	// trace.Print("%d, %s, %s, %s", i, set.GetDesc(), set.GetValue(), set.GetAllValues())
	set.ModifyValue(b)
	st := this.ComposeLine(i)
	this.w.SetLineForDialog(i, st, i, 0)
	// trace.Print("%d, %s, %s, %s", i, set.GetDesc(), set.GetValue(), set.GetAllValues())
// trace.End() // t //
	this.w.PrintWindowAnyway()
}

/*
Sets also the default values which may be overridden by the saved settings. Must be called before loading the saved settings.
*/
func (this *cSettings) InitSettings() {
// trace.BeginSilent("Set.InitSettings") // t //

	this.RepeatMode.values = []string{SET_REPEAT_OFF, SET_REPEAT_1, SET_REPEAT_ALL}
	this.RepeatMode.Pos = 0
	this.RepeatMode.desc = "Repeat playlist"
	this.RepeatMode._fWhenModified = func() {
		Display.AddLineToPrint(1, 4, "").Flush()
	}
	this.allSettings = append(this.allSettings, &this.RepeatMode)

	// this.Shuffle.Value = false
	// this.Shuffle.values = [2]string{SET_YES, SET_NO}
	// this.Shuffle.desc = "Shuffle playlist"
	// this.allSettings = append(this.allSettings, &this.Shuffle)

	// this.ShowRemainingTime.Value = false
	// this.ShowRemainingTime.values = [2]string{SET_YES, SET_NO}
	// this.ShowRemainingTime.desc = "Show remaining time"
	// this.allSettings = append(this.allSettings, &this.ShowRemainingTime)

	this.Equalizer.Pos = 0
	this.Equalizer.desc = "Equalizer"
	this.Equalizer._fWhenModified = func() {
		Mpv.AudiofilterRemoveAll()
		Mpv.AudiofilterAdd(Equalizer.GetEqusetString(UNSET, false))
		Equalizer.keys.iActive = this.Equalizer.Pos
	}
	this.allSettings = append(this.allSettings, &this.Equalizer)

	// this.RestoreSession.Value = false
	// this.RestoreSession.values = [2]string{SET_YES, SET_NO}
	// this.RestoreSession.desc = "Restore last session"
	// this.allSettings = append(this.allSettings, &this.RestoreSession)

	this.SleepTime.values = []uint16{200, 400, 600, 800, 1000}
	this.SleepTime.Pos = 1
	this.SleepTime.desc = "MusicBar delay in milliseconds"
	this.allSettings = append(this.allSettings, &this.SleepTime)

	this.Bro1Folding.values = []string{"Artist", "Composer"}
	this.Bro1Folding.Pos = 0
	this.Bro1Folding.desc = "Browser1: Folding, 1st Field"
	this.allSettings = append(this.allSettings, &this.Bro1Folding)

	this.Bro2Folding.values = []string{"Artist", "Composer"}
	this.Bro2Folding.Pos = 1
	this.Bro2Folding.desc = "Browser2: Folding, 1st Field"
	this.allSettings = append(this.allSettings, &this.Bro2Folding)

	this.OpenInBrowserApple.Value = false
	this.OpenInBrowserApple.values = [2]string{SET_YES, SET_NO}
	this.OpenInBrowserApple.desc = "Open in Browser at Apple.com"
	this.allSettings = append(this.allSettings, &this.OpenInBrowserApple)

	this.OpenInBrowserGoogle.Value = true
	this.OpenInBrowserGoogle.values = [2]string{SET_YES, SET_NO}
	this.OpenInBrowserGoogle.desc = "Open in Browser at Google.com"
	this.allSettings = append(this.allSettings, &this.OpenInBrowserGoogle)

	this.OpenInBrowserYoutube.Value = false
	this.OpenInBrowserYoutube.values = [2]string{SET_YES, SET_NO}
	this.OpenInBrowserYoutube.desc = "Open in Browser at Youtube.com"
	this.allSettings = append(this.allSettings, &this.OpenInBrowserYoutube)

	// this.ShowNumbers.Value = false
	// this.ShowNumbers.values = [2]string{SET_YES, SET_NO}
	// this.ShowNumbers.desc = "Show color numbers"
	// this.allSettings = append(this.allSettings, &this.ShowNumbers)

	this.PlayerSelectActive.Value = false
	this.PlayerSelectActive.values = [2]string{SET_YES, SET_NO}
	this.PlayerSelectActive.desc = "Player: select the active one"
	this.allSettings = append(this.allSettings, &this.PlayerSelectActive)

	this.BlinkPlaying.Value = false
	this.BlinkPlaying.values = [2]string{SET_YES, SET_NO}
	this.BlinkPlaying.desc = "Playing line is blinking"
	this.allSettings = append(this.allSettings, &this.BlinkPlaying)
	this.BlinkPlaying._fWhenModified = func() {
		if this.BlinkPlaying.Value {
			styleActive.fontStyle = FONT_BLINK
			styleSelActive.fontStyle = FONT_BLINK
		} else {
			styleActive.fontStyle = FONT_NORMAL
			styleSelActive.fontStyle = FONT_NORMAL
		}
	}

	this.PlaySecs.Pos = 0
	this.PlaySecs.values = []string{SET_NO, "10", "20", "30"}
	this.PlaySecs.desc = "Play only the first N seconds"
	this.allSettings = append(this.allSettings, &this.PlaySecs)

	this.AudioPitchCorrection.Value = false
	this.AudioPitchCorrection.values = [2]string{SET_YES, SET_NO}
	this.AudioPitchCorrection.desc = "Audio Pitch Correction"
	this.AudioPitchCorrection._fWhenModified = func() {
// trace.PrintGreen("AudioPitchCorrection.fWhenModified") // t //
		Mpv.SetAudioPitchCorrection(this.AudioPitchCorrection.Value)
	}
	this.allSettings = append(this.allSettings, &this.AudioPitchCorrection)

	// trace.Print("%v", this.allSettings[0])
// trace.End() // t //
}

// type t__setts struct {

// }

func (this *cSettings) marshall(filename string, x any) {
	fd, _ := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	defer fd.Close()
	by, _ := json.MarshalIndent(x, "", "  ")
	fmt.Fprint(fd, string(by))
}

func (this *cSettings) unmarshall(filename string, x any) {
	by, _ := os.ReadFile(filename)
	// trace.Print("::::::::::::::::::::::::::::: json: %s", string(by))
	json.Unmarshal([]byte(by), &x)
	// trace.Print("set unmarshal:\n%+v", x)
}

func (this *cSettings) SaveSettings() {
// trace.Begin("Set.SaveSettings") // t //
	// fd, _ := os.OpenFile("a/settings.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	// defer fd.Close()
	// fdIndent, _ := os.OpenFile("a/settings_indented.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	// defer fdIndent.Close()

	Set.LastSectionLeftWindow = MainScreen.LeftSection.GetActiveWindow().pModel.GetLongName()
	Set.LastSectionRightWindow = MainScreen.pRightSection.GetActiveWindow().pModel.GetLongName()

	// Set.LastSectionRightWindow = Screen.SectionView1.activeLayer.activeWindow.pModel.GetLongName()
	// Set.LastSectionLeftWindow = Screen.SectionLeft.activeLayer.activeWindow.pModel.GetLongName()

	if k := Play.keys.GetActive(); k != nil {
		Set.LastPlayedFile = k.v
		Set.LastTime = Mpv.GetFileTime()
	} else {
		Set.LastPlayedFile = ""
	}
	this.marshall(FILE_SETTINGS_UIPL, &ui_pl)
	this.marshall(FILE_SETTINGS_SETS, &Set)
	this.marshall(FILE_SETTINGS_STYLES, &AllMaps.allStyleSettings)
	this.marshall(FILE_SETTINGS_EQU, &Equalizer.EquSets)
// trace.EndAdd("Set.SaveSettings") // t //
}

func (this *cSettings) LoadSettings() {
// trace.BeginSilent("Set.LoadSettings") // t //
	this.unmarshall(FILE_SETTINGS_UIPL, &ui_pl)
	this.unmarshall(FILE_SETTINGS_SETS, &Set)
	this.unmarshall(FILE_SETTINGS_STYLES, &AllMaps.allStyleSettings)
	this.unmarshall(FILE_SETTINGS_EQU, &Equalizer.EquSets)
	// trace.Print("Set: %+v", Set)
	// trace.Print("Equalizer.EquSets: %+v", Equalizer.EquSets)
	if this.BlinkPlaying.Value {
		styleActive.fontStyle = FONT_BLINK
		styleSelActive.fontStyle = FONT_BLINK
	}
	Equalizer.keys.iActive = this.Equalizer.Pos
// trace.End() // t //
}

//////	t__SetSt

func (this *cSetSt) GetDesc() string {
	return this.desc
}

func (this *cSetSt) GetValue() string {
	return this.values[this.Pos]
}

func (this *cSetSt) GetAllValues() string {
	var st string
	for i, v := range this.values {
		st += v
		if i < len(this.values)-1 {
			st += " "
		}
	}
	return st
}

func (this *cSetSt) ModifyValue(b bool) {
	if b == true {
		this.Pos++
		if this.Pos >= uint16(len(this.values)) {
			this.Pos = 0
		}
	} else {
		if this.Pos > 0 {
			this.Pos--
		} else {
			this.Pos = uint16(len(this.values) - 1)
		}
	}
	if this._fWhenModified != nil {
		this._fWhenModified()
	}
}

//////

func (this *cSetStEqu) GetDesc() string {
	return this.desc
}

func (this *cSetStEqu) GetValue() string {
	return Equalizer.EquSets[this.Pos].Name
}

func (this *cSetStEqu) GetAllValues() string {
	return ""
}

func (this *cSetStEqu) ModifyValue(plus bool) {
	if plus == true {
		if this.Pos < uint16(len(Equalizer.EquSets)-1) {
			this.Pos++
		}
	} else {
		if this.Pos > 0 {
			this.Pos--
		}
	}
	if this._fWhenModified != nil {
		this._fWhenModified()
	}
}

//////	t__SetUi

func (this *cSetUi) GetDesc() string {
	return this.desc
}

func (this *cSetUi) GetValue() string {
	return fmt.Sprintf("%d", this.values[this.Pos])
}

func (this *cSetUi) GetUintValue() uint16 {
	return this.values[this.Pos]
}

func (this *cSetUi) GetAllValues() string {

	return fmt.Sprintf("%d...%d", this.values[0], this.values[len(this.values)-1])
}

func (this *cSetUi) ModifyValue(b bool) {
	if b == true {
		this.Pos++
		if this.Pos >= uint16(len(this.values)) {
			this.Pos = 0
		}
	} else {
		if this.Pos > 0 {
			this.Pos--
		} else {
			this.Pos = uint16(len(this.values) - 1)
		}
	}
}

//////	t__SetBo

func (this *cSetBo) GetDesc() string {
	return this.desc
}

func (this *cSetBo) GetValue() string {
	if this.Value {
		return this.values[0]
	}
	return this.values[1]
}

func (this *cSetBo) ModifyValue(b bool) {
	this.Value = !this.Value
	if this._fWhenModified != nil {
		this._fWhenModified()
	}
}

func (this *cSetBo) GetAllValues() string {
	return this.values[0] + " " + this.values[1]
}

type mSettings struct {
	mBasic
}

type cSettings struct {
	///
	consoleModeOn          bool
	LastSectionLeftWindow  string
	LastSectionRightWindow string
	LastPlayedFile         string
	LastTime               uint32
	allSettings            []iSet

	RepeatMode cSetSt
	// Shuffle              cSetBo
	// ShowRemainingTime    cSetBo
	Equalizer cSetStEqu
	// RestoreSession       cSetBo
	SleepTime            cSetUi
	AudioPitchCorrection cSetBo
	// MusicBarRune         cSetSt
	Bro1Folding          cSetSt
	Bro2Folding          cSetSt
	OpenInBrowserGoogle  cSetBo
	OpenInBrowserApple   cSetBo
	OpenInBrowserYoutube cSetBo
	// ShowNumbers          cSetBo
	PlaySecs           cSetSt
	BlinkPlaying       cSetBo
	PlayerSelectActive cSetBo
}

type iSet interface {
	GetDesc() string
	GetValue() string
	GetAllValues() string
	ModifyValue(bool)
}

type cSetSt struct {
	// value  string
	values         []string
	desc           string
	Pos            uint16
	_fWhenModified func()
}

type cSetStEqu struct {
	// value  string
	desc           string
	Pos            uint16
	_fWhenModified func()
}

type cSetUi struct {
	// value  uint16
	values []uint16
	desc   string
	Pos    uint16
}

type cSetBo struct {
	Value          bool
	values         [2]string
	desc           string
	_fWhenModified func()
}

