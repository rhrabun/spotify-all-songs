package main

import (
	"context"
	"log"

	"github.com/zmb3/spotify/v2"
)

func chunkSlice(slice []spotify.ID, chunkSize int) [][]spotify.ID {
	var chunks [][]spotify.ID
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func SyncSongs(client *spotify.Client, ctx context.Context, songsToAdd []spotify.ID, songsToRemove []spotify.ID) {
	if len(songsToAdd) != 0 {
		if len(songsToAdd) > 50 {
			log.Println("Starting chunking")
			chunkedSongsToAdd := chunkSlice(songsToAdd, 49)
			for _, chunk := range chunkedSongsToAdd {
				addErr := client.AddTracksToLibrary(ctx, chunk...)
				if addErr != nil {
					log.Fatal(addErr)
				}
			}
		}
		addErr := client.AddTracksToLibrary(ctx, songsToAdd...)
		if addErr != nil {
			log.Fatal(addErr)
		}

		log.Printf("Added %d songs to All Songs playlist\n", len(songsToAdd))
	}

	if len(songsToRemove) != 0 {
		removeErr := client.RemoveTracksFromLibrary(ctx, songsToRemove...)
		if removeErr != nil {
			log.Fatal(removeErr)
		}

		log.Printf("Removed %d songs from All Songs playlist\n", len(songsToRemove))
	}
}

func getLikedTracks(client *spotify.Client, ctx context.Context) []spotify.ID {
	var likedTracks []spotify.SavedTrack
	res, err := client.CurrentUsersTracks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for page := 1; ; page++ {
		likedTracks = append(likedTracks, res.Tracks...)

		err := client.NextPage(ctx, res)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	var trackIds []spotify.ID
	for _, item := range likedTracks {
		trackIds = append(trackIds, item.ID)
	}

	return trackIds
}

func GetPlaylistTracks(client *spotify.Client, ctx context.Context, playlistId spotify.ID) []spotify.ID {
	var tracks []spotify.PlaylistItem
	res, err := client.GetPlaylistItems(ctx, playlistId)
	if err != nil {
		log.Fatal(err)
	}

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
