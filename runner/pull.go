package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type PullProgress struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

type PullResponse struct {
	Status         string       `json:"status"`
	ProgressDetail PullProgress `json:"progressDetail"`
	Progress       string       `json:"progress"`
	ID             string       `json:"id"`
}

func pullImage(ctx context.Context, cli client.APIClient, image string) error {
	pullStats, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	needsNewline := false
	decoder := json.NewDecoder(pullStats)
	for decoder.More() {
		var progress PullResponse
		decoder.Decode(&progress)

		switch progress.Status {
		case "Downloading":
			needsNewline = true
			fmt.Printf("\rDownloading %s", progress.Progress)

		case "Extracting":
			needsNewline = true
			fmt.Printf("\rExtracting %s", progress.Progress)

		default:
			if needsNewline {
				fmt.Println()
				needsNewline = false
			}
			fmt.Println(progress.Status)
		}
	}

	return nil
}
