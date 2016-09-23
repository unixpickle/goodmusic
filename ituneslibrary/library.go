package ituneslibrary

import (
	"errors"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/DHowett/go-plist"
)

const libraryXMLPath = "Music/iTunes/iTunes Music Library.xml"

// Library returns play information about the iTunes
// library of the current user.
//
// This only works on Mac OS X, and it requires the
// "Share iTunes Library XML with other applications"
// setting to be checked in the advanced preferences
// tab of iTunes.
type Library struct {
	artists []string
	plays   map[string]int
}

// LoadLibrary loads the user's iTunes library in its
// current state.
func LoadLibrary() (*Library, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(user.HomeDir, libraryXMLPath)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	plistData := map[string]interface{}{}
	if _, err := plist.Unmarshal(data, plistData); err != nil {
		return nil, err
	}
	tracks, ok := plistData["Tracks"].(map[string]interface{})
	if !ok {
		return nil, errors.New("library plist: missing or invalid tracks list")
	}
	res := &Library{
		plays: map[string]int{},
	}
	for _, trackObj := range tracks {
		track, ok := trackObj.(map[string]interface{})
		if !ok {
			return nil, errors.New("library plist: unexpected track type")
		}
		artistName, ok1 := track["Artist"].(string)
		playCount, _ := track["Play Count"].(uint64)
		if !ok1 {
			continue
		}
		clean := cleanArtistName(artistName)
		if _, ok := res.plays[clean]; !ok {
			res.artists = append(res.artists, clean)
		}
		res.plays[clean] += int(playCount)
	}
	return res, nil
}

// Artists returns the list of artists in the iTunes
// library, as of when the library was loaded.
func (l *Library) Artists() []string {
	return l.artists
}

// Count returns the total number of times the user has
// played a song by the given artist, as of when the
// library was loaded.
func (l *Library) Count(artist string) int {
	return l.plays[artist]
}

func cleanArtistName(name string) string {
	return strings.ToLower(name)
}
