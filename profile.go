package goodmusic

import "math"

// A Profile represents somebody's music preferences by
// recording all of the users in a Dataset who like
// similar artists.
//
// More specifically, each user x in a Dataset is mapped
// to the number of artists which x likes that the
// Profile's user also likes.
type Profile map[int64]int

// LibraryProfile constructs a user's Profile by running
// a Library through a Dataset, counting how many times
// each user in the Dataset "agrees" with the Library on
// artist preferences.
//
// The minPlays argument specifies how many plays an
// artist in the Library must have to be considered while
// looking for corresponding artists in the Dataset.
func LibraryProfile(l Library, d Dataset, minPlays int) Profile {
	res := Profile{}
	for artistInfo := range d.Artists() {
		plays := l.Plays(artistInfo.Name)
		if plays > minPlays {
			for _, user := range artistInfo.Users {
				res[user]++
			}
		}
	}
	return res
}

// Correlation computes a measure of correlation between
// the given ArtistInfo and the Profile.
// A higher correlation means that the the user with the
// profile is more likely to enjoy the artist.
func (p Profile) Correlation(a ArtistInfo) float64 {
	var dotSum float64
	for _, u := range a.Users {
		dotSum += float64(p[u])
	}
	return dotSum / math.Sqrt(float64(len(a.Users)))
}
