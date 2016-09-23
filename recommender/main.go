package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/unixpickle/goodmusic"
	"github.com/unixpickle/goodmusic/ituneslibrary"
	"github.com/unixpickle/goodmusic/spotifylibrary"
)

const (
	ApplicationArg = 1
	DataArg        = 2
	MinCountArg    = 3
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "application lastfm_data min_count")
		fmt.Fprintln(os.Stderr, "\nAvailable applications:")
		fmt.Fprintln(os.Stderr, " spotify (min_count refers to song count)")
		fmt.Fprintln(os.Stderr, " itunes  (min_count refers to play count)")
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	minCount, err := strconv.Atoi(os.Args[MinCountArg])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid min_count:", os.Args[MinCountArg])
		os.Exit(1)
	}

	log.Println("Loading music library...")
	var library goodmusic.Library
	switch os.Args[ApplicationArg] {
	case "spotify":
		token, err := spotifylibrary.Auth()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Spotify authorization failed:", err)
			os.Exit(1)
		}
		library, err = spotifylibrary.LoadLibrary(token)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Spotify load failed:", err)
			os.Exit(1)
		}
	case "itunes":
		library, err = ituneslibrary.LoadLibrary()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to load iTunes library:", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Unknown application:", os.Args[ApplicationArg])
		os.Exit(1)
	}

	log.Println("Loading dataset...")
	dataset, err := goodmusic.LoadLastfmDataset(os.Args[DataArg])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load dataset:", err)
		os.Exit(1)
	}

	log.Println("Building music profile...")
	profile := goodmusic.LibraryProfile(library, dataset, minCount)

	log.Println("Computing artist correlations...")
	var list correlationList
	for artist := range dataset.Artists() {
		if library.Count(artist.Name) > 0 {
			continue
		}
		cor := profile.Correlation(artist)
		list.artists = append(list.artists, artist.Name)
		list.correlations = append(list.correlations, cor)
	}

	log.Println("Sorting correlations...")
	sort.Sort(&list)

	log.Println("Recommending...")
	for i, artist := range list.artists {
		cor := list.correlations[i]
		fmt.Printf("%s (correlation of %.02f%%)\n", artist, cor*100)
	}
}

type correlationList struct {
	artists      []string
	correlations []float64
}

func (c *correlationList) Len() int {
	return len(c.artists)
}

func (c *correlationList) Swap(i, j int) {
	c.artists[i], c.artists[j] = c.artists[j], c.artists[i]
	c.correlations[i], c.correlations[j] = c.correlations[j], c.correlations[i]
}

func (c *correlationList) Less(i, j int) bool {
	return c.correlations[i] > c.correlations[j]
}
