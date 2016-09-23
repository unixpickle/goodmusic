package goodmusic

// A Library returns basic information about a user's
// music library.
type Library interface {
	// Artists returns a list of artist names in the
	// music library.
	Artists() []string

	// Count returns some metric of how much the user has
	// listened to the artist.
	// This might be number of plays, some sort of rating,
	// or the number of songs by the artist.
	Count(artist string) int
}
