package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
)

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
