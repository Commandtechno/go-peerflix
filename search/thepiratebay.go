package search

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type TPBResult struct {
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	Size     string `json:"size"`
	Seeders  string `json:"seeders"`
	Leechers string `json:"leechers"`
}

func SearchThePirateBay(query string) []*Torrent {
	res, err := http.Get("https://apibay.org/q.php?q=" + url.QueryEscape(query))
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	var tpbResults []*TPBResult
	if err := json.NewDecoder(res.Body).Decode(&tpbResults); err != nil {
		panic(err)
	}

	var torrents []*Torrent
	for _, tpbResult := range tpbResults {
		torrents = append(torrents, &Torrent{
			Name:     tpbResult.Name,
			Size:     prettyBytes(toInt(tpbResult.Size)),
			Seeders:  toInt(tpbResult.Seeders),
			Leechers: toInt(tpbResult.Leechers),
			MagnetURI: func() string {
				return "magnet:?xt=urn:btih:" + tpbResult.InfoHash
			},
		})
	}

	return torrents
}
