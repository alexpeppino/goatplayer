package main

import (
	"fmt"
	"slices"
	"sort"
)

func (this *mContainer) SelectActive_o() {
// trace.Begin_n(this, "SelectActive") // t //
	Input.DisplayInput()
	for i := this.keys.iSelected + 1; i < uint16(len(this.keys.sl)); i++ {
		if this.keys.sl[i].isActiveFilter {
			this.SelectKeyAndLine(i)
			this.w.PrintWindow()
// trace.ReturnAdd("i: %d", i) // t //
			return
		}
	}
	for i := uint16(0); i < this.keys.iSelected; i++ {
		if this.keys.sl[i].isActiveFilter {
			this.SelectKeyAndLine(i)
			this.w.PrintWindow()
// trace.ReturnAdd("i: %d", i) // t //
			return
		}
	}
// trace.End() // t //
}

func (this *mContainer) SelectKey_o() {
// trace.BeginEnd_n(this, "SelectKey_o {mContainer}") // t //
	Term.StatusLine.SetContent(this.keys.GetSelected().v)
}

func (this *mContainer) SetOuter() {
	this.mBasic.outer = this
	this.outer = this
}

func (this *mContainer) InitModel(name string, sFocus string, helpfile string) {
	if this.outer == nil {
		this.SetOuter()
	}
	this.mBasic.InitModel(name, sFocus, helpfile)
	AllMaps.mapContainers[name] = this
	AllMaps.allContainers = append(AllMaps.allContainers, this)
}

func (this *mContainer) PrepareModel() {
	this.mBasic.PrepareModel()
}

func (this *mContainer) ApplyAllFilters() {

}

func (this *mContainer) UpdateAllValues() {
	// for _, i := range this.pBrowserWsp.dbIndexes.sl {

	// }
}

func (this *mContainer) ActivateFilterByTRow(iTRow uint16) {
// trace.BeginAdd_n(this, "ActivateFilterByTRow", "iTRow: %d", iTRow) // t //
	iLine := this.w.CalcLineIndexByTRow(iTRow)
	this.ActivateFilter(iLine)
// trace.End() // t //
}

func (this *mContainer) ActivateFilter(iLine uint16) {
	this.keys.iActive = UNSET
	// iLine := this.w.CalcLineIndexByTRow(iTRow)
// trace.BeginAdd_n(this, "ActivateFilter", "iLine: %d", iLine) // t //
	p := this.w.lines.Get(iLine)
	if p == nil || p.iKey == UNSET {
		// trace.N_Error(this, "SelectLineByTRow, invalid index %d", iRow)
// trace.ReturnAdd("invalid index") // t //
		return
	}
	///	Toggle new filter
	this.ToggleActiveFilter(p.iKey)
	// toggleIsTrue := this.ToggleActiveFilter(p.iKey)
	this.w.PrintLine(iLine)
	///
	this.FilterAndPrint__new()
// trace.End() // t //
}

func (this *mContainer) DontFilterAndPrint() {
// trace.Begin_n(this, "DontFilterAndPrint") // t //
	browser := this.pBrowserWsp
	browser.View.LoadAndPrepare(&this.pBrowserWsp.wspDbIndexes)
	browser.View.w.PrintWindowAnyway()
	browser.GenresCon.LoadAllValues(false, false)
	browser.ArtistsCon.LoadAllValues(false, false)
	browser.AlbumsCon.LoadAllValues(false, false)
	browser.AlbumArtistsCon.LoadAllValues(false, false)
	browser.ComposersCon.LoadAllValues(false, false)
	browser.YearsCon.LoadAllValues(false, false)
	this.w.BringWindowToFront(true)
// trace.End() // t //
}

func (this *mContainer) FilterAndPrint__new() {
// trace.Begin_n(this, "FilterAndPrint__new") // t //
	var dbIndexes tDbIndexes
	// var allContainersWithFilters []*mContainer
	browser := this.pBrowserWsp
	this.w.PrintStatusLine()
	browser.View.Reset()

	nGenres := browser.GenresCon.SetValidValues()
	nArtists := browser.ArtistsCon.SetValidValues()
	nAlbums := browser.AlbumsCon.SetValidValues()
	nAlbumArtists := browser.AlbumArtistsCon.SetValidValues()
	nComposers := browser.ComposersCon.SetValidValues()
	nYears := browser.YearsCon.SetValidValues()
	tot := nGenres + nArtists + nAlbums + nAlbumArtists + nComposers + nYears
// trace.Print("Total valid values in all containers: %d", tot) // t //

	/// Case: There is no filtering in all containers
	if tot == 0 {
		this.DontFilterAndPrint()
		return
		// browser.View.LoadAndPrepare(&this.pBrowserWsp.wspDbIndexes)
		// browser.View.w.PrintWindowAnyway()
		// browser.GenresCon.LoadAllValues(false, false)
		// browser.ArtistsCon.LoadAllValues(false, false)
		// browser.AlbumsCon.LoadAllValues(false, false)
		// browser.AlbumArtistsCon.LoadAllValues(false, false)
		// browser.ComposersCon.LoadAllValues(false, false)
		// browser.YearsCon.LoadAllValues(false, false)
		// this.w.BringWindowToFront(true)
		// trace.ReturnAdd("There is no filtering in all containers")
	}

	var mapValidGenres map[string]uint16
	var mapValidArtists map[string]uint16
	var mapValidAlbums map[string]uint16
	var mapValidAlbumArtists map[string]uint16
	var mapValidComposers map[string]uint16
	var mapValidYears map[string]uint16
	mapValidGenres = make(map[string]uint16)
	mapValidArtists = make(map[string]uint16)
	mapValidAlbums = make(map[string]uint16)
	mapValidAlbumArtists = make(map[string]uint16)
	mapValidComposers = make(map[string]uint16)
	mapValidYears = make(map[string]uint16)

	var mapPresentArtists map[string]uint16
	var mapPresentAlbums map[string]uint16
	var mapPresentAlbumArtists map[string]uint16
	var mapPresentComposers map[string]uint16
	var mapPresentYears map[string]uint16
	mapPresentArtists = make(map[string]uint16)
	mapPresentAlbums = make(map[string]uint16)
	mapPresentAlbumArtists = make(map[string]uint16)
	mapPresentComposers = make(map[string]uint16)
	mapPresentYears = make(map[string]uint16)

	var mapAllArtists map[string]uint16
	var mapAllAlbums map[string]uint16
	var mapAllAlbumArtists map[string]uint16
	var mapAllComposers map[string]uint16
	var mapAllYears map[string]uint16
	mapAllArtists = make(map[string]uint16)
	mapAllAlbums = make(map[string]uint16)
	mapAllAlbumArtists = make(map[string]uint16)
	mapAllComposers = make(map[string]uint16)
	mapAllYears = make(map[string]uint16)

	// var validGenres []string
	// var validArtists []string
	// var validAlbums []string
	// var validAlbumArtists []string
	// var validComposers []string
	// var validYears []string

	// var presentAlbums []string
	// var presentAlbumArtists []string
	// var presentComposers []string
	// var presentYears []string

	var allArtists []string
	var allAlbums []string
	var allAlbumArtists []string
	var allComposers []string
	var allYears []string

	// var valid3rdLevel, validAlbum, validAlbumArtist, ValidComposer, validYear bool

	/// Case: There is filtering in the containers
// trace.Print("There is filtering in the containers (%d valid values)", tot) // t //
	for _, i := range browser.wspDbIndexes.sl {
		audiof := &AF.sl[i]
		if nGenres == 0 || browser.GenresCon.ApplyFilterOR(audiof) {
			// trace.Print("genres, valid audiof: %s", audiof.filename)
			if nGenres != 0 {
				mapValidGenres[audiof.genre]++
			}
			mapAllArtists[audiof.artist]++
			if nArtists == 0 || browser.ArtistsCon.ApplyFilterOR(audiof) {
				// trace.Print("artists, valid audiof: %s", audiof.filename)
				if nArtists != 0 {
					mapValidArtists[audiof.artist]++
				}
				mapAllAlbums[audiof.album]++
				mapAllAlbumArtists[audiof.albumArtist]++
				mapAllComposers[audiof.composer]++
				mapAllYears[audiof.year]++

				if nAlbums == 0 || browser.AlbumsCon.ApplyFilterOR(audiof) {
					if nAlbumArtists == 0 || browser.AlbumArtistsCon.ApplyFilterOR(audiof) {
						if nComposers == 0 || browser.ComposersCon.ApplyFilterOR(audiof) {
							if nYears == 0 || browser.YearsCon.ApplyFilterOR(audiof) {
// trace.Print("valid audiof: %s", audiof.filename) // t //
								if nAlbums != 0 {
									mapValidAlbums[audiof.album]++
								}
								if nAlbumArtists != 0 {
									mapValidAlbumArtists[audiof.albumArtist]++
								}
								if nComposers != 0 {
									mapValidComposers[audiof.composer]++
								}
								if nYears != 0 {
									mapValidYears[audiof.year]++
								}
								dbIndexes.AppendIndex(uint16(i))
							}
						}
					}
				}
			}
		}
	}
// trace.Print("dbIndexes:  %#v", dbIndexes) // t //

	if len(dbIndexes.sl) == 0 {
		Term.StatusLine.PrintMessage("Filters have been reset (empty result)")
		this.DontFilterAndPrint()
// trace.Return() // t //
		return
	}

	for _, u := range dbIndexes.sl {
		audiof := &AF.sl[u]
		mapPresentArtists[audiof.artist]++
		mapPresentAlbums[audiof.album]++
		mapPresentAlbumArtists[audiof.albumArtist]++
		mapPresentComposers[audiof.composer]++
		mapPresentYears[audiof.year]++
	}

	// dump.DumpMapStringUint("all artists", &mapAllArtists)
	// dump.DumpMapStringUint("all albums", &mapAllAlbums)
	// dump.DumpMapStringUint("all albumArtists", &mapAllAlbumArtists)
	// dump.DumpMapStringUint("all composers", &mapAllComposers)
	// dump.DumpMapStringUint("all years", &mapAllYears)
	// dump.DumpMapStringUint("valid genres", &mapValidGenres)
	// dump.DumpMapStringUint("valid artists", &mapValidArtists)
	// dump.DumpMapStringUint("valid albums", &mapValidAlbums)
	// dump.DumpMapStringUint("valid albumArtists", &mapValidAlbumArtists)
	// dump.DumpMapStringUint("valid composers", &mapValidComposers)
	// dump.DumpMapStringUint("valid years", &mapValidYears)
	// dump.DumpMapStringUint("present albums", &mapPresentAlbums)
	// dump.DumpMapStringUint("present albumArtists", &mapPresentAlbumArtists)
	// dump.DumpMapStringUint("present composers", &mapPresentComposers)
	// dump.DumpMapStringUint("present years", &mapPresentYears)

	for name := range mapAllArtists {
		allArtists = append(allArtists, name)
	}
	sort.Strings(allArtists)

	for name := range mapAllAlbums {
		allAlbums = append(allAlbums, name)
	}
	sort.Strings(allAlbums)

	for name := range mapAllAlbumArtists {
		allAlbumArtists = append(allAlbumArtists, name)
	}
	sort.Strings(allAlbumArtists)

	for name := range mapAllComposers {
		allComposers = append(allComposers, name)
	}
	sort.Strings(allComposers)

	for name := range mapAllYears {
		allYears = append(allYears, name)
	}
	sort.Strings(allYears)

	if this == &browser.GenresCon {
		// if len(mapValidArtists) > 0 {
		browser.ArtistsCon.LoadFromSliceAndMap(&allArtists, &mapAllArtists)
		// } else {
		// 	browser.ArtistsCon.LoadFromSliceAndMap(&allArtists, &mapPresentArtists)
		// }
	} else if this == &browser.ArtistsCon {
		iSelected := this.keys.iSelected
		// if len(mapValidArtists) > 0 {
		browser.ArtistsCon.LoadFromSliceAndMap(&allArtists, &mapAllArtists)
		// } else {
		// 	browser.ArtistsCon.LoadFromSliceAndMap(&allArtists, &mapPresentArtists)
		// }
		this.SelectKeyAndLine(iSelected)
		this.w.iOldPage = UNSET
	}
	iSelected := this.keys.iSelected
	browser.View.LoadAndPrepare(&dbIndexes)
	browser.View.w.PrintWindowAnyway()
	browser.AlbumsCon.LoadFromSliceAndMap3rdLevel(&allAlbums, &mapValidAlbums, &mapPresentAlbums)
	browser.AlbumArtistsCon.LoadFromSliceAndMap3rdLevel(&allAlbumArtists, &mapValidAlbumArtists, &mapPresentAlbumArtists)
	browser.ComposersCon.LoadFromSliceAndMap3rdLevel(&allComposers, &mapValidComposers, &mapPresentComposers)
	browser.YearsCon.LoadFromSliceAndMap3rdLevel(&allYears, &mapValidYears, &mapPresentYears)
	this.SelectKeyAndLine(iSelected)

	this.w.BringWindowToFront(false)
	Term.Flush()
// trace.End() // t //
}

func (this *mContainer) LoadFromSliceAndMap(allList *[]string, mapValidOrPresent *map[string]uint16) {
// trace.Begin_n(this, "LoadFromSliceAndPrint") // t //
// dump.DumpMapStringUint("valid artists", mapValidOrPresent) // t //
	this.keys.iActive = UNSET
	this.Reset()
	this.keys.sl = make([]tTableKey, len(*allList))
	this.w.lines.sl = make([]tTableLine, len(*allList))
	for i, name := range *allList {
		key := &this.keys.sl[i]
		line := &this.w.lines.sl[i]
		n := (*mapValidOrPresent)[name]
		// trace.Print("name: %s, u: %d", name, (*ma)[name])
		var t string
		if n > 0 {
			t = fmt.Sprintf("  %s (%d)", name, n)
		} else {
			t = fmt.Sprintf("  %s", name)
		}
		key.v = name
		key.iLine = uint16(i)
		line.iKey = uint16(i)
		line.text = t
		line.normal = styleNormal.EchoStyle(t, this.w.nCols)
		line.selected = styleSelected.EchoStyle(t, this.w.nCols)
	}
// trace.Print("filterValues %s", this.validValues) // t //
	for _, value := range this.validValues {
		i := this.keys.GetIndexByFilename(value)
		if key := this.keys.Get(i); key != nil {
			key.isActiveFilter = true
		}
	}
	this.AfterLoad()
	this.w.SetEmptyLine()
	// this.PrepareView()
// trace.End() // t //
}

func (this *mContainer) LoadFromSliceAndMap3rdLevel(allList *[]string, mapValid, mapPresent *map[string]uint16) {
// trace.Begin_n(this, "LoadFromSliceAndMap3rdLevel") // t //
	this.keys.iActive = UNSET
	this.Reset()
	this.keys.sl = make([]tTableKey, len(*allList))
	this.w.lines.sl = make([]tTableLine, len(*allList))
	if len((*mapValid)) > 0 {
		for i, name := range *allList {
			key := &this.keys.sl[i]
			line := &this.w.lines.sl[i]
			if v, ok := (*mapValid)[name]; ok {
// trace.Print("%s is valid", name) // t //
				t := fmt.Sprintf("  %s (%d)", name, v)
				key.v = name
				key.isActiveFilter = true
				key.iLine = uint16(i)
				line.iKey = uint16(i)
				line.text = t
				line.normal = styleNormal.EchoStyle(t, this.w.nCols)
				line.selected = styleSelected.EchoStyle(t, this.w.nCols)
			} else {
				key.v = name
// trace.Print("%s is not valid", name) // t //
				key.iLine = uint16(i)
				line.iKey = uint16(i)
				name = "  " + name
				line.text = name
				line.normal = styleNormal.EchoStyle(name, this.w.nCols)
				line.selected = styleSelected.EchoStyle(name, this.w.nCols)
			}
		}
	} else {
		for i, name := range *allList {
			key := &this.keys.sl[i]
			line := &this.w.lines.sl[i]
			if v, ok := (*mapPresent)[name]; ok {
				key.v = name
// trace.Print("%s is present", name) // t //
				t := fmt.Sprintf("  %s (%d)", name, v)
				key.iLine = uint16(i)
				line.iKey = uint16(i)
				line.text = t
				line.normal = styleNormal.EchoStyle(t, this.w.nCols)
				line.selected = styleSelected.EchoStyle(t, this.w.nCols)
			} else {
				key.v = name
// trace.Print("%s is not present", name) // t //
				key.iLine = uint16(i)
				line.iKey = uint16(i)
				name = "  " + name
				line.text = name
				line.normal = styleNormalLight.EchoStyle(name, this.w.nCols)
				line.selected = styleSelectedLight.EchoStyle(name, this.w.nCols)
			}
		}
	}
	this.AfterLoad()
	this.w.SetEmptyLine()
	// this.PrepareView()
// trace.End() // t //
}

func (this *mContainer) FilterAndPrint(toggleIsTrue bool) {
	this.FilterAndPrint__new()
// trace.Begin_n(this, "FilterAndPrint") // t //
	var dbIndexes tDbIndexes
	var allContainersWithFilters []*mContainer
	var monitoring bool = true
	browser := this.pBrowserWsp
	this.w.PrintStatusLine()
	browser.View.Reset()
	///	Set all filters
	switch this {
	case &browser.GenresCon:
		/// 1° level filter
// trace.Print("1° level filter") // t //
		this.keys.iActive = UNSET
		// browser.Artists.filterValues = nil
		// browser.Albums.filterValues = nil
		this.SetValidValues()
		if len(this.validValues) == 0 {
			/// There are no filter values
		} else {
			/// There are filter values
			allContainersWithFilters = append(allContainersWithFilters, this)
			a := this.GetFieldValuesForFiltered(GetAudiofArtist)
			if a == nil {
				browser.ArtistsCon.LoadFromSliceAndPrint(&browser.ArtistsCon.allValues)
			} else {
				browser.ArtistsCon.LoadFromSliceAndPrint(a)
			}
		}
	case &browser.ArtistsCon:
		/// 2° level filter
// trace.Print("2° level filter") // t //
		// browser.Albums.filterValues = nil
		var containers []*mContainer
		containers = append(containers, &browser.GenresCon, &browser.ArtistsCon)
		for _, container := range containers {
			container.SetValidValues()
			if len(container.validValues) > 0 {
				allContainersWithFilters = append(allContainersWithFilters, container)
			}
		}
	case &browser.AlbumsCon, &browser.AlbumArtistsCon, &browser.ComposersCon, &browser.YearsCon:
		/// 3° level filter
// trace.Print("3° level filter") // t //
		var containers []*mContainer
		containers = append(containers, &browser.GenresCon, &browser.ArtistsCon, &browser.AlbumsCon, &browser.AlbumArtistsCon, &browser.ComposersCon, &browser.YearsCon)
		for _, container := range containers {
			container.SetValidValues()
			if len(container.validValues) > 0 {
				allContainersWithFilters = append(allContainersWithFilters, container)
			}
		}
	}
// trace.Print("allFilters: %s", GetLongNames(allContainersWithFilters)) // t //
	///	Apply filters to indexes; load and prepare the View
	// trace.Print("this.pBrowserWsp.dbIndexes:  %+v", this.pBrowserWsp.dbIndexes)
	if len(allContainersWithFilters) == 0 {
		/// Load all audiofiles
		browser.View.LoadAndPrepare(&this.pBrowserWsp.wspDbIndexes)
		if monitoring {
			Monitor.AddLine("no filters")
			Monitor.UpdateFileWithContent()
		}
	} else {
		if monitoring {
			for _, container := range allContainersWithFilters {
				Monitor.AddLine("-- " + container.name)
				for _, v := range container.validValues {
					Monitor.AddLine(v)
				}
				container.keys.TraceActiveFilters()
			}
			Monitor.UpdateFileWithContent()
		}
		/// Apply all filters
		for _, i := range this.pBrowserWsp.wspDbIndexes.sl {
			audiof := &AF.sl[i]
			// trace.Print("audiof is %d %s", i, audiof.title)
			isValid := false
			for _, container := range allContainersWithFilters {
				if container.ApplyFilterOR(audiof) {
					isValid = true
				} else {
					isValid = false
					break
				}
			}
			if isValid {
// trace.Print("valid audiof: %s", audiof.title) // t //
				dbIndexes.sl = append(dbIndexes.sl, uint16(i))
			}

		}
		// browser.View.dbIndexes.CopyFrom(&dbIndexes)
		browser.View.LoadAndPrepare(&dbIndexes)
	}
	browser.View.w.PrintWindowAnyway()
	/// Update containers
	switch this {
	case &browser.GenresCon:
		browser.ArtistsCon.LoadFromIndexesAndPrint(&dbIndexes, toggleIsTrue, true)
		browser.AlbumsCon.LoadFromIndexesAndPrint(&dbIndexes, toggleIsTrue, true)
		browser.AlbumArtistsCon.LoadFromIndexesAndPrint(&dbIndexes, toggleIsTrue, false)
		browser.ComposersCon.LoadFromIndexesAndPrint(&dbIndexes, toggleIsTrue, false)
		browser.YearsCon.LoadFromIndexesAndPrint(&dbIndexes, toggleIsTrue, false)
	case &browser.ArtistsCon:
		browser.AlbumsCon.LoadFromIndexesAndPrint(&dbIndexes, toggleIsTrue, true)
		browser.AlbumArtistsCon.LoadFromIndexesAndPrint(&dbIndexes, toggleIsTrue, false)
		browser.ComposersCon.LoadFromIndexesAndPrint(&dbIndexes, toggleIsTrue, false)
		browser.YearsCon.LoadFromIndexesAndPrint(&dbIndexes, toggleIsTrue, false)
	case &browser.AlbumsCon:
	case &browser.AlbumArtistsCon:
	case &browser.ComposersCon:
	}
	///
	Term.Flush()
// trace.End() // t //
}

func (this *mContainer) LoadAllValues(doReactivateFilters bool, doPrint bool) {
	this.Reset()
	this.keys.sl = nil
	this.w.lines.sl = nil
	for _, value := range this.allValues {
		this.AppendKey(value)
	}
	this.PrepareView()
	if doPrint {
		this.w.PrintWindowAnyway()
	}
}

/*
If dbIndexes.sl is nil, loads all values
*/
func (this *mContainer) LoadFromIndexesAndPrint(dbIndexes *tDbIndexes, doReactivateFilters bool, doPrint bool) {
// trace.Begin_n(this, "LoadFromIndexesAndPrint") // t //
// trace.Print("indexes: %v", *dbIndexes) // t //
	this.keys.iActive = UNSET
	this.Reset()
	if dbIndexes.sl == nil {
		for _, value := range this.allValues {
			this.AppendKey(value)
		}
	} else {
		var mapA map[string]uint16
		mapA = make(map[string]uint16)
		for _, i := range dbIndexes.sl {
			audiof := &AF.sl[i]
			mapA[this._fGetAudiofValue(audiof)]++
		}
		// trace.Print("map albums: %v", mapAlbums)
		var aList []string
		for aName := range mapA {
			aList = append(aList, aName)
		}
		sort.Strings(aList)
		for _, name := range aList {
			this.AppendKey(name)
		}
	}
	// if reactivateFilters {
	// 	this.ReactivateFilters()
	// } else {
	// 	this.filterValues = nil
	// }
	this.PrepareView()
	if doPrint {
		this.w.PrintWindowAnyway()
	}
// trace.End() // t //
}

func (this *mContainer) LoadFromSliceAndPrint(list *[]string) {
// trace.Begin_n(this, "LoadFromSliceAndPrint") // t //
	this.keys.iActive = UNSET
	this.Reset()
	for _, name := range *list {
		this.AppendKey(name)
	}
// trace.Print("filterValues %s", this.validValues) // t //
	for _, value := range this.validValues {
		i := this.keys.GetIndexByFilename(value)
		this.keys.Get(i).isActiveFilter = true
	}
	this.PrepareView()
	this.w.PrintWindowAnyway()
// trace.End() // t //
}

/* Returns true is this.filterValues contains the audiof value. */
func (this *mContainer) ApplyFilterOR(audiof *aAudiofile) bool {
	b := slices.Contains(this.validValues, this._fGetAudiofValue(audiof))
	// trace.N_BeginEnd(this, "ApplyFilterOr, %v", b)
	return b
}

/* Sets this.filterValues. Returns number of values. */
func (this *mContainer) SetValidValues() uint16 {
// trace.BeginSilent_n(this, "SetValidValues") // t //
	this.validValues = nil
	for i := range this.keys.sl {
		k := &this.keys.sl[i]
		if k.isActiveFilter {
			this.validValues = append(this.validValues, k.v)
		}
	}
// trace.AddToPrint("all valid values: %#v", this.validValues) // t //
// trace.End() // t //
	return uint16(len(this.validValues))
}

// func (this *mContainer) AppendFilterValue(value string) {
// 	trace.AddToPrint(value)
// 	trace.BeginEnd_n(this, "AppendFilterValue")
// 	this.filterValues = append(this.filterValues, value)
// }

func (this *mContainer) PrepareView() {
	w := &this.w
	// trace.BeginSilent_n(this, "PrepareView")
// trace.Begin_n(this, "PrepareView") // t //
	w.lines.Make(int(this.keys.Len()))
	// trace.N_Begin(this, "PrepareLines")
	for i := range this.keys.Len() {
		key := &this.keys.sl[i]
		v := "  " + key.v
		// v = " " + v

		//!!
		// if this == &this.pBrowserWsp.AlbumsCon {
		// 	AdjustWidth(&v, 60)
		// 	v = fmt.Sprintf("%-30s  %s", v, this.pBrowserWsp.AlbumsCon.sa1.v[i])
		// } else if this == &this.pBrowserWsp.YearsCon {
		// 	v = fmt.Sprintf("%-20s   %d", v, key.i)
		// }

		AdjustWidth(&v, w.nCols-2)
		w.SetLine(uint16(i), v, uint16(i), 0)
	}
	//	Empty line
	w.SetEmptyLine()
	this.AfterLoad()
// trace.EndAdd("%d lines", w.lines.Len()) // t //
}

func (this *mContainer) UseFoundItems_o() {
// trace.Print("%s use selected: %s, %s", this.GetLongName(), Find.GetSelected(), FindDialog.focusKeys.activeKey) // t //
	iKey := this.keys.GetIndexByFilename(Find.GetSelected())
// trace.Print("iLine: %d", iKey) // t //
	this.ActivateFilter(iKey)
	this.keys.iSelected = iKey
}

func (this *mContainer) GetKeysFromDB() {
	var sl []string
	AF.GetFieldSortedUnique(&sl, this.fKeyFilter)
	// ...
}

// Called by an AFList when in single-filter mode. Put the contents in ‘sl’.
func (this *mContainer) GetFileListFromDB(sl *[]string) {
	AF.Query(sl, this.fFileListFilter)
}

// { Artists >> load to Playlist }
func (this *mContainer) LoadTo(that *mFilelist) {
// trace.BeginEndAdd_n(this, "LoadTo", "%s", that.name) //@1 // t //
	this.keys.iActive = this.keys.iSelected
	this.w.UpdateActiveLines()
	this.w.pLayer.pSection.PrintFocusButtons(true)
// trace.Print("%d %d", this.keys.iActive, this.keys.iSelected) // t //
	// Screen.FileContainers.actLayer = this.w.layer
	switch this {
	case &this.pBrowserWsp.GenresCon:
		that.LoadGenre(this.keys.GetSelected().v)
	// case &PFLists.mContainer:
	// 	that.LoadPFList(this.keys.GetSelected().v)
	case &this.pBrowserWsp.ArtistsCon:
		that.LoadArtist()
	case &this.pBrowserWsp.AlbumsCon:
	case &this.pBrowserWsp.YearsCon:
		that.LoadYear(this.pBrowserWsp.YearsCon.keys.GetSelected().v)
	}
	// that.PrepareLines()
	that.w.PrintWindow()
	// Folded.Load()
}

type mContainer struct {
	mBasic
	outer            miContainerOuter
	validValues      []string
	allValues        []string // all values for the browser dbIndexes
	sa1              cStringArray
	sa2              cStringArray
	fKeyFilter       f_FieldFilter
	fFileListFilter  f_FilelistFilter
	_fGetAudiofValue func(*aAudiofile) string
	// filters          cFilters
}

