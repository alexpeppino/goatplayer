package musicscanner // Single-instance package

import (
	"fmt"
	"os"
	"path/filepath"

	tag "github.com/dhowden/tag"
)

var (
	FoldersN  uint16
	DirsN     uint16
	nMp3      uint16
	nM4a      uint16
	nOgg      uint16
	FileList  []string // all audiofiles found
	ValidDirs []TValidDir

	ValidFoldersN     uint16   // deprecated
	ValidFolders      []string // deprecated
	ValidFolderFilesN []uint16 // deprecated
	Ext               map[string]uint16
)

// go get github.com/wtolson/go-taglib

// Scan music folder recursively. List all folders containing audio files recording the number of files found.
// List all audio files.
// Get metadata for each audio file.

// Just checks if ext is a valid extension.
func IsValidExtension(s string) bool {
	_, ok := Ext[s]
	return ok
}

func InitMusicScanner() {
	Ext = make(map[string]uint16)
	Ext[".mp3"] = 0
	Ext[".m4a"] = 0
	Ext[".ogg"] = 0
	Ext[".aac"] = 0
}

func Reset() {
	ValidDirs = nil
	DirsN = 0
	FoldersN = 0
	nMp3 = 0
	nM4a = 0
	nOgg = 0
	ValidFoldersN = 0
	ValidFolders = nil
	ValidFolderFilesN = nil
	// fileList = nil
	// for v := range ext {
	// 	ext[v] = 0
	// }
}

// Scans dir and subdirs, fills the list of dirs (ValidDirs) even with dirs that have no audiofiles.
// musicscanner.ValidDirs => updated, musicscanner.nDirs => updated
func FindValidDirs(dir string, subdir string, level uint16) {
	var findValidDirs func(dir string, subdir string, level uint16) uint16
// trace.Begin("musicscanner.FindValidDirs") // t //
	// Recursive function.
	findValidDirs = func(dir string, subdir string, level uint16) uint16 {
		if level == 0 {
// trace.Print("ScanForValidDirs, dir: %s", dir) // t //
		}
		var nFilesInDir uint16
		var path string
		if subdir != "" {
			path = dir + "/" + subdir
		} else {
			path = dir
		}
		// fmt.Printf("Scanning %s...\n", dir)
		entries, err := os.ReadDir(path)
		check(err)
		var filename string
		//  Checks if the directory contains audio files.
		for _, entry := range entries {
			ext := filepath.Ext(entry.Name())
			if IsValidExtension(ext) {
				Ext[ext]++
				nFilesInDir++
			}
			// fmt.Println(ent.Type(), ent.Name(), typ.String()[0], ext)
		}
		var nnn int
		ValidDirs = append(ValidDirs, TValidDir{path, subdir, nFilesInDir, level})
		nnn = len(ValidDirs) - 1
		// trace.Print("dir %s %d", path, sca.ValidDirs[nnn].n)

		//  Checks for subdirectories and scans them.
		for _, entry := range entries {
			// typ := ent.Type()
			filename = entry.Name()
			if entry.IsDir() {
				DirsN++
				// fmt.Println("is dir")
				ValidDirs[nnn].N += findValidDirs(path, filename, level+1)
			}
		}
		// trace.Print("%v", sca.ValidDirs)
		// trace.Print("dir %s %d", path, sca.ValidDirs[nnn].n)
		return ValidDirs[nnn].N
		// sort.Slice(sca.ValidDirs, func(i, j int) bool {
		// 	return strings.ToLower(sca.ValidDirs[j].dir) > strings.ToLower(sca.ValidDirs[i].dir)
		// })
	}
	findValidDirs(dir, subdir, level)
// trace.EndAdd("musicscanner.FindValidDirs, found %d dirs, %v", len(ValidDirs), Ext) // t //
}

func WriteDirsToFile() {
	fd, _ := os.OpenFile("a/validDirs.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	// trace.Err(err)
	defer fd.Close()
	for i := range ValidDirs {
		v := ValidDirs[i]
		fmt.Fprintf(fd, "%-120s%-80s%03d\n", v.DirPath, v.Basename, v.N)
	}
}

func WriteFilelistToFile() {
	fd, _ := os.OpenFile("a/validFiles.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	// trace.Err(err)
	defer fd.Close()
	for i := range FileList {
		v := FileList[i]
		fmt.Fprintf(fd, "%s\n", v)
	}
}

// Searches 'dirPath' recursively for all audiofiles.
// musicscanner.fileList => updated
func FindFiles(dirPath string) {
// trace.BeginAdd("musicscanner.FindFiles", "dirPath: '%s'", dirPath) // t //
	var scanForFiles func(dirPath string)
	FileList = nil

	/// The recursive function
	scanForFiles = func(dirPath string) {
		var nFilesInDir uint16
		// fmt.Printf("Scanning %s...\n", dir)
		entries, err := os.ReadDir(dirPath)
		// trace.Print("dirPath is %s, %d entries found", dirPath, len(entries))
		check(err)
		var entryName string
		for _, entry := range entries {
			// typ := ent.Type()
			entryName = entry.Name()
			entryPath := dirPath + "/" + entryName
			// info, _ := entry.Info()
			// if info.Mode()&os.ModeSymlink != 0 {
			// 	target, _ := filepath.EvalSymlinks(entryPath)
			// 	trace.Print("found symlink: %s -> %s", entryPath, target)
			// 	// entryPath, _ = filepath.EvalSymlinks(entryPath)
			// }
			ext := filepath.Ext(entryName)
			// trace.Print("%s, entry: %s (%s)", dirPath, entryName, ext)
			if entry.IsDir() {
				DirsN++
				// fmt.Println("is dir")
				scanForFiles(entryPath)
			} else if IsValidExtension(ext) {
				Ext[ext]++
				FileList = append(FileList, entryPath)
				nFilesInDir++
			}
			// fmt.Println(ent.Type(), ent.Name(), typ.String()[0], ext)
		}
		if nFilesInDir > 0 {
			ValidFolders = append(ValidFolders, dirPath)
			ValidFolderFilesN = append(ValidFolderFilesN, nFilesInDir)
			ValidFoldersN++
		}
	}

	///
	scanForFiles(dirPath)
// trace.EndAdd("musicscanner.FindFiles, found %d files", len(FileList)) // t //
}

// https://dev.to/moseeh_52/understanding-osstat-vs-oslstat-in-go-file-and-symlink-handling-3p5d

// Contains all metadata retrieved by 'tag.ReadFrom'.
var MetaTags tag.Metadata

// Sets 'MetaTags' for the given filename after calling 'tag.ReadFrom'.
func GetMetadata(filename string) {
	fd, _ := os.Open(filename)
	// log.Println(err)
	MetaTags, _ = tag.ReadFrom(fd)
	// log.Println(err, MetaTags)
}

func Metadata(file *os.File) {
	meta, err := tag.ReadFrom(file)
	if err != nil {
		return
	}
	fmt.Println(meta.Title())
	// fmt.Println(meta.Album())
	// fmt.Println(meta.Artist())
	// fmt.Println(meta.Track())
}

// func Test () {

//		ScanForValidDirs("/home/alex/Music")
//		for i, v := range ValidFolders {
//			fmt.Println(i, v)
//		}
//		fmt.Println("valid folders:", validFolders_n)
//		fmt.Println("folders:", folders_n)
//		fmt.Println("mp3 tot:", mp3_tot)
//		fmt.Println("m4a tot:", m4a_tot)
//		fmt.Println("ogg tot:", ogg_tot)
//		fmt.Println("audiofiles tot:", mp3_tot + m4a_tot + ogg_tot)
//	}

type TValidDir struct {
	DirPath  string
	Basename string
	N        uint16
	Level    uint16
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

