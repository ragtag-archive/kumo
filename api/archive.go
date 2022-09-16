package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ragtag-archive/kumo/util"
)

type RagtagVideo struct {
	ChannelName string `json:"channel_name"`
	ChannelId   string `json:"channel_id"`
	VideoId     string `json:"video_id"`
}

type esSearchResult struct {
	Took     int64 `json:"took"`
	TimedOut bool  `json:"timed_out"`
	Hits     struct {
		Total struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		Hits []struct {
			Index  string      `json:"_index"`
			Id     string      `json:"_id"`
			Score  float64     `json:"_score"`
			Source RagtagVideo `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// GetArchivedItems returns a list of video IDs on Ragtag Archive for a given
// channel ID.
func GetArchivedItems(ctx context.Context, client *http.Client, baseUrl string, channelId string) (util.Set[string], error) {
	// Fetch the channel's videos
	url := fmt.Sprintf("%s/api/v1/search?size=10000&channel_id=%s", baseUrl, channelId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		if err == context.Canceled {
			return nil, err
		}
		return nil, fmt.Errorf("failed to fetch channel videos: %w", err)
	}

	// Read the response body as JSON
	var result esSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Extract the video IDs
	videos := make(util.Set[string])
	for _, hit := range result.Hits.Hits {
		videos[hit.Source.VideoId] = struct{}{}
	}

	return videos, nil
}
