# Spotify-all-songs

Basically, this is a copy of my Python script that adds/removes songs in All Songs playlist, based on other playlists.
Rewritten in Go in attempt to learn it


### Build
`go install .`

### Run
*Export install dir to PATH via `export PATH=$PATH:$(dirname $(go list -f '{{.Target}}' .))`*

`SPOTIFY_ID="<id>" SPOTIFY_SECRET="<Secret>" spotify-all-songs`


### TODO's
* Try caching auth token to file(the way Python script does it), so it doesn't require to login via browser each time
