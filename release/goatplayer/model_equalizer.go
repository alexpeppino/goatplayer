package main

import "fmt"

func (this *mEqualizer) SetOuter() {
	this.outer = this
	this.mBasic.outer = this
}

func (this *mEqualizer2) SetOuter() {
	this.outer = this
	this.mBasic.outer = this
}

func (this *mEqualizer) ResetValues() {
	this.EquSets[0].Name = "EquSet 1"
	this.EquSets[0].Gains = [10]int16{0, 4, 5, 4, 0, -3, -2, -1, -1, -1}
	this.EquSets[1].Name = "EquSet 2"
	this.EquSets[1].Gains = [10]int16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	this.EquSets[2].Name = "EquSet 3"
	this.EquSets[2].Gains = [10]int16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	this.EquSets[3].Name = "Classical"
	this.EquSets[3].Gains = [10]int16{0, 0, 0, 0, 0, 0, -7, -7, -7, -10}
	this.EquSets[4].Name = "Club"
	this.EquSets[4].Gains = [10]int16{0, 0, 8, 6, 6, 6, 3, 0, 0, 0}
	this.EquSets[5].Name = "Dance"
	this.EquSets[5].Gains = [10]int16{10, 7, 2, 0, 0, -6, -7, -7, 0, 0}
	this.EquSets[6].Name = "Full bass"
	this.EquSets[6].Gains = [10]int16{-8, 10, 10, 6, 2, -4, -8, -10, -11, -11}
	this.EquSets[7].Name = "Full bass and treble"
	this.EquSets[7].Gains = [10]int16{7, 6, 0, -7, -5, 2, 8, 11, 12, 12}
	this.EquSets[8].Name = "Full treble"
	this.EquSets[8].Gains = [10]int16{-10, -10, -10, -4, 2, 11, 16, 16, 16, 17}
	this.EquSets[9].Name = "Headphones"
	this.EquSets[9].Gains = [10]int16{5, 11, 6, -3, -2, 2, 5, 10, 13, 14}
	this.EquSets[10].Name = "Large Hall"
	this.EquSets[10].Gains = [10]int16{10, 10, 6, 6, 0, -5, -5, -5, 0, 0}
	this.EquSets[11].Name = "Live"
	this.EquSets[11].Gains = [10]int16{-5, 0, 4, 6, 6, 6, 4, 2, 2, 2}
	this.EquSets[12].Name = "Party"
	this.EquSets[12].Gains = [10]int16{7, 7, 0, 0, 0, 0, 0, 0, 7, 7}
	this.EquSets[13].Name = "Pop"
	this.EquSets[13].Gains = [10]int16{-2, 5, 7, 8, 6, 0, -2, -2, -2, -2}
	this.EquSets[14].Name = "Reggae"
	this.EquSets[14].Gains = [10]int16{0, 0, 0, -6, 0, 6, 6, 0, 0, 0}
	this.EquSets[15].Name = "Rock"
	this.EquSets[15].Gains = [10]int16{8, 5, -6, -8, -3, 4, 9, 11, 11, 11}
	this.EquSets[16].Name = "Ska"
	this.EquSets[16].Gains = [10]int16{-2, -5, -4, 0, 4, 6, 9, 10, 11, 10}
	this.EquSets[17].Name = "Soft"
	this.EquSets[17].Gains = [10]int16{5, 2, 0, -2, 0, 4, 8, 10, 11, 12}
	this.EquSets[18].Name = "Soft rock"
	this.EquSets[18].Gains = [10]int16{4, 4, 2, 0, -4, -6, -3, 0, 2, 9}
	this.EquSets[19].Name = "Techno"
	this.EquSets[19].Gains = [10]int16{8, 6, 0, -6, -5, 0, 8, 10, 10, 9}

}

func (this *mEqualizer) InitModel(name string, focusKey string, helpfile string) {
	if this.outer == nil {
		this.SetOuter()
		this.m2.SetOuter()
	}
	this.mBasic.InitModel(name, focusKey, helpfile)
	// this._fSelectKey = func() {
	// 	this.PrepareW2()
	// 	this.w2.PrintWindowAnyway()
	// }
}

func (this *mEqualizer) PrepareModel() {
	this.mBasic.PrepareModel()
	// this._fSelectKey = func() {
	// 	this.PrepareW2()
	// 	this.w2.PrintWindowAnyway()
	// }
}

func (this *mEqualizer) SelectKey_o() {
// trace.Begin_n(this, "SelectKey_o") // t //
	this.PrepareW2()
	this.m2.w.PrintWindowAnyway()
// trace.End() // t //
}

func (this *mEqualizer) GetEqusetStringOneFreq(iEquset, iFreq uint16) string {
	e := &this.EquSets[iEquset]
	return fmt.Sprintf("@e%d:equalizer=f=%s:t=q:w=1:g=%d", iFreq, frequencies[iFreq], e.Gains[iFreq])
}

/*
whichIndex: UNSET => use value in settings; UNSET-1 => use selected in window. doAddAF: true = add "--af".
*/
func (this *mEqualizer) GetEqusetString(whichIndex uint16, doAddAF bool) string {
	var s2 string
	switch whichIndex {
	case UNSET:
		whichIndex = Set.Equalizer.Pos
	case UNSET - 1:
		whichIndex = this.keys.iSelected
	}
	e := &this.EquSets[whichIndex]
	le := len(e.Gains)
	if doAddAF {
		s2 = "--af="
	}
	for j := range le {
		s2 += fmt.Sprintf("@e%d:equalizer=f=%s:t=q:w=1:g=%d", j, frequencies[j], e.Gains[j])
		if j < le-1 {
			s2 += ","
		}
	}
	// trace.Print("equalizer: %s", s2)
	// @e0:equalizer=f=36:t=q:w=1:g=0
	return s2
}

func (this *mEqualizer) PrepareW2Line(iLine uint16) {
	iEqu := this.keys.iSelected
	gain := this.EquSets[iEqu].Gains[iLine]
	var s2 string
	var pallina = "●"
	var barra = "┊"
	if gain < 0 {
		s2 = fmt.Sprintf("%*s%*s", 15+gain, pallina, 0-gain, barra)
	} else if gain > 0 {
		s2 = fmt.Sprintf("%*s%*s", 15, barra, gain, pallina)
	} else {
		s2 = fmt.Sprintf("%15s", pallina)
	}
	s := fmt.Sprintf("%7s   %+5d             %s", frequencies_[iLine], gain, s2)
	this.m2.w.SetLineForDialog(iLine, s, iLine, 0)
}

func (this *mEqualizer) PrepareW2() {
	w2 := &this.m2.w
	this.m2.keys.iActive = UNSET
	equset := &this.EquSets[this.keys.iSelected]
	le := len(equset.Gains)
	this.m2.keys.Make(le)
	w2.lines.Make(le)
	for i := range le {
		this.PrepareW2Line(uint16(i))
	}
}

func (this *mEqualizer) PrepareView() {
	this.Reset()
	this.m2.Reset()
	this.keys.iActive = Set.Equalizer.Pos
	this.SelectKeyAndLine(0)
	this.m2.SelectKeyAndLine(0)
	w := &this.w
	le := len(this.EquSets)
	this.keys.Make(le)
	w.lines.Make(le)
	for i := range le {
		e := &this.EquSets[i]
		if e.Name == "" {
			e.Name = fmt.Sprintf("equset%d", i+1)
		}
		// e.Freq = [10]uint16{}
		w.SetLineForDialog(uint16(i), "  "+e.Name, uint16(i), 0)
	}
	this.SelectKeyAndLine(0)
	w.SetEmptyLineForDialog()
	w.columnsCaption = stylePageColCaption.EchoStyle("  Preset", this.w.fGetWidth())
	//	w2
	w2 := &this.m2.w
	w2.columnsCaption = stylePageColCaption.EchoStyle("   Freq.     Gain                       - ┊ +", this.m2.w.fGetWidth())
	this.PrepareW2()
	w2.SetEmptyLineForDialog()
}

func (this *mEqualizer) FreqM1() {
	if this.m2.w.lines.iSelected > 0 {
		this.m2.w.lines.iSelected--
	}
	this.m2.w.PrintWindowAnyway()
}

func (this *mEqualizer) FreqP1() {
	if this.m2.w.lines.iSelected < 9 {
		this.m2.w.lines.iSelected++
	}
	this.m2.w.PrintWindowAnyway()
}

func (this *mEqualizer) GainM1() {
	w2 := &this.m2.w
	iEquset := this.keys.iSelected
	iFreq := w2.lines.iSelected
	if this.EquSets[iEquset].Gains[iFreq] > -15 {
		this.EquSets[iEquset].Gains[iFreq]--
		if iEquset == Set.Equalizer.Pos {
			Mpv.AudiofilterRemove(fmt.Sprintf("@e%d", iFreq))
			Mpv.AudiofilterAdd(this.GetEqusetStringOneFreq(Set.Equalizer.Pos, iFreq))
		}
	}
	this.PrepareW2Line(w2.lines.iSelected)
	w2.PrintWindowAnyway()
}

func (this *mEqualizer) GainP1() {
	w2 := &this.m2.w
	iEquset := this.keys.iSelected
	iFreq := w2.lines.iSelected
	if this.EquSets[iEquset].Gains[iFreq] < 15 {
		this.EquSets[iEquset].Gains[iFreq]++
		if iEquset == Set.Equalizer.Pos {
			Mpv.AudiofilterRemove(fmt.Sprintf("@e%d", iFreq))
			Mpv.AudiofilterAdd(this.GetEqusetStringOneFreq(Set.Equalizer.Pos, iFreq))
		}
	}
	this.PrepareW2Line(w2.lines.iSelected)
	w2.PrintWindowAnyway()
}

type mEqualizer struct {
	mBasic
	m2          mEqualizer2
	paletteEcho string
	EquSets     [20]tEquset
	// EquSetFilters [10]string
}

type mEqualizer2 struct {
	mBasic
}

var frequencies = []string{"36", "64", "125", "250", "500", "1000", "2000", "4000", "8000", "16000"}
var frequencies_ = []string{"36", "64", "125", "250", "500", "1K", "2K", "4K", "8K", "16K"}

type tEquset struct {
	Name  string
	Gains [10]int16
	echo  string
	// Items [20]tEquItem
}

// type tEquItem struct {
// 	Gain int16
// 	Q    uint16
// 	Freq uint16
// }

