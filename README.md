# musicthread-spotify-importer

Import Spotify playlists into [MusicThread](https://musicthread.app/).

## Build
```bash
go build -o spotify-importer
```

## Run
```bash
./spotify-importer \
  -spotify_client="spotify api client id" \
  -spotify_secret="spotify api client secret" \
  -thread="musicthread thread id" \
  -token="musicthread api token" \
  -playlist="spotify playlist id"
```
