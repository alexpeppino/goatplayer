package main

import (
	"fmt"
	"os"
	"path/filepath"

	tag "github.com/dhowden/tag"
)

// go get github.com/wtolson/go-taglib

// Scan music folder recursively. List all folders containing audio files recording the number of files found.
// List all audio files.
// Get metadata for each audio file.

// Just checks if ext is a valid extension.
func (this *cMusicScanner) IsValidExtension(s string) bool {
	_, ok := this.ext[s]
	return ok
}

func (this *cMusicScanner) InitMusicScanner() {
	this.ext = make(map[string]uint16)
	this.ext[".mp3"] = 0
	this.ext[".m4a"] = 0
	this.ext[".ogg"] = 0
	this.ext[".aac"] = 0
}

func (this *cMusicScanner) Reset() {
	this.validDirs = nil
	this.nDirs = 0
	this.nFolders = 0
	this.nMp3 = 0
	this.nM4a = 0
	this.nOgg = 0
	this.nValidFolders = 0
	this.validFolders = nil
	this.nValidFolderFiles = nil
	// this.fileList = nil
	// for v := range this.ext {
	// 	this.ext[v] = 0
	// }
}

// Scans dir and subdirs, fills the list of dirs (ValidDirs) even with dirs that have no audiofiles.
// MusicScanner.validDirs => updated, MusicScanner.nDirs => updated
func (this *cMusicScanner) FindValidDirs(dir string, subdir string, level uint16) {
	var findValidDirs func(dir string, subdir string, level uint16) uint16
// trace.Begin("MusicScanner.FindValidDirs") // t //
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
			if this.IsValidExtension(ext) {
				this.ext[ext]++
				nFilesInDir++
			}
			// fmt.Println(ent.Type(), ent.Name(), typ.String()[0], ext)
		}
		var nnn int
		MusicScanner.validDirs = append(MusicScanner.validDirs, tValidDir{path, subdir, nFilesInDir, level})
		nnn = len(MusicScanner.validDirs) - 1
		// trace.Print("dir %s %d", path, sca.ValidDirs[nnn].n)

		//  Checks for subdirectories and scans them.
		for _, entry := range entries {
			// typ := ent.Type()
			filename = entry.Name()
			if entry.IsDir() {
				MusicScanner.nDirs++
				// fmt.Println("is dir")
				MusicScanner.validDirs[nnn].n += findValidDirs(path, filename, level+1)
			}
		}
		// trace.Print("%v", sca.ValidDirs)
		// trace.Print("dir %s %d", path, sca.ValidDirs[nnn].n)
		return MusicScanner.validDirs[nnn].n
		// sort.Slice(sca.ValidDirs, func(i, j int) bool {
		// 	return strings.ToLower(sca.ValidDirs[j].dir) > strings.ToLower(sca.ValidDirs[i].dir)
		// })
	}
	findValidDirs(dir, subdir, level)
// trace.EndAdd("MusicScanner.FindValidDirs, found %d dirs, %v", len(MusicScanner.validDirs), this.ext) // t //
}

func (this *cMusicScanner) WriteDirsToFile() {
	fd, _ := os.OpenFile("a/validDirs.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	// trace.Err(err)
	defer fd.Close()
	for i := range this.validDirs {
		v := this.validDirs[i]
		fmt.Fprintf(fd, "%-120s%-80s%03d\n", v.dirPath, v.basename, v.n)
	}
}

func (this *cMusicScanner) WriteFilelistToFile() {
	fd, _ := os.OpenFile("a/validFiles.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	// trace.Err(err)
	defer fd.Close()
	for i := range this.fileList {
		v := this.fileList[i]
		fmt.Fprintf(fd, "%s\n", v)
	}
}

// Searches 'dirPath' recursively for all audiofiles.
// MusicScanner.fileList => updated
func (this *cMusicScanner) FindFiles(dirPath string) {
// trace.BeginAdd("MusicScanner.FindFiles", "dirPath: '%s'", dirPath) // t //
	var scanForFiles func(dirPath string)
	this.fileList = nil

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
				MusicScanner.nDirs++
				// fmt.Println("is dir")
				scanForFiles(entryPath)
			} else if this.IsValidExtension(ext) {
				this.ext[ext]++
				MusicScanner.fileList = append(MusicScanner.fileList, entryPath)
				nFilesInDir++
			}
			// fmt.Println(ent.Type(), ent.Name(), typ.String()[0], ext)
		}
		if nFilesInDir > 0 {
			MusicScanner.validFolders = append(MusicScanner.validFolders, dirPath)
			MusicScanner.nValidFolderFiles = append(MusicScanner.nValidFolderFiles, nFilesInDir)
			MusicScanner.nValidFolders++
		}
	}

	///
	scanForFiles(dirPath)
// trace.EndAdd("MusicScanner.FindFiles, found %d files", len(this.fileList)) // t //
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

// 	ScanForValidDirs("/home/alex/Music")
// 	for i, v := range ValidFolders {
// 		fmt.Println(i, v)
// 	}
// 	fmt.Println("valid folders:", validFolders_n)
// 	fmt.Println("folders:", folders_n)
// 	fmt.Println("mp3 tot:", mp3_tot)
// 	fmt.Println("m4a tot:", m4a_tot)
// 	fmt.Println("ogg tot:", ogg_tot)
// 	fmt.Println("audiofiles tot:", mp3_tot + m4a_tot + ogg_tot)
// }

type cMusicScanner struct {
	nFolders  uint16
	nDirs     uint16
	nMp3      uint16
	nM4a      uint16
	nOgg      uint16
	fileList  []string // all audiofiles found
	validDirs []tValidDir

	nValidFolders     uint16   // deprecated
	validFolders      []string // deprecated
	nValidFolderFiles []uint16 // deprecated
	ext               map[string]uint16
}

type tValidDir struct {
	dirPath  string
	basename string
	n        uint16
	level    uint16
}

