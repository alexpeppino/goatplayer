package main

import (
	AF "goatplayer/internal/audiofiles"
	"goatplayer/internal/table"
	"sort"
)

func (this *mPlayerWsp) InitPlayerWsp(name, publicName string) {
	this.name = name
	// this.outer = this
	// this.mWorkspace.outer = this
	this.pFilelist = &this.Player.mAbsFilelist
	this.publicName = publicName
	MainScreen.mapWorkspaces[this.name] = &this.maWorkspace
	this.Player.pPlayerW = this
	this.Player.pWorkspace = &this.maWorkspace
	this.AudioPropCon.pWorkspace = &this.maWorkspace
	// this.Param.pWorkspace = &this.maWorkspace
	// this.PFListsCon.pWorkspace = &this.maWorkspace
}

// func (this *mEditorWsp) InitEditorWsp(name, publicName string) {
// 	this.name = name
// 	// this.outer = this
// 	// this.mWorkspace.outer = this
// 	this.publicName = publicName
// 	MainScreen.mapWorkspaces[this.name] = &this.maWorkspace
// 	this.Edit.pEditorWsp = this
// 	this.Edit.pWorkspace = &this.maWorkspace
// 	this.Param.pWorkspace = &this.maWorkspace
// 	this.PFLists.pWorkspace = &this.maWorkspace
// 	// this.FileProp.pBrowserWsp = &BrowserWsp1
// }

func (this *mBrowserWsp) InitBrowserWsp(name, publicName string) {
	this.name = name
	this.Dirs.workingDir = AF.MusicDir
	// trace.BeginSilent_n(this, "InitBrowserWsp") // t //
	this.publicName = publicName
	MainScreen.mapWorkspaces[this.name] = &this.maWorkspace
	this.wspDbIndexes.CopyFrom(&AF.AllDbIndexes)

	// trace.AddToPrint("dbIndexes len: %d", len(this.wspDbIndexes.Sl)) // t //
	// trace.AddToPrint("this.dbIndexes:  %+v", this.dbIndexes)

	this.outer = this
	this.maWorkspace.outer = this
	this.pFilelist = &this.View.mAbsFilelist

	this.View.pBrowserWsp = this
	this.GenresCon.pBrowserWsp = this
	this.ArtistsCon.pBrowserWsp = this
	this.AlbumsCon.pBrowserWsp = this
	this.YearsCon.pBrowserWsp = this
	this.AlbumArtistsCon.pBrowserWsp = this
	this.ComposersCon.pBrowserWsp = this
	// this.CommentsCon.pBrowserWsp = this
	this.AudioPropCon.pBrowserWsp = this
	this.Dirs.pBrowserWsp = this

	this.View.pWorkspace = &this.maWorkspace
	this.GenresCon.pWorkspace = &this.maWorkspace
	this.ArtistsCon.pWorkspace = &this.maWorkspace
	this.AlbumsCon.pWorkspace = &this.maWorkspace
	this.YearsCon.pWorkspace = &this.maWorkspace
	this.AlbumArtistsCon.pWorkspace = &this.maWorkspace
	this.ComposersCon.pWorkspace = &this.maWorkspace
	// this.CommentsCon.pWorkspace = &this.maWorkspace
	this.AudioPropCon.pWorkspace = &this.maWorkspace
	this.Dirs.pWorkspace = &this.maWorkspace
	// trace.End() // t //
}

/*
Called after loading audiofiles. Prepares the content of all containers for browser, using dbIndexes as list of audiofile indexes.
*/
func (this *mBrowserWsp) PrepareContainers(dbIndexes *AF.GDbIndexes) {
	// trace.BeginSilent_n(this, "PrepareContainers") // t //

	genres := &this.GenresCon
	artists := &this.ArtistsCon
	albums := &this.AlbumsCon
	albumArtists := &this.AlbumArtistsCon
	composers := &this.ComposersCon
	years := &this.YearsCon
	// comments := &this.CommentsCon

	mapGenres := make(map[string]int)
	mapArtists := make(map[string]int)
	mapAlbums := make(map[string]int)
	mapAlbumArtists := make(map[string]int)
	mapComposers := make(map[string]int)
	mapYears := make(map[string]int)
	mapComments := make(map[string]int)

	genres.Reset()
	artists.Reset()
	albums.Reset()
	albumArtists.Reset()
	composers.Reset()
	// comments.Reset()
	years.Reset()
	genres.allValues = nil
	artists.allValues = nil
	albums.allValues = nil
	albumArtists.allValues = nil
	composers.allValues = nil
	// comments.allValues = nil
	years.allValues = nil

	for _, i := range dbIndexes.Sl {
		audiof := &AF.Sl[i]
		mapGenres[audiof.Genre]++
		mapArtists[audiof.Artist]++
		mapAlbums[audiof.Album]++
		mapAlbumArtists[audiof.AlbumArtist]++
		mapComposers[audiof.Composer]++
		mapComments[audiof.Comment]++
		mapYears[audiof.Year]++
	}

	// trace.Print("mapGenres: %+v", mapGenres)
	// trace.Print("mapArtists: %+v", mapArtists)
	// trace.Print("mapAlbums: %+v", mapAlbums)
	// trace.Print("mapAlbumArtists: %+v", mapAlbumArtists)
	// trace.Print("mapComposers: %+v", mapComposers)
	// trace.Print("mapComments: %+v", mapComments)
	// trace.Print("mapYears: %+v", mapYears)

	for k := range mapGenres {
		genres.keys.Append(&table.TTableKey{V: k, LineI: 0, DbIndex: 0, IsActiveFilter: false})
		genres.allValues = append(genres.allValues, k)
	}
	sort.Strings(genres.allValues)
	sort.Slice(genres.keys.Sl, func(i, j int) bool {
		return genres.keys.Sl[i].V < genres.keys.Sl[j].V
	})

	for k := range mapArtists {
		artists.keys.Append(&table.TTableKey{V: k, LineI: 0, DbIndex: 0, IsActiveFilter: false})
		artists.allValues = append(artists.allValues, k)
	}
	sort.Strings(artists.allValues)
	sort.Slice(artists.keys.Sl, func(i, j int) bool {
		return artists.keys.Sl[i].V < artists.keys.Sl[j].V
	})

	for k := range mapAlbums {
		albums.keys.Append(&table.TTableKey{V: k, LineI: 0, DbIndex: 0, IsActiveFilter: false})
		albums.allValues = append(albums.allValues, k)
	}
	sort.Strings(albums.allValues)
	sort.Slice(albums.keys.Sl, func(i, j int) bool {
		return albums.keys.Sl[i].V < albums.keys.Sl[j].V
	})

	for k := range mapAlbumArtists {
		albumArtists.keys.Append(&table.TTableKey{V: k, LineI: 0, DbIndex: 0, IsActiveFilter: false})
		albumArtists.allValues = append(albumArtists.allValues, k)
	}
	sort.Strings(albumArtists.allValues)
	sort.Slice(albumArtists.keys.Sl, func(i, j int) bool {
		return albumArtists.keys.Sl[i].V < albumArtists.keys.Sl[j].V
	})

	for k := range mapComposers {
		composers.keys.Append(&table.TTableKey{V: k, LineI: 0, DbIndex: 0, IsActiveFilter: false})
		composers.allValues = append(composers.allValues, k)
	}
	sort.Strings(composers.allValues)
	sort.Slice(composers.keys.Sl, func(i, j int) bool {
		return composers.keys.Sl[i].V < composers.keys.Sl[j].V
	})

	// for k := range mapComments {
	// 	comments.keys.Append(&tTableKey{k, 0, 0, false})
	// 	comments.allValues = append(comments.allValues, k)
	// }
	// sort.Strings(comments.allValues)
	// sort.Slice(comments.keys.Sl, func(i, j int) bool {
	// 	return comments.keys.Sl[i].V < comments.keys.Sl[j].v
	// })

	for k := range mapYears {
		years.keys.Append(&table.TTableKey{V: k, LineI: 0, DbIndex: 0, IsActiveFilter: false})
		years.allValues = append(years.allValues, k)
	}
	sort.Strings(years.allValues)
	sort.Slice(years.keys.Sl, func(i, j int) bool {
		return years.keys.Sl[i].V < years.keys.Sl[j].V
	})

	// trace.Print("genres.allValues:  %+v", genres.allValues)
	// trace.Print("artists.allValues:  %+v", artists.allValues)
	// trace.Print("albums.allValues:  %+v", albums.allValues)
	// trace.Print("albumArtists.allValues:  %+v", albumArtists.allValues)
	// trace.Print("composers.allValues:  %+v", composers.allValues)
	// trace.Print("comments.allValues:  %+v", comments.allValues)
	// trace.Print("years.allValues:  %+v", years.allValues)

	// trace.End() // t //
}

func (this *maWorkspace) GetLongName() string { return this.name }

/* short=Wsp */
type maWorkspace struct {
	name          string
	pRightSection *vSection
	sections      []*vSection
	focusKeys     uFocusKeys /* focus-to-model keys */
	outer         miWorkspace
	publicName    string
	AudioPropCon  mAudioProperties
	pFilelist     *mAbsFilelist
}

type mBrowserWsp struct {
	maWorkspace
	View            mFilelist
	GenresCon       mContainer
	ArtistsCon      mContainer
	AlbumsCon       mContainer
	YearsCon        mContainer
	AlbumArtistsCon mContainer
	ComposersCon    mContainer
	// CommentsCon     mContainer
	Dirs         mDirs
	wspDbIndexes AF.GDbIndexes
}

type mPlayerWsp struct {
	maWorkspace
	Player mPlaylist
	// PFListsCon mPFListsCon
	// Param      mParameters
}

// type mEditorWsp struct {
// 	maWorkspace
// 	Edit    mEditor
// 	PFLists mPFLists
// 	// PFListsCon mPFLists
// }

type miWorkspace interface {
	// Dump_o() /*c*/
}
