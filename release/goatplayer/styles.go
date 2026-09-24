package main

import (
	"fmt"
)

func (this *cStyle) EchoStyle(s string, width uint16) string {
	var textColor, backColor uint16
	var em string
	if this.textColor == UNSET {
		textColor = this.otherStyleTextColor.textColor
	} else {
		textColor = this.textColor
	}
	if this.backColor == UNSET {
		backColor = this.otherStyleBackColor.backColor
	} else {
		backColor = this.backColor
	}
	color := fmt.Sprintf("\x1b[38;5;%d;48;5;%dm", textColor, backColor)
	// s = strings.TrimSpace(s)
	if width > 0 {
		AdjustWidthSeparated(&s, &em, width)
		em = fmt.Sprintf("%s%s%s", color, em, FONT_RESET)
	}
	s = fmt.Sprintf("%s%s%s%s", color, this.fontStyle, s, FONT_RESET)
	if this.align == ALIGN_LEFT {
		return s + em
	} else {
		return em + s
	}
}

func AdjustWidthSeparated(str, em *string, finalLen uint16) {
	var le = RuneLen(str)
	// var le = uint16(len(*str))
	if le > finalLen {
		*str = string([]rune(*str)[0:finalLen])
	} else {
		*em = fmt.Sprintf("%-*s", finalLen-le, *em)
	}
}

/* Don't use with fontstyle underlined or inherited colors (use Echo instead). */
func (this *cStyle) EchoStyleSimple(s string, width uint16) string {
	AdjustWidthAligned(&s, width, this.align)
	return fmt.Sprintf("\x1b[38;5;%d;48;5;%dm%s%s%s", this.textColor, this.backColor, this.fontStyle, s, FONT_RESET)
}

/* Adds a blank before and after. Doesn't change the width. */
func (this *cStyle) EchoStyleFormatOnly(s string, addBlank bool) string {
	var textColor, backColor uint16
	if this.textColor == UNSET {
		textColor = this.otherStyleTextColor.textColor
	} else {
		textColor = this.textColor
	}
	if this.backColor == UNSET {
		backColor = this.otherStyleBackColor.backColor
	} else {
		backColor = this.backColor
	}
	if addBlank {
		s = " " + s + " "
	}
	return fmt.Sprintf("\x1b[38;5;%d;48;5;%dm%s%s%s", textColor, backColor, this.fontStyle, s, FONT_RESET)
}

func (this *cStyle) EchoTwoStyles(s, right string, width uint16, otherStyle *cStyle) string {
	var leRight = RuneLen(&right)
	s1 := this.EchoStyle(s, width-leRight)
	s2 := otherStyle.EchoStyle(right, leRight)
	return s1 + s2
}

/* Returns only the echo of the colors. */
func (this *cStyle) JustTheColors() string {
	var textColor, backColor uint16
	if this.textColor == UNSET {
		textColor = this.otherStyleTextColor.textColor
	} else {
		textColor = this.textColor
	}
	if this.backColor == UNSET {
		backColor = this.otherStyleBackColor.backColor
	} else {
		backColor = this.backColor
	}
	return fmt.Sprintf("\x1b[38;5;%d;48;5;%dm", textColor, backColor)
}

func styleBold(s string) string {
	return DC_BOLD + s + DC_RESET
}

func styleItalic(s string) string {
	return DC_ITALIC + s + DC_RESET
}

func styleUnderline(s string) string {
	return DC_UNDERLINE + s + DC_RESET
}

func (this *cStyle) ImportSettings() {
	if styleSet, ok := AllMaps.mapStyleSettings[this.name]; ok {
		/// Read the style from settings
		this.textColor = styleSet.TextColor
		this.backColor = styleSet.BackColor
	} else {
		/// Add the style
		AllMaps.allStyleSettings = append(AllMaps.allStyleSettings, cStyleSettings{this.name, this.textColor, this.backColor})
		AllMaps.mapStyleSettings[this.name] = &AllMaps.allStyleSettings[len(AllMaps.allStyleSettings)-1]
	}
	AllMaps.mapStyle[this.name] = this
}

func (this *cStyle) InitStyle() {
	this.defTextColor = this.textColor
	this.defBackColor = this.backColor
	this.ImportSettings()
}

func InitStyles() {
	AllMaps.mapStyleSettings = make(map[string]*cStyleSettings)
	for i := range AllMaps.allStyleSettings {
		s := &AllMaps.allStyleSettings[i]
		AllMaps.mapStyleSettings[s.Name] = s
	}
	AllMaps.mapStyle = make(map[string]*cStyle)

	AllMaps.allConfigurableStyles = append(AllMaps.allConfigurableStyles,
		&styleNormal,
		&styleSelected,
		&styleActive,
		&styleActFilter,
		&styleSelActFilter,
		&styleNormalLight,
		&styleSelectedLight,

		&styleFolder1,
		&styleFolder2,
		&styleButton_SelFrame,
		&styleSelButton_SelFrame,
		&styleSelButton_Frame,
		&styleDialogButton_Frame,
		&styleDialogSelButton_Frame,

		&styleDisplay,
		&styleDisplayMessage,
		&styleProgressBar,
		&styleWindowStatusLine,
		&styleStatusLine,
		&styleStatusLineMessage,
		&styleFindPattern,

		&styleCommandPanel,
		&stylePageColored,
		&styleSidepanelKey)

	for _, style := range AllMaps.allConfigurableStyles {
		style.InitStyle()
	}

	styleProgressBarInv = &cStyle{
		styleProgressBar.backColor,
		styleProgressBar.textColor,
		0, 0, ALIGN_LEFT, FONT_NORMAL, "", nil, nil}
}

var (
	/// Configurable

	styleNormal        = cStyle{231, 234, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Normal line", nil, nil}
	styleSelected      = cStyle{UNSET, 242, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Selected line", &styleNormal, nil}
	styleActive        = cStyle{226, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Active line", nil, &styleNormal}
	styleActFilter     = cStyle{16, 86, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Active filter", nil, nil}
	styleSelActFilter  = cStyle{9, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Selected-active filter", nil, &styleActFilter}
	styleNormalLight   = cStyle{231, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Normal, light", nil, &styleNormal}
	styleSelectedLight = cStyle{231, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Selected, light", nil, &styleSelected}

	styleFolder1               = cStyle{136, 235, 0, 0, ALIGN_LEFT, FONT_ITALIC, "Folder, level 1", nil, nil}
	styleFolder2               = cStyle{137, UNSET, 0, 0, ALIGN_LEFT, FONT_ITALIC, "Folder, level 2", nil, &styleFolder1}
	styleButton_SelFrame       = cStyle{232, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Button, Sel Frame", nil, &styleSelButton_SelFrame}
	styleSelButton_SelFrame    = cStyle{229, 131, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Sel Button, Sel Frame", nil, nil} // 137
	styleSelButton_Frame       = cStyle{UNSET, 240, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Sel Button, Frame", &styleSelButton_SelFrame, nil}
	styleDialogButton_Frame    = cStyle{6, 15, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Dialog Button", nil, nil}
	styleDialogSelButton_Frame = cStyle{16, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Dialog Sel Button", nil, &styleDialogButton_Frame}

	styleDisplay           = cStyle{195, 23, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Display", nil, nil}
	styleDisplayMessage    = cStyle{227, UNSET, 0, 0, ALIGN_LEFT, FONT_BLINK, "Display, message", nil, &styleDisplay}
	styleProgressBar       = cStyle{228, 21, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Progress Bar", nil, nil}
	styleWindowStatusLine  = cStyle{4, 235, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Window, status line", nil, nil}
	styleStatusLine        = cStyle{245, 234, 0, 0, ALIGN_LEFT, FONT_ITALIC, "Status line", nil, nil}
	styleStatusLineMessage = cStyle{227, UNSET, 0, 0, ALIGN_LEFT, FONT_BLINK, "Status line, message", nil, &styleStatusLine}
	styleFindPattern       = cStyle{231, 23, 0, 0, ALIGN_LEFT, FONT_BOLD, "Find, pattern", nil, nil}

	styleCommandPanel      = cStyle{250, 234, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Sidepanel", nil, nil}
	styleSidepanelKey      = cStyle{231, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Sidepanel Key", nil, &styleCommandPanel}
	styleCommandPanelTitle = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_BOLD, "Sidepanel title", &styleSidepanelKey, &styleCommandPanel}
	stylePageColored       = cStyle{45, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Page Underlined text", nil, &styleNormal}

	/// Inverted

	styleProgressBarInv *cStyle

	/// Not configurable

	styleCommandPanelItalic = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_ITALIC, "", &styleCommandPanel, &styleCommandPanel}
	styleSelActive          = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "", &styleActive, &styleSelected}
	styleNormal_            = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_ITALIC, "", &styleNormal, &styleNormal}
	styleSelected_          = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_ITALIC, "", &styleSelected, &styleSelected}
	styleButton_Frame       = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "", &styleButton_SelFrame, &styleSelButton_Frame}

	// styleActive_    = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_ITALIC, "Active line italic", nil, nil}
	// styleSelActive_ = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_ITALIC, "Selected-active line italic", nil, nil}

	styleFrameText  = cStyle{16, 66, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Frame text", nil, nil}
	styleError      = cStyle{3, 16, 0, 0, ALIGN_LEFT, FONT_BOLD, "Frame text", nil, nil}
	styleGrey       = cStyle{244, 238, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Grey text", nil, nil}
	styleTimeBar    = cStyle{15, 16, 0, 0, ALIGN_LEFT, FONT_NORMAL, "MusicBar", nil, nil}
	styleBlackWhite = cStyle{15, 235, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Black-White", nil, nil}

	styleNormalInv = cStyle{238, 231, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Normal line inv", nil, nil}
	styleMargins   = cStyle{16, 196, 0, 0, ALIGN_LEFT, FONT_BOLD, "", nil, nil}

	// stylePage           = cStyle{16, 195, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Page", nil, nil}
	stylePageStatusLine = cStyle{4, 235, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Page StatusLine", nil, nil}
	styleEmptyLine      = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_NORMAL, "Empty line", &styleNormal, &styleNormal}
	stylePageColCaption = cStyle{246, UNSET, 0, 0, ALIGN_LEFT, FONT_ITALIC, "Columns caption", nil, &styleNormal}

	stylePageBold      = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_BOLD, "Page Bold text", &styleNormal, &styleNormal}
	stylePageItalic    = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_ITALIC, "Page Italic text", &styleNormal, &styleNormal}
	stylePageUnderline = cStyle{UNSET, UNSET, 0, 0, ALIGN_LEFT, FONT_UNDERLINE, "Page Underlined text", &styleNormal, &styleNormal}
)

var (
	traceStyleNormal = cTraceStyle{231, FONT_NORMAL}
	traceStyleBegin  = cTraceStyle{157, FONT_BOLD}
	traceStyleEnd    = cTraceStyle{157, FONT_NORMAL}
	traceStyleError  = cTraceStyle{3, FONT_NORMAL}
	traceStylePlayer = cTraceStyle{4, FONT_NORMAL}
	traceStyleGreen  = cTraceStyle{83, FONT_NORMAL}
	traceStyleFocus  = cTraceStyle{80, FONT_NORMAL}
)

type cStyleSettings struct {
	Name      string
	TextColor uint16
	BackColor uint16
}

// textColor/backColor are nil when inherited from otherStyle.
type cStyle struct {
	textColor           uint16
	backColor           uint16
	defTextColor        uint16
	defBackColor        uint16
	align               uint16
	fontStyle           string
	name                string
	otherStyleTextColor *cStyle
	otherStyleBackColor *cStyle
}

var (
// ////	Styles
)

func (Pen *t__Pencils) InitPencils() {
	// Pen.Normal = fmt.Sprintf(COL, CLR_WHITE, CLR_GREY_23)
	Pen.Normal = fmt.Sprintf(CLR, 231, 235)
	Pen.NormalInv = fmt.Sprintf(CLR, 237, 231)
	Pen.Selected = fmt.Sprintf(CLR, 16, 228)
	Pen.Active = fmt.Sprintf(CLR, 160, 235)
	Pen.SelectedActive = fmt.Sprintf(CLR, 160, 228)
	Pen.StatusLine = fmt.Sprintf(CLR, 4, 235)
	Pen.Frame = fmt.Sprintf(CLR, 252, 235)
	Pen.EmptyFrame = fmt.Sprintf(CLR, 235, 66)
	Pen.Tab = fmt.Sprintf(CLR, 48, 235)
	Pen.SelectedTab = fmt.Sprintf(CLR, 235, 48)
	Pen.Grey = fmt.Sprintf(CLR, 244, 235)
	Pen.ButtonUnpressed = fmt.Sprintf(CLR, 16, 251)
	Pen.ButtonPressed = fmt.Sprintf(CLR, 16, 244)
	Pen.Error = fmt.Sprintf(CLR, 9, 235)
}

