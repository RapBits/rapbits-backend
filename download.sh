#!/bin/bash

# folder to dowload media into
mkdir "$1"
cd "$1"

# download youtube video
youtube-dl -o "$1" -f 'bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best' "$2"

# slice video
snippetMP4FileName="$1_snippet.mp4"
ffmpeg -i  "$1.mp4" -ss "$3" -t "$4" "$snippetMP4FileName"

# extract audio
snippetMP3FileName="$1_snippet.mp3"
ffmpeg -i "$snippetMP4FileName" -f mp3 -ab 100000 -vn "$snippetMP3FileName"

# download album cover
curl "$5" > "$1"
ext=''
case $(file -b "$1") in
    *ASCII*) ext='.txt' ;;
    *JPEG*)  ext='.jpg' ;;
    *PDF*)   ext='.pdf' ;;
    *PNG*) ext='.png' ;;
    *SVG*) ext='.svg' ;;
    *GIF*) ext='.gif' ;;
    *BMP*) ext='.bmp' ;;
    *) continue ;;
esac
mv "$1" "$1${ext}"

# upload files to media server 

# delete files from server
cd ../
rm -rf "$1"
