package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/ragtag-archive/kumo/api"
	"github.com/ragtag-archive/kumo/config"
	"github.com/ragtag-archive/kumo/util"
	"github.com/robfig/cron/v3"
	"golang.org/x/sync/semaphore"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var ytService *youtube.Service

// getPlaylistId returns the ID of the playlist for the given channel ID.
func getPlaylistId(channelId string) (string, error) {
	if !strings.HasPrefix(channelId, "UC") {
		return "", fmt.Errorf("Invalid channel ID: %s", channelId)
	}

	return "UU" + channelId[2:], nil
}

func doChannel(
	ctx context.Context, cfg *config.Config, client *http.Client,
	channel *api.Channel, presetName string,
) {
	if !strings.EqualFold(channel.CronPreset, presetName) {
		return
	}

	playlistId, err := getPlaylistId(channel.Id)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[%s] Error getting playlist ID: %s", channel.Name, err)
		}
		return
	}

	// Fetch the channel's videos from YouTube
	youtubeIds, err := api.GetPlaylistItems(ytService, playlistId)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[%s] Error fetching videos: %s", channel.Name, err)
		}
		return
	}

	// Fetch the channel's videos from the archive
	archivedIds, err := api.GetArchivedItems(ctx, client, cfg.Archive.ArchiveURL, channel.Id)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[%s] Error fetching archived videos: %s", channel.Name, err)
		}
		return
	}

	// Find the videos that need to be archived
	idsToArchive := util.SetToSlice(util.SetDifference(youtubeIds, archivedIds))

	log.Printf("[%s] YT: %d, Archive: %d, New: %d",
		channel.Name, len(youtubeIds), len(archivedIds), len(idsToArchive))

	// Add the videos to the archive queue
	sem := semaphore.NewWeighted(cfg.App.MaxConcurrency)
	for _, id := range idsToArchive {
		if err := sem.Acquire(ctx, 1); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[%s] Error acquiring semaphore: %s", channel.Name, err)
			}
			break
		}
		go func(id string) {
			defer sem.Release(1)

			// Send a PUT request to the archive queue
			req, err := http.NewRequestWithContext(
				ctx, http.MethodPut, cfg.Archive.QueueURL, strings.NewReader(id))
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Printf("[%s:%s] Error creating request: %s", channel.Name, id, err)
				}
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Printf("[%s:%s] Error sending request: %s", channel.Name, id, err)
				}
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				if !errors.Is(err, context.Canceled) {
					log.Printf("[%s:%s] Error sending request: %s", channel.Name, id, resp.Status)
				}
				return
			}
		}(id)
	}

	// Wait for all the channels to finish
	if err := sem.Acquire(context.Background(), int64(8)); err != nil {
		log.Printf("Error waiting for channels: %s", err)
		return
	}
}

func runCronJob(cfg *config.Config, ctx context.Context, presetName string) {
	log.Printf("Running cron job %s", presetName)
	defer log.Printf("Finished cron job %s", presetName)

	log.Printf("Fetching YouTube channels")
	client := &http.Client{}
	channels, err := api.GetChannels(ctx, client, cfg.Archive.ChannelsListURL)
	if err != nil {
		log.Printf("Error fetching YouTube channels: %s", err)
		return
	}
	log.Printf("Found %d YouTube channels", len(channels))

	// Process the channels
	sem := semaphore.NewWeighted(cfg.App.MaxConcurrency)
	for _, channel := range channels {
		if err := sem.Acquire(ctx, 1); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[%s] Error acquiring semaphore: %s", channel.Name, err)
			}
			break
		}
		go func(channel api.Channel) {
			defer sem.Release(1)
			doChannel(ctx, cfg, client, &channel, presetName)
		}(channel)
	}

	// Wait for all the channels to finish
	if err := sem.Acquire(context.Background(), int64(8)); err != nil {
		log.Printf("Error waiting for channels: %s", err)
		return
	}
}

func main() {
	log.Println("kumo")

	log.Println("Loading configuration")
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading configuration: %s", err))
	}

	log.Println("Initializing YouTube service")
	ctx, cancel := context.WithCancel(context.Background())
	ytService, err = youtube.NewService(ctx, option.WithAPIKey(cfg.YouTube.ApiKey))
	if err != nil {
		panic(fmt.Errorf("Error initializing YouTube service: %s", err))
	}

	log.Println("Registering cron jobs")
	c := cron.New()
	wg := &sync.WaitGroup{}
	for name, spec := range cfg.CronPresets {
		func(name string) {
			log.Printf("Registering %s: %s", name, spec)
			c.AddFunc(spec, func() {
				wg.Add(1)
				defer wg.Done()
				runCronJob(cfg, ctx, name)
			})
		}(name)
	}

	// Start the cron jobs
	log.Println("Starting cron jobs")
	c.Start()

	log.Println("Preloading channels list")
	channels, err := api.GetChannels(ctx, &http.Client{}, cfg.Archive.ChannelsListURL)
	if err != nil {
		panic(fmt.Errorf("Error preloading channels list: %s", err))
	}
	log.Printf("Found %d channels", len(channels))

	// Wait for a signal to stop
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	log.Println("Ready")
	<-ch

	log.Println("Interrupted, shutting down")
	cancel()
	c.Stop()
	wg.Wait()

	log.Println("Bye")
}
