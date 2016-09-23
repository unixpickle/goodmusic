# goodmusic

This project aims to use music listener databases to find new music for you to enjoy.

# Datasets

Currently, I have been using the [Last.fm 360K](http://www.dtic.upf.edu/~ocelma/MusicRecommendationDataset/lastfm-360K.html) dataset. I plan to use abstractions that allow for other datasets as well.

# Usage

If you have a Mac and use iTunes as your music player, you are all set. Otherwise, file an issue on this repository and I (or somebody else) may add support for your setup.

First, download the [Last.fm music dataset](http://www.dtic.upf.edu/~ocelma/MusicRecommendationDataset/lastfm-360K.html). While this is running, make sure you have [Go installed](https://golang.org/doc/install) and a GOPATH configured.

To download the tool, run:

```
$ go get github.com/unixpickle/goodmusic/recommender
```

Now you can run the command as follows, and you will see similar output:

```
$ $GOPATH/bin/recommender /path/to/usersha1-artmbid-artname-plays.tsv 100 | more
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
