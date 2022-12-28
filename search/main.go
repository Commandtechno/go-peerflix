package search

import "sort"

type Torrent struct {
	Name      string
	Size      string
	Seeders   int
	Leechers  int
	MagnetURI func() string
}

func Search(query string) []*Torrent {
	var torrents []*Torrent
	torrents = append(torrents, Search1337x(query)...)
	torrents = append(torrents, SearchThePirateBay(query)...)

	sort.Slice(torrents, func(i, j int) bool {
		// compare seeders
		return torrents[i].Seeders > torrents[j].Seeders
	})

	if len(torrents) > 10 {
		torrents = torrents[:10]
	}

	sort.Slice(torrents, func(i, j int) bool {
		// compare seeder to leecher ratio
		return float64(torrents[i].Seeders)/float64(torrents[i].Leechers) > float64(torrents[j].Seeders)/float64(torrents[j].Leechers)
	})

	return torrents
}
