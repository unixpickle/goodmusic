package spotifylibrary

import (
	"encoding/json"
	"net/http"
	"strings"
)

const myTracksAPI = "https://api.spotify.com/v1/me/tracks"

// A Library stores information about a user's saved
// music on Spotify.
type Library struct {
	artists []string
	counts  map[string]int
}

// LoadLibrary loads the saved music from a user's
// Spotify account, given an access token.
func LoadLibrary(key string) (*Library, error) {
	artists, errs := requestTracks(key)
	res := &Library{counts: map[string]int{}}
	for artist := range artists {
		if _, ok := res.counts[artist]; !ok {
			res.artists = append(res.artists, artist)
		}
		res.counts[artist]++
	}
	return res, <-errs
}

// Artists returns a list of artist names from the
// user's Spotify library.
func (l *Library) Artists() []string {
	return l.artists
}

// Count returns the number of songs by the given
// artist in the user's Spotify library.
func (l *Library) Count(artist string) int {
	return l.counts[artist]
}

type tracksResponse struct {
	Next  string       `json:"next"`
	Items []tracksItem `json:"items"`
}

type tracksItem struct {
	Track struct {
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
	} `json:"track"`
}

func requestTracks(key string) (artists <-chan string, errs <-chan error) {
	artistsChan := make(chan string, 1)
	errsChan := make(chan error, 1)
	go func() {
		defer close(artistsChan)
		defer close(errsChan)
		req, err := http.NewRequest("GET", myTracksAPI, nil)
		if err != nil {
			errsChan <- err
			return
		}
		req.Header.Set("Authorization", "Bearer "+key)
		for {
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errsChan <- err
				return
			}
			dec := json.NewDecoder(resp.Body)
			var respData tracksResponse
			err = dec.Decode(&respData)
			resp.Body.Close()
			if err != nil {
				errsChan <- err
				return
			}
			for _, item := range respData.Items {
				for _, artist := range item.Track.Artists {
					artistsChan <- canonicalizeAritstName(artist.Name)
				}
			}
			if respData.Next == "" {
				break
			}
			req, err = http.NewRequest("GET", respData.Next, nil)
			if err != nil {
				errsChan <- err
				return
			}
			req.Header.Set("Authorization", "Bearer "+key)
		}
	}()
	return artistsChan, errsChan
}

func canonicalizeAritstName(name string) string {
	return strings.ToLower(name)
}
