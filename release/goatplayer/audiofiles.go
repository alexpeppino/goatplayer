package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"

	taglib "github.com/wtolson/go-taglib"
)

var mutex1 sync.Mutex

var nNew, nMissing uint16

func (this *aAudiofiles) CheckDb() {
// trace.Begin("AF.CheckDb") // t //
	var ss string
	reader := bufio.NewReader(os.Stdin)
	filePath := FILE_AUDIOFILES
	/// Check if FILE_AUDIOFILES exists
	fmt.Println()

	if _, err := os.Stat(filePath); err == nil {
		/// AudifilesData exists
		fmt.Println("Database files already exist.")
		fmt.Println("Loading data from files...")
		this.LoadDataFromFile()
		this.PrepareMapFilenameIndex()

		fmt.Println("Comparing database files with filesystem...")
		this.Compare()
		/// Check for New files from the FS
		ll := len(this.notInFile)
		nNew = uint16(ll)
		if ll == 0 {
			fmt.Println("No new files found.")
		} else {
			fmt.Println("New files found:")
			for _, filename := range this.notInFile {
				fmt.Println(filename)
			}
			fmt.Printf("There are %d new audiofiles. Add them to the database? ", ll)
			ss = ReadLineFromConsole()
			fmt.Printf("hai scritto: %s\n", ss)
			fmt.Printf("adding new files...\n")
			/// Add filenames to the file
			for _, filename := range this.notInFile {
				// fmt.Printf("len before: %d\n", len(this.sl))
				this.AddAudiofileFromFileSystem(-1, filename)
				// fmt.Printf("new len: %d\n", len(this.sl))
				// trace.Print("new audiof: %+v", this.sl[len(this.sl)-1])
			}
			this.SortAudiofiles(SortAudiofilesAATT)
			this.SaveDataToFile()
		}

		/// Check for missing files
		fmt.Println("Searching for missing files...")
		ll = len(this.notInFS)
		nMissing = uint16(ll)
		if ll == 0 {
			/// no missing files
			fmt.Println("No missing files found.")
		} else {
			/// there are missing files
			fmt.Println("Missing files found, delete from the database...")
			for _, filename := range this.notInFS {
				fmt.Println("deleting file:", filename)
				fmt.Println("len before:", len(this.sl))
				this.DeleteAudiofile(filename)
				fmt.Println("len after:", len(this.sl))
			}
			fmt.Println("Saving data to file...")
			this.SaveDataToFile()
		}
		this.PrepareMapFilenameIndex()
		this.PrepareAllDbIndexes()

	} else if os.IsNotExist(err) {
		/// File does not exist. Create the file.
		fmt.Println("Database files do not exist. Searching the filesystem...")
		this.LoadDataFromFilesystem()
		this.PrepareMapFilenameIndex()
		this.PrepareAllDbIndexes()
		fmt.Println("Saving data to file...")
		this.SaveDataToFile()
		/// No need to compare. return
		// ss, _ = reader.ReadString('\n')
	} else {
		fmt.Println("Error checking file:", err)
	}

	/// { At this point, AF.sl and audiofiles.txt are ready. }

	///

	fmt.Printf("Press Enter to continue.\n")
	ss, _ = reader.ReadString('\n')
// trace.End() // t //
}

/* Saves to file the audiofiles database. Order of fields is important. */
func (this *aAudiofiles) SaveDataToFile() {
// trace.Begin("AF.SaveDataToFile") // t //
	fd, _ := os.OpenFile(FILE_AUDIOFILES, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	defer fd.Close()
	ss := "" +
		"artist\t" + // 00
		"album\t" + // 01
		"iTrack\t" + // 02
		"title\t" + // 03
		"genre\t" + // 04
		"filename\t" + // 05
		"year\t" + // 06
		"modTime\t" + // 07
		"nTrack\t" + // 08
		"albumArtist\t" + // 09
		"composer\t" + // 10
		"comment\t" + // 11
		"channels\t" + // 12
		"sampleRate\t" + // 13
		"duration\t" + // 14
		"bitRate" // 15
	fmt.Fprintf(fd, "%s\n", ss)
	for i := 0; i < len(this.sl); i++ {
		audiof := &this.sl[i]
		fmt.Fprintf(fd, "%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%d\t%s\n",
			audiof.artist,      // 00
			audiof.album,       // 01
			audiof.iTrack,      // 02
			audiof.title,       // 03
			audiof.genre,       // 04
			audiof.filename,    // 05
			audiof.year,        // 06
			audiof.modTime,     // 07
			audiof.nTrack,      // 08
			audiof.albumArtist, // 09
			audiof.composer,    // 10
			audiof.comment,     // 11
			audiof.channels,    // 12
			audiof.sampleRate,  // 13
			audiof.duration,    // 14
			audiof.bitRate)     // 15
	}
// trace.EndAdd("AF.SaveDataToFile") // t //
}

/*  */
func (this *aAudiofiles) LoadDataFromFile() {
// trace.BeginSilent("AF.LoadDataFromFile") // t //
	fd, _ := os.OpenFile(FILE_AUDIOFILES, os.O_RDONLY, 0664)
	defer fd.Close()
	AF.sl = nil
	var audiof aAudiofile
	sc := bufio.NewScanner(fd)
	var dbIndex uint16
	sc.Scan()
	for sc.Scan() {
		//	For each line of the file
		// trace.Print(sc.Text())
		fields := strings.Split(sc.Text(), "\t")
		audiof.artist = fields[0]
		audiof.album = fields[1]
		audiof.iTrack = utils.str2uint(fields[2])
		audiof.title = fields[3]
		audiof.genre = fields[4]
		audiof.filename = fields[5]
		audiof.year = fields[6]
		audiof.modTime = fields[7]
		audiof.nTrack = utils.str2uint(fields[8])
		audiof.albumArtist = fields[9]
		audiof.composer = fields[10]
		audiof.comment = fields[11]
		audiof.channels = fields[12]
		audiof.sampleRate = fields[13]
		audiof.duration = utils.str2uint32(fields[14])
		audiof.bitRate = fields[15]
		audiof.dbIndex = dbIndex
		dbIndex++
		AF.sl = append(AF.sl, audiof)
	}
	// this.PrepareFilenameMap()
	sc.Err()
	// le := len(this.sl)
	// this.PrepareAllDbIndexes()
	// this.relIndexDbSl = make([]uint16, le)
	// this.Sort(SortAudiofilesAATT)
	// this.PrepareContainers(&BrowserWsp1, &this.allDbIndexes)
	// this.PrepareContainers(&BrowserWsp2, &this.allDbIndexes)
// trace.Print("read %d audiofiles", len(AF.sl)) // t //
	// trace2.Print("%+v", AF.sl[0])
	// trace2.Print("%+v", AF.sl[200])
	// trace2.Print("%+v", AF.sl[400])
	// trace2.Print("%+v", AF.sl[600])
	// trace2.Print("%+v", AF.sl[800])
// trace.End() // t //
}

/*
Sort the database lines in this order: Artist, Album, Track number, Title
*/
func SortAudiofilesAATT(i, j int) bool {
	audiof_i := &AF.sl[i]
	audiof_j := &AF.sl[j]
	if audiof_i.artist != audiof_j.artist {
		return audiof_j.artist > audiof_i.artist
	} else {
		if audiof_i.album != audiof_j.album {
			return audiof_j.album > audiof_i.album
		} else {
			if audiof_i.iTrack != audiof_j.iTrack {
				return audiof_j.iTrack > audiof_i.iTrack
			} else {
				return audiof_j.title > audiof_i.title
			}
		}
	}
}

func SortDbIndexesCATT(i, j int) bool {
	audiof_i := &AF.sl[i]
	audiof_j := &AF.sl[j]
	if audiof_i.composer != audiof_j.composer {
		return audiof_j.composer > audiof_i.composer
	} else {
		if audiof_i.album != audiof_j.album {
			return audiof_j.album > audiof_i.album
		} else {
			if audiof_i.iTrack != audiof_j.iTrack {
				return audiof_j.iTrack > audiof_i.iTrack
			} else {
				return audiof_j.title > audiof_i.title
			}
		}
	}
}

/*
Loads audio-file durations from database. If it doesn't find the duration in the database for a file, gets it directly from the audiofile and adds it to the database.
*/
// func SaveAllDurationsToFile_go() {
// 	mapFilenameDuration = make(map[string]uint32)
// 	fd, _ := os.OpenFile("a/db-durations.txt", os.O_CREATE, 0664)
// 	sca := bufio.NewScanner(fd)
// 	/// Read all durarions from the database
// 	trace.GoPrint("SaveAllDurationsToFile >>")
// 	mutex1.Lock()
// 	for sca.Scan() {
// 		fields := strings.Split(sca.Text(), "\t")
// 		// trace.Print("%s %s", fields[0], fields[1])
// 		mapFilenameDuration[fields[0]] = utils.str2uint32(fields[1])
// 	}
// 	sca.Err()
// 	mutex1.Unlock()
// 	trace.GoPrint("GetDurations unlock")
// 	fd, _ = os.OpenFile("a/db-durations.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
// 	defer fd.Close()
// 	/// Add durations to the database.
// 	for key := range AF.mapFilenameIndex {
// 		if _, ok := mapFilenameDuration[key]; !ok {
// 			fd2, err := taglib.Read(key)
// 			if err != nil {
// 			}
// 			u1 := uint32(fd2.Length() / 10000000)
// 			mapFilenameDuration[key] = u1
// 			fmt.Fprintf(fd, "%s\t%d\n", key, u1)
// 		}
// 	}
// 	trace.GoPrint("SaveAllDurationsToFile <<")
// }

/* Get indexrel. */
func (this *aAudiofiles) GetIndexByFilename(filename string) uint16 {
	var i uint16
	var ok bool
	if i, ok = this.mapFilenameIndex[filename]; !ok {
// trace.Error("filename not found") // t //
		return UNSET
	}
	return i
	// // nIndexRel := uint16(len(this.relIndexDbSl))
	// if i >= nIndexRel {
	// 	trace.Error("indexRel index out of range: %d %d", i, nIndexRel)
	// }
	// iIndexRel := this.relIndexDbSl[i]
	// nAudiof := uint16(len(this.sl))
	// if iIndexRel >= nAudiof {
	// 	trace.Error("audiof index out of range: %d/%d", iIndexRel, nAudiof)
	// }
	// // return &this.audiof[this.mapFilenameIndex[filename]]
	// // trace.BeginEnd("AF.GetIndexByFilename, filename: %s, index; ", filename, iIndexRel)
	// return iIndexRel
}

/*
Give a filename, get the pointer to a database line.
*/
func (this *aAudiofiles) GetAudiofByFilename(filename string) *aAudiofile {

	// trace.Print("AF.GetAudiofByFilename, %s", filename)
	iIndexRel := this.GetIndexByFilename(filename)
	nAudiof := uint16(len(this.sl))
	if iIndexRel >= nAudiof {
// trace.Error("audiof index out of range: %d %d", iIndexRel, nAudiof) // t //
	}
	// return &this.audiof[this.mapFilenameIndex[filename]]
	return &this.sl[iIndexRel]
}

/* Get audiofile by database index. */
func (this *aAudiofiles) GetAudiof(i uint16) *aAudiofile {
	return &this.sl[i]
}

/* Get sl index by db index. */
func (this *aAudiofiles) GetIndex(i uint16) uint16 {
	return this.relIndexDbSl[i]
}

/* To be called after load from file. Compares filelist from file with filelist from filesystem. */
func (this *aAudiofiles) Compare() {
// trace.Begin("AF.Compare") // t //
	// var isChanged bool = false
	var filenamesFS map[string]struct{}
	filenamesFS = make(map[string]struct{})
	var n int
	MusicScanner.FindFiles(AF.musicDir)
	for _, filenameFS := range MusicScanner.fileList {
		filenamesFS[filenameFS] = struct{}{}
	}
	/// Files which are not in the filesystem (missing)
	for filenameFile := range this.mapFilenameIndex {
		if _, ok := filenamesFS[filenameFile]; !ok {
// trace.Print("file not found in filesystem: '%s'", filenameFile) // t //
			// Report.AddReportLine("  %s", filenameFile)
			this.notInFS = append(this.notInFS, filenameFile)
			n++
			// isChanged = true
		}
	}
	// Report.AddReportLine("Audiofiles found in the database but not in the filesystem: %d", n)
	n = 0
	/// Diff: files which are not in the db
	for filenameFS := range filenamesFS {
		if _, ok := this.mapFilenameIndex[filenameFS]; !ok {
// trace.Print("file not found in file: '%s'", filenameFS) // t //
			// Report.AddReportLine("  %s", filenameFS)
			this.notInFile = append(this.notInFile, filenameFS)
			n++
			// isChanged = true
		}
	}
	// Report.AddReportLine("Audiofiles found in the filesystem but not in the database: %d", n)

	/// Remove filenames from the file

	// if isChanged {
	// 	this.Sort(sortAATT)
	// }
	// if len(notInFile) > 0 {
	// }
// trace.EndAdd("AF.Compare") // t //
}

/*
Called after loading audiofiles. Prepares the content of all containers for browser, using dbIndexes as list of audiofile indexes.
*/
func (this *aAudiofiles) PrepareContainers(browser *mBrowserWsp, dbIndexes *tDbIndexes) {
// trace.BeginSilent("AF.PrepareContainers") // t //
// trace.AddToPrint("browser: %s", browser.name) // t //

	genres := &browser.GenresCon
	artists := &browser.ArtistsCon
	albums := &browser.AlbumsCon
	albumArtists := &browser.AlbumArtistsCon
	composers := &browser.ComposersCon
	years := &browser.YearsCon
	// comments := &browser.CommentsCon

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

	for _, i := range dbIndexes.sl {
		audiof := &this.sl[i]
		mapGenres[audiof.genre]++
		mapArtists[audiof.artist]++
		mapAlbums[audiof.album]++
		mapAlbumArtists[audiof.albumArtist]++
		mapComposers[audiof.composer]++
		mapComments[audiof.comment]++
		mapYears[audiof.year]++
	}

	// trace.Print("mapGenres: %+v", mapGenres)
	// trace.Print("mapArtists: %+v", mapArtists)
	// trace.Print("mapAlbums: %+v", mapAlbums)
	// trace.Print("mapAlbumArtists: %+v", mapAlbumArtists)
	// trace.Print("mapComposers: %+v", mapComposers)
	// trace.Print("mapComments: %+v", mapComments)
	// trace.Print("mapYears: %+v", mapYears)

	for k := range mapGenres {
		genres.keys.Append(&tTableKey{k, 0, 0, false})
		genres.allValues = append(genres.allValues, k)
	}
	sort.Strings(genres.allValues)
	sort.Slice(genres.keys.sl, func(i, j int) bool {
		return genres.keys.sl[i].v < genres.keys.sl[j].v
	})

	for k := range mapArtists {
		artists.keys.Append(&tTableKey{k, 0, 0, false})
		artists.allValues = append(artists.allValues, k)
	}
	sort.Strings(artists.allValues)
	sort.Slice(artists.keys.sl, func(i, j int) bool {
		return artists.keys.sl[i].v < artists.keys.sl[j].v
	})

	for k := range mapAlbums {
		albums.keys.Append(&tTableKey{k, 0, 0, false})
		albums.allValues = append(albums.allValues, k)
	}
	sort.Strings(albums.allValues)
	sort.Slice(albums.keys.sl, func(i, j int) bool {
		return albums.keys.sl[i].v < albums.keys.sl[j].v
	})

	for k := range mapAlbumArtists {
		albumArtists.keys.Append(&tTableKey{k, 0, 0, false})
		albumArtists.allValues = append(albumArtists.allValues, k)
	}
	sort.Strings(albumArtists.allValues)
	sort.Slice(albumArtists.keys.sl, func(i, j int) bool {
		return albumArtists.keys.sl[i].v < albumArtists.keys.sl[j].v
	})

	for k := range mapComposers {
		composers.keys.Append(&tTableKey{k, 0, 0, false})
		composers.allValues = append(composers.allValues, k)
	}
	sort.Strings(composers.allValues)
	sort.Slice(composers.keys.sl, func(i, j int) bool {
		return composers.keys.sl[i].v < composers.keys.sl[j].v
	})

	// for k := range mapComments {
	// 	comments.keys.Append(&tTableKey{k, 0, 0, false})
	// 	comments.allValues = append(comments.allValues, k)
	// }
	// sort.Strings(comments.allValues)
	// sort.Slice(comments.keys.sl, func(i, j int) bool {
	// 	return comments.keys.sl[i].v < comments.keys.sl[j].v
	// })

	for k := range mapYears {
		years.keys.Append(&tTableKey{k, 0, 0, false})
		years.allValues = append(years.allValues, k)
	}
	sort.Strings(years.allValues)
	sort.Slice(years.keys.sl, func(i, j int) bool {
		return years.keys.sl[i].v < years.keys.sl[j].v
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

func (this *aAudiofiles) DeleteAudiofile(filename string) {
	var i uint16
	i = this.GetIndexByFilename(filename)
// trace.Print("found at index %d - %s", i, this.sl[i].filename) // t //
	this.sl = slices.Delete(this.sl, int(i), int(i)+1)
}

/*
Called while loading audiofiles data.
i: index, -1 => append new item
*/
func (this *aAudiofiles) AddAudiofileFromFileSystem(i int, filename string) {
	var genre, title, artist, album string
	var track_i, track_n int
	if i == -1 {
		this.sl = append(this.sl, aAudiofile{})
		i = len(this.sl) - 1
		// this.allDbIndexes = append(this.allDbIndexes, uint16(i))
	}
	//	Create a database line for the audiofile
	audiof := &this.sl[i]
	audiof.dbIndex = uint16(i)
	// this.allDbIndexes[i] = uint16(i)
	this.mapFilenameIndex[filename] = uint16(i)
	//	File's modification time
	info, _ := os.Stat(filename)
	modTime := info.ModTime()
	audiof.modTime = modTime.Format("2006-01-02 15:04:05")
	//	Metadata
	GetMetadata(filename)
	this.sl[i].dbIndex = uint16(i)
	this.sl[i].filename = filename
	this.sl[i].dir, this.sl[i].file = filepath.Split(filename)
	if MetaTags != nil {
		// trace.Print("AlbumArtist %s", MetaTags.AlbumArtist())
		// trace.Print("Composer %s", MetaTags.Composer())
		// trace.Print("Comment %s", MetaTags.Comment())
		if genre = MetaTags.Genre(); genre == "" {
			genre = UNKNOWN_GENRE
		}
		audiof.genre = genre
		///
		if artist = MetaTags.Artist(); artist == "" {
			artist = UNKNOWN_ARTIST
		}
		audiof.artist = artist
		///
		if album = MetaTags.Album(); album == "" {
			album = UNKNOWN_ALBUM
		}
		audiof.album = album
		///
		if title = MetaTags.Title(); title == "" {
			title = UNKNOWN_TITLE
		}
		audiof.title = title
		///
		audiof.albumArtist = MetaTags.AlbumArtist()
		audiof.composer = MetaTags.Composer()
		audiof.comment = MetaTags.Comment()
		///
		track_i, track_n = MetaTags.Track()
		audiof.iTrack = uint16(track_i)
		audiof.nTrack = uint16(track_n)
		audiof.year = fmt.Sprintf("%d", MetaTags.Year())
		// if audiof.artist == "The Kinks" {
		// 	trace.Print("Year   %d   %s", MetaTags.Year(), audiof.title)
		// }
	} else {
		//	If there is no metadata for the file
		audiof.genre = UNKNOWN_GENRE
		audiof.artist = UNKNOWN_ARTIST
		audiof.album = UNKNOWN_ALBUM
		audiof.title = UNKNOWN_TITLE
		audiof.iTrack = 0
		audiof.nTrack = 0
		audiof.year = ""
	}
	/// taglib
	fd, err := taglib.Read(filename)
	if err == nil {
		if audiof.artist == "Pink Floyd" {
// trace.Print("Year   %d   %s", fd.Year(), audiof.title) // t //
		}
		// audiof.year = fmt.Sprintf("%d", fd.Year())
		audiof.bitRate = utils.int2str(fd.Bitrate())
		audiof.sampleRate = utils.int2str(fd.Samplerate())
		audiof.channels = utils.int2str(fd.Channels())
		audiof.duration = uint32(fd.Length() / 10000000)
	}
}

func (this *aAudiofiles) PrepareMapFilenameIndex() {
	this.mapFilenameIndex = make(map[string]uint16)
	var audiof *aAudiofile
	for i := range this.sl {
		audiof = &this.sl[i]
		// audiof.i = uint16(i) // set the index of the db line
		this.mapFilenameIndex[audiof.filename] = uint16(i)
	}
// trace.BeginEndAdd("PrepareMapFilenameIndex", "len: %d", len(this.mapFilenameIndex)) // t //
}

/*
After doing a scan of all music directories, reads the tags of all the audiofiles found and fills the database (this.audiof). Sorts the results.
*/
func (this *aAudiofiles) LoadDataFromFilesystem() {
// trace.Begin("AF.LoadDataFromFilesystem") // t //
	/// Get from the filesystem the filenames of all the audiofiles
	MusicScanner.FindFiles(AF.musicDir)
	le := uint16(len(MusicScanner.fileList))
	fmt.Println("Audiofiles found:", le)
	this.sl = make([]aAudiofile, le)
	// this.allDbIndexes = make([]uint16, le)
	this.mapFilenameIndex = make(map[string]uint16)
	/// For each filename
	carriageReturn_eraseLine := "\r\x1b[0K"
	for i, filename := range MusicScanner.fileList {
		this.AddAudiofileFromFileSystem(i, filename)
		fmt.Printf("%sAdding audiofile: %d / %d", carriageReturn_eraseLine, i+1, le)
	}
	fmt.Println()
	// this.PrepareAllDbIndexes()
	this.SortAudiofiles(SortAudiofilesAATT)
	// this.PrepareFilenameMap()
	// this.PrepareContainers(&BrowserWsp1, &this.allDbIndexes)
	// this.PrepareContainers(&BrowserWsp2, &this.allDbIndexes)
	// logFile.WriteString("----- database -----\n")
// trace.Print("len(this.audiof): %d", len(this.sl)) //@1 // t //
	// trace.Print("len(this.mapFilenameIndex2): %d", len(this.mapFilenameIndex2)) //@1
// trace.Print("Genres.keys.Len(): %d", BrowserWsp1.GenresCon.keys.Len())   //@1 // t //
// trace.Print("Artists.keys.Len(): %d", BrowserWsp1.ArtistsCon.keys.Len()) //@1 // t //
// trace.Print("Albums.keys.Len(): %d", BrowserWsp1.AlbumsCon.keys.Len())   //@1 // t //
// trace.Print("Genres.keys.Len(): %d", BrowserWsp2.GenresCon.keys.Len())   //@1 // t //
// trace.Print("Artists.keys.Len(): %d", BrowserWsp2.ArtistsCon.keys.Len()) //@1 // t //
// trace.Print("Albums.keys.Len(): %d", BrowserWsp2.AlbumsCon.keys.Len())   //@1 // t //
// trace.End() // t //
}

func (this *aAudiofiles) PrepareAllDbIndexes() {
	le := len(this.sl)
	this.allDbIndexes.sl = make([]uint16, le)
	for i := range le {
		this.allDbIndexes.sl[i] = uint16(i)
	}
// trace.BeginEndAdd("PrepareAllDbIndexes", "len: %d", len(this.allDbIndexes.sl)) // t //
}

/*
Applies the sorting function f to the slice and prepares the dbIndexes.
this.sl => sorted
*/
func (this *aAudiofiles) SortAudiofiles(f f_SortSlice) {
// trace.BeginEnd("AF.Sort") // t //
	sort.Slice(this.sl, f)
	var i uint16
	le := uint16(len(this.sl))
	//	Update the index relation table. indexRel is a short way to find the audiof line inside the newly ordered database according to the real index.
	// indexRel: db index => sl index
	// this.relIndexDbSl = make([]uint16, le)
	for i = range le {
		this.sl[i].dbIndex = i
	}
	// trace2.Print("sl:0, db:%d, sl[0]: %s, db[0]: %s", this.relIndexDbSl[0], AF.sl[0].filename, AF.sl[this.relIndexDbSl[0]].filename)
	// trace2.Print("indexRel: %+v", this.relIndexDbSl)
}

/* func 'compare' selects the database lines, the lisf of filenames found are stored in 'filenames'. */
func (this *aAudiofiles) Query(filenames *[]string, compare f_FilelistFilter) {
	for i := 0; i < len(this.sl); i++ {
		audiof := &this.sl[i]
		if compare(audiof, i) {
			*filenames = append(*filenames, audiof.filename)
		}
	}
	// trace.BeginEnd("DB.Query, found %d audiof", len(*sl))
}

/*
This Query uses an index selection slice to give the old selection and to get the new selection.
func 'compare' selects the database lines.
*/
func (this *aAudiofiles) QueryI(I *[]uint16, compare f_FilelistFilter) {
	var inI *[]uint16
	var outI []uint16
	if I != nil {
		inI = I
	} else {
		inI = &this.filteredIndexes
	}
	for i := 0; i < len(*inI); i++ {
		audiof := this.GetAudiof((*inI)[i])
		if compare(audiof, i) {
			outI = append(outI, uint16(i))
		}
	}
	if I != nil {
		I = &outI
	} else {
		this.filteredIndexes = make([]uint16, len(outI))
		copy(this.filteredIndexes, outI)
	}
	// trace.BeginEnd("DB.Query, found %d audiof", len(*sl))
}

/* indexis are db-indexis. audiofile db is sorted. new selected indexis are returned sorted as the db at that moment. */
func (this *aAudiofiles) QueryRel(dbIndexis *[]uint16, filter f_FilelistFilterUi) {
	var outDbIndexis []uint16
	var iDb uint16
	for _, iDb = range *dbIndexis {
		audiof := this.GetAudiof(iDb)
		if filter(audiof, iDb) {
			outDbIndexis = append(outDbIndexis, uint16(iDb))
		}
	}
	*dbIndexis = outDbIndexis
	// trace.BeginEnd("DB.Query, found %d audiof", len(*sl))
}

/*
Get a sorted unique list of all the values of a field searching the whole database.
func 'fieldValue' must return the value of that field.
*/
func (this *aAudiofiles) GetFieldSortedUnique(results *[]string, fieldValue f_FieldFilter) {
	var fieldMap map[string]int
	fieldMap = make(map[string]int)
	for i := range len(this.sl) {
		audiof := &this.sl[i]
		fieldMap[fieldValue(audiof)]++
	}
	for st := range fieldMap {
		*results = append(*results, st)
	}
	sort.Strings(*results)
	// trace.BeginEnd("DB.GetFieldSortedUnique, %d filenames, found %d audiof", len(*filenames), len(*results))
}

/*
Give the filenames, get results (sorted and unique values of one field, only for the database lines of those filenames).
compare must return the value of that field.
*/
func (this *aAudiofiles) GetFieldByFilenames(filenames *[]string, results *[]string, compare f_FieldFilter) {
	var fieldMap map[string]int
	fieldMap = make(map[string]int)
	for _, filename := range *filenames {
		audiof := this.GetAudiofByFilename(filename)
		fieldMap[compare(audiof)]++
	}
	for st := range fieldMap {
		*results = append(*results, st)
	}
	sort.Strings(*results)
	// trace.BeginEnd("DB.GetFieldSortedUnique, %d filenames, found %d audiof", len(*filenames), len(*results))
}

type aAudiofile struct {
	title       string
	genre       string
	artist      string
	album       string
	year        string
	albumArtist string
	composer    string
	comment     string
	iTrack      uint16
	filename    string
	dbIndex     uint16 // database line index
	modTime     string
	nTrack      uint16
	channels    string
	sampleRate  string
	duration    uint32
	bitRate     string

	track_is string
	track_ns string
	dir      string
	file     string
	volume   uint16
}

type aAudiofiles struct {
	sl                  []aAudiofile
	mapFilenameIndex    map[string]uint16 // Useful to search an index by filename
	selectedIndexes     tDbIndexes
	filteredIndexes     []uint16   // selected db lines
	allDbIndexes        tDbIndexes // the complete set of db lines
	relIndexDbSl        []uint16   // index relation table (db-line index to db-slice index)
	notInFile           []string
	notInFS             []string
	doSortAndSaveAtQuit bool
	musicDir            string
}

