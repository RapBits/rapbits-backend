# Rapbits Backend

## Tech
- Go
- MySQL - song information is stored here
- youtube-dl - used to download songs
- ffmpeg - used to extract audio from video

## API Endpoints

- `GET` /song/:songId
- `GET` /songs?limit=<limit>&offset=<offset>
- `POST` /upload

## Run
To run the project run:

```
go run github.com/rapbits-backend
```


