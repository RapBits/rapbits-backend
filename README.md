# Rapbits Backend

The backend for rapbits is a mainly a vanilla Go application.

## Tech
- Go
- MySQL - song information is stored here
- youtube-dl - used to download songs
- ffmpeg - used to extract audio from video

## API Endpoints

- `GET` /song/:songId - returns a song where songId is the alphanumeric key for a specific song

Example of request:
```
/song/007ts5zek20iPw02n2997EXI1xSj7z
```

- `GET` /songs?offset=<offset> - returns 30 songs where offset is the offset to begin retrieving the songs from 

Example of request:
```
/songs?offset=30
```

- `POST` /upload - post request to upload song

Example of post body:
```
{
    "songName": "some song",
    "artist": "john doe",
    "url": "https://www.youtube.com/watch?v=123",
    "startTime": "60",
    "endTime": "120",
    "lyric": "foo bar",
    "albumCover": "https://i.ytimg.com/an_webp/dQw4w9WgXcQ/mqdefault_6s.webp?du=3000&sqp=CLzn4egF&rs=AOn4CLABNWTXa2_TJImw8lQYbmDyt1HtpA"
}
```

## Run
To run the project run:

```
go run github.com/rapbits-backend
```


