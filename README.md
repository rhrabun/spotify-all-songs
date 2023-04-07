# Spotify-all-songs

Basically, this is a copy of my script that adds/removes songs in All Songs playlist, based on other playlists.
This is my attempt to learn Go


### Build
`go install .`

### Run
*Export install dir to PATH via `export PATH=$PATH:$(dirname $(go list -f '{{.Target}}' .))`*

`SPOTIFY_ID="<id>" SPOTIFY_SECRET="<Secret>" go run spotify-all-songs`


### TODO's
* Try caching auth token to file(the way Python script does it), so it doesn't require to login via browser each time
* Add concurrent execution for getting songs from each playlist to try to reduce execution time
