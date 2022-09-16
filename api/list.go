package api

import (
	"bufio"
	"context"
	"net/http"
	"strings"
	"time"
)

// Channel represents a single YouTube channel
type Channel struct {
	AgencyName string `json:"agency_name"`
	BatchName  string `json:"batch_name"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	CronPreset string `json:"cron_preset"`
}

// Channel cache
var cache []Channel
var cacheLastModified time.Time

// GetChannels returns a list of channels to process
func GetChannels(ctx context.Context, client *http.Client, youtubeChannelsListURL string) ([]Channel, error) {
	// Check if the cache is still valid
	if cacheLastModified.After(time.Now().Add(-5 * time.Minute)) {
		return cache, nil
	}

	// Fetch the channels list
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, youtubeChannelsListURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body line-by-line
	channels := make([]Channel, 0)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into fields
		fields := strings.Split(line, "\t")
		if len(fields) < 5 {
			continue
		}

		// Parse the fields
		channels = append(channels, Channel{
			AgencyName: strings.TrimSpace(fields[0]),
			BatchName:  strings.TrimSpace(fields[1]),
			Id:         strings.TrimSpace(fields[2]),
			Name:       strings.TrimSpace(fields[3]),
			CronPreset: strings.TrimSpace(fields[4]),
		})
	}

	// Update the cache
	cache = channels
	cacheLastModified = time.Now()

	return channels, nil
}
