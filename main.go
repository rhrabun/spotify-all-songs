package main

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/zmb3/spotify/v2"
)

func main() {
	client, ctx := SpotifyAuthenticate()

	playlists := GetPlaylists(client, ctx)

	if playlists != nil {
		// This is more efficient loop to check filter playlists
		// https://stackoverflow.com/a/20551116
		i := 0
		for _, item := range playlists {
			if item.Owner.DisplayName == "Afandi_bobo" && !strings.HasPrefix(item.Name, "_") {
				playlists[i] = item
				i++
			}
		}
		// Rewrite slice in-place to get rid of not needed values
		playlists = playlists[:i]
	} else {
		log.Fatal("Couldn't find any playlists")
	}

	// Get songs from Liked playlist
	likedTracks := getLikedTracks(client, ctx)

	// Get songs from other playlists
	wg := sync.WaitGroup{}
	mut := sync.Mutex{}
	var otherSongs []spotify.ID

	log.Println("Gettings songs from other playlists")
	// For each playlist run concurrent function to extract songs and add them to original slice by pointer
	for _, item := range playlists {
		wg.Add(1)
		go func(client *spotify.Client, ctx context.Context, item spotify.SimplePlaylist, otherSongs *[]spotify.ID) {
			playlistSongs := GetPlaylistTracks(client, ctx, item.ID)

			mut.Lock()
			*otherSongs = append(*otherSongs, playlistSongs...)
			mut.Unlock()

			wg.Done()
		}(client, ctx, item, &otherSongs)
	}
	wg.Wait()

	songsToAdd := compare(otherSongs, likedTracks)
	songsToRemove := compare(likedTracks, otherSongs)

	if len(songsToAdd) == 0 && len(songsToRemove) == 0 {
		log.Println("Nothing to do here")
	} else {
		SyncSongs(client, ctx, songsToAdd, songsToRemove)
	}
}

// difference returns the elements in `a` that aren't in `b`.
func compare(a []spotify.ID, b []spotify.ID) []spotify.ID {
	// Create a map out of slice b, cuz map allows to check if item exists
	mb := make(map[spotify.ID]struct{}, len(b))
	for _, item := range b {
		mb[item] = struct{}{}
	}

	// Adding found items to map to avoid duplicates
	diff := make(map[spotify.ID]struct{})
	for _, item := range a {
		if _, found := mb[item]; !found {
			diff[item] = struct{}{}
		}
	}

	// Now converting to slice
	ids := make([]spotify.ID, len(diff))
	i := 0
	for k := range diff {
		ids[i] = k
		i++
	}

	return ids
}
