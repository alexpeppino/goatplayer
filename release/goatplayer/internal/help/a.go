package help

import (
	"goatplayer/internal/K"
	"sort"
)

var MepTitles map[string]*string
var AllTitles []string

func Init() {
	MepTitles = make(map[string]*string)
	MepTitles[K.HELPFILE_BROWSERS] = &Browsers
	MepTitles[K.HELPFILE_COLORS] = &Colors
	MepTitles[K.HELPFILE_DIRS] = &Dirs
	MepTitles[K.HELPFILE_DOWNLOAD] = &Download
	MepTitles[K.HELPFILE_EQUALIZER] = &Equalizer
	MepTitles[K.HELPFILE_FIND] = &Find
	MepTitles[K.HELPFILE_HELP] = &Help
	MepTitles[K.HELPFILE_LISTS] = &Lists
	MepTitles[K.HELPFILE_PARAMETERS] = &Parameters
	MepTitles[K.HELPFILE_PLAYER] = &Player
	MepTitles[K.HELPFILE_PROGRESS_BARS] = &ProgressBars
	MepTitles[K.HELPFILE_PROPERTIES] = &Properties
	MepTitles[K.HELPFILE_SECTIONS] = &Sections
	MepTitles[K.HELPFILE_SETTINGS] = &Settings
	MepTitles[K.HELPFILE_TAG_FIELDS] = &TagFields
	MepTitles[K.HELPFILE_TERMINOLOGY] = &Terminology

	for s := range MepTitles {
		AllTitles = append(AllTitles, s)
	}
	sort.Strings(AllTitles)
}

