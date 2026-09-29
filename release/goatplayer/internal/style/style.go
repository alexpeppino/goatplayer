package style

import (
	"fmt"
	"goatplayer/internal/K"
	"goatplayer/internal/small/utils"
)

var Pen t__Pencils
var MapStyle map[string]*GStyle // map: style name > style reference
var AllConfigurableStyles []*GStyle
var AllStyleSettings []GStyleSettings
var MapStyleSettings map[string]*GStyleSettings

func (this *GStyle) EchoStyle(s string, width uint16) string {
	var textColor, backColor uint16
	var em string
	if this.TextColor == K.UNSET {
		textColor = this.OtherStyleTextColor.TextColor
	} else {
		textColor = this.TextColor
	}
	if this.BackColor == K.UNSET {
		backColor = this.OtherStyleBackColor.BackColor
	} else {
		backColor = this.BackColor
	}
	color := fmt.Sprintf("\x1b[38;5;%d;48;5;%dm", textColor, backColor)
	// s = strings.TrimSpace(s)
	if width > 0 {
		utils.AdjustWidthSeparated(&s, &em, width)
		em = fmt.Sprintf("%s%s%s", color, em, K.FONT_RESET)
	}
	s = fmt.Sprintf("%s%s%s%s", color, this.FontStyle, s, K.FONT_RESET)
	if this.Align == K.ALIGN_LEFT {
		return s + em
	} else {
		return em + s
	}
}

/* Don't use with fontstyle underlined or inherited colors (use Echo instead). */
func (this *GStyle) EchoStyleSimple(s string, width uint16) string {
	utils.AdjustWidthAligned(&s, width, this.Align)
	return fmt.Sprintf("\x1b[38;5;%d;48;5;%dm%s%s%s", this.TextColor, this.BackColor, this.FontStyle, s, K.FONT_RESET)
}

/* Adds a blank before and after. Doesn't change the width. */
func (this *GStyle) EchoStyleFormatOnly(s string, addBlank bool) string {
	var textColor, backColor uint16
	if this.TextColor == K.UNSET {
		textColor = this.OtherStyleTextColor.TextColor
	} else {
		textColor = this.TextColor
	}
	if this.BackColor == K.UNSET {
		backColor = this.OtherStyleBackColor.BackColor
	} else {
		backColor = this.BackColor
	}
	if addBlank {
		s = " " + s + " "
	}
	return fmt.Sprintf("\x1b[38;5;%d;48;5;%dm%s%s%s", textColor, backColor, this.FontStyle, s, K.FONT_RESET)
}

func (this *GStyle) EchoTwoStyles(s, right string, width uint16, otherStyle *GStyle) string {
	var leRight = utils.RuneLen(&right)
	s1 := this.EchoStyle(s, width-leRight)
	s2 := otherStyle.EchoStyle(right, leRight)
	return s1 + s2
}

/* Returns only the echo of the colors. */
func (this *GStyle) JustTheColors() string {
	var textColor, backColor uint16
	if this.TextColor == K.UNSET {
		textColor = this.OtherStyleTextColor.TextColor
	} else {
		textColor = this.TextColor
	}
	if this.BackColor == K.UNSET {
		backColor = this.OtherStyleBackColor.BackColor
	} else {
		backColor = this.BackColor
	}
	return fmt.Sprintf("\x1b[38;5;%d;48;5;%dm", textColor, backColor)
}

func styleBold(s string) string {
	return K.DC_BOLD + s + K.DC_RESET
}

func styleItalic(s string) string {
	return K.DC_ITALIC + s + K.DC_RESET
}

func styleUnderline(s string) string {
	return K.DC_UNDERLINE + s + K.DC_RESET
}

func (this *GStyle) ImportSettings() {
	if styleSet, ok := MapStyleSettings[this.Name]; ok {
		/// Read the style from settings
		this.TextColor = styleSet.TextColor
		this.BackColor = styleSet.BackColor
	} else {
		/// Add the style
		AllStyleSettings = append(AllStyleSettings, GStyleSettings{this.Name, this.TextColor, this.BackColor})
		MapStyleSettings[this.Name] = &AllStyleSettings[len(AllStyleSettings)-1]
	}
	MapStyle[this.Name] = this
}

func (this *GStyle) InitStyle() {
	this.DefTextColor = this.TextColor
	this.DefBackColor = this.BackColor
	this.ImportSettings()
}

func InitStyles() {
	MapStyleSettings = make(map[string]*GStyleSettings)
	for i := range AllStyleSettings {
		s := &AllStyleSettings[i]
		MapStyleSettings[s.Name] = s
	}
	MapStyle = make(map[string]*GStyle)

	AllConfigurableStyles = append(AllConfigurableStyles,
		&Normal,
		&Selected,
		&Active,
		&ActFilter,
		&SelActFilter,
		&NormalLight,
		&SelectedLight,

		&Folder1,
		&Folder2,
		&Button_SelFrame,
		&SelButton_SelFrame,
		&SelButton_Frame,
		&DialogButton_Frame,
		&DialogSelButton_Frame,

		&Display,
		&DisplayMessage,
		&ProgressBar,
		&WindowStatusLine,
		&StatusLine,
		&StatusLineMessage,
		&FindPattern,

		&CommandPanel,
		&PageColored,
		&SidepanelKey)

	for _, style := range AllConfigurableStyles {
		style.InitStyle()
	}

	ProgressBarInv = &GStyle{
		ProgressBar.BackColor,
		ProgressBar.TextColor,
		0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "", nil, nil}
}

var (
	/// Configurable

	Normal        = GStyle{231, 234, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Normal line", nil, nil}
	Selected      = GStyle{K.UNSET, 242, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Selected line", &Normal, nil}
	Active        = GStyle{226, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Active line", nil, &Normal}
	ActFilter     = GStyle{16, 86, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Active filter", nil, nil}
	SelActFilter  = GStyle{9, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Selected-active filter", nil, &ActFilter}
	NormalLight   = GStyle{231, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Normal, light", nil, &Normal}
	SelectedLight = GStyle{231, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Selected, light", nil, &Selected}

	Folder1               = GStyle{136, 235, 0, 0, K.ALIGN_LEFT, K.FONT_ITALIC, "Folder, level 1", nil, nil}
	Folder2               = GStyle{137, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_ITALIC, "Folder, level 2", nil, &Folder1}
	Button_SelFrame       = GStyle{232, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Button, Sel Frame", nil, &SelButton_SelFrame}
	SelButton_SelFrame    = GStyle{229, 131, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Sel Button, Sel Frame", nil, nil} // 137
	SelButton_Frame       = GStyle{K.UNSET, 240, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Sel Button, Frame", &SelButton_SelFrame, nil}
	DialogButton_Frame    = GStyle{6, 15, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Dialog Button", nil, nil}
	DialogSelButton_Frame = GStyle{16, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Dialog Sel Button", nil, &DialogButton_Frame}

	Display           = GStyle{195, 23, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Display", nil, nil}
	DisplayMessage    = GStyle{227, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_BLINK, "Display, message", nil, &Display}
	ProgressBar       = GStyle{228, 21, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Progress Bar", nil, nil}
	WindowStatusLine  = GStyle{4, 235, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Window, status line", nil, nil}
	StatusLine        = GStyle{245, 234, 0, 0, K.ALIGN_LEFT, K.FONT_ITALIC, "Status line", nil, nil}
	StatusLineMessage = GStyle{227, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_BLINK, "Status line, message", nil, &StatusLine}
	FindPattern       = GStyle{231, 23, 0, 0, K.ALIGN_LEFT, K.FONT_BOLD, "Find, pattern", nil, nil}

	CommandPanel      = GStyle{250, 234, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Sidepanel", nil, nil}
	SidepanelKey      = GStyle{231, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Sidepanel Key", nil, &CommandPanel}
	CommandPanelTitle = GStyle{K.UNSET, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_BOLD, "Sidepanel title", &SidepanelKey, &CommandPanel}
	PageColored       = GStyle{45, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Page Underlined text", nil, &Normal}

	/// Inverted

	ProgressBarInv *GStyle

	/// Not configurable

	CommandPanelItalic = GStyle{K.UNSET, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_ITALIC, "", &CommandPanel, &CommandPanel}
	SelActive          = GStyle{K.UNSET, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "", &Active, &Selected}
	Normal_            = GStyle{K.UNSET, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_ITALIC, "", &Normal, &Normal}
	Selected_          = GStyle{K.UNSET, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_ITALIC, "", &Selected, &Selected}
	Button_Frame       = GStyle{K.UNSET, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "", &Button_SelFrame, &SelButton_Frame}

	// Active_    = cStyle{K.UNSET, K.UNSET, 0, 0, ALIGN_LEFT, K.FONT_ITALIC, "Active line italic", nil, nil}
	// SelActive_ = cStyle{K.UNSET, K.UNSET, 0, 0, ALIGN_LEFT, K.FONT_ITALIC, "Selected-active line italic", nil, nil}

	FrameText  = GStyle{16, 66, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Frame text", nil, nil}
	Error      = GStyle{3, 16, 0, 0, K.ALIGN_LEFT, K.FONT_BOLD, "Frame text", nil, nil}
	Grey       = GStyle{244, 238, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Grey text", nil, nil}
	TimeBar    = GStyle{15, 16, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "MusicBar", nil, nil}
	BlackWhite = GStyle{15, 235, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Black-White", nil, nil}

	NormalInv = GStyle{238, 231, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Normal line inv", nil, nil}
	Margins   = GStyle{16, 196, 0, 0, K.ALIGN_LEFT, K.FONT_BOLD, "", nil, nil}

	// Page           = cStyle{16, 195, 0, 0, ALIGN_LEFT, K.FONT_NORMAL, "Page", nil, nil}
	PageStatusLine = GStyle{4, 235, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Page StatusLine", nil, nil}
	EmptyLine      = GStyle{K.UNSET, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_NORMAL, "Empty line", &Normal, &Normal}
	PageColCaption = GStyle{246, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_ITALIC, "Columns caption", nil, &Normal}

	PageBold      = GStyle{K.UNSET, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_BOLD, "Page Bold text", &Normal, &Normal}
	PageItalic    = GStyle{K.UNSET, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_ITALIC, "Page Italic text", &Normal, &Normal}
	PageUnderline = GStyle{K.UNSET, K.UNSET, 0, 0, K.ALIGN_LEFT, K.FONT_UNDERLINE, "Page Underlined text", &Normal, &Normal}
)

type GStyleSettings struct {
	Name      string
	TextColor uint16
	BackColor uint16
}

// textColor/backColor are nil when inherited from otherStyle.
type GStyle struct {
	TextColor           uint16
	BackColor           uint16
	DefTextColor        uint16
	DefBackColor        uint16
	Align               uint16
	FontStyle           string
	Name                string
	OtherStyleTextColor *GStyle
	OtherStyleBackColor *GStyle
}

func (Pen *t__Pencils) InitPencils() {
	// Pen.Normal = fmt.Sprintf(COL, CLR_WHITE, CLR_GREY_23)
	Pen.Normal = fmt.Sprintf(K.CLR, 231, 235)
	Pen.NormalInv = fmt.Sprintf(K.CLR, 237, 231)
	Pen.Selected = fmt.Sprintf(K.CLR, 16, 228)
	Pen.Active = fmt.Sprintf(K.CLR, 160, 235)
	Pen.SelectedActive = fmt.Sprintf(K.CLR, 160, 228)
	Pen.StatusLine = fmt.Sprintf(K.CLR, 4, 235)
	Pen.Frame = fmt.Sprintf(K.CLR, 252, 235)
	Pen.EmptyFrame = fmt.Sprintf(K.CLR, 235, 66)
	Pen.Tab = fmt.Sprintf(K.CLR, 48, 235)
	Pen.SelectedTab = fmt.Sprintf(K.CLR, 235, 48)
	Pen.Grey = fmt.Sprintf(K.CLR, 244, 235)
	Pen.ButtonUnpressed = fmt.Sprintf(K.CLR, 16, 251)
	Pen.ButtonPressed = fmt.Sprintf(K.CLR, 16, 244)
	Pen.Error = fmt.Sprintf(K.CLR, 9, 235)
}

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

// const (
// 	K.ALIGN_LEFT = iota
// 	_ALIGN_RIGHT
// 	K.CLR = "\x1b[38;5;%d;48;5;%dm" // color
// )

// const (
// 	K.FONT_NORMAL         = ""
// 	K.FONT_BOLD           = "\x1b[1m"
// 	FONT_BOLD_UNDERLINE = "\x1b[1m\x1b[4m"
// 	K.FONT_ITALIC         = "\x1b[3m"
// 	FONT_UNDERLINE      = "\x1b[4m"
// 	FONT_RESET          = "\x1b[0m"
// 	FONT_BLINK          = "\x1b[5m"
// )

// const (
// 	DC_RESET     = "\x1b[0m"
// 	DC_BOLD      = "\x1b[1m"
// 	DC_ITALIC    = "\x1b[3m"
// 	DC_UNDERLINE = "\x1b[4m"
// )
// const (
// 	K.UNSET       uint16 = 65535
// 	NEVERSET    uint16 = 65534
// 	OUT_OF_PAGE uint16 = 65535
// 	UNDEF              = "UNDEFINED"
// 	K.UNSET_S            = "K.UNSET"
// 	NONE               = "NONE"
// )

