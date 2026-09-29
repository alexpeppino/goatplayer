package trace

import (
	"fmt"
	"goatplayer/internal/K"
)

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

