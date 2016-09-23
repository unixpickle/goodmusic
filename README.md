# goodmusic

This project aims to use music listener databases to find new music for you to enjoy.

# Datasets

Currently, I have been using the [Last.fm 360K](http://www.dtic.upf.edu/~ocelma/MusicRecommendationDataset/lastfm-360K.html) dataset. I plan to use abstractions that allow for other datasets as well.

# Usage

This project currently supports the following music setups:

 * iTunes on Mac OS X
 * Spotify on any operating system

If you use a different music player, you may file an issue on this repository and I (or somebody else) may add support for your setup.

First, download the [Last.fm music dataset](http://www.dtic.upf.edu/~ocelma/MusicRecommendationDataset/lastfm-360K.html). While this is running, make sure you have [Go installed](https://golang.org/doc/install) and a GOPATH configured.

To download the tool, run:

```
$ go get github.com/unixpickle/goodmusic/recommender
```

Now, follow the instructions on [iTunes](#using-itunes) or [Spotify](#using-spotify) below.

## Using iTunes

Run the command as follows, and you will see similar output:

```
$ $GOPATH/bin/recommender itunes /path/to/usersha1-artmbid-artname-plays.tsv 100 | more
2016/09/22 19:35:40 Loading music library...
2016/09/22 19:35:41 Loading dataset...
2016/09/22 19:36:06 Building music profile...
2016/09/22 19:36:06 Computing artist correlations...
2016/09/22 19:36:09 Sorting correlations...
2016/09/22 19:36:09 Recommending...
the killers (correlation of 77279.27%)
the beatles (correlation of 65038.06%)
britney spears (correlation of 63985.37%)
red hot chili peppers (correlation of 63947.08%)
...
```

The last argument (100 in this case) specifies how many plays an artist needs in iTunes in order to be considered. I have found that the default of 100 works fine, since most of my favorite artists have hundreds of plays and I have few artists in iTunes that I don't enjoy.

## Using Spotify

Run the command as follows, and you will see similar output:

```
$ go run *.go spotify ~/Downloads/lastfm-dataset-360K/usersha1-artmbid-artname-plays.tsv 1
2016/09/23 19:12:21 Loading music library...
Please navigate to:
https://accounts.spotify.com/authorize?client_id=dac4cf3243db4bd1b055efc7af84be72&redirect_uri=http%3A%2F%2Flocalhost%3A14505%2Fspotifydone&response_type=token&scope=user-library-read
2016/09/23 19:12:28 Loading dataset...
2016/09/23 19:13:06 Building music profile...
2016/09/23 19:13:06 Computing artist correlations...
2016/09/23 19:13:09 Sorting correlations...
2016/09/23 19:13:09 Recommending...
the killers (correlation of 77279.27%)
the beatles (correlation of 65038.06%)
britney spears (correlation of 63985.37%)
red hot chili peppers (correlation of 63947.08%)
...
```

The last argument (1 in this case) specifies how many songs you must have by a given artist for that artist to be considered.
