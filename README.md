# GoatPlayer

GoatPlayer is a music library's manager and audiofiles player.
It is written in Go.
It works in the gnome-terminal.
It uses mpv as a backend (via socket)

GOAT = Go + Audiofiles + Terminal

Other commands and libraries you need to run GoatPlayer:

- mpv
- socat
- libtag static library

Install libtag:

`$ sudo apt install libtagc0-dev`

This is a new project. Any help and collaboration will be welcome.
Download and try it out if you want to help me with its development.

All this thing started in February 2026 as a bash script. I did like using mpv from the command line, so I started writing a bash script to manage playlists. After some time I decided that bash was not good for this project; I looked aroung and choose the Go language. I like Go and I think it has been a good choice. 

I use:

- Ubuntu 24.04.5 LTS
- go version go1.25.8 linux/amd64
- Visual Studio Code
- gnome-terminal

I'm not an expert with Ubuntu, nor with the terminal, nor with releasing a software. I like to share what I came to till now. I hope someone will help me in the development. I do believe in this project.

This is a music library's manager similar in some aspects to Rhythmbox: it has columns where you can select (activate) tag-fields. But here there are not just 3 columns, there are 6: Genres, Artists, Albums, Years, Album-artists and Composers. All these columns can be used to filter the audiofiles.

Then, the window where the audiofiles are listed (view) is not the usual one: it GROUPS the audiofiles by artist/album or by composer/album. There is one line with Artist/Album and then, below, there are lines with just the title of the audiofiles (on the right side there are the track number and the duration too). I think this is very nice too see and much better than the usual view. It is similar to how cmus does it, but it is not the same.

## release

In this directory there is the source code of the release version (no logging).

## develop

In this directory there is the source code of the development (debug) version. 

s(It prints a colorful and detailed logfile I can watch in another terminal using `tail -F`. I removed this from the release version.)

## install

Here I put the zip file containing the executable file (goatplayer) and some required text files.

In the Wiki you will find a documentation about the usage and the software.

## Go packages used

- github.com/dhowden/tag
- github.com/wtolson/go-taglib
