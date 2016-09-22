package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/unixpickle/goodmusic"
	"github.com/unixpickle/goodmusic/ituneslibrary"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "lastfm_data min_plays")
		os.Exit(1)
	}

	minPlays, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid min_plays:", os.Args[2])
		os.Exit(1)
	}

	log.Println("Loading music library...")
	library, err := ituneslibrary.LoadLibrary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load iTunes library:", err)
		os.Exit(1)
	}

	log.Println("Loading dataset...")
	dataset, err := goodmusic.LoadLastfmDataset(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load dataset:", err)
		os.Exit(1)
	}

	log.Println("Building music profile...")
	profile := goodmusic.LibraryProfile(library, dataset, minPlays)

	log.Println("Computing artist correlations...")
	var list correlationList
	for artist := range dataset.Artists() {
		if library.Plays(artist.Name) > 0 {
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
