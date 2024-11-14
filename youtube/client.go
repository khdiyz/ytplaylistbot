package youtube

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeClient struct {
	Service *youtube.Service
}

func NewYouTubeClient(apiKey string) *YouTubeClient {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create YouTube client: %v", err)
	}
	return &YouTubeClient{Service: service}
}

func (yt *YouTubeClient) GetPlaylistVideos(playlistID string) ([]*youtube.PlaylistItem, error) {
	var videos []*youtube.PlaylistItem
	pageToken := ""

	for {
		call := yt.Service.PlaylistItems.List([]string{"snippet"}).PlaylistId(playlistID).MaxResults(50).PageToken(pageToken)
		response, err := call.Do()
		if err != nil {
			return nil, err
		}
		videos = append(videos, response.Items...)
		if response.NextPageToken == "" {
			break
		}
		pageToken = response.NextPageToken
	}

	return videos, nil
}
