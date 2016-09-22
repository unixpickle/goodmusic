package goodmusic

// A Library returns basic information about a user's
// music library.
type Library interface {
	// Artists returns a list of artist names in the
	// music library.
	Artists() []string

	// Plays returns the collective number of times the
	// given artist's songs have been played.
	Plays(artist string) int
}
