package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/docker/docker/api/types"
)

type PullProgress struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type PullResponse struct {
	Status         string       `json:"status"`
	ProgressDetail PullProgress `json:"progressDetail"`
	Progress       string       `json:"progress"`
	ID             string       `json:"id"`
}

// Pull an image
func (d *DockerClient) Pull(image string) error {
	pullStats, err := d.cli.ImagePull(d.ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	downloads := make(map[string]int)
	downloadMaxs := make(map[string]int)
	done := 0
	total := 0

	needsNewline := false
	decoder := json.NewDecoder(pullStats)
	for decoder.More() {
		var progress PullResponse
		decoder.Decode(&progress)

		switch progress.Status {
		case "Downloading":
			downloads[progress.ID] = progress.ProgressDetail.Current
			downloadMaxs[progress.ID] = progress.ProgressDetail.Total
			needsNewline = true
			fmt.Printf("\rDownloading: %v (%v/%v complete)               ", percs(downloads, downloadMaxs), done, total)

		case "Download complete":
			done++
			downloads[progress.ID] = downloadMaxs[progress.ID]

		case "Pulling fs layer":
			total++

		case "Waiting":
		case "Verifying Checksum":
		case "Extracting":

		default:
			if total == 0 && done == total {
				if needsNewline {
					fmt.Println()
					needsNewline = false
				}
				fmt.Println(progress.Status)
			}
		}
	}

	return nil
}

func percs(statuses, maxes map[string]int) string {
	status := ""
	var groups [][]string
	for id, stat := range statuses {
		perc := float64(stat) / float64(maxes[id])

		if perc != 1 {
			groups = append(groups, []string{
				id,
				strconv.Itoa(int(perc * 100)),
			})
		}
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i][0] < groups[j][0]
	})

	for _, group := range groups {
		if status != "" {
			status += " | "
		}
		status += group[1] + "%"
	}
	return status
}

func sum(statuses map[string]int) int {
	sum := 0
	for _, val := range statuses {
		sum += val
	}
	return sum
}

// PullAll images for registered languages
func (d *DockerClient) PullAll() error {
	for _, def := range languageDefs {
		if err := d.Pull(def.Image); err != nil {
			return err
		}
	}
	return nil
}
