package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/zmb3/spotify/v2"
)

func main() {
	client, ctx := SpotifyAuthenticate()

	playlists := GetPlaylists(client, ctx)

	var allSongsId spotify.ID
	if playlists != nil {
		// This is more efficient loop to check filter playlists
		// https://stackoverflow.com/a/20551116
		i := 0
		for _, item := range playlists {
			if item.Name == "All Songs" {
				fmt.Printf("Found playlist %s with ID: %s\n", item.Name, item.ID)
				allSongsId = item.ID
			} else if item.Owner.DisplayName == "Afandi_bobo" && !strings.HasPrefix(item.Name, "_") {
				playlists[i] = item
				i++
			}
		}
		// Rewrite slice in-place to get rid of not needed values
		playlists = playlists[:i]
	} else {
		log.Fatal("Couldn't find any playlists")
	}

	// Get songs from All songs playlist
	allSongsTracks := GetPlaylistTracks(client, ctx, allSongsId)

	// Get songs from other playlists
	var otherSongs []spotify.ID
	for _, item := range playlists {
		log.Printf("Getting songs from playlist %s\n", item.Name)
		playlistSongs := GetPlaylistTracks(client, ctx, item.ID)
		otherSongs = append(otherSongs, playlistSongs...)
	}

	songsToAdd := compare(otherSongs, allSongsTracks)
	songsToRemove := compare(allSongsTracks, otherSongs)

	if len(songsToAdd) == 0 && len(songsToRemove) == 0 {
		log.Println("Nothing to do here")
	} else {
		SyncSongs(client, ctx, songsToAdd, songsToRemove, allSongsId)
	}
}

func SyncSongs(client *spotify.Client, ctx context.Context, songsToAdd []spotify.ID, songsToRemove []spotify.ID, playlistId spotify.ID) {
	if len(songsToAdd) != 0 {
		_, addErr := client.AddTracksToPlaylist(ctx, spotify.ID(playlistId), songsToAdd...)
		if addErr != nil {
			fmt.Printf("%+v", addErr)
		}

		fmt.Printf("Added %d songs to All Songs playlist\n", len(songsToAdd))
	}

	if len(songsToRemove) != 0 {
		_, removeErr := client.RemoveTracksFromPlaylist(ctx, spotify.ID(playlistId), songsToRemove...)
		if removeErr != nil {
			log.Fatal(removeErr)
		}

		fmt.Printf("Removed %d songs from All Songs playlist\n", len(songsToRemove))
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

func GetPlaylistTracks(client *spotify.Client, ctx context.Context, playlistId spotify.ID) []spotify.ID {
	var tracks []spotify.PlaylistItem
	res, err := client.GetPlaylistItems(ctx, playlistId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Total of %d tracks", res.Total)

	for page := 1; ; page++ {
		tracks = append(tracks, res.Items...)

		err := client.NextPage(ctx, res)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	var trackIds []spotify.ID
	for _, item := range tracks {
		trackIds = append(trackIds, item.Track.Track.ID)
	}

	return trackIds
}

func GetPlaylists(client *spotify.Client, ctx context.Context) []spotify.SimplePlaylist {
	var playlists []spotify.SimplePlaylist

	res, err := client.CurrentUsersPlaylists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Total of %d playlists", res.Total)

	for page := 1; ; page++ {
		playlists = append(playlists, res.Playlists...)

		err := client.NextPage(ctx, res)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	return playlists
}
