package main

import "fmt"

/*
There are 2 kinds of "f" functions:
1.  fName  : outer method (for an embedding type) called inside basic.Name()
2.  _fName : defined for each variable and defined in this file
*/
func _f() {
	BrowserWsp1.GenresCon._fGetBrowserString = func() string {
		return fmt.Sprintf("%s music", BrowserWsp1.GenresCon.keys.GetSelected().v)
	}
	BrowserWsp1.ArtistsCon._fGetBrowserString = func() string {
		return BrowserWsp1.ArtistsCon.keys.GetSelected().v
	}
	BrowserWsp1.AlbumsCon._fGetBrowserString = func() string {
		i := BrowserWsp1.AlbumsCon.keys.iSelected
		return fmt.Sprintf("%s %s", BrowserWsp1.AlbumsCon.keys.sl[i].v, BrowserWsp1.AlbumsCon.sa1.Get(i))
	}
	BrowserWsp1.YearsCon._fGetBrowserString = func() string {
		return fmt.Sprintf("%s best songs", BrowserWsp1.YearsCon.keys.GetSelected().v)
	}
	BrowserWsp1.AlbumArtistsCon._fGetBrowserString = func() string {
		return BrowserWsp1.AlbumArtistsCon.keys.GetSelected().v
	}
	BrowserWsp1.ComposersCon._fGetBrowserString = func() string {
		return BrowserWsp1.ComposersCon.keys.GetSelected().v
	}
	// BrowserWsp1.CommentsCon._fGetBrowserString = func() string {
	// 	return BrowserWsp1.CommentsCon.keys.GetSelected().v
	// }

	BrowserWsp2.GenresCon._fGetBrowserString = func() string {
		return fmt.Sprintf("%s music", BrowserWsp2.GenresCon.keys.GetSelected().v)
	}
	BrowserWsp2.ArtistsCon._fGetBrowserString = func() string {
		return BrowserWsp2.ArtistsCon.keys.GetSelected().v
	}
	BrowserWsp2.AlbumsCon._fGetBrowserString = func() string {
		i := BrowserWsp2.AlbumsCon.keys.iSelected
		return fmt.Sprintf("%s %s", BrowserWsp2.AlbumsCon.keys.sl[i].v, BrowserWsp2.AlbumsCon.sa1.Get(i))
	}
	BrowserWsp2.YearsCon._fGetBrowserString = func() string {
		return fmt.Sprintf("%s best songs", BrowserWsp2.YearsCon.keys.GetSelected().v)
	}
	BrowserWsp2.AlbumArtistsCon._fGetBrowserString = func() string {
		return BrowserWsp2.AlbumArtistsCon.keys.GetSelected().v
	}
	BrowserWsp2.ComposersCon._fGetBrowserString = func() string {
		return BrowserWsp2.ComposersCon.keys.GetSelected().v
	}
	// BrowserWsp2.CommentsCon._fGetBrowserString = func() string {
	// 	return BrowserWsp2.CommentsCon.keys.GetSelected().v
	// }
	BrowserWsp1.GenresCon._fGetAudiofValue = GetAudiofGenre
	BrowserWsp1.ArtistsCon._fGetAudiofValue = GetAudiofArtist
	BrowserWsp1.AlbumsCon._fGetAudiofValue = GetAudiofAlbum
	BrowserWsp1.AlbumArtistsCon._fGetAudiofValue = GetAudiofAlbumArtist
	BrowserWsp1.ComposersCon._fGetAudiofValue = GetAudiofComposer
	BrowserWsp1.YearsCon._fGetAudiofValue = GetAudiofYear

	BrowserWsp2.GenresCon._fGetAudiofValue = GetAudiofGenre
	BrowserWsp2.ArtistsCon._fGetAudiofValue = GetAudiofArtist
	BrowserWsp2.AlbumsCon._fGetAudiofValue = GetAudiofAlbum
	BrowserWsp2.AlbumArtistsCon._fGetAudiofValue = GetAudiofAlbumArtist
	BrowserWsp2.ComposersCon._fGetAudiofValue = GetAudiofComposer
	BrowserWsp2.YearsCon._fGetAudiofValue = GetAudiofYear

	// EditorWsp.PFLists._fOpen = func() {
	// 	Edit.LoadList(EditorWsp.PFLists.keys.GetSelected().v)
	// 	Edit.PrepareIndexes()
	// 	Edit.LoadAndPrepare(nil)
	// 	Edit.w.PrintWindowAnyway()
	// }

	// BrowserWsp1.View._fSetFocusToModel = func() {
	// 	Term.AddStringWithPos(54, 0, "BrowserWsp1")
	// 	Term.Flush()
	// }
	// BrowserWsp2.View._fSetFocusToModel = func() {
	// 	Term.AddStringWithPos(54, 0, "BrowserWsp2")
	// 	Term.Flush()
	// }
	// PlayerWsp.ActivePL._fSetFocusToModel = func() {
	// 	Term.AddStringWithPos(54, 0, "PlayerWsp   ")
	// 	Term.Flush()
	// }
}

