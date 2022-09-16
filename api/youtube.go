package api

import (
	"time"

	"github.com/ragtag-archive/kumo/util"
	"google.golang.org/api/youtube/v3"
)

// playlistLastFullFetch keeps track of when the last full fetch was performed
// for each playlist.
var playlistLastFullFetch map[string]time.Time = make(map[string]time.Time)

// GetPlaylistItems returns a list of video IDs in the given playlist. It will
// fetch the full playlist if it hasn't been fetched in the last 72 hours.
func GetPlaylistItems(service *youtube.Service, playlistId string) (util.Set[string], error) {
	// Check if we need to fetch the full playlist
	needFullFetch := playlistLastFullFetch[playlistId].Before(time.Now().Add(-72 * time.Hour))

	videos := make(util.Set[string])
	pageToken := ""
	for {
		// Fetch the next page of videos
		call := service.PlaylistItems.List([]string{"snippet"}).
			PlaylistId(playlistId).
			MaxResults(50)
		if pageToken != "" {
			call.PageToken(pageToken)
		}
		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		// Add the videos to the set
		for _, item := range resp.Items {
			videos[item.Snippet.ResourceId.VideoId] = struct{}{}
		}

		// Check if we need to fetch the full playlist
		if !needFullFetch {
			break
		}

		// Check if we've reached the end of the playlist
		if resp.NextPageToken == "" {
			// Update the last full fetch time
			playlistLastFullFetch[playlistId] = time.Now()
			break
		}

		// Update the page token
		pageToken = resp.NextPageToken
	}

	return videos, nil
}
