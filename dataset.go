package goodmusic

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ArtistInfo includes metadata and listener info about
// a musical artist.
type ArtistInfo struct {
	// Name is the human-readable artist name.
	Name string

	// Users is a list of user IDs who like this artist.
	Users []int64
}

// Dataset contains music listener data.
type Dataset interface {
	// Artists returns a new channel of all the artists
	// in the dataset.
	Artists() <-chan ArtistInfo
}

type lastfmDataset struct {
	artists map[string][]int64
}

// LoadLastfmDataset loads the dataset from
// http://www.dtic.upf.edu/~ocelma/MusicRecommendationDataset/lastfm-360K.html
// using the file named usersha1-artmbid-artname-plays.tsv.
func LoadLastfmDataset(playsTSVPath string) (Dataset, error) {
	f, err := os.Open(playsTSVPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	res := &lastfmDataset{artists: map[string][]int64{}}

	var userId int64
	var userHash string
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 4 {
			return nil, fmt.Errorf("expected %d columns but got %d", 4, len(parts))
		}
		if userHash != parts[0] {
			userId++
			userHash = parts[0]
		}
		res.artists[parts[2]] = append(res.artists[parts[2]], userId)
	}

	return res, nil
}

func (l *lastfmDataset) Artists() <-chan ArtistInfo {
	res := make(chan ArtistInfo, 1)

	go func() {
		for artist, ids := range l.artists {
			res <- ArtistInfo{Name: artist, Users: ids}
		}
		close(res)
	}()

	return res
}
