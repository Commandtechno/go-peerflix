package search

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getMagnetURI(url string) func() string {
	return func() string {
		doc, err := loadDocUrl("https://1377x.to" + url)
		if err != nil {
			panic(err)
		}

		return doc.Find("a[href^=\"magnet\"]").AttrOr("href", "")
	}
}

func Search1337x(query string) []*Torrent {
	doc, err := loadDocUrl("https://1377x.to/search/" + url.PathEscape(query) + "/1/")
	if err != nil {
		panic(err)
	}

	var torrents []*Torrent
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		torrents = append(torrents, &Torrent{
			Name:      trimLen(strings.TrimSpace(s.Find(".name").Text()), 100),
			Size:      strings.TrimSpace(s.Find(".size").Text()),
			Seeders:   toInt(s.Find(".seeds").Text()),
			Leechers:  toInt(s.Find(".leeches").Text()),
			MagnetURI: getMagnetURI(s.Find("a[href^=\"/torrent/\"]").AttrOr("href", "")),
		})
	})

	return torrents
}
