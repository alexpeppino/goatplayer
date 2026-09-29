package audiofiles // Single-instance package

import (
	"bufio"
	"fmt"
	"goatplayer/internal/K"
	"goatplayer/internal/input"
	"goatplayer/internal/musicscanner"
	"goatplayer/internal/small/utils"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/wtolson/go-taglib"
)

var (
	Sl                  []TAudiofile
	MepFilenameIndex    map[string]uint16 // Useful to search an index by filename
	SelectedIndexes     GDbIndexes
	FilteredIndexes     []uint16   // selected db lines
	AllDbIndexes        GDbIndexes // the complete set of db lines
	RelIndexDbSl        []uint16   // index relation table (db-line index to db-slice index)
	NotInFile           []string
	NotInFS             []string
	DoSortAndSaveAtQuit bool
	MusicDir            string
)

var NewN, MissingN uint16
var Mutex1 sync.Mutex

type TAudiofile struct {
	Title       string
	Genre       string
	Artist      string
	Album       string
	Year        string
	AlbumArtist string
	Composer    string
	Comment     string
	TrackI      uint16
	Filename    string
	DbIndex     uint16 // database line index
	ModTime     string
	TrackN      uint16
	ChannelsN   string
	SampleRate  string
	Duration    uint32
	BitRate     string

	Track_is string
	Track_ns string
	Dir      string
	File     string
	Volume   uint16
}

func CheckDb() error {
// trace.Begin("CheckDb") // t //
	var ss string
	reader := bufio.NewReader(os.Stdin)
	filePath := K.FILE_AUDIOFILES

	/// Check if FILE_AUDIOFILES exists

	fmt.Println()

	if _, err := os.Stat(filePath); err == nil {

		/// AudifilesData exists

		fmt.Println("Database files already exist.")
		fmt.Println("Loading data from files...")
		if err := LoadDataFromFile(); err != nil {
			return err
		}
		PrepareMepFilenameIndex()

		fmt.Println("Comparing database files with filesystem...")
		Compare()

		/// Check for New files from the FS

		ll := len(NotInFile)
		NewN = uint16(ll)
		if ll == 0 {
			fmt.Println("No new files found.")
		} else {
			fmt.Println("New files found:")
			for _, filename := range NotInFile {
				fmt.Println(filename)
			}
			fmt.Printf("There are %d new audiofiles. Add them to the database? ", ll)
			ss = input.ReadLineFromConsole()
			fmt.Printf("hai scritto: %s\n", ss)
			fmt.Printf("adding new files...\n")

			/// Add filenames to the file

			for _, filename := range NotInFile {
				// fmt.Printf("len before: %d\n", len(sl))
				AddAudiofileFromFileSystem(-1, filename)
				// fmt.Printf("new len: %d\n", len(sl))
				// trace.Print("new audiof: %+v", sl[len(sl)-1])
			}
			SortAudiofiles(SortAudiofilesAATT)
			SaveDataToFile()
		}

		/// Check for missing files

		fmt.Println("Searching for missing files...")
		ll = len(NotInFS)
		MissingN = uint16(ll)
		if ll == 0 {
			/// no missing files
			fmt.Println("No missing files found.")
		} else {

			/// There are missing files

			fmt.Println("Missing files found, delete from the database...")
			for _, filename := range NotInFS {
				fmt.Println("deleting file:", filename)
				fmt.Println("len before:", len(Sl))
				DeleteAudiofile(filename)
				fmt.Println("len after:", len(Sl))
			}
			fmt.Println("Saving data to file...")
			SaveDataToFile()
		}
		PrepareMepFilenameIndex()
		PrepareAllDbIndexes()

	} else if os.IsNotExist(err) {

		/// File does not exist. Create the file.

		fmt.Println("Database files do not exist. Searching the filesystem...")
		LoadDataFromFilesystem()
		PrepareMepFilenameIndex()
		PrepareAllDbIndexes()
		fmt.Println("Saving data to file...")
		SaveDataToFile()

		/// Save data, reset, load data from file (this is done to check now the integrity of data)
		fmt.Println("Checking integrity of data files...")
		Sl = nil
		if err := LoadDataFromFile(); err != nil {
			return err
		}
		PrepareMepFilenameIndex()
		PrepareAllDbIndexes()

		/// No need to compare. return
		// ss, _ = reader.ReadString('\n')
	} else {
		fmt.Println("Error checking file:", err)
	}

	/// { At this point, sl and audiofiles.txt are ready. }

	///

	fmt.Printf("Press Enter to continue.\n")
	ss, _ = reader.ReadString('\n')
// trace.End() // t //
	return nil
}

/* Saves to file the audiofiles database. Order of fields is important. */
func SaveDataToFile() {
// trace.Begin("SaveDataToFile") // t //
	fd, _ := os.OpenFile(K.FILE_AUDIOFILES, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, K.PERM_OPEN_FILE)
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
	for i := 0; i < len(Sl); i++ {
		audiof := &Sl[i]
		fmt.Fprintf(fd, "%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%d\t%s\n",
			audiof.Artist,      // 00
			audiof.Album,       // 01
			audiof.TrackI,      // 02
			audiof.Title,       // 03
			audiof.Genre,       // 04
			audiof.Filename,    // 05
			audiof.Year,        // 06
			audiof.ModTime,     // 07
			audiof.TrackN,      // 08
			audiof.AlbumArtist, // 09
			audiof.Composer,    // 10
			audiof.Comment,     // 11
			audiof.ChannelsN,   // 12
			audiof.SampleRate,  // 13
			audiof.Duration,    // 14
			audiof.BitRate)     // 15
	}
// trace.EndAdd("SaveDataToFile") // t //
}

/*  */
func LoadDataFromFile() error {
// trace.BeginSilent("LoadDataFromFile") // t //
	fd, _ := os.OpenFile(K.FILE_AUDIOFILES, os.O_RDONLY, K.PERM_OPEN_FILE)
	defer fd.Close()
	Sl = nil
	var audiof TAudiofile
	sc := bufio.NewScanner(fd)
	var dbIndex uint16
	var iLine uint16 = 1
	sc.Scan() // captions
	for sc.Scan() {
		iLine++
		/// For each line of the file
		// trace.Print(sc.Text())
		str1 := sc.Text()
		fields := strings.Split(str1, "\t")
		if len(fields) < 16 {
			return fmt.Errorf("Error: database line #%d contains %d fields and not 16. Line:\n%s", iLine, len(fields), str1)
		}
		audiof.Artist = fields[0]
		audiof.Album = fields[1]
		audiof.TrackI = utils.Str2uint(fields[2])
		audiof.Title = fields[3]
		audiof.Genre = fields[4]
		audiof.Filename = fields[5]
		audiof.Year = fields[6]
		audiof.ModTime = fields[7]
		audiof.TrackN = utils.Str2uint(fields[8])
		audiof.AlbumArtist = fields[9]
		audiof.Composer = fields[10]
		audiof.Comment = fields[11]
		audiof.ChannelsN = fields[12]
		audiof.SampleRate = fields[13]
		audiof.Duration = utils.Str2uint32(fields[14])
		audiof.BitRate = fields[15]
		audiof.DbIndex = dbIndex
		dbIndex++
		Sl = append(Sl, audiof)
	}
	// PrepareFilenameMap()
	sc.Err()
	// le := len(sl)
	// PrepareAllDbIndexes()
	// relIndexDbSl = make([]uint16, le)
	// Sort(SortAudiofilesAATT)
	// PrepareContainers(&BrowserWsp1, &allDbIndexes)
	// PrepareContainers(&BrowserWsp2, &allDbIndexes)
// trace.Print("read %d audiofiles", len(Sl)) // t //
	// trace2.Print("%+v", sl[0])
	// trace2.Print("%+v", sl[200])
	// trace2.Print("%+v", sl[400])
	// trace2.Print("%+v", sl[600])
	// trace2.Print("%+v", sl[800])
// trace.End() // t //
	return nil
}

/*
Sort the database lines in this order: Artist, Album, Track number, Title
*/
func SortAudiofilesAATT(i, j int) bool {
	audiof_i := &Sl[i]
	audiof_j := &Sl[j]
	if audiof_i.Artist != audiof_j.Artist {
		return audiof_j.Artist > audiof_i.Artist
	} else {
		if audiof_i.Album != audiof_j.Album {
			return audiof_j.Album > audiof_i.Album
		} else {
			if audiof_i.TrackI != audiof_j.TrackI {
				return audiof_j.TrackI > audiof_i.TrackI
			} else {
				return audiof_j.Title > audiof_i.Title
			}
		}
	}
}

func SortDbIndexesCATT(i, j int) bool {
	audiof_i := &Sl[i]
	audiof_j := &Sl[j]
	if audiof_i.Composer != audiof_j.Composer {
		return audiof_j.Composer > audiof_i.Composer
	} else {
		if audiof_i.Album != audiof_j.Album {
			return audiof_j.Album > audiof_i.Album
		} else {
			if audiof_i.TrackI != audiof_j.TrackI {
				return audiof_j.TrackI > audiof_i.TrackI
			} else {
				return audiof_j.Title > audiof_i.Title
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
// 	for key := range mapFilenameIndex {
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
func GetIndexByFilename(filename string) uint16 {
	var i uint16
	var ok bool
	if i, ok = MepFilenameIndex[filename]; !ok {
// trace.Error("filename not found") // t //
		return K.UNSET
	}
	return i
	// // nIndexRel := uint16(len(relIndexDbSl))
	// if i >= nIndexRel {
	// 	trace.Error("indexRel index out of range: %d %d", i, nIndexRel)
	// }
	// iIndexRel := relIndexDbSl[i]
	// nAudiof := uint16(len(sl))
	// if iIndexRel >= nAudiof {
	// 	trace.Error("audiof index out of range: %d/%d", iIndexRel, nAudiof)
	// }
	// // return &audiof[mapFilenameIndex[filename]]
	// // trace.BeginEnd("GetIndexByFilename, filename: %s, index; ", filename, iIndexRel)
	// return iIndexRel
}

/*
Give a filename, get the pointer to a database line.
*/
func GetAudiofByFilename(filename string) *TAudiofile {

	// trace.Print("GetAudiofByFilename, %s", filename)
	iIndexRel := GetIndexByFilename(filename)
	nAudiof := uint16(len(Sl))
	if iIndexRel >= nAudiof {
// trace.Error("audiof index out of range: %d %d", iIndexRel, nAudiof) // t //
	}
	// return &audiof[mapFilenameIndex[filename]]
	return &Sl[iIndexRel]
}

/* Get audiofile by database index. */
func GetAudiof(i uint16) *TAudiofile {
	return &Sl[i]
}

/* Get sl index by db index. */
func GetIndex(i uint16) uint16 {
	return RelIndexDbSl[i]
}

/* To be called after load from file. Compares filelist from file with filelist from filesystem. */
func Compare() {
// trace.Begin("Compare") // t //
	// var isChanged bool = false
	var filenamesFS map[string]struct{}
	filenamesFS = make(map[string]struct{})
	var n int
	musicscanner.FindFiles(MusicDir)
	for _, filenameFS := range musicscanner.FileList {
		filenamesFS[filenameFS] = struct{}{}
	}
	/// Files which are not in the filesystem (missing)
	for filenameFile := range MepFilenameIndex {
		if _, ok := filenamesFS[filenameFile]; !ok {
// trace.Print("file not found in filesystem: '%s'", filenameFile) // t //
			// Report.AddReportLine("  %s", filenameFile)
			NotInFS = append(NotInFS, filenameFile)
			n++
			// isChanged = true
		}
	}
	// Report.AddReportLine("Audiofiles found in the database but not in the filesystem: %d", n)
	n = 0
	/// Diff: files which are not in the db
	for filenameFS := range filenamesFS {
		if _, ok := MepFilenameIndex[filenameFS]; !ok {
// trace.Print("file not found in file: '%s'", filenameFS) // t //
			// Report.AddReportLine("  %s", filenameFS)
			NotInFile = append(NotInFile, filenameFS)
			n++
			// isChanged = true
		}
	}
	// Report.AddReportLine("Audiofiles found in the filesystem but not in the database: %d", n)

	/// Remove filenames from the file

	// if isChanged {
	// 	Sort(sortAATT)
	// }
	// if len(notInFile) > 0 {
	// }
// trace.EndAdd("Compare") // t //
}

func DeleteAudiofile(filename string) {
	var i uint16
	i = GetIndexByFilename(filename)
// trace.Print("found at index %d - %s", i, Sl[i].Filename) // t //
	Sl = slices.Delete(Sl, int(i), int(i)+1)
}

/*
Called while loading audiofiles data.
i: index, -1 => append new item
*/
func AddAudiofileFromFileSystem(i int, filename string) {
	var genre, title, artist, album string
	var track_i, track_n int
	if i == -1 {
		Sl = append(Sl, TAudiofile{})
		i = len(Sl) - 1
		// allDbIndexes = append(allDbIndexes, uint16(i))
	}
	//	Create a database line for the audiofile
	audiof := &Sl[i]
	audiof.DbIndex = uint16(i)
	// allDbIndexes[i] = uint16(i)
	MepFilenameIndex[filename] = uint16(i)
	//	File's modification time
	info, _ := os.Stat(filename)
	modTime := info.ModTime()
	audiof.ModTime = modTime.Format("2006-01-02 15:04:05")
	//	Metadata
	musicscanner.GetMetadata(filename)
	Sl[i].DbIndex = uint16(i)
	Sl[i].Filename = filename
	Sl[i].Dir, Sl[i].File = filepath.Split(filename)
	if musicscanner.MetaTags != nil {
		// trace.Print("AlbumArtist %s", MetaTags.AlbumArtist())
		// trace.Print("Composer %s", MetaTags.Composer())
		// trace.Print("Comment %s", MetaTags.Comment())
		if genre = musicscanner.MetaTags.Genre(); genre == "" {
			genre = K.UNKNOWN_GENRE
		}
		audiof.Genre = genre
		///
		if artist = musicscanner.MetaTags.Artist(); artist == "" {
			artist = K.UNKNOWN_ARTIST
		}
		audiof.Artist = artist
		///
		if album = musicscanner.MetaTags.Album(); album == "" {
			album = K.UNKNOWN_ALBUM
		}
		audiof.Album = album
		///
		if title = musicscanner.MetaTags.Title(); title == "" {
			title = K.UNKNOWN_TITLE
		}
		audiof.Title = title
		///
		audiof.AlbumArtist = musicscanner.MetaTags.AlbumArtist()
		audiof.Composer = musicscanner.MetaTags.Composer()
		audiof.Comment = musicscanner.MetaTags.Comment()
		///
		track_i, track_n = musicscanner.MetaTags.Track()
		audiof.TrackI = uint16(track_i)
		audiof.TrackN = uint16(track_n)
		audiof.Year = fmt.Sprintf("%d", musicscanner.MetaTags.Year())
		// if audiof.Artist == "The Kinks" {
		// 	trace.Print("Year   %d   %s", MetaTags.Year(), audiof.Title)
		// }
	} else {
		//	If there is no metadata for the file
		audiof.Genre = K.UNKNOWN_GENRE
		audiof.Artist = K.UNKNOWN_ARTIST
		audiof.Album = K.UNKNOWN_ALBUM
		audiof.Title = K.UNKNOWN_TITLE
		audiof.TrackI = 0
		audiof.TrackN = 0
		audiof.Year = ""
	}
	/// taglib
	fd, err := taglib.Read(filename)
	if err == nil {
		if audiof.Artist == "Pink Floyd" {
// trace.Print("Year   %d   %s", fd.Year(), audiof.Title) // t //
		}
		// audiof.Year = fmt.Sprintf("%d", fd.Year())
		audiof.BitRate = utils.Int2str(fd.Bitrate())
		audiof.SampleRate = utils.Int2str(fd.Samplerate())
		audiof.ChannelsN = utils.Int2str(fd.Channels())
		audiof.Duration = uint32(fd.Length() / 10000000)
	}
}

func PrepareMepFilenameIndex() {
	MepFilenameIndex = make(map[string]uint16)
	var audiof *TAudiofile
	for i := range Sl {
		audiof = &Sl[i]
		// audiof.i = uint16(i) // set the index of the db line
		MepFilenameIndex[audiof.Filename] = uint16(i)
	}
// trace.BeginEndAdd("PrepareMapFilenameIndex", "len: %d", len(MepFilenameIndex)) // t //
}

func GetAudiofGenre(audiof *TAudiofile) string       { return audiof.Genre }
func GetAudiofArtist(audiof *TAudiofile) string      { return audiof.Artist }
func GetAudiofAlbum(audiof *TAudiofile) string       { return audiof.Album }
func GetAudiofAlbumArtist(audiof *TAudiofile) string { return audiof.AlbumArtist }
func GetAudiofComposer(audiof *TAudiofile) string    { return audiof.Composer }
func GetAudiofYear(audiof *TAudiofile) string        { return audiof.Year }

/*
After doing a scan of all music directories, reads the tags of all the audiofiles found and fills the database (audiof). Sorts the results.
*/
func LoadDataFromFilesystem() {
// trace.Begin("LoadDataFromFilesystem") // t //
	/// Get from the filesystem the filenames of all the audiofiles
	musicscanner.FindFiles(MusicDir)
	le := uint16(len(musicscanner.FileList))
	fmt.Println("Audiofiles found:", le)
	Sl = make([]TAudiofile, le)
	// allDbIndexes = make([]uint16, le)
	MepFilenameIndex = make(map[string]uint16)
	/// For each filename
	carriageReturn_eraseLine := "\r\x1b[0K"
	for i, filename := range musicscanner.FileList {
		AddAudiofileFromFileSystem(i, filename)
		fmt.Printf("%sAdding audiofile: %d / %d", carriageReturn_eraseLine, i+1, le)
	}
	fmt.Println()
	// PrepareAllDbIndexes()
	SortAudiofiles(SortAudiofilesAATT)
	// PrepareFilenameMap()
	// PrepareContainers(&BrowserWsp1, &allDbIndexes)
	// PrepareContainers(&BrowserWsp2, &allDbIndexes)
	// logFile.WriteString("----- database -----\n")
// trace.Print("len(audiof): %d", len(Sl)) //@1 // t //
	// trace.Print("len(mapFilenameIndex2): %d", len(mapFilenameIndex2)) //@1

	// trace.Print("Genres.keys.Len(): %d", BrowserWsp1.GenresCon.keys.Len())   //@1
	// trace.Print("Artists.keys.Len(): %d", BrowserWsp1.ArtistsCon.keys.Len()) //@1
	// trace.Print("Albums.keys.Len(): %d", BrowserWsp1.AlbumsCon.keys.Len())   //@1
	// trace.Print("Genres.keys.Len(): %d", BrowserWsp2.GenresCon.keys.Len())   //@1
	// trace.Print("Artists.keys.Len(): %d", BrowserWsp2.ArtistsCon.keys.Len()) //@1
	// trace.Print("Albums.keys.Len(): %d", BrowserWsp2.AlbumsCon.keys.Len())   //@1
// trace.End() // t //
}

func PrepareAllDbIndexes() {
	le := len(Sl)
	AllDbIndexes.Sl = make([]uint16, le)
	for i := range le {
		AllDbIndexes.Sl[i] = uint16(i)
	}
// trace.BeginEndAdd("PrepareAllDbIndexes", "len: %d", len(AllDbIndexes.Sl)) // t //
}

/*
Applies the sorting function f to the slice and prepares the dbIndexes.
sl => sorted
*/
func SortAudiofiles(f f_int_int_bool) {
// trace.BeginEnd("Sort") // t //
	sort.Slice(Sl, f)
	var i uint16
	le := uint16(len(Sl))
	//	Update the index relation table. indexRel is a short way to find the audiof line inside the newly ordered database according to the real index.
	// indexRel: db index => sl index
	// relIndexDbSl = make([]uint16, le)
	for i = range le {
		Sl[i].DbIndex = i
	}
	// trace2.Print("sl:0, db:%d, sl[0]: %s, db[0]: %s", relIndexDbSl[0], sl[0].filename, sl[relIndexDbSl[0]].filename)
	// trace2.Print("indexRel: %+v", relIndexDbSl)
}

/* func 'compare' selects the database lines, the lisf of filenames found are stored in 'filenames'. */
func Query(filenames *[]string, compare f_audiof_int_bool) {
	for i := 0; i < len(Sl); i++ {
		audiof := &Sl[i]
		if compare(audiof, i) {
			*filenames = append(*filenames, audiof.Filename)
		}
	}
	// trace.BeginEnd("DB.Query, found %d audiof", len(*sl))
}

/*
Give the filenames, get results (sorted and unique values of one field, only for the database lines of those filenames).
compare must return the value of that field.
*/
func GetFieldByFilenames(filenames *[]string, results *[]string, compare f_audiof_string) {
	var fieldMap map[string]int
	fieldMap = make(map[string]int)
	for _, filename := range *filenames {
		audiof := GetAudiofByFilename(filename)
		fieldMap[compare(audiof)]++
	}
	for st := range fieldMap {
		*results = append(*results, st)
	}
	sort.Strings(*results)
	// trace.BeginEnd("DB.GetFieldSortedUnique, %d filenames, found %d audiof", len(*filenames), len(*results))
}

func MakeFunc_SortCATT(dbIndexes *GDbIndexes) f_int_int_bool {
	return func(i, j int) bool {
		audiof_i := &Sl[dbIndexes.Sl[i]]
		audiof_j := &Sl[dbIndexes.Sl[j]]
		if audiof_i.Composer != audiof_j.Composer {
			return audiof_j.Composer > audiof_i.Composer
		} else {
			if audiof_i.Album != audiof_j.Album {
				return audiof_j.Album > audiof_i.Album
			} else {
				if audiof_i.TrackI != audiof_j.TrackI {
					return audiof_j.TrackI > audiof_i.TrackI
				} else {
					return audiof_j.Title > audiof_i.Title
				}
			}
		}
	}
}

type f_int_int_bool func(i int, j int) bool
type f_audiof_int_bool func(audiof *TAudiofile, i int) bool
type f_audiof_uint_bool func(audiof *TAudiofile, i uint16) bool
type f_audiof_string func(audiof *TAudiofile) string
type f_audiof_bool func(audiof *TAudiofile) bool

